package validator

import (
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var (
	valid            *validator.Validate
	uni              *ut.UniversalTranslator
	enTrans, zhTrans ut.Translator
)

func init() {
	var found bool
	valid = validator.New() //验证器初始化
	ent := en.New()
	uni = ut.New(ent, ent)
	enTrans, found = uni.GetTranslator("en")
	if !found {
		panic("en translation not found")
	}
	zht := zh.New()
	uni = ut.New(zht, zht)
	zhTrans, found = uni.GetTranslator("zh")
	if !found {
		panic("zh translation not found")
	}
	//注册英文翻译
	err := en_translations.RegisterDefaultTranslations(valid, enTrans)
	if err != nil {
		panic(err)
	}
	//注册中文翻译
	err = zh_translations.RegisterDefaultTranslations(valid, zhTrans)
	if err != nil {
		panic(err)
	}
	//自定义校验函数
	// ....
	fmt.Println("register success")
}

// ValidateData 验证方法
func ValidateData(l string, dataStruct interface{}) error {
	trans := enTrans
	//前端传的字段 通过tag json获取
	valid.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	if l == "zh" {
		trans = zhTrans
		// 中文需要字段国际化 注册字段标签
		valid.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			return name
		})
	}
	err := valid.Struct(dataStruct)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range errs {
				errStr := fe.Translate(trans)
				//可以根据平台的业务错误码进行封装业务错误
				fmt.Println(errStr)
			}
		}
		return err
	}
	return nil
}
