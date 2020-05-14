package utils

import (
	"encoding/json"

	"github.com/astaxie/beego"
)

type HTTPData struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func ReturnHTTPSuccess(this *beego.Controller, val interface{}) {
	rtnData := HTTPData{
		ErrNo:  0,
		ErrMsg: "",
		Data:   val,
	}

	data, err := json.Marshal(rtnData)
	if err != nil {
		this.Data["json"] = err
	} else {
		this.Data["json"] = json.RawMessage(string(data))
	}
}

func GetHTTPRtnJsonData(errCode int, errMsg string) interface{} {
	rtnData := HTTPData{
		ErrNo:  errCode,
		ErrMsg: errMsg,
		Data:   nil,
	}
	data, _ := json.Marshal(rtnData)
	return json.RawMessage(string(data))
}
