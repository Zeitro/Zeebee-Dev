package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	go_validator "github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	zh2 "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
)

const (
	LocalZh = "zh"
	LocalEn = "en"
)

type Validator struct {
	validator *go_validator.Validate
	trans     ut.Translator
	local     string
}

func NewValidator(local string) (v *Validator, err error) {
	var trans ut.Translator
	validator := go_validator.New()

	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	switch local {
	case LocalZh:
		zh := zh.New()
		uni := ut.New(zh)
		trans, _ = uni.GetTranslator(local)
		err = zh2.RegisterDefaultTranslations(validator, trans)
	case LocalEn:
		en := en.New()
		uni := ut.New(en)
		trans, _ = uni.GetTranslator(local)
		err = en2.RegisterDefaultTranslations(validator, trans)
	default:
		en := en.New()
		uni := ut.New(en)
		trans, _ = uni.GetTranslator(LocalEn)
		err = en2.RegisterDefaultTranslations(validator, trans)
	}
	return &Validator{
		validator: validator,
		trans:     trans,
		local:     local,
	}, nil
}

func (v *Validator) Validate(i interface{}) (err error) {
	validatorErr := v.validator.Struct(i)
	if validatorErr != nil && v.trans != nil {
		for _, err := range validatorErr.(go_validator.ValidationErrors) {
			return errors.New(err.Translate(v.trans))
		}
	}
	return validatorErr
}

func GetEnValidator() (v *Validator, err error) {
	return NewValidator(LocalEn)
}
