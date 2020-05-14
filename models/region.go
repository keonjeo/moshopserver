package models

import (
	"github.com/astaxie/beego/orm"
)

func GetRegionName(regionId int) string {
	o := orm.NewOrm()
	var region NideshopRegion
	o.QueryTable(&NideshopRegion{}).Filter("id", regionId).One(&region)

	return region.Name
}
