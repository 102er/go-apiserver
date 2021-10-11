package errors

import (
	E "errors"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"strings"
)

type Error interface {
	// Generic error interface
	error
	Code() ErrCode
	Message() string
	I18nMessage(l string) string
	Cause() error
	Causes() []error
	Data() map[string]interface{}
	String() string
	ResponseErrType() ResponseErrType
	SetResponseType(r ResponseErrType) Error
	Component() ErrComponent
	SetComponent(c ErrComponent) Error
	Retryable() bool
	SetRetryable() Error
	AppendCause() Error
}

type GoError struct {
	code         ErrCode
	message      string
	data         map[string]interface{}
	causes       []error
	component    ErrComponent
	responseType ResponseErrType
	retryable    bool
	appendCause  bool
}

type ErrComponent string

const (
	ErrService ErrComponent = "service"
	ErrRepo    ErrComponent = "repository"
	ErrLib     ErrComponent = "library"
)

type ResponseErrType string

const (
	BadRequest                 ResponseErrType = "BadRequest"
	Unauthorized               ResponseErrType = "Unauthorized"
	AlreadyExists              ResponseErrType = "AlreadyExists"
	Forbidden                  ResponseErrType = "Forbidden"                  //没有权限访问资源
	NotFound                   ResponseErrType = "NotFound"                   //资源不存在
	InternalServer             ResponseErrType = "InternalServer"             //请求未完成，服务异常
	UpstreamServiceUnavailable ResponseErrType = "UpstreamServiceUnavailable" //请求未完成，服务器从上游服务器收到一个无效响应
	ServiceUnavailable         ResponseErrType = "ServiceUnavailable"         //服务不可用
	GatewayTimeout             ResponseErrType = "GatewayTimeout"             //请求超时
)

func NewGoError(code ErrCode, data map[string]interface{}, cause error) Error {
	e := GoError{
		code:   code,
		data:   data,
		causes: []error{cause},
	}
	if cause != nil {
		e.appendCause = true
	}
	return &e
}

func (e *GoError) Error() string {
	s := fmt.Sprintf("%d:%s", e.code, e.Message())
	if e.appendCause {
		s += getCauses(e.causes)
	}
	return s
}

func I18nCode(err error) ErrCode {
	type causer interface {
		Cause() error
	}
	type coder interface {
		Code() ErrCode
	}
	for err != nil {
		code, ok := err.(coder)
		if ok {
			return code.Code()
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return -1
}
func I18nMessage(err error, lang string) string {
	type messager interface {
		I18nMessage(l string) string
	}
	type causer interface {
		Cause() error
	}
	for err != nil {
		ge, ok := err.(messager)
		if ok {
			return ge.I18nMessage(lang)
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	if err != nil {
		return err.Error()
	}
	return "unknown"
}

func Is(err, target error) bool {
	return E.Is(err, target)
}

func (e *GoError) Code() ErrCode {
	return e.code
}
func (e *GoError) Message() string {
	return localize(e.code, "en", e.data)
}
func (e *GoError) I18nMessage(l string) string {
	return localize(e.code, l, e.data)
}
func (e *GoError) Cause() error {
	if len(e.causes) > 1 {
		return e.causes[0]
	}
	return nil
}

func (e *GoError) Causes() []error {
	return e.causes
}

func (e *GoError) Data() map[string]interface{} {
	return e.data
}

func (e *GoError) String() string {
	return e.Error()
}

func (e *GoError) ResponseErrType() ResponseErrType {
	return e.responseType
}
func (e *GoError) SetResponseType(r ResponseErrType) Error {
	e.responseType = r
	return e
}
func (e *GoError) Component() ErrComponent {
	return e.component
}
func (e *GoError) SetComponent(c ErrComponent) Error {
	e.component = c
	return e
}
func (e *GoError) Retryable() bool {
	return e.retryable
}
func (e *GoError) SetRetryable() Error {
	e.retryable = true
	return e
}

// AppendCause will append the error cause to the error string
func (e *GoError) AppendCause() Error {
	e.appendCause = true
	return e
}
func getCauses(errors []error) string {
	var s strings.Builder
	for _, err := range errors {
		s.WriteString(err.Error())
		s.WriteString("; ")
	}
	return s.String()
}

func localize(id ErrCode, lang string, data map[string]interface{}) string {
	l := language.English
	if lang == "zh" || lang == "zh_CN" {
		l = language.Chinese
	}
	template := Templates(lang, id)
	bundle := i18n.NewBundle(l)
	localizer := i18n.NewLocalizer(bundle, l.String())
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    string(rune(id)),
			Other: template,
		},
		TemplateData: data,
	})
}
