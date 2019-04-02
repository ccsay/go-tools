package qiniu

import (
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/api.v7/auth/qbox"
	"fmt"
	"context"
	"github.com/liuchonglin/go-utils"
	"qiniupkg.com/x/errors.v7"
	"net/http"
	"io/ioutil"
	"bytes"
	"github.com/liuchonglin/go-utils/timeutil"
	"github.com/liuchonglin/go-utils/stringutil"
	"path"
)

const (
	bucket       = "bucket1"
	accessKey    = "Vjo6qJFqGBFw64xkK9PB7Jto81VvxL9TlT2DEJvI"
	secretKey    = "vypqulvH-Y9n4ND8uKOY7999omxt106n_lBdvIKa"
	emptyMessage = "%s is empty"
)

var (
	keyEmptyError      = errors.New(fmt.Sprintf(emptyMessage, "key"))
	filePathEmptyError = errors.New(fmt.Sprintf(emptyMessage, "filePath"))
	upTokenEmptyError  = errors.New(fmt.Sprintf(emptyMessage, "upToken"))
	fileUrlEmptyError  = errors.New(fmt.Sprintf(emptyMessage, "fileUrl"))
)

// 文件上传响应
type FileUploadResp struct {
	Key    string `json:"key"`
	Hash   string `json:"hash"`
	Fsize  int    `json:"fsize"`
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// 提供了对资源进行管理的操作
var bucketManager *storage.BucketManager

// 文件上传，资源管理等配置
var cfg = &storage.Config{
	Zone: &storage.ZoneHuabei,
}

// 鉴权
var mac *qbox.Mac

// 表单上传对象
var formUploader *storage.FormUploader

func init() {
	mac = qbox.NewMac(accessKey, secretKey)
	formUploader = storage.NewFormUploader(cfg)
	bucketManager = storage.NewBucketManager(mac, cfg)
}

// 获取上传凭证
func GetUpToken() string {
	// 自定义凭证有效期（Expires 单位为秒，为上传凭证的有效时间）
	putPolicy := storage.PutPolicy{
		Scope:      bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
		// 2小时有效期（默认有效期为1个小时）
		Expires: 2 * 60 * 60,
	}
	return putPolicy.UploadToken(mac)
}

// 根据文件路径上传文件，自动生成Key及文件名称
func UploadByFilePathAutoKey(filePath string, upToken string) (*FileUploadResp, error) {
	return UploadByFilePath(filePath, "", upToken)
}

// 根据文件路径上传文件
func UploadByFilePath(filePath string, key string, upToken string) (*FileUploadResp, error) {
	if utils.IsEmpty(filePath) {
		return nil, filePathEmptyError
	}
	if utils.IsEmpty(upToken) {
		return nil, upTokenEmptyError
	}
	if utils.IsEmpty(key) {
		key = timeutil.GetCurrentTime(timeutil.NoDivFormatTime) + stringutil.RandomString(10) + path.Ext(filePath)
	}
	ret := &FileUploadResp{}
	if err := formUploader.PutFile(context.Background(), ret, upToken, key, filePath, nil); err != nil {
		return nil, err
	}
	return ret, nil
}

// 根据链接上传文件，自动生成Key及文件名称
func UploadByUrlAutoKey(resURL string, upToken string) (*FileUploadResp, error) {
	return UploadByUrl(resURL, "", upToken)
}

// 根据链接上传文件
func UploadByUrl(resURL string, key string, upToken string) (*FileUploadResp, error) {
	if utils.IsEmpty(resURL) {
		return nil, fileUrlEmptyError
	}

	resp, err := http.Get(resURL)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if utils.IsEmpty(key) {
		key = timeutil.GetCurrentTime(timeutil.NoDivFormatTime) + stringutil.RandomString(10) + path.Ext(resURL)
	}

	ret := &FileUploadResp{}
	if err = formUploader.Put(context.Background(), ret, upToken, key, bytes.NewReader(data), int64(len(data)), &storage.PutExtra{}); err != nil {
		return nil, err
	}
	return ret, nil
}

/*
// 根据链接上传文件，自动生成Key及文件名称
func UploadByUrlAutoKey(resURL string) (*storage.FetchRet, error) {
	return UploadByUrl(resURL, "")
}

// 根据链接上传文件
func UploadByUrl(resURL string, key string) (*storage.FetchRet, error) {
	if utils.IsEmpty(resURL) {
		return nil, fileUrlEmptyError
	}
	var putRet storage.FetchRet
	var err error
	if utils.IsEmpty(key) {
		// 不指定保存的key，默认用文件hash作为文件名
		putRet, err = bucketManager.FetchWithoutKey(resURL, bucket)
	} else {
		putRet, err = bucketManager.Fetch(resURL, bucket, key)
	}

	if err != nil {
		return nil, err
	}
	return &putRet, nil
}
*/

// 根据key获取文件信息
func GetFileInfo(key string) (*storage.FileInfo, error) {
	if utils.IsEmpty(key) {
		return nil, keyEmptyError
	}
	info, err := bucketManager.Stat(bucket, key)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// 根据key删除文件
func DelFile(key string) error {
	if utils.IsEmpty(key) {
		return keyEmptyError
	}
	return bucketManager.Delete(bucket, key)
}

// 设置文件存活时间
func SetFileLifeTime(key string, days int) error {
	if utils.IsEmpty(key) {
		return keyEmptyError
	}
	return bucketManager.DeleteAfterDays(bucket, key, days)
}
