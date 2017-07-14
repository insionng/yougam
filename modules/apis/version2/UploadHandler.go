package version2

import (
	"fmt"
	"github.com/insionng/makross"
	"github.com/insionng/makross/jwt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"github.com/insionng/yougam/helper"
)

type UploadError struct {
	Error string `json:"error"`
}

type Files struct {
	Files []string `json:"files"`
}

func PostUploadHandler(self *makross.Context) error {
	var uploadError UploadError
	uploadError.Error = "ErrServiceUnavailable"
	var fz = new(Files)

	parent := self.Args("parent").String() //格式 parent="topic:2",创建时id为0，须要服务端在前端完成提交流程后再进行二次修正操作
	fmt.Println(parent)

	//single file
	uploadFile, serr := self.FormFile("uploadFile")
	if serr != nil {

		// Multipart form
		form, merr := self.MultipartForm()
		if merr != nil {
			uploadError.Error = fmt.Sprintf("PostUploadHandler() self.FormFile() Error:(%v) and self.MultipartForm() Error:(%v)", serr, merr)
			return self.JSON(uploadError)
		}

		files := form.File["uploadFiles"]
		for _, file := range files {
			// Source
			srcFile, err := file.Open()
			if err != nil {
				uploadError.Error = fmt.Sprintf("PostUploadHandler() self.MultipartForm() file.Open() (%v)", err)
				return self.JSON(uploadError)
			}
			defer srcFile.Close()

			// Copy
			fz, err = ioProcess(srcFile, file.Filename, fz, self)
			if err != nil {
				uploadError.Error = fmt.Sprintf("PostUploadHandler() ioProcess() Error:%v", err)
				return self.JSON(uploadError)
			}
		}
		if len(files) == len(fz.Files) {
			return self.JSON(fz)
		}

		uploadError.Error = fmt.Sprintf("PostUploadHandler() self.FormFile() Error:%v", serr)
		return self.JSON(uploadError)
	} else {
		srcFile, err := uploadFile.Open()
		if err != nil {
			uploadError.Error = "PostUploadHandler() uploadFile.Open() Errors"
			return self.JSON(uploadError)
		}
		defer srcFile.Close()

		// Copy
		fz, err = ioProcess(srcFile, uploadFile.Filename, fz, self)
		if err != nil {
			uploadError.Error = fmt.Sprintf("PostUploadHandler() ioProcess() Error:%v", err)
			return self.JSON(uploadError)
		}
		return self.JSON(fz)
	}

}

func ioProcess(srcFile multipart.File, uploadFilename string, files *Files, self *makross.Context) (*Files, error) {
	var uploadError UploadError
	uploadError.Error = "ErrServiceUnavailable"

	claims := jwt.GetMapClaims(self)
	var IsSigned bool
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			uploadError.Error = "尚未登录"
		}
		IsSigned = true
	}

	if !IsSigned {
		return nil, self.JSON(uploadError)
	}

	targetFolder := helper.FileStorageDir + "file"
	ext := strings.ToLower(path.Ext(uploadFilename))
	filename := helper.MD5(time.Now().String()) + ext
	dirPath := fmt.Sprintf("%v/%v", targetFolder, fmt.Sprintf("%v/%v", uid, time.Now().Format("03/04/")))

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		uploadError.Error = fmt.Sprintf("PostUploadHandler() ioProcess() os.MkdirAll() (%v)", err)
		return nil, self.JSON(uploadError)
	}

	tempPath := dirPath + filename
	f, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		uploadError.Error = fmt.Sprintf("PostUploadHandler() ioProcess() os.OpenFile() (%v)", err)
		return nil, self.JSON(uploadError)
	}
	defer f.Close()

	if _, err := io.Copy(f, srcFile); err != nil {
		uploadError.Error = fmt.Sprintf("PostUploadHandler() ioProcess() io.Copy() (%v)", err)
		return nil, self.JSON(uploadError)
	}
	if !(ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif") {
		filehash, _ := helper.Filehash(helper.URL2local(tempPath), nil)
		fname := helper.EncryptHash(filehash+strconv.Itoa(int(uid)), nil)
		finalPath := dirPath + fname + ext

		if err := os.Rename(helper.URL2local(tempPath), helper.URL2local(finalPath)); err != nil {
			log.Println("重命名文件发生错误：", err)
		} else {
			os.Remove(helper.URL2local(tempPath))
		}
		if fhashed, _ := helper.Filehash(helper.URL2local(finalPath), nil); helper.ValidateHash(fname, fhashed+strconv.Itoa(int(uid))) {
			s := fmt.Sprintf("/file/%v", finalPath[8:])
			files.Files = append(files.Files, s)
			return files, nil
		} else {
			uploadError.Error = "PostUploadHandler() ioProcess() 校验文件不正确"

			if e := os.Remove(helper.URL2local(finalPath)); e != nil {
				log.Println(fmt.Sprintf("PostUploadHandler() ioProcess() 清除错误文件 %v 出错:%v", finalPath, e))
			}

			return nil, self.JSON(uploadError)
		}
	} else {

		inputFile := tempPath
		outputFile := fmt.Sprintf("%s%v-%s", dirPath, uid, filename)
		outputSize := "768x0"
		outputAlign := "center"
		background := "white"
		staticPath := "./themes/" + helper.Theme() + "/static"
		watermarkFile := staticPath + "/img/watermark.png"

		//所有上传的图片都会被缩略处理
		err = helper.Thumbnail("resize", inputFile, outputFile, outputSize, outputAlign, background)
		if err != nil {
			uploadError.Error = fmt.Sprintf("PostUploadHandler() ioProcess() 生成缩略图出错：%v", err)
			if err := os.Remove(helper.URL2local(inputFile)); err != nil {
				os.Remove(helper.URL2local(outputFile))
				uploadError.Error = fmt.Sprintf("PostUploadHandler() ioProcess() 删除缩略图出错：%v", err)
			}
			return nil, self.JSON(uploadError)
		} else {

			os.Remove(helper.URL2local(inputFile))

			helper.Watermark(watermarkFile, outputFile, outputFile, "SouthEast")

			filehash, _ := helper.Filehash(helper.URL2local(outputFile), nil)
			fname := helper.EncryptHash(filehash+strconv.Itoa(int(uid)), nil)
			finalPath := dirPath + fname + ext

			if err := os.Rename(helper.URL2local(outputFile), helper.URL2local(finalPath)); err != nil {
				log.Println("重命名文件发生错误：", err)
			} else {
				os.Remove(helper.URL2local(outputFile))
			}

			//文件权限校验 通过说明文件上传转换过程中没发生错误
			//首先读取被操作文件的hash值 和 用户请求中的文件hash值  以及 用户当前id的string类型  进行验证
			if fhashed, _ := helper.Filehash(helper.URL2local(finalPath), nil); helper.ValidateHash(fname, fhashed+strconv.Itoa(int(uid))) {
				s := fmt.Sprintf("/file/%v", finalPath[8:])
				files.Files = append(files.Files, s)
				return files, nil
			} else {
				uploadError.Error = "PostUploadHandler() ioProcess() 校验图片不正确"
				if e := os.Remove(helper.URL2local(finalPath)); e != nil {
					log.Println(fmt.Sprintf("PostUploadHandler() ioProcess() 清除错误文件 %v 出错:%v", finalPath, e))
				}
				return nil, self.JSON(uploadError)
			}
		}
	}
}
