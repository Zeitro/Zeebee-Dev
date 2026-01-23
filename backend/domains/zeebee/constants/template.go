package constants

import "fmt"

type TemplateCategory string

const (
	TemplateCategoryEmail    TemplateCategory = "Email"
	TemplateCategoryContract TemplateCategory = "Contract"
)

func (c TemplateCategory) Validate() error {
	if c != TemplateCategoryEmail && c != TemplateCategoryContract {
		return fmt.Errorf("Invalid Template Category:%v", c)
	}
	return nil
}
