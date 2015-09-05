package controllers

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var allowImageType = "gif|jpeg|jpg|png|bmp"
var allowFileType = "rar|doc|docx|zip|pdf|txt|swf|mkv|avi|rm|rmvb|mpeg|mpg|ogg|mov|wmv|mp4|webm"
var website = "/"
var uploadfile = "attach/files/"
var uploadvideo = "attach/vedio/"
var uploadimage = "attach/images/"

var configJson []byte // 当客户端请求 /ueditor/go/controller?action=config 返回的json内容

func init() {
	file, err := os.Open("static/editor/ueditor/go/config.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer file.Close()
	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(file)

	configJson = buf.Bytes()
}

type UEController struct {
	beego.Controller
}

//入口方法
func (this *UEController) UEditor() {

	action := this.Ctx.Input.Query("action")
	switch {
	case action == "config":
		Config(&this.Controller)
	case action == "uploadimage" || action == "uploadscrawl" || action == "catchimage":
		UploadImage(&this.Controller)
	case action == "uploadvideo":
		UploadVedio(&this.Controller)
	case action == "uploadfile":
		UploadFile(&this.Controller)
	case action == "listimage":
		ImageManager(&this.Controller)
	case action == "listfile":
		FilesManager(&this.Controller)
	}
}

func Config(this *beego.Controller) {

	this.Ctx.WriteString(string(configJson))
	this.StopRun()
}

//上传图片
func UploadImage(this *beego.Controller) {
	f, h, _ := this.GetFile("upfile")
	filename := h.Filename
	f.Close() //关闭，减少缓存
	//获取扩展名
	ext := filename[strings.LastIndex(filename, ".")+1:]

	if !strings.Contains(allowImageType, ext) {
		fmt.Println(filename)
		this.Ctx.WriteString("{\"state\":\"FAILED\"}")
		this.StopRun()
	}
	newname := strconv.FormatInt(time.Now().Unix(), 10) + "_" + filename
	err := this.SaveToFile("upfile", uploadimage+newname)
	state := "SUCCESS"
	if err != nil {
		fmt.Println(err)
		state = "FAILED"
	}
	url := website + uploadimage + newname
	this.Ctx.WriteString("{\"state\": \"" + state + "\", \"url\": \"" + url + "\", \"title\": \"\",\"original\": \"" + filename + "\"}")
	this.StopRun()
}

//上传视频
func UploadVedio(this *beego.Controller) {
	f, h, _ := this.GetFile("upfile")
	filename := h.Filename
	f.Close() //关闭，减少缓存
	index := strings.LastIndex(filename, ".")
	filetype := ""
	if index == -1 {
		this.Ctx.WriteString("{\"state\":\"FAILED\"}")
		this.StopRun()
	}
	filetype = filename[index:]
	ext := filetype[1:]
	if !strings.Contains(allowFileType, ext) {
		this.Ctx.WriteString("{\"state\":\"FAILED\"}")
		this.StopRun()
	}
	newname := strconv.FormatInt(time.Now().Unix(), 10) + "_" + filename
	err := this.SaveToFile("upfile", uploadvideo+newname)
	state := "SUCCESS"
	if err != nil {
		fmt.Println(err)
		state = "FAILED"
	}
	url := website + uploadvideo + newname
	this.Ctx.WriteString("{\"state\": \"" + state + "\", \"url\": \"" + url + "\", \"title\": \"\",\"original\": \"" + filename + "\"}")
	this.StopRun()
}

//上传文件
func UploadFile(this *beego.Controller) {
	f, h, _ := this.GetFile("upfile")
	filename := h.Filename
	f.Close() //关闭，减少缓存

	index := strings.LastIndex(filename, ".")
	filetype := ""
	if index == -1 {
		this.Ctx.WriteString("{\"state\":\"FAILED\"}")
		this.StopRun()
	}
	filetype = filename[index:]
	ext := filetype[1:]
	if !strings.Contains(allowFileType, ext) {
		this.Ctx.WriteString("{\"state\":\"FAILED\"}")
		this.StopRun()
	}
	newname := strconv.FormatInt(time.Now().Unix(), 10) + "_" + filename
	err := this.SaveToFile("upfile", uploadfile+newname)
	state := "SUCCESS"
	if err != nil {
		fmt.Println(err)
		state = "FAILED"
	}
	url := website + uploadfile + newname

	this.Ctx.WriteString("{\"state\": \"" + state + "\", \"url\": \"" + url + "\", \"title\": \"\",\"original\": \"" + filename + "\"}")
	this.StopRun()
}

//图片管理
func ImageManager(this *beego.Controller) {
	strRet := ""
	callbackjson := "{\"state\": \"SUCCESS\",\"list\": ["
	total := 0
	err := filepath.Walk(uploadimage, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		ext := path[strings.LastIndex(path, ".")+1:]
		if strings.Contains(allowImageType, ext) {
			strRet += (path + "ue_separate_ue")
			callbackjson += "{\"url\": \"/" + uploadimage + f.Name() + "\"},"
			total++
			fmt.Println("allow:", path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	callbackjson += "],\"start\": 0,\"total\": " + strconv.Itoa(total) + "}"
	fmt.Println(strRet)
	this.Ctx.WriteString(callbackjson)
	this.StopRun()
}

//文件管理
func FilesManager(this *beego.Controller) {
	strRet := ""
	callbackjson := "{\"state\": \"SUCCESS\",\"list\": ["
	total := 0
	err := filepath.Walk(uploadfile, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		ext := path[strings.LastIndex(path, ".")+1:]
		if strings.Contains(allowFileType, ext) {
			strRet += (path + "ue_separate_ue")
			callbackjson += "{\"url\": \"/" + uploadfile + f.Name() + "\"},"
			total++
			fmt.Println("allow:", path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	callbackjson += "],\"start\": 0,\"total\": " + strconv.Itoa(total) + "}"
	fmt.Println(strRet)
	this.Ctx.WriteString(callbackjson)
	this.StopRun()
}
