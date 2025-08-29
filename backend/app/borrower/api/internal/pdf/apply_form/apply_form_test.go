package apply_form

import (
	"os"
	"testing"
)

func TestApplyForm_Generate(t *testing.T) {
	output := "apply_form_test.pdf"
	f, err := os.Create(output)
	if err != nil {
		t.Errorf("os.Create() error = %v", err)
	}
	defer f.Close()
	form, err := NewApplyForm()
	if err != nil {
		t.Errorf("NewApplyForm() error = %v", err)
	}
	err = form.Generate(f)
	if err != nil {
		t.Errorf("Generate() error = %v", err)
	}
	t.Logf("file %s created", output)
}
