package resource

import (
	"sync"
	"time"

	"github.com/alileza/gebet/config"
	"github.com/alileza/gebet/resource/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidType   = errors.New("invaid type")
	ErrNotFound      = errors.New("not")
	ErrInvalidParams = errors.New("invalid")
)

type Resource interface{}

func HTTP(i Resource) *http.Client {
	return i.(*http.Client)
}

type Manager struct {
	resources []*config.Resource
	cache     sync.Map
}

func NewManager(cfgs []*config.Resource) *Manager {
	return &Manager{resources: cfgs}
}

func (mgr *Manager) Get(name string) (Resource, error) {
	for _, resource := range mgr.resources {
		if resource.Name == name {
			client, ok := mgr.cache.Load(resource)
			if ok {
				return client, nil
			}

			switch resource.Type {
			case "http":
				return mgr.http(resource)
			default:
				return nil, ErrInvalidType
			}
		}
	}
	return nil, ErrNotFound
}

func (mgr *Manager) http(cfg *config.Resource) (interface{}, error) {
	opts := &http.Options{}
	for key, val := range cfg.Params {
		switch key {
		case "base_url":
			opts.BaseURL = val
		case "timeout":
			timeout, err := time.ParseDuration(val)
			if err != nil {
				return nil, errors.Wrap(err, "timeout: get http client, invalid params value")
			}
			opts.Timeout = timeout
		default:
			return nil, errors.New(key + ": invalid params")
		}
	}
	client := http.New(opts)

	mgr.cache.Store(cfg, client)

	return client, nil
}