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
const HttpProxy = "http_proxy"
const HttpsProxy = "https_proxy"

//环境变量文件名称
const BashProfileFileName = ".bash_profile"

//环境变量文件全路径
var BashProfilePath string

//http代理地址
var httpProxy string

//https代理地址
var httpsProxy string

//http proxy环境变量设置
var ExportHttpProxy string

//https proxy环境变量设置
var ExportHttpsProxy string

//add proxy
var a bool

//remove proxy
var r bool

//show proxy config
var s bool

func init() {
	flag.BoolVar(&a, "a", false, "add:增加代理")
	flag.BoolVar(&r, "r", false, "remove:移除代理")
	flag.BoolVar(&s, "s", false, "show:显示代理配置")
	flag.StringVar(&httpProxy, "http", "http://127.0.0.1:1087", "http代理地址")
	flag.StringVar(&httpsProxy, "https", "http://127.0.0.1:1087", "https代理地址")
	flag.Parse()
}

func main() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}
	BashProfilePath = userHome + "/" + BashProfileFileName
	ExportHttpProxy = EXPORT + " " + HttpProxy + "=" + httpProxy + "\n"
	ExportHttpsProxy = EXPORT + " " + HttpsProxy + "=" + httpsProxy + "\n"
	//增加代理
	if a {
		addProxy()
	}
	//移除代理
	if r {
		removeProxy()
	}
	//显示代理配置
	if s {
		showProxyConfig()
	}
	if !a && !r && !s {
		fmt.Println("-？ for help")
	}
}

///增加代理
func addProxy() {
	//追加写入
	f, err := os.OpenFile(BashProfilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	content := string(b)
	if strings.Contains(content, HttpProxy) || strings.Contains(content, HttpsProxy) {
		showProxyConfig()
		fmt.Println("代理设置已经存在")
		return
	}
	_, err = io.WriteString(f, ExportHttpProxy)
	if err != nil {
		panic(err)
	}
	_, err = io.WriteString(f, ExportHttpsProxy)
	if err != nil {
		panic(err)
	}
	cmdSuccess := execCommand("source", []string{BashProfilePath})
	if cmdSuccess {
		fmt.Print(ExportHttpProxy)
		fmt.Print(ExportHttpsProxy)
		fmt.Println("执行成功")
	} else {
		fmt.Println("执行失败")
	}
}

///移除代理
func removeProxy() {
	exist := checkFileIsExist(BashProfilePath)
	if !exist {
		fmt.Println("代理设置不存在")
		return
	}
	//覆盖写入
	f, err := os.OpenFile(BashProfilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
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
		if strings.Contains(line, HttpProxy) || strings.Contains(line, HttpsProxy) {
			continue
		}
		content = content + line
	}
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	}
	cmdSuccess := execCommand("source", []string{BashProfilePath})
	if cmdSuccess {
		fmt.Println("执行成功")
	} else {
		fmt.Println("执行失败")
	}
}

///显示代理配置
func showProxyConfig() {
	exist := checkFileIsExist(BashProfilePath)
	if !exist {
		fmt.Println("代理设置不存在")
		return
	}
	//读取文件
	f, err := os.OpenFile(BashProfilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var hasProxyConfig = false
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if strings.Contains(line, HttpProxy) || strings.Contains(line, HttpsProxy) {
			fmt.Println(line)
			hasProxyConfig = true
		}
	}
	if !hasProxyConfig {
		fmt.Println("代理设置不存在")
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
