package models

import (
	"github.com/astaxie/beego/orm"
)

func IsUserHasCollect(userId, typeId, valueId int) int {
	o := orm.NewOrm()
	var collect NideshopCollect

	err := o.QueryTable(&NideshopCollect{}).Filter("type_id", typeId).Filter("value_id", valueId).Filter("user_id", userId).One(&collect)
	if err == nil {
		return 1
	}
	return 0
}
