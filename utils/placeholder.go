package utils

import (
	"fmt"
	"regexp"
	"strings"
)

/**
做变量的提取，和获取
*/
var FactorPlaceholderReplacer = &TemplateFactorPlaceholderReplacer{}

type TemplateFactorPlaceholderReplacer struct{}

// 找到所有的占位符${placeholder}
func (*TemplateFactorPlaceholderReplacer) FindPlaceholder(template string) ([]string, error) {
	// 不包含占位符的直接return
	if template == "" || !strings.Contains(template, "${") || !strings.Contains(template, "}") {
		return nil, nil
	}

	// 正则匹配所有占位符变量名
	reg, err := regexp.Compile(`\${([\p{Han}\w:._-]+)}`)
	if err != nil {
		return nil, err
	}
	matchedVariables := reg.FindAllString(template, -1)
	if matchedVariables == nil {
		return nil, nil
	}
	placeholderNames := make([]string, 0, len(matchedVariables))
	for _, matched := range matchedVariables {
		placeholderNames = append(placeholderNames, matched[2:len(matched)-1])
	}
	return placeholderNames, nil
}

// 将template中の${var}变量替换掉
func (*TemplateFactorPlaceholderReplacer) Replace(template string, params map[string]interface{}) (val string, err error) {
	// 参数为空也不用替换了
	if params == nil || len(params) <= 0 {
		return template, nil
	}
	// 不包含占位符的直接return
	if template == "" || !strings.Contains(template, "${") || !strings.Contains(template, "}") {
		return template, nil
	}

	// 正则匹配所有占位符变量名
	reg, err := regexp.Compile(`\${([\p{Han}\w:._-]+)}`)
	if err != nil {
		return
	}
	matchedVariables := reg.FindAllString(template, -1)
	if matchedVariables == nil {
		return template, nil
	}
	for _, matched := range matchedVariables {
		variableName := matched[2 : len(matched)-1]
		if variable, exists := params[variableName]; exists {
			template = strings.ReplaceAll(template, matched, fmt.Sprintf("%+v", variable))
		}
	}
	return template, nil
}
