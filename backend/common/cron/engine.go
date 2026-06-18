package cron

import (
	"fmt"

	"github.com/robfig/cron/v3"

	"zeebee-backend/common/mapx"
)

type Engine struct {
	engine  *cron.Cron
	taskMap *mapx.Map[string, cron.EntryID]
}

func NewEngine() *Engine {
	return &Engine{
		engine:  cron.New(cron.WithSeconds()),
		taskMap: mapx.NewMap[string, cron.EntryID](),
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
	if name == "" {
		return fmt.Errorf("task name cannot be empty")
	}
	if spec == "" {
		return fmt.Errorf("task %s: spec cannot be empty", name)
	}
	if cmd == nil {
		return fmt.Errorf("task %s: cmd cannot be nil", name)
	}

	if _, existed := e.taskMap.Get(name); existed {
		return fmt.Errorf("task %s already exists", name)
	}

	id, err := e.engine.AddFunc(spec, cmd)
	if err != nil {
		return fmt.Errorf("task %s: invalid spec %q: %w", name, spec, err)
	}
	e.taskMap.Set(name, id)
	return nil
}

// Unregister removes a named task. Returns error if the task does not exist.
func (e *Engine) Unregister(name string) error {
	id, existed := e.taskMap.Get(name)
	if !existed {
		return fmt.Errorf("task %s does not exist", name)
	}
	e.engine.Remove(id)
	e.taskMap.Delete(name)
	return nil
}

// IsExisted reports whether a named task is registered.
func (e *Engine) IsExisted(name string) bool {
	_, existed := e.taskMap.Get(name)
	return existed
}
