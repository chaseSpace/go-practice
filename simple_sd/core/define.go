package core

import (
	"context"
	"errors"
	"fmt"
)

type ServiceDiscovery interface {
	Name() string
	Register(instance ServiceInstance) error
	Deregister(service, addr string) error
	Discovery(ctx context.Context, service string, lastHash string) ([]ServiceInstance, string, error)
}

// ServiceInstance 表示注册的单个实例
type ServiceInstance struct {
	Service  string
	IsUDP    bool // TCP by default
	Host     string
	Port     int
	Metadata map[string]string

	fails int
}

func (s ServiceInstance) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s ServiceInstance) Check() error {
	if s.Service == "" {
		return errors.New("ServiceInstance must have service name")
	}
	if s.Host == "" || s.Port < 1 {
		return errors.New("ServiceInstance must have valid address and port")
	}
	return nil
}