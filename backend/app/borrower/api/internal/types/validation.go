package types

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// registerValidation 存储所有的验证方法
// key: 验证标签名称，value: 对应的验证函数
var registerValidation = map[string]validator.Func{
	"mobile": validateMobile, // 手机号验证
	"idcard": validateIdCard, // 身份证号验证
}

// 初始化并注册所有验证方法
func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			name = strings.ToLower(fld.Name)
		}
		return name
	})
	// 遍历注册所有验证方法
	for tag, handler := range registerValidation {
		_ = validate.RegisterValidation(tag, handler)
	}
}

// 验证手机号
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 使用正则表达式验证中国大陆手机号(13,14,15,16,17,18,19开头的11位数字)
	match, _ := regexp.MatchString("^1[3-9]\\d{9}$", mobile)
	return match
}

// 验证身份证号
func validateIdCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()
	// 支持15位或18位身份证号
	match, _ := regexp.MatchString("(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)", idCard)
	return match
}
