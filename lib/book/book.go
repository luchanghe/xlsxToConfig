package book

import (
	Config "awesomeProject/lib/config"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/orderedmap"
	"github.com/tealeg/xlsx"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Body struct {
	head *Head
	data any
}

type Head struct {
	enName            string //英文字段
	cnName            string //中文字段
	fType             string //字段类型
	relationSheetName string
	relationKey       string
}
type PageData struct {
	name     string
	bodyRows []*BodyRow
}

type BodyRow struct {
	bodies []*Body
}

type Book struct {
	name    string
	dataMap *orderedmap.OrderedMap
}

func NewBook(name string) *Book {
	return &Book{name: name}
}
func (b *Book) IncludeExcel(file *xlsx.File) {
	fileTable := map[string]*xlsx.Sheet{}
	for _, sheet := range file.Sheets {
		fileTable[sheet.Name] = sheet
	}
	file.Sheets[0].Rows[2].Cells[0].Value = "id"
	b.dataMap = createDataMap(create(fileTable, file.Sheets[0], "", ""))
}

func createDataMap(data *PageData) *orderedmap.OrderedMap {
	jsonMap := orderedmap.New()
	for i, row := range data.bodyRows {
		bodyMap := orderedmap.New()
		keyStr := ""
		for _, body := range row.bodies {
			head := body.head
			switch head.fType {
			case "array":
				bodyMap.Set(head.enName, createDataMap(body.data.(*PageData)))
			case "int", "string", "float":
				bodyMap.Set(head.enName, body.data)
			case "json":
				if len(body.data.(string)) == 0 {
					body.data = "{}"
				}

				if body.data.(string)[0] == '[' && body.data.(string)[len(body.data.(string))-1] == ']' {
					var jsonArray []interface{}
					err := json.Unmarshal([]byte(body.data.(string)), &jsonArray)
					if err != nil {
						panic(fmt.Sprintf("json解析失败,%s=>%s::%v\n", data.name, head.enName, err))
					}
					bodyMap.Set(head.enName, jsonArray)
				} else {
					jsonTempMap := orderedmap.New()
					err := jsonTempMap.UnmarshalJSON([]byte(body.data.(string)))
					if err != nil {
						panic(fmt.Sprintf("json解析失败,%s=>%s::%v\n", data.name, head.enName, err))
					}
					bodyMap.Set(head.enName, jsonTempMap)
				}
			case "id":
				s, err := strconv.Atoi(body.data.(string))
				if err == nil {
					bodyMap.Set(head.enName, s)
				} else {
					bodyMap.Set(head.enName, body.data)
				}
				keyStr = head.enName
			}
		}
		if keyStr != "" {
			vKey, _ := bodyMap.Get(keyStr)
			switch vKey.(type) {
			case string:
				jsonMap.Set(vKey.(string), bodyMap)
			case int:
				jsonMap.Set(strconv.Itoa(vKey.(int)), bodyMap)
			}
		} else {
			jsonMap.Set(strconv.Itoa(i), bodyMap)
		}
	}
	return jsonMap
}

func create(table map[string]*xlsx.Sheet, sheet *xlsx.Sheet, relationKey string, relationVal string) *PageData {
	relationIndex := -1
	data := &PageData{name: sheet.Name}
	var heads []*Head
	//xlsx第一行获取字段名
	for _, cell := range sheet.Rows[0].Cells {
		if cell.String() == "" {
			break
		}
		enName := strings.Split(cell.String(), "@")
		//第一行,存储的是字段名
		heads = append(heads, &Head{enName: enName[len(enName)-1]})
	}
	//xlsx第二行获取中文备注
	for i := 0; i < len(heads); i++ {
		heads[i].cnName = sheet.Rows[1].Cells[i].String()
	}
	//xlsx第三行获取字段类型
	for i := 0; i < len(heads); i++ {
		heads[i].fType = sheet.Rows[2].Cells[i].String()
		if heads[i].fType == "array" {
			pattern := `(\w*):#(\w*).(\w*)`
			// 编译正则表达式
			re := regexp.MustCompile(pattern)
			// 查找匹配的字符串
			matches := re.FindStringSubmatch(heads[i].enName)
			// 如果找到匹配项
			if len(matches) >= 4 {
				heads[i].enName = matches[1]
				heads[i].relationSheetName = matches[2]
				heads[i].relationKey = matches[3]
			}
		} else {
			if heads[i].enName == relationKey {
				relationIndex = i
			}
		}
	}
	//循环后面的每一行
	for _, row := range sheet.Rows[3:] {
		if len(row.Cells) == 0 {
			break
		}
		if relationIndex != -1 && row.Cells[relationIndex].String() != relationVal {
			continue
		}
		//基于头标循环每一格
		rowTemp := &BodyRow{}
		for i := 0; i < len(heads); i++ {
			var fieldVal any
			var ceilString string
			if len(row.Cells) > i {
				ceilString = row.Cells[i].String()
			} else {
				ceilString = ""
			}
			switch heads[i].fType {
			case "array":
				temp := create(table, table[heads[i].relationSheetName], heads[i].relationKey, row.Cells[0].String())
				fieldVal = temp
			case "string":
				fieldVal = ceilString
			case "id":
				fieldVal = ceilString
			case "json":
				fieldVal = ceilString
			case "float":
				temp, _ := strconv.ParseFloat(ceilString, 64)
				fieldVal = temp
			case "int":
				temp, _ := strconv.Atoi(ceilString)
				fieldVal = temp
			}
			rowTemp.bodies = append(rowTemp.bodies, &Body{
				head: heads[i],
				data: fieldVal,
			})
		}
		data.bodyRows = append(data.bodyRows, rowTemp)
	}
	return data
}

func (b *Book) write(d []byte, s string) {
	c := Config.GetConfig()
	dir := c.OutPut + s + "/"
	_, err := os.Stat(dir)
	if err != nil {
		// 如果目录不存在，则创建它
		if os.IsNotExist(err) {
			err2 := os.MkdirAll(dir, 0755)
			if err2 != nil {
				panic(fmt.Sprintf("创建目录出错: %v\n", err2))
			}
		}
	}
	fileLink := dir + b.name + "." + s
	file, err := os.Create(fileLink)
	if err != nil {
		panic(fmt.Sprintf("创建文件失败:%v\n", err))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// 写入数据
	_, err = file.Write(d)
	if err != nil {
		panic(fmt.Sprintf("写入文件失败:%v\n", err))
	}
}
