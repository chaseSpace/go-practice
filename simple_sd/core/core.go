package core

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type SimpleSd struct {
	mu               sync.RWMutex
	store            map[string]*Service
	newServiceNotify sync.Cond
}

var Sd ServiceDiscovery = &SimpleSd{store: map[string]*Service{}}

func (s *SimpleSd) Name() string {
	return "SimpleSd"
}

var (
	ErrInstanceNotRegistered = errors.New("instance not registered")
)

const DiscoveryInterval = time.Millisecond * 200

func (s *SimpleSd) Register(instance ServiceInstance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	serv := s.store[instance.Service]
	if serv == nil {
		serv = newService(instance.Service)
		s.store[instance.Service] = serv
	}
	err := serv.Add(&instance)
	return err
}

func (s *SimpleSd) Deregister(service, addr string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if serv := s.store[service]; serv != nil {
		return serv.Remove(addr)
	}
	return errors.Wrap(ErrInstanceNotRegistered, fmt.Sprintf("service: %s, addr: %s", service, addr))
}

func (s *SimpleSd) Discovery(ctx context.Context, service string, lastHash string) ([]ServiceInstance, string, error) {
	for {
		select {
		case <-ctx.Done():
			instances, currHash := s.getInstances(service)
			return instances, currHash, nil

		case <-time.After(DiscoveryInterval): // refresh interval
			instances, currHash := s.getInstances(service)
			if currHash != lastHash {
				return instances, currHash, nil
			}
		}
	}
}

func (s *SimpleSd) getInstances(service string) ([]ServiceInstance, string) {
	s.mu.RLock()
	serv := s.store[service]
	s.mu.RUnlock()

	if serv == nil {
		return nil, ""
	} else {
		list, currHash := serv.InstanceList()
		return list, currHash
	}
}

type Service struct {
	service string
	mu      sync.RWMutex
	smap    map[string]*ServiceInstance

	nowHash string
	quit    chan struct{}
}

const healthCheckInterval = time.Second * 5
const healthCheckMaxFails = 2

func newService(service string) *Service {
	s := &Service{service: service, smap: make(map[string]*ServiceInstance)}
	go s.healthCheck()
	return s
}

func (s *Service) Stop() {
	close(s.quit)
}

func (s *Service) healthCheck() {
	for {
		select {
		case <-s.quit:
			Sdlogger.Debug("Service: %s health check stopped", s.service)
			return
		case <-time.After(healthCheckInterval):
			__inst, _ := s.InstanceList()
			for _, ins := range __inst {
				if ins.IsUDP {
					pass := netDialTest(ins.Addr(), 3, time.Second)
					if pass {
						if ins.fails > 0 {
							ins.fails = 0
							Sdlogger.Debug("Service: %s, instance %s health recovered", s.service, ins.Addr())
						}
					} else {
						ins.fails++
						if ins.fails >= healthCheckMaxFails {
							// remove instance
							_ = s.Remove(ins.Addr())
							Sdlogger.Debug("Service: %s, instance %s was removed", s.service, ins.Addr())
						}
					}
				}
			}
		}
	}
}

func (s *Service) Add(instance *ServiceInstance) error {
	s.mu.Lock()
	s.smap[instance.Addr()] = instance
	s.mu.Unlock()
	s.resetHash()
	return nil
}

func (s *Service) Remove(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.smap[addr] == nil {
		return fmt.Errorf("instance: %s not found", addr)
	}
	delete(s.smap, addr)
	s.resetHash()
	return nil
}

func (s *Service) resetHash() {
	var __inst []ServiceInstance
	for _, instance := range s.smap {
		__inst = append(__inst, *instance)
	}
	s.nowHash = CalInstanceHash(__inst)
}

func (s *Service) InstanceList() (list []ServiceInstance, hash string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, instance := range s.smap {
		list = append(list, *instance)
	}
	hash = s.nowHash
	return
}
