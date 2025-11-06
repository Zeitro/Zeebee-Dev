package email

import (
	"bytes"
	"html/template"
	"sync"

	"github.com/pkg/errors"
)

// global template manager
var defaultTemplates *Templates

// ErrEmailTemplateNotFound not found template  error
var ErrEmailTemplateNotFound = errors.New("email template not found")

func init() {
	defaultTemplates = &Templates{
		mux:       sync.RWMutex{},
		templates: make(map[string]*template.Template),
	}
}

// RegisterWithText register template with text content
func RegisterTemplateWithText(name string, tplContent string) (err error) {
	return defaultTemplates.RegisterWithText(name, tplContent)
}

// GetTemplate return template by name
func GetTemplate(name string) (tpl *template.Template, err error) {
	return defaultTemplates.GetTemplate(name)
}

// Parse return parsed template content
func ParseTemplate(name string, data map[string]interface{}) (content string, err error) {
	return defaultTemplates.Parse(name, data)
}

// Templates  template manager
type Templates struct {
	mux       sync.RWMutex
	templates map[string]*template.Template
}

// RegisterWithContent register template with text content
func (t *Templates) RegisterWithText(name string, tplContent string) (err error) {
	t.mux.Lock()
	defer t.mux.Unlock()
	var tpl *template.Template
	tpl, err = template.New(name).Parse(tplContent)
	if err != nil {
		return
	}
	t.templates[name] = tpl
	return err
}

// GetTemplate return template by
func (t *Templates) GetTemplate(name string) (tpl *template.Template, err error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	var ok bool
	tpl, ok = t.templates[name]
	if !ok {
		return nil, ErrEmailTemplateNotFound
	}
	return tpl, nil
}

// Parse return parsed template content
func (t *Templates) Parse(name string, data map[string]interface{}) (content string, err error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	var tpl *template.Template
	tpl, err = t.GetTemplate(name)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
