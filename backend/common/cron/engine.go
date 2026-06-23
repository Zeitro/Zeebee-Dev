package cron

import (
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

type Engine struct {
	engine  *cron.Cron
	mu      sync.RWMutex
	taskMap map[string]cron.EntryID
}

func NewEngine() *Engine {
	return &Engine{
		engine:  cron.New(cron.WithSeconds()),
		taskMap: make(map[string]cron.EntryID),
	}
}

func (e *Engine) Start() {
	e.engine.Start()
}

func (e *Engine) Stop() {
	e.engine.Stop()
}

// Register registers a named task. Returns error on duplicate name or invalid spec.
func (e *Engine) Register(name string, spec string, cmd func()) (err error) {
	if err = validateTask(name, spec, cmd); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	return e.registerLocked(name, spec, cmd)
}

// Unregister removes a named task. Returns error if the task does not exist.
func (e *Engine) Unregister(name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.unregisterLocked(name)
}

// IsExisted reports whether a named task is registered.
func (e *Engine) IsExisted(name string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, existed := e.taskMap[name]
	return existed
}

// registerLocked adds the task; caller must hold e.mu.
func (e *Engine) registerLocked(name string, spec string, cmd func()) (err error) {
	if _, existed := e.taskMap[name]; existed {
		return fmt.Errorf("task %s already exists", name)
	}
	id, err := e.engine.AddFunc(spec, cmd)
	if err != nil {
		return fmt.Errorf("task %s: invalid spec %q: %w", name, spec, err)
	}
	e.taskMap[name] = id
	return nil
}

// unregisterLocked removes the task; caller must hold e.mu.
func (e *Engine) unregisterLocked(name string) error {
	id, existed := e.taskMap[name]
	if !existed {
		return fmt.Errorf("task %s does not exist", name)
	}
	e.engine.Remove(id)
	delete(e.taskMap, name)
	return nil
}

func validateTask(name string, spec string, cmd func()) error {
	if name == "" {
		return fmt.Errorf("task name cannot be empty")
	}
	if spec == "" {
		return fmt.Errorf("task %s: spec cannot be empty", name)
	}
	if cmd == nil {
		return fmt.Errorf("task %s: cmd cannot be nil", name)
	}
	return nil
}
