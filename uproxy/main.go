package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const EXPORT = "export"
const HTTP_PROXY = "http_proxy"
const HTTPS_PROXY = "https_proxy"

//环境变量文件名称
const BASH_PROFILE_FILE_NAME = ".bash_profile"

//环境变量文件全路径
var BASH_PROFILE_PATH string
//操作类型 a:增加代理 r:移除代理
var processType string
//http代理地址
var httpProxy string
//https代理地址
var httpsProxy string

//http proxy环境变量设置
var EXPORT_HTTP_PROXY string
//https proxy环境变量设置
var EXPORT_HTTPS_PROXY string

func cmd() {
	flag.StringVar(&processType, "p", "add", "操作类型 add:增加代理 remove:移除代理")
	flag.StringVar(&httpProxy, "http", "http://127.0.0.1:1087", "http代理地址")
	flag.StringVar(&httpsProxy, "https", "http://127.0.0.1:1087", "https代理地址")
	flag.Parse()
}

func main() {
	cmd()
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}
	BASH_PROFILE_PATH = userHome + "/" + BASH_PROFILE_FILE_NAME
	EXPORT_HTTP_PROXY = EXPORT + " " + HTTP_PROXY + "=" + httpProxy + "\n"
	EXPORT_HTTPS_PROXY = EXPORT + " " + HTTPS_PROXY + "=" + httpsProxy + "\n"
	switch processType {
	case "add":
		//增加代理
		addProxy()
		break
	case "remove":
		//移除代理
		removeProxy()
		break
	default:
		//参数错误
		fmt.Println("参数错误 -p a:增加代理 r:移除代理")
	}

}

///增加代理
func addProxy() {
	//追加写入
	f, err := os.OpenFile(BASH_PROFILE_PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	content := string(b)
	if strings.Contains(content, HTTP_PROXY) || strings.Contains(content, HTTPS_PROXY) {
		fmt.Println("代理设置已经存在")
		return
	}
	_, err = io.WriteString(f, EXPORT_HTTP_PROXY)
	if err != nil {
		panic(err)
	}
	_, err = io.WriteString(f, EXPORT_HTTPS_PROXY)
	if err != nil {
		panic(err)
	}
	cmdSuccess := execCommand("source", []string{BASH_PROFILE_PATH})
	if cmdSuccess {
		fmt.Println("执行成功")
	} else {
		fmt.Println("执行失败")
	}
}

///移除代理
func removeProxy() {
	exist := checkFileIsExist(BASH_PROFILE_PATH)
	if !exist {
		fmt.Println("代理设置不存在")
		return
	}
	//覆盖写入
	f, err := os.OpenFile(BASH_PROFILE_PATH, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var content string
	for {
		line, err := r.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if strings.Contains(line, HTTP_PROXY) || strings.Contains(line, HTTPS_PROXY) {
			continue
		}
		content = content + line
	}
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	}
	cmdSuccess := execCommand("source", []string{BASH_PROFILE_PATH})
	if cmdSuccess {
		fmt.Println("执行成功")
	} else {
		fmt.Println("执行失败")
	}
}

///判断文件是否存在  存在返回 true 不存在返回false
func checkFileIsExist(file string) bool {
	var exist = true
	if _, err := os.Stat(file); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

///执行cmd
func execCommand(commandName string, params []string) bool {
	cmd := exec.Command(commandName, params...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}
	cmd.Wait()
	return true
}
