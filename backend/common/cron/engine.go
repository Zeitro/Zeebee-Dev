package cron

import (
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

var defaultEngine = NewEngine()

func Start() {
	defaultEngine.Start()
}
func Stop() {
	defaultEngine.Stop()
}
func Register(name string, spec string, cmd func()) error {
	return defaultEngine.Register(name, spec, cmd)
}
func Unregister(name string) error {
	return defaultEngine.Unregister(name)
}
func Replace(name string, spec string, cmd func()) error {
	return defaultEngine.Replace(name, spec, cmd)
}
func IsExisted(name string) bool {
	return defaultEngine.IsExisted(name)
}

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

// Replace replaces a named task atomically. Returns error on invalid
// name/spec or nil cmd. The new entry is added before the old one is removed,
// so an invalid spec leaves any existing task intact.
func (e *Engine) Replace(name string, spec string, cmd func()) (err error) {
	if err = validateTask(name, spec, cmd); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	oldID, existed := e.taskMap[name]
	id, err := e.engine.AddFunc(spec, cmd)
	if err != nil {
		return fmt.Errorf("task %s: invalid spec %q: %w", name, spec, err)
	}
	e.taskMap[name] = id
	if existed {
		e.engine.Remove(oldID)
	}
	return nil
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
