package locker

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

type Locker struct {
	lock   *consul.Lock
	config *Config
}

type Config struct {
	ConsulAddress string
	Key           string
	Value         string
}

func NewLocker(c *Config) (*Locker, error) {
	consulConfig := consul.DefaultConfig()
	consulConfig.Address = c.ConsulAddress

	consulClient, err := consul.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("error while creating consul client for '%s': %w", c.ConsulAddress, err)
	}

	options := consul.LockOptions{
		Key:   c.Key,
		Value: []byte(c.Value),
	}
	consulLock, err := consulClient.LockOpts(&options)
	if err != nil {
		return nil, fmt.Errorf("error while creating consul lock for key '%s': %w", c.Key, err)
	}

	return &Locker{
		lock:   consulLock,
		config: c,
	}, nil
}

func (l *Locker) Lock(abortCh <-chan struct{}) (<-chan struct{}, error) {
	resultCh, err := l.lock.Lock(abortCh)
	if err != nil {
		return nil, fmt.Errorf("error while locking key '%s': %w", l.config.Key, err)
	}
	return resultCh, nil
}

func (l *Locker) Unlock() error {
	err := l.lock.Unlock()
	if err != nil {
		return fmt.Errorf("error while unlocking key '%s': %w", l.config.Key, err)
	}
	return nil
}
