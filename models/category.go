package models

import (
	"github.com/astaxie/beego/orm"
	"moshopserver/utils"
)

func GetChildCategoryId(categoryId int) []int64 {
	o := orm.NewOrm()
	var childIds []orm.Params
	o.QueryTable(&NideshopCategory{}).Filter("parent_id", categoryId).Limit(10000).Values(&childIds, "id")
	childIntIds := utils.ExactMapValues2Int64Array(childIds, "Id")
	return childIntIds
}

func GetCategoryWhereIn(categoryId int) []int64 {
	childCategoryIds := GetChildCategoryId(categoryId)
	childCategoryIds = append(childCategoryIds, int64(categoryId))
	return childCategoryIds
}
