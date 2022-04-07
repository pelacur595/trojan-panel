package api

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"strconv"
)

var validate *validator.Validate

func InitValidator() {
	// Validate为单例对象
	validate = validator.New()
	// 注册校验对象
	_ = validate.RegisterValidation("validateStr", validateStr)
	_ = validate.RegisterValidation("validatePositiveInt", validatePositiveInt)
}

// 字符串必须是字母和数字的组合
func validateStr(f validator.FieldLevel) bool {
	str := f.Field().String()
	reg := "^[A-Za-z0-9]+$"
	compile := regexp.MustCompile(reg)
	return compile.MatchString(str)
}

// 正整数
func validatePositiveInt(f validator.FieldLevel) bool {
	str := strconv.FormatUint(f.Field().Uint(), 10)
	reg := "^[1-9]\\d*$"
	compile := regexp.MustCompile(reg)
	return compile.MatchString(str)
}
