package protokit

import (
	"regexp"
	"strings"

	"github.com/hoisie/mustache"
)

func camelToSnake(queryPath string) string {
	// 使用正则表达式将大写字母前面插入下划线，并将字符串转换为小写
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snakeCase := re.ReplaceAllString(queryPath, "${1}_${2}")
	// 将字符串全部转为小写
	snakeCase = strings.ToLower(snakeCase)
	return snakeCase
}

func standardQueryPath(queryPath string, snakeCase bool, mapping map[string]string) string {

	pp := make(map[string]interface{})
	for k, v := range mapping {
		pp[k] = v
	}

	queryPath = mustache.Render(queryPath, pp)

	// 保证开头有且只有一个/
	for strings.HasPrefix(queryPath, "/") {
		queryPath = strings.TrimPrefix(queryPath, "/")
	}
	queryPath = "/" + queryPath
	// 移除结尾的/
	for strings.HasSuffix(queryPath, "/") {
		queryPath = strings.TrimSuffix(queryPath, "/")
	}
	// 移除全部空格
	queryPath = strings.ReplaceAll(queryPath, " ", "")
	// 全部snake case
	if snakeCase {
		queryPath = camelToSnake(queryPath)
	}
	if queryPath == "" {
		queryPath = "/"
	}
	return queryPath
}
