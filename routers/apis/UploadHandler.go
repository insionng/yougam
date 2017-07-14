package apis

import (
	"errors"
	"fmt"
	"github.com/insionng/makross"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func PostUploadHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)

	TplNames := "editor-ajax-result"
	uid := int64(0)
	if (_usr_ != nil) && okay {
		uid = _usr_.Id
	} else {
		self.Flash.Error("UploadPostHandler获取UID错误0!")
		self.Set("result", "UploadPostHandler获取UID错误0!")
		self.Set("resultcode", "failed")
		return self.Render(TplNames)

	}

	file, e := self.FormFile("uploadfile")
	if e != nil {
		self.Flash.Error("UploadPostHandler获取文件错误1!")
		self.Set("result", "UploadPostHandler() self.GetFile() Errors")
		self.Set("resultcode", "failed")
		return self.Render(TplNames)

	}

	targetFolder := helper.FileStorageDir + "file"

	if file.Filename == "" {

		self.Flash.Error("UploadPostHandler获取文件错误3!")

		self.Set("result", " ")
		self.Set("resultcode", "failed")
		return self.Render(TplNames)

	}

	ext := strings.ToLower(path.Ext(file.Filename))
	filename := helper.MD5(time.Now().String()) + ext
	dirpath := fmt.Sprintf("%v/%v", targetFolder, fmt.Sprintf("%v/%v", uid, time.Now().Format("03/04/")))

	_err := os.MkdirAll(dirpath, 0755)
	if _err != nil {
		self.Flash.Error(_err.Error())

		self.Set("result", " ")
		self.Set("resultcode", "failed")
		return self.Render(TplNames)
	}

	path := dirpath + filename
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {

		self.Flash.Error("UploadPostHandler获取文件错误2!")

		self.Set("result", " ")
		self.Set("resultcode", "failed")
		return self.Render(TplNames)

	}
	defer f.Close()

	src, err := file.Open()
	if err != nil {
		self.Set("result", err)
		self.Set("resultcode", "failed")
		return self.Render(TplNames)

	}
	defer src.Close()

	// Copy
	if _, err = io.Copy(f, src); err != nil {
		self.Set("result", err)
		self.Set("resultcode", "failed")
		return self.Render(TplNames)

	}

	input_file := path
	output_file := fmt.Sprintf("%s%v-%s", dirpath, uid, filename)
	output_size := "768x0"
	output_align := "center"
	background := "white"
	staticPath := "./themes/" + helper.Theme() + "/static"
	watermark_file := staticPath + "/img/watermark.png"

	//所有上传的图片都会被缩略处理
	err = helper.Thumbnail("resize", input_file, output_file, output_size, output_align, background)
	if err != nil {

		self.Flash.Error(fmt.Sprintf("UploadPostHandler生成缩略图出错:%v", err))

		if e := os.Remove(helper.URL2local(input_file)); e != nil {
			os.Remove(helper.URL2local(output_file))
			self.Flash.Error(fmt.Sprintf("UploadPostHandler生成缩略图出错:%v", e))
		}

		self.Set("result", err)
		self.Set("resultcode", "failed")
		return self.Render(TplNames)

	} else {

		os.Remove(helper.URL2local(input_file))

		helper.Watermark(watermark_file, output_file, output_file, "SouthEast")

		filehash, _ := helper.Filehash(helper.URL2local(output_file), nil)
		fname := helper.EncryptHash(filehash+strconv.Itoa(int(uid)), nil)
		newpath := dirpath + fname + ext

		if err := os.Rename(helper.URL2local(output_file), helper.URL2local(newpath)); err != nil {
			log.Println("重命名文件发生错误：", err)
		} else {
			os.Remove(helper.URL2local(output_file))
		}

		//文件权限校验 通过说明文件上传转换过程中没发生错误
		//首先读取被操作文件的hash值 和 用户请求中的文件hash值  以及 用户当前id的string类型  进行验证

		if fhashed, _ := helper.Filehash(helper.URL2local(newpath), nil); helper.ValidateHash(fname, fhashed+strconv.Itoa(int(uid))) {

			self.Set("result", "file_uploaded")
			self.Set("resultcode", "ok")
			self.Set("file_name", newpath[6:])
			return self.Render(TplNames)

		} else {

			self.Flash.Error("UploadPostHandler校验图片不正确!")

			self.Set("result", " ")
			self.Set("resultcode", "failed")

			if e := os.Remove(helper.URL2local(newpath)); e != nil {
				self.Flash.Error(fmt.Sprintf("UploadPostHandler清除错误文件 %v 出错:%v", newpath, e))
			}
		}

	}

	return errors.New("UploadPostHandler Errors")

}
