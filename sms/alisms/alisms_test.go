package alisms

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	utils "github.com/liuchonglin/go-utils"
)

func TestNewAliSms(t *testing.T) {
	tests := []struct {
		name    string
		want    *AliSms
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAliSms()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAliSms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAliSms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAliSms_SendSms(t *testing.T) {
	aliSms, err := NewAliSms()
	if err != nil {
		panic(err)
	}

	type fields struct {
		Client *sdk.Client
	}
	type args struct {
		phoneNumber   string
		signName      string
		templateCode  string
		templateParam string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Client: aliSms.Client,
			},
			args: args{
				phoneNumber:   "17710660309",
				signName:      "亖堂小镇",
				templateCode:  "SMS_99430013",
				templateParam: "{\"code\":\"你老公叫你回来吃饭了，知道不\"}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AliSms{
				Client: tt.fields.Client,
			}
			isSuccess, resp, err := a.SendSms(tt.args.phoneNumber, tt.args.signName, tt.args.templateCode, tt.args.templateParam)
			if err != nil {
				t.Errorf("a.SendSms() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println("send sms status=", isSuccess)
			utils.PrintlnJson(resp)
		})
	}
}

func TestAliSms_Mns(t *testing.T) {
	aliSms, err := NewAliSms()
	if err != nil {
		panic(err)
	}
	type fields struct {
		Client *sdk.Client
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:"ok",
			fields:fields{
				Client:aliSms.Client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AliSms{
				Client: tt.fields.Client,
			}
			a.Mns()
		})
	}
}
