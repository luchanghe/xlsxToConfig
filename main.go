package main

import (
	"awesomeProject/lib/book"
	"awesomeProject/lib/config"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 4096)
			runtime.Stack(stack, false)
			fmt.Println(fmt.Sprintf("[%s]Recovered from panic: %s\n%s\n", time.Now().Format(time.DateTime), err, stack))
		}
	}()
	c := Config.GetConfig()
	// 调用 filepath.Walk() 函数来遍历目录
	err := filepath.Walk(c.Input, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("访问路径出错 %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileSlice := strings.Split(info.Name(), ".")
		if len(fileSlice) != 2 {
			return nil
		}
		if fileSlice[1] != "xlsx" {
			return nil
		}
		name := fileSlice[0]
		xlFile, err := xlsx.OpenFile(path)
		if err != nil {
			panic(err)
		}
		bookSrt := book.NewBook(name)
		bookSrt.IncludeExcel(xlFile)
		bookSrt.WriteToJson()
		bookSrt.WriteToPhp()
		return nil
	})
	// 错误处理
	if err != nil {
		panic(err)
	}

}
