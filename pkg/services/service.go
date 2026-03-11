package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"theSone/internal/application"
)

type Service interface {
	Init() error
	Run(ctx context.Context) error
	Stop() error
}

type Manager struct {
	log      application.Logger
	services []Service
}

func NewManager(log application.Logger) *Manager {
	return &Manager{log: log}
}

func (m *Manager) AddService(services ...Service) {
	m.services = append(m.services, services...)
}

func (m *Manager) Run(ctx context.Context) error {
	m.log.Info("Starting services...")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, s := range m.services {
		if err := s.Init(); err != nil {
			m.log.Error("Failed to initialize service", "error", err)
			m.Stop()
			return err
		}
		go func(svc Service) {
			if err := svc.Run(ctx); err != nil {
				m.log.Error("Service stopped with error", "error", err)
				cancel()
			}
		}(s)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		m.log.Info("Shutdown signal received")
	case <-ctx.Done():
		m.log.Info("Context cancelled")
	}

	m.Stop()
	return nil
}

func (m *Manager) Stop() {
	m.log.Info("Stopping services...")
	for _, s := range m.services {
		if err := s.Stop(); err != nil {
			m.log.Error("Error stopping service", "error", err)
		}
	}
}
