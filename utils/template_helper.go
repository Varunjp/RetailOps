package utils

import "html/template"

func TemplateFuncs() template.FuncMap{
	return template.FuncMap{
		"inc":inc,
		"dec":dec,
	}
}

func inc(x int)int{ return x +1}
func dec(x int)int{ return x-1}