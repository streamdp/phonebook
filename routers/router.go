package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"phonebook/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/reset", &controllers.MainController{}, "get:Reset")
	beego.Router("/updateRecord", &controllers.RecordController{})
	beego.Router("/v1/third", &controllers.OptionsController{})
	beego.Router("/v1/get_record", &controllers.RecordController{}, "get:GetRecord")
	beego.Router("/deleteRecord", &controllers.RecordController{}, "get:DeleteRecord")
	beego.Router("/v1/delete_many", &controllers.RecordController{}, "post:DeleteManyRecords")
	beego.Router("/back", &controllers.RecordController{}, "get:GoBack")
}
