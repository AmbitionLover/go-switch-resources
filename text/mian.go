package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var FileInfo []os.FileInfo
	var err error
	relativePath := "./img"

	if FileInfo, err = ioutil.ReadDir(relativePath); err != nil {
		fmt.Println("读取 img 文件夹出错")
		return
	}

	for _, fileInfo := range FileInfo {
		fmt.Println(fileInfo.Name())
	}
}
