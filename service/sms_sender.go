package service

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

var clientNow *sms.Client

func Init() {
	credential := common.NewCredential(
		viper.GetString("tx_secret_id"),
		viper.GetString("tx_secret_key"),
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	cpf.SignMethod = "HmacSHA1"
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)
	clientNow = client
}

func TxSendSms(phone string, code string) error {
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(viper.GetString("tx_app_id"))
	request.SignName = common.StringPtr(viper.GetString("tx_sign_name"))
	request.TemplateId = common.StringPtr(viper.GetString("tx_template_id"))
	request.TemplateParamSet = common.StringPtrs([]string{code})
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + phone})
	response, err := clientNow.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(response.Response)
	// 打印返回的json字符串
	fmt.Printf("%s", b)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}
