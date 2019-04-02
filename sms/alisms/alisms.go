package alisms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/liuchonglin/go-utils"
	"errors"
	"encoding/json"
	"encoding/base64"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dybaseapi/mns"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dybaseapi"
	"time"
	"github.com/liuchonglin/go-utils/stringutil"
	"strconv"
	"fmt"
)

const (
	accessKeyId  = "LTAInUcF65NaKR37"
	accessSecret = "Qf1EpQzp7ywN1EYe4l4OSuWVXZN3cR"
	mnsDomain    = "1943695596114318.mns.cn-hangzhou.aliyuncs.com"
	post         = "POST"
	https        = "https"
	domain       = "dysmsapi.aliyuncs.com"
	version      = "2017-05-25"
)

var (
	phoneNumberIsNilError  = errors.New("phoneNumber is nil")
	signNameIsNilError     = errors.New("signName is nil")
	templateCodeIsNilError = errors.New("templateCode is nil")
	sendDateIsNilError     = errors.New("sendDate is nil")
	bizIdIsNilError        = errors.New("bizId is nil")
)

type AliSms struct {
	Client *sdk.Client
}

// 发送短信响应信息
type SendSmsResponse struct {
	// 状态码的描述
	Message string `json:"Message"`
	// 请求ID
	RequestId string `json:"RequestId"`
	// 发送回执ID，可根据该ID在接口QuerySendDetails中查询具体的发送状态
	BizId string `json:"BizId"`
	// 请求状态码
	// 1.返回OK代表请求成功
	// 2.其他错误码详见错误码列表(https://help.aliyun.com/document_detail/101346.html?spm=a2c4g.11186623.2.14.3f53202aotJb4F)
	Code string `json:"Code"`
}

// 查询短信发送详细响应信息
type QuerySendDetailsResponse struct {
	// 总数
	TotalCount string `json:"totalCount"`
	// 状态码的描述
	Message string `json:"message"`
	// 请求ID
	RequestId string `json:"requestId"`
	// 请求状态码
	// 1.返回OK代表请求成功
	// 2.其他错误码详见错误码列表(https://help.aliyun.com/document_detail/101346.html?spm=a2c4g.11186623.2.14.3f53202aotJb4F)
	Code string `json:"code"`
	// 短信发送详情数组
	SmsSendDetailDTOs []SmsSendDetailDTO `json:"SmsSendDetailDTOs"`
}

// 短信发送详情
type SmsSendDetailDTO struct {
	// 转发给运营商的时间
	SendDate string `json:"sendDate"`
	// 收到运营商回执的时间
	ReceiveDate string `json:"receiveDate"`
	// 发送状态
	SendStatus int `json:"sendStatus"`
	// 错误码
	ErrCode string `json:"errCode"`
	// 模板代码
	TemplateCode string `json:"templateCode"`
	// 发送内容
	Content string `json:"content"`
	// 手机号码
	PhoneNum string `json:"phoneNum"`
}

type MnsResponse struct {
	// 发送时间
	SendTime string `json:"send_time"`
	// 接收时间
	ReceiveTime string `json:"receive_time"`
	// 是否成功
	Success bool `json:"success"`
	// 错误信息
	ErrMsg string `json:"err_msg"`
	// 错误码
	ErrCode string `json:"err_code"`
	// 手机号码
	PhoneNumber string `json:"phone_number"`
	// 短信大小 140字节算一条短信，短信长度超过140字节时会拆分成多条短信发送
	SmsSize string `json:"sms_size"`
	// 发送回执ID，可根据该ID在接口QuerySendDetails中查询具体的发送状态
	BizId string `json:"biz_id"`
}

func NewAliSms() (*AliSms, error) {
	client, err := sdk.NewClientWithAccessKey("default", accessKeyId, accessSecret)
	if err != nil {
		return nil, err
	}
	return &AliSms{Client: client}, nil
}

func (a *AliSms) SendSms(phoneNumber, signName, templateCode, templateParam string) (bool, *SendSmsResponse, error) {
	if utils.IsEmpty(phoneNumber) {
		return false, nil, phoneNumberIsNilError
	}
	if utils.IsEmpty(signName) {
		return false, nil, signNameIsNilError
	}
	if utils.IsEmpty(templateCode) {
		return false, nil, templateCodeIsNilError
	}
	request := requests.NewCommonRequest()
	// 设置请求方式
	request.Method = post
	// 设置协议
	request.Scheme = https
	// 指定域名则不会寻址
	request.Domain = domain
	// 指定产品版本
	request.Version = version
	// 指定接口名
	request.ApiName = "SendSms"
	// 设置参数值
	request.QueryParams["PhoneNumbers"] = phoneNumber
	request.QueryParams["SignName"] = signName
	request.QueryParams["TemplateCode"] = templateCode
	request.QueryParams["TemplateParam"] = templateParam

	response, err := a.Client.ProcessCommonRequest(request)
	if err != nil {
		return false, nil, err
	}

	sendSmsResponse := &SendSmsResponse{}
	if err := json.Unmarshal([]byte(response.GetHttpContentString()), sendSmsResponse); err != nil {
		return false, nil, err
	}

	return response.IsSuccess(), sendSmsResponse, nil
}

func (a *AliSms) GetSmsDetails(phoneNumber, sendDate, bizId string) (*QuerySendDetailsResponse, error) {
	if utils.IsEmpty(bizId) {
		return nil, bizIdIsNilError
	}
	return a.GetSmsDetailsList(phoneNumber, sendDate, bizId, 1, 1)
}

func (a *AliSms) GetSmsDetailsList(phoneNumber, sendDate, bizId string, pageSize int, currentPage int) (*QuerySendDetailsResponse, error) {
	if utils.IsEmpty(phoneNumber) {
		return nil, phoneNumberIsNilError
	}
	if utils.IsEmpty(sendDate) {
		return nil, sendDateIsNilError
	}
	request := requests.NewCommonRequest()
	request.Method = post
	request.Scheme = https
	request.Domain = domain
	request.Version = version
	request.ApiName = "QuerySendDetails"
	request.QueryParams["PhoneNumber"] = phoneNumber
	request.QueryParams["SendDate"] = sendDate
	request.QueryParams["PageSize"] = strconv.Itoa(pageSize)
	request.QueryParams["CurrentPage"] = strconv.Itoa(currentPage)
	if !utils.IsEmpty(bizId) {
		request.QueryParams["BizId"] = bizId
	}
	response, err := a.Client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}
	querySendDetailsResponse := &QuerySendDetailsResponse{}
	if err := json.Unmarshal([]byte(response.GetHttpContentString()), querySendDetailsResponse); err != nil {
		return nil, err
	}

	return querySendDetailsResponse, nil
}

func (a *AliSms) Mns() {
	regionId := "cn-hangzhou"
	endpoints.AddEndpointMapping(regionId, "Dybaseapi", "dybaseapi.aliyuncs.com")

	// 创建client实例
	client, err := dybaseapi.NewClientWithAccessKey(regionId, accessKeyId, accessSecret)
	if err != nil {
		// 异常处理
		panic(err)
	}

	queueName := "Alicom-Queue-1307199738792642-SmsReport"
	// 需要接收的消息类型
	// 短信回执：SmsReport，短信上行：SmsUp
	messageType := "SmsReport"

	var token *dybaseapi.MessageTokenDTO
	for {
		if token == nil || stringutil.StringToInt64(token.ExpireTime)-time.Now().Unix() > 2*60 {
			// 创建API请求并设置参数
			request := dybaseapi.CreateQueryTokenForMnsQueueRequest()
			request.MessageType = messageType
			request.QueueName = queueName
			// 发起请求并处理异常
			response, err := client.QueryTokenForMnsQueue(request)
			if err != nil {
				// 异常处理
				panic(err)
			}

			token = &response.MessageTokenDTO
		}

		mnsClient, err := mns.NewClientWithStsToken(
			regionId,
			token.AccessKeyId,
			token.AccessKeySecret,
			token.SecurityToken,
		)

		if err != nil {
			panic(err)
		}
		mnsRequest := mns.CreateBatchReceiveMessageRequest()
		mnsRequest.Domain = mnsDomain
		mnsRequest.QueueName = queueName
		//mnsRequest.ConnectTimeout = 30 * time.Second
		mnsRequest.NumOfMessages = "10"
		// 当队列中有消息时，请求立即返回；
		// 当队列中没有消息时，请求在MNS服务器端挂5秒钟，在这期间，有消息写入队列，请求会立即返回消息，5秒后，请求返回队列没有消息；
		mnsRequest.WaitSeconds = "5"

		mnsResponse, err := mnsClient.BatchReceiveMessage(mnsRequest)
		if err != nil {
			continue
		}
		receiptHandles := make([]string, len(mnsResponse.Message))
		for i, message := range mnsResponse.Message {
			messageBody, decodeErr := base64.StdEncoding.DecodeString(message.MessageBody)
			if decodeErr != nil {
				panic(decodeErr)
			}
			//TODO 业务处理
			//fmt.Println("messageBody+++++++++++", string(messageBody))
			mnsResponse := &MnsResponse{}
			if err := json.Unmarshal(messageBody, mnsResponse); err != nil {
				panic(err)
			}
			utils.PrintlnJson(mnsResponse)
			fmt.Println("==============================================")
			receiptHandles[i] = message.ReceiptHandle
		}
		if len(receiptHandles) > 0 {
			mnsDeleteRequest := mns.CreateBatchDeleteMessageRequest()
			mnsDeleteRequest.Domain = mnsDomain
			mnsDeleteRequest.QueueName = queueName
			mnsDeleteRequest.SetReceiptHandles(receiptHandles)
			//_, err = mnsClient.BatchDeleteMessage(mnsDeleteRequest) // 取消注释将删除队列中的消息
			if err != nil {
				panic(err)
			}
		}

	}
}
