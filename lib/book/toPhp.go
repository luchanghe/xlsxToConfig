package book

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/orderedmap"
	"strconv"
	"strings"
)

type id string

func (b *Book) WriteToPhp() {
	// 将Go map 转换为PHP数组的字符串形式
	phpArray := goMapToPHPArray(b.dataMap)
	var buffer bytes.Buffer
	buffer.WriteString("<?php\n$CONF=")
	buffer.WriteString(phpArray)
	buffer.WriteString(";")
	b.write(buffer.Bytes(), "php")
}

func goMapToPHPArray(data *orderedmap.OrderedMap) string {
	var builder strings.Builder
	builder.WriteString("array(")
	for _, key := range data.Keys() {
		val, _ := data.Get(key)
		// 检查 str1 是否为整数字符串
		pre := `"%s" => %s, `
		if _, err := strconv.Atoi(key); err == nil {
			pre = `%s => %s, `
		}
		builder.WriteString(fmt.Sprintf(pre, key, goValueToPHPValue(val)))
	}
	builder.WriteString(")")
	return builder.String()
}

// goValueToPHPValue 将Go值转换为PHP值的字符串形式
func goValueToPHPValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v)
	case int, int64:
		return fmt.Sprintf("%v", v)
	case float64:
		return fmt.Sprintf("%v", v)
	case *orderedmap.OrderedMap:
		return goMapToPHPArray(v)
	default:
		return "" // 不支持的类型，返回空字符串
	}
}
