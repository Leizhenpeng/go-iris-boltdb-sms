package service

import (
	"github.com/spf13/viper"
	"leizhenpeng/go-iris-boltdb-sms/model"
	"regexp"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type SmsCore interface {
	IfExist(phone string) bool
	GenCode(phone string) (string, error)
	ValidPhone(phone string) bool
	ValidCode(phone string, code string) bool
	ClearAll()
	Total() int
}

type SmsCoreImpl struct {
}

func (s SmsCoreImpl) ClearAll() {
	model.DbNow.Flush()
}

func (s SmsCoreImpl) IfExist(phone string) bool {
	get, err := model.DbNow.Get(phone)
	if err != nil {
		return false
	}
	if get == nil {
		return false
	}
	return true
}

func (s SmsCoreImpl) Total() int {
	i, i2 := model.DbNow.Len()
	if i2 != nil {
		return 0
	}
	return i
}

func genWay() string {
	s, err := gonanoid.Generate("0123456789", 6)
	if err != nil {
		return ""
	}
	return s
}

func (s SmsCoreImpl) GenCode(phone string) (string, error) {
	sNew := genWay()
	err := model.DbNow.Set(phone, []byte(sNew))
	if err != nil {
		return "", err
	}
	model.DbNow.Expire(phone, viper.GetInt("EXPIRE"))
	err = TxSendSms(phone, sNew)
	if err != nil {
		model.DbNow.Del(phone)
		return "", err
	}
	return sNew, nil
}

func (s SmsCoreImpl) ValidPhone(phone string) bool {
	//check if phone is valid by regex
	rePhone := `^1[3-9]\d{9}$`
	return regexp.MustCompile(rePhone).MatchString(phone)

}

func (s SmsCoreImpl) ValidCode(phone string, code string) bool {
	if !s.ValidPhone(phone) {
		return false
	}
	if !s.IfExist(phone) {
		return false
	}
	get, err := model.DbNow.Get(phone)
	if err != nil {
		return false
	}
	if string(get) == code {
		//del phone
		err := model.DbNow.Del(phone)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

var (
	Sms SmsCore = SmsCoreImpl{}
)
