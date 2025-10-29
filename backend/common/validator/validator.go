package validator

import (
	"reflect"
	"regexp"
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
	v = &Validator{
		validator: validator,
		trans:     trans,
		local:     local,
	}
	err = v.init()
	return v, err
}

func (v *Validator) init() (err error) {
	err = v.registerRules()
	if err != nil {
		return err
	}
	err = v.registerTrans()
	if err != nil {
		return err
	}
	return err
}

func (v *Validator) Validate(i interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.WithStack(errors.Errorf("validate error:%v", e))
		}
	}()
	validatorErr := v.validator.Struct(i)
	if validatorErr != nil && v.trans != nil {
		for _, err := range validatorErr.(go_validator.ValidationErrors) {
			return errors.New(err.Translate(v.trans))
		}
	}
	return validatorErr
}

func (v *Validator) registerTrans() (err error) {
	err = v.validator.RegisterTranslation("password", v.trans, func(ut ut.Translator) error {
		return ut.Add("password", "{0} must contain uppercase and lowercase letters, numbers, and at least 8 digits", true)
	}, func(ut ut.Translator, fe go_validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})
	return err
}

func (v *Validator) registerRules() (err error) {
	err = v.validator.RegisterValidation("password", validatePwd)
	return err
}

func GetEnValidator() (v *Validator, err error) {
	return NewValidator(LocalEn)
}

func validatePwd(fl go_validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	return hasLower && hasUpper && hasNumber
}
