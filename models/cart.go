package models

import "github.com/astaxie/beego/orm"

func ClearBuyGoods(userId int) {
	o := orm.NewOrm()
	o.QueryTable(&NideshopCart{}).Filter("user_id", userId).Filter("session_id", 1).Filter("checked", 1).Delete()
}
