package core

import (
	"flag"
	"fmt"
)

// Interface implemented by service providers
type ServiceProvider interface {
	Register(c Container)
}

// Interface implemented by bootable providers
type BootableProvider interface {
	// Provider method called by application on bootstrap
	Boot(container Container) error
}

// Interface implemented by Shutdown providers
type ShutdownProvider interface {
	// Provider method called by application on Shutdown phase
	Shutdown(container Container)
}

// Interface implemented by Reconfigure providers
type ReconfigurableProvider interface {
	// Provider method called by SIGHUP signal
	Reconfigure(container Container)
}

type ExtendFn func(old interface{}, c Container) interface{}

// Service container
type Container interface {
	Set(key string, val interface{})
	Factory(key string, fn func(c Container) interface{})
	Protect(key string, fn func(c Container) interface{})
	Has(key string) bool
	Get(key string) (interface{}, error)
	MustGet(key string) interface{}
	Extend(key string, fn ExtendFn) error
	MustExtend(key string, fn ExtendFn)
}

// Application entry point
type App struct {
	providers []interface{}
	items     map[string]interface{}
	instances map[string]interface{}
	protected map[string]bool
	factories map[string]bool
}

// Create application instance
func NewApp() *App {
	return &App{
		items:     make(map[string]interface{}),
		instances: make(map[string]interface{}),
		protected: make(map[string]bool),
		factories: make(map[string]bool),
	}
}

// Register service
func (a *App) Set(key string, val interface{}) {
	a.items[key] = val
}

// Register factory
func (a *App) Factory(key string, fn func(c Container) interface{}) {
	a.factories[key] = true
	a.Set(key, fn)
}

// Register protected value
func (a *App) Protect(key string, fn func(c Container) interface{}) {
	a.protected[key] = true
	a.Set(key, fn)
}

// Returns true if app has specified service or factory. False otherwise.
func (a *App) Has(key string) bool {
	_, ok := a.items[key]

	return ok
}

// Returns specified service or factory
func (a *App) Get(key string) (interface{}, error) {
	item, ok := a.items[key]
	if !ok {
		return nil, fmt.Errorf("identifier '%s' is not defined", key)
	}

	var obj interface{}
	if a.isServiceDefinition(item) {
		itemFn := item.(func(c Container) interface{})
		protected := a.isProtected(key)
		if protected {
			obj = itemFn
		} else if instance, exists := a.instances[key]; exists {
			obj = instance
		} else {
			obj = itemFn(a)
			if !a.isFactory(key) {
				a.instances[key] = obj
			}
		}
	} else {
		obj = item
	}

	return obj, nil
}

// Returns specified service or factory. Panics if it's not defined
func (a *App) MustGet(key string) interface{} {
	val, err := a.Get(key)
	if err != nil {
		panic(err)
	}

	return val
}

// Extends specified entry
func (a *App) Extend(key string, fn ExtendFn) error {
	orig, exists := a.items[key]
	if !exists {
		return fmt.Errorf("identifier '%s' is not defined", key)
	}

	if !a.isServiceDefinition(orig) {
		return fmt.Errorf("idenfier '%s' does not contains service definition", key)
	}

	callable := orig.(func(c Container) interface{})
	a.items[key] = func(c Container) interface{} {
		return fn(callable(a), a)
	}

	return nil
}

// Extends specified entry. Panics on error
func (a *App) MustExtend(key string, fn ExtendFn) {
	err := a.Extend(key, fn)
	if err != nil {
		panic(err)
	}
}

// Register provider
func (a *App) Register(provider interface{}) {
	sProvider, suitable := provider.(ServiceProvider)
	if suitable {
		sProvider.Register(a)
	}

	a.providers = append(a.providers, provider)
}

// Bootstrap application
func (a *App) Boot() error {
	flag.Parse()

	for _, provider := range a.providers {
		bProvider, suitable := provider.(BootableProvider)
		if suitable {
			err := bProvider.Boot(a)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Shutdown application
func (a *App) Shutdown() {
	// in reverse order
	for i := len(a.providers) - 1; i >= 0; i-- {
		provider := a.providers[i]

		sProvider, suitable := provider.(ShutdownProvider)
		if suitable {
			sProvider.Shutdown(a)
		}
	}
}

// Reconfigure application
func (a *App) Reconfigure() error {
	flag.Parse()

	for _, provider := range a.providers {
		rProvider, suitable := provider.(ReconfigurableProvider)
		if suitable {
			rProvider.Reconfigure(a)
		}
	}

	return nil
}

func (a *App) isServiceDefinition(val interface{}) bool {
	_, ok := val.(func(c Container) interface{})

	return ok
}

func (a *App) isProtected(key string) bool {
	_, ok := a.protected[key]

	return ok
}

func (a *App) isFactory(key string) bool {
	_, ok := a.factories[key]

	return ok
}
