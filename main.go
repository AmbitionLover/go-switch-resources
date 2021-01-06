package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type BaseInfo struct {
	Versions string `yaml:"versions"`
	Ip       string `yaml:"ip"`
	Host     string `yaml:"host"`
}

const (
	resources = "private-resources/"
)

var inputName string
var versions = []string{"4.0", "4.1.1", "4.1.2", "4.1.3", "4.1.4", "4.2.0", "4.2.1", "4.2.2", "4.2.3"}

func main() {

	before(resources)

	fmt.Println("请输入你需要切换分支全称,输入q退出程序")

	for {
		info := BaseInfo{}
		conf, err := info.GetConf()
		if err != nil {
			fmt.Printf("错误信息：%s\n", err)
			continue
		}
		destPath := conf.Versions
		if len(destPath) == 0 {
			destPath = "private-java"
		}

		fmt.Printf("移动的主目录是 %s\n", string(destPath))

		// 移动前端资源到对应目录
		moveResources(destPath + "/oss/src/main/resources")
		moveResources(destPath + "/sign/src/main/resources")

		fmt.Println("移动完成 输入Q退出 ")
	}
}

/**
移动资源
*/
func moveResources(destPath string) {
	RemoveContents(destPath + "/static")
	RemoveContents(destPath + "/templates")
	CopyDir(resources+inputName, "/static", destPath)
	CopyDir(resources+inputName, "/templates", destPath)
}

/**
1、检查private-resources是否存在
2、检查粘贴的文件夹是否在。oss的
*/
func before(relativePath string) {
	fmt.Println("前端资源替换程序\t version：1.0\t author：Ambi")
	if exists, _ := pathExists(relativePath); !exists {
		fmt.Println("private-resources 文件夹不存在")
		panic("退出程序")
	}
	fmt.Println("资源替换文件夹检查成功")
}

func (c *BaseInfo) GetConf() (*BaseInfo, error) {

	_, err := fmt.Scan(&inputName)
	if err == io.EOF || inputName == "Q" || inputName == "q" {
		os.Exit(1)
	}

	relativePath := resources + inputName

	if exists, _ := pathExists(relativePath); !exists {
		return nil, errors.New(inputName + " 文件夹不存在")
	}

	if exists, _ := pathExists(relativePath + "/static"); !exists {
		return nil, errors.New(inputName + " static 文件夹不存在")
	}

	if exists, _ := pathExists(relativePath + "/templates"); !exists {
		return nil, errors.New(inputName + " templates 文件夹不存在")
	}

	yamlFile, err := ioutil.ReadFile(relativePath + "/config.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c, err
}

/**
 * 拷贝文件夹,同时拷贝文件夹中的文件
 * @param srcPath  		需要拷贝的文件夹路径: D:/test
 * @param destPath		拷贝到的位置: D:/backup/
 */
func CopyDir(srcPath string, dir string, destPath string) error {

	fmt.Printf("开始进行 %s 目录移动\n", dir)
	srcPath += dir

	//检测目录正确性
	if srcInfo, err := os.Stat(srcPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("srcPath不是一个正确的目录！")
			fmt.Println(e.Error())
			return e
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !destInfo.IsDir() {
			e := errors.New("destInfo不是一个正确的目录！")
			fmt.Println(e.Error())
			return e
		}
	}
	////加上拷贝时间:不用可以去掉
	//destPath = destPath + "_" + time.Now().Format("20060102150405")

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath+dir, -1)
			//fmt.Println("复制文件:" + path + " 到 " + destNewPath)
			copyFile(path, destNewPath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
	return err
}

//生成目录并拷贝文件
func copyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()
	//分割path目录
	destSplitPathDirs := strings.Split(dest, "/")

	//检测时候存在目录
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b, _ := pathExists(destSplitPath)
			if b == false {
				fmt.Println("创建目录:" + destSplitPath)
				//创建目录
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

//检测文件夹路径时候存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
删除文件夹内容
*/
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
