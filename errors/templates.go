package errors

// i18N bundle file can be used to templatizing the error string
var templates = map[ErrCode]string{
	BadRequestFiled: "field {{.kind}} is required",
}

var templateZh = map[ErrCode]string{
	BadRequestFiled: "字段 {{.kind}} 必填",
}

func GetTemplates(lan string, code ErrCode) string {
	if lan == "zh" || lan == "zh_CN" {
		return templateZh[code]
	}
	return templates[code]
}
