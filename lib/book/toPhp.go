package book

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/orderedmap"
	"strconv"
	"strings"
)

func (b *Book) WriteToPhp() {
	// 将Go map 转换为PHP数组的字符串形式
	phpArray := goMapToPHPArray(b.dataMap, 1, "    ")
	var buffer bytes.Buffer
	buffer.WriteString("<?php\n$CONF = ")
	buffer.WriteString(phpArray)
	buffer.WriteString(";")
	b.write(buffer.Bytes(), "php")
}

func goMapToPHPArray(data *orderedmap.OrderedMap, depth int, indent string) string {
	var builder strings.Builder
	builder.WriteString("array(\n")
	dst := createDepthBuilder(depth, indent)
	for _, key := range data.Keys() {
		val, _ := data.Get(key)
		// 检查 str1 是否为整数字符串
		pre := `"%s" => %s, `
		if _, err := strconv.Atoi(key); err == nil {
			pre = `%s => %s, `
		}
		builder.WriteString(dst.String())
		builder.WriteString(fmt.Sprintf(pre+"\n", key, goValueToPHPValue(val, depth, indent)))
	}
	builder.WriteString(dst.String())
	builder.WriteString(")")
	return builder.String()
}

// goValueToPHPValue 将Go值转换为PHP值的字符串形式
func goValueToPHPValue(value interface{}, depth int, indent string) string {
	dst := createDepthBuilder(depth, indent)
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case int, int64:
		return fmt.Sprintf("%v", v)
	case float64:
		return fmt.Sprintf("%v", v)
	case *orderedmap.OrderedMap:
		return goMapToPHPArray(v, depth+1, indent)
	case []interface{}:
		var builder strings.Builder
		dst2 := createDepthBuilder(depth+1, indent)
		builder.WriteString("array(\n")
		for _, item := range v {
			s := fmt.Sprintf("%v", item)
			if _, err := strconv.Atoi(s); err == nil {
				builder.WriteString(dst2.String())
				builder.WriteString(s)
				builder.WriteString(",\n")
			} else {
				builder.WriteString(dst2.String())
				builder.WriteString("\"")
				builder.WriteString(s)
				builder.WriteString("\"")
				builder.WriteString(",\n")
			}
		}
		builder.WriteString(dst.String())
		builder.WriteString(")")
		return builder.String()
	default:
		return "" // 不支持的类型，返回空字符串
	}
}

func createDepthBuilder(depth int, indent string) strings.Builder {
	var dst strings.Builder
	for i := 0; i < depth; i++ {
		dst.WriteString(indent)
	}
	return dst
}
