package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"moshopserver/models"
	"moshopserver/utils"
)

type AddressController struct {
	beego.Controller
}

type AddressListRtnJson struct {
	models.NideshopAddress
	ProviceName  string `json:"provice_name"`
	CityName     string `json:"city_name"`
	DistrictName string `json:"district_name"`
	FullRegion   string `json:"full_region"`
}

func (c *AddressController) Address_List() {

	var addresses []models.NideshopAddress
	o := orm.NewOrm()
	o.QueryTable(&models.NideshopAddress{}).Filter("user_id", getLoginUserId()).All(&addresses)

	rtnaddress := make([]AddressListRtnJson, 0)

	for _, val := range addresses {

		provinceName := models.GetRegionName(val.ProvinceId)
		cityName := models.GetRegionName(val.CityId)
		distinctName := models.GetRegionName(val.DistrictId)
		rtnaddress = append(rtnaddress, AddressListRtnJson{
			NideshopAddress: val,
			ProviceName:     provinceName,
			CityName:        cityName,
			DistrictName:    distinctName,
			FullRegion:      provinceName + cityName + distinctName,
		})

	}

	utils.ReturnHTTPSuccess(&c.Controller, rtnaddress)
	c.ServeJSON()

}
func (c *AddressController) Address_Detail() {
	id := c.GetString("id")
	intId := utils.String2Int(id)

	var address models.NideshopAddress
	var val AddressListRtnJson

	o := orm.NewOrm()
	err := o.QueryTable(&models.NideshopAddress{}).Filter("id", intId).Filter("user_id", getLoginUserId()).One(&address)
	if err != orm.ErrNoRows {

		provinceName := models.GetRegionName(address.ProvinceId)
		cityName := models.GetRegionName(address.CityId)
		distinctName := models.GetRegionName(address.DistrictId)
		val = AddressListRtnJson{
			NideshopAddress: address,
			ProviceName:     provinceName,
			CityName:        cityName,
			DistrictName:    distinctName,
			FullRegion:      provinceName + cityName + distinctName,
		}
	}
	utils.ReturnHTTPSuccess(&c.Controller, val)
	c.ServeJSON()
}

type AddressSaveBody struct {
	Address    string `json:"address"`
	CityId     int    `json:"city_id"`
	DistrictId int    `json:"district_id"`
	IsDefault  bool   `json:"is_default"`
	Mobile     string `json:"mobile"`
	Name       string `json:"name"`
	ProvinceId int    `json:"province_id"`
	AddressId  int    `json:"address_id"`
}

func (c *AddressController) Address_Save() {

	var asb AddressSaveBody
	body := c.Ctx.Input.RequestBody
	json.Unmarshal(body, &asb)

	address := asb.Address
	name := asb.Name
	mobile := asb.Mobile
	provinceId := asb.ProvinceId
	cityId := asb.CityId
	distinctId := asb.DistrictId
	isDefault := asb.IsDefault
	addressId := asb.AddressId
	userId := getLoginUserId()
	var intisdefault int
	if isDefault {
		intisdefault = 1
	} else {
		intisdefault = 0
	}

	intCityId := cityId
	intProvinceId := provinceId
	intDistinctId := distinctId

	addressData := models.NideshopAddress{
		Address:    address,
		CityId:     intCityId,
		DistrictId: intDistinctId,
		ProvinceId: intProvinceId,
		Name:       name,
		Mobile:     mobile,
		UserId:     userId,
		IsDefault:  intisdefault,
	}

	var intId int64
	o := orm.NewOrm()
	if addressId == 0 {
		id, err := o.Insert(&addressData)
		if err == nil {
			intId = id
		}
	} else {
		o.QueryTable(&models.NideshopAddress{}).Filter("id", intId).Filter("user_id", userId).Update(orm.Params{
			"is_default": 0,
		})
	}

	if isDefault {
		_, err := o.Raw("UPDATE nideshop_address SET is_default = 0 where id <> ? and user_id = ?", intId, userId).Exec()
		if err == nil {
			//res.RowsAffected()
			//fmt.Println("mysql row affected nums: ", num)
		}
	}
	var addressInfo models.NideshopAddress
	o.QueryTable(&models.NideshopAddress{}).Filter("id", intId).One(&addressInfo)

	utils.ReturnHTTPSuccess(&c.Controller, addressInfo)
	c.ServeJSON()
}

func (c *AddressController) Address_Delete() {
	addressId := c.GetString("id")
	intAddressId := utils.String2Int(addressId)

	o := orm.NewOrm()
	o.QueryTable(&models.NideshopAddress{}).Filter("id", intAddressId).Filter("user_id", getLoginUserId()).Delete()

	return
}
