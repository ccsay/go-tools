package jsonWebToken

import (
	"testing"
	"fmt"
	"reflect"
)

func TestJwt_CreateToken(t *testing.T) {
	type args struct {
		data map[string]interface{}
	}
	tests := []struct {
		name    string
		j       *jsonWebToken
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			j:    &jsonWebToken{},
			args: args{
				data: map[string]interface{}{"name": "admin", "phone": "13300220033", "role": "admin"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.CreateToken(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Jwt.CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("Jwt.CreateToken() token = %v \n", got)
		})
	}
}

func TestJwt_ParseToken(t *testing.T) {
	j := &jsonWebToken{}
	data := map[string]interface{}{"name": "admin", "phone": "13300220033", "role": "admin"}
	token, err := j.CreateToken(data)
	if err != nil {
		panic(err)
	}

	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		j       *jsonWebToken
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "ok",
			j:    &jsonWebToken{},
			args: args{
				tokenString: token,
			},
			want:    data,
			wantErr: false,
		}, {
			name: "format error",
			j:    &jsonWebToken{},
			args: args{
				tokenString: token + "a",
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "is expired",
			j:    &jsonWebToken{},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6ImFkbWluIiwicGhvbmUiOiIxMzMwMDIyMDAzMyIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTU1MTA5NTU5NSwiaXNzIjoibWFydGluIiwibmJmIjoxNTUxMDk0NTk0fQ.xFt_nOdRc_ZU4K7XoMKkPRj6XID866FUK0xYGQAa1Jw",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.ParseToken(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("Jwt.ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Jwt.ParseToken() = %v, want %v", got, tt.want)
			} else {
				fmt.Printf("Jwt.ParseToken() claims = %v\n", got)
			}
		})
	}
}

func TestJwt_RefreshToken(t *testing.T) {
	j := &jsonWebToken{}
	data := map[string]interface{}{"name": "admin", "phone": "13300220033", "role": "admin"}
	token, err := j.CreateToken(data)
	if err != nil {
		panic(err)
	}

	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		j       *jsonWebToken
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			j:    &jsonWebToken{},
			args: args{
				tokenString: token,
			},
			wantErr: false,
		}, {
			name: "format error",
			j:    &jsonWebToken{},
			args: args{
				tokenString: token + "a",
			},
			wantErr: true,
		}, {
			name: "is expired",
			j:    &jsonWebToken{},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6ImFkbWluIiwicGhvbmUiOiIxMzMwMDIyMDAzMyIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTU1MTA5NTU5NSwiaXNzIjoibWFydGluIiwibmJmIjoxNTUxMDk0NTk0fQ.xFt_nOdRc_ZU4K7XoMKkPRj6XID866FUK0xYGQAa1Jw",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.RefreshToken(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("Jwt.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("Jwt.CreateToken() token = %v \n", got)
		})
	}
}
