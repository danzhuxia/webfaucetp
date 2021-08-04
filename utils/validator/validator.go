package validator

import (
	"reflect"

	"webfaucetp/utils/errmsg"

	"github.com/go-playground/locales/zh_Hans_CN"
	untrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/sirupsen/logrus"
)

func Validate(data interface{}) (string, int) {
	validate := validator.New()
	uni := untrans.New(zh_Hans_CN.New())
	trans, _ := uni.GetTranslator("zh_Hans_CN")

	err := zhtrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		logrus.Errorln("err:", err)
	}
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		return label
	})

	err = validate.Struct(data)
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			return v.Translate(trans), errmsg.ERROR
		}
	}
	return "", errmsg.SUCCESS
}
