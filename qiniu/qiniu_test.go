package qiniu

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/qiniu/api.v7/storage"
	"path"
)

func TestGetUpToken(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUpToken()
			fmt.Println(got)
		})
	}
}

func TestUploadByFilePathAutoKey(t *testing.T) {
	upToken := GetUpToken()
	type args struct {
		filePath string
		upToken  string
	}
	tests := []struct {
		name    string
		args    args
		want    *FileUploadResp
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				filePath: "./a6efce1b9d16fdfac33cef92b38f8c5494ee7b7d.jpg",
				upToken:  upToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UploadByFilePathAutoKey(tt.args.filePath, tt.args.upToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadByFilePathAutoKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("UploadByFilePathAutoKey() = %v\n", got)
		})
	}
}

func TestUploadByFilePath(t *testing.T) {
	type args struct {
		filePath string
		key      string
		upToken  string
	}
	tests := []struct {
		name    string
		args    args
		want    *FileUploadResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UploadByFilePath(tt.args.filePath, tt.args.key, tt.args.upToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadByFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UploadByFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUploadByUrlAutoKey(t *testing.T) {
	upToken := GetUpToken()
	type args struct {
		resURL  string
		upToken string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				resURL:  "http://static2.jihaoba.com/7niu/upload/58a53abac6009.jpg",
				upToken: upToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UploadByUrlAutoKey(tt.args.resURL, upToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadByUrlAutoKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fmt.Printf("UploadByUrlAutoKey() = %v\n", got)
		})
	}
}

func TestUploadByUrl(t *testing.T) {
	type args struct {
		resURL  string
		key     string
		upToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *storage.FetchRet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UploadByUrl(tt.args.resURL, tt.args.key, tt.args.upToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadByUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UploadByUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFileInfo(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    *storage.FileInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileInfo(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelFile(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DelFile(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("DelFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetFileLifeTime(t *testing.T) {
	type args struct {
		key  string
		days int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetFileLifeTime(tt.args.key, tt.args.days); (err != nil) != tt.wantErr {
				t.Errorf("SetFileLifeTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
	fmt.Println(path.Ext("./a6efce1b9d16fdfac33cef92b38f8c5494ee7b7d.jpg"))
}
