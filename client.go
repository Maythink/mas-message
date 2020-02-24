package mas

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Req struct {
	host   string
	params string
}

//普通短信
type NorClient struct {
	EcName    string `json:"ecName"`  //企业名称
	ApId      string `json:"apId"`    //接口账号用户名
	Mobiles   string `json:"mobiles"` //每批次限5000个号码
	Content   string `json:"content"`
	Sign      string `json:"sign"`
	AddSerial string `json:"addSerial"`
	Mac       string `json:"mac"`
}

//模板短信
type TmpClient struct {
	EcName     string `json:"ecName"`     //企业名称
	ApId       string `json:"apId"`       //接口账号用户名
	TemplateId string `json:"templateId"` //模板ID
	Mobiles    string `json:"mobiles"`    //每批次限5000个号码
	Params     string `json:"params"`     //模板变量
	Sign       string `json:"sign"`
	AddSerial  string `json:"addSerial"`
	Mac        string `json:"mac"`
}

type VivoTokenPar struct {
}

func NewTmpClient(ecName, apId, secretKey, templateId, mobiles, params, sign, addSerial string) (*Req, error) {
	if len(strings.Split(mobiles, ",")) > 5000 {
		return nil, errors.New("发送手机号码超限")
	}
	mac := Md5(ecName + apId + secretKey + templateId + mobiles + params + sign + addSerial)
	vc := &TmpClient{
		EcName:     ecName,
		ApId:       apId,
		TemplateId: templateId,
		Mobiles:    mobiles,
		Params:     params,
		Sign:       sign,
		AddSerial:  addSerial,
		Mac:        mac,
	}
	res, err := json.Marshal(vc)
	if err != nil {
		return nil, err
	}
	req := &Req{
		host:   ProductionHost,
		params: Base64_Encode(string(res)),
	}
	return req, err
}

func Md5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	s := hex.EncodeToString(md5Ctx.Sum(nil))
	return s
}

func Base64_Encode(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

//----------------------------------------Sender----------------------------------------//
// 发送模板短信到手机
func (r *Req) SendTmpMessage() (*Result, error) {
	res, err := r.doPost(r.host+TmpURL, []byte(r.params))
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, errors.New(result.Rspcod)
	}
	return &result, nil
}

func handleResponse(response *http.Response) ([]byte, error) {
	defer func() {
		_ = response.Body.Close()
	}()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *Req) doPost(url string, formData []byte) ([]byte, error) {
	var result []byte
	var req *http.Request
	var resp *http.Response
	var err error

	req, err = http.NewRequest("POST", url, bytes.NewReader(formData))
	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	tryTime := 0
tryAgain:
	resp, err = client.Do(req)
	if err != nil {
		tryTime += 1
		if tryTime < PostRetryTimes {
			goto tryAgain
		}
		return nil, err
	}
	result, err = handleResponse(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("network error")
	}
	return result, nil
}

func (r *Req) doGet(url string, params string) ([]byte, error) {
	var result []byte
	var req *http.Request
	var resp *http.Response
	var err error
	req, err = http.NewRequest("GET", url+params, nil)
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	result, err = handleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}
