package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"moshopserver/models"
	"moshopserver/services"

	"github.com/astaxie/beego"
	"moshopserver/utils"
)

type AuthController struct {
	beego.Controller
}

type AuthLoginBody struct {
	Code     string               `json:"code"`
	UserInfo services.ResUserInfo `json:"userInfo"`
}

func (c *AuthController) Auth_LoginByWeixin() {
	var alb AuthLoginBody
	body := c.Ctx.Input.RequestBody

	err := json.Unmarshal(body, &alb)
	if err != nil {
		fmt.Println("Auth_LoginByWeixin => Fail to json.Unmarshal, err: ", err)
		data := utils.GetHTTPRtnJsonData(500, fmt.Sprintf("Auth_LoginByWeixin => Fail to json.Unmarshal, err: %v", err))
		c.Ctx.Output.JSON(data, true, false)
		return
	}

	clientIP := c.Ctx.Input.IP()
	fmt.Println("+++++++++++++++++++++++++")
	fmt.Println("Auth_LoginByWeixin => alb: ", alb)
	fmt.Println("Auth_LoginByWeixin => clientIP: ", clientIP)
	fmt.Println("+++++++++++++++++++++++++")

	userInfo := services.Login(alb.Code, alb.UserInfo)
	if userInfo == nil {
		fmt.Println("Auth_LoginByWeixin => Fail to get userInfo")
		data := utils.GetHTTPRtnJsonData(500, "Auth_LoginByWeixin => Fail to get userInfo")
		c.Ctx.Output.JSON(data, true, false)
		return
	}

	var user models.NideshopUser
	o := orm.NewOrm()
	err = o.QueryTable(&models.NideshopUser{}).Filter("weixin_openid", userInfo.OpenID).One(&user)
	if err == orm.ErrNoRows {
		newUser := models.NideshopUser{Username: utils.GetUUID(), Password: "", RegisterTime: utils.GetTimestamp(),
			RegisterIp: clientIP, Mobile: "", WeixinOpenid: userInfo.OpenID, Avatar: userInfo.AvatarUrl, Gender: userInfo.Gender,
			Nickname: userInfo.NickName}
		o.Insert(&newUser)
		o.QueryTable(&models.NideshopUser{}).Filter("weixin_openid", userInfo.OpenID).One(&user)
	}

	userInfoMap := make(map[string]interface{})
	userInfoMap["id"] = user.Id
	userInfoMap["username"] = user.Username
	userInfoMap["nickname"] = user.Nickname
	userInfoMap["gender"] = user.Gender
	userInfoMap["avatar"] = user.Avatar
	userInfoMap["birthday"] = user.Birthday

	user.LastLoginIp = clientIP
	user.LastLoginTime = utils.GetTimestamp()

	_, err = o.Update(&user)
	if err != nil {
		fmt.Println("Auth_LoginByWeixin => Fail to Update user, err: ", err)
		data := utils.GetHTTPRtnJsonData(500, fmt.Sprintf("Auth_LoginByWeixin => Fail to Update user, err: %v", err))
		c.Ctx.Output.JSON(data, true, false)
		return
	}

	sessionKey := services.Create(utils.Int2String(user.Id))
	fmt.Println("Auth_LoginByWeixin => sessionKey: " + sessionKey)

	rtnInfo := make(map[string]interface{})
	rtnInfo["token"] = sessionKey
	rtnInfo["userInfo"] = userInfo

	utils.ReturnHTTPSuccess(&c.Controller, rtnInfo)
	c.ServeJSON()
}
