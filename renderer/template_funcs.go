package renderer

import (
	"strings"
	"text/template"
)

var templateFuncs = template.FuncMap{
	"Contains":         strings.Contains,
	"Replace":          strings.Replace,
	"ToUpper":          strings.ToUpper,
	"ToLower":          strings.ToLower,
	"SafeComputerName": SafeComputerName,
}

// SafeComputerName modifies the specified string to make it Windows computer
// name valid
func SafeComputerName(name string) string {
	if len(name) <= 0 {
		return "computername"
	}
	invalidChars := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, s := range invalidChars {
		name = strings.Replace(name, s, "", -1)
	}
	i := len(name)
	if i > 15 {
		i = 15
	}
	return name[0:i]
}
