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

func (s *SimpleSd) Discovery(ctx context.Context, service string, lastHash string, block bool) ([]ServiceInstance, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Millisecond * 500): // refresh interval
			s.mu.RLock()
			serv := s.store[service]
			s.mu.RUnlock()

			if serv == nil {
				if lastHash == "" || !block {
					Sdlogger.Debug("SimpleSd: returned with lastHash:%s and block:%v", lastHash, block)
					return nil, nil
				}
			} else {
				list := serv.InstanceList()
				if lastHash == "" || lastHash != calInstanceHash(list) || !block {
					Sdlogger.Debug("SimpleSd: returned with lastHash:%s and block:%v", lastHash, block)
					return list, nil
				}
			}
		}
	}
}

type Service struct {
	service string
	mu      sync.RWMutex
	smap    map[string]*ServiceInstance
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
			for _, ins := range s.InstanceList() {
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
	return nil
}

func (s *Service) Remove(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.smap[addr] == nil {
		return fmt.Errorf("instance: %s not found", addr)
	}
	delete(s.smap, addr)
	return nil
}

func (s *Service) InstanceList() (list []ServiceInstance) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, instance := range s.smap {
		list = append(list, *instance)
	}
	return
}
