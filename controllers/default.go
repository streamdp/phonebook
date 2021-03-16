package controllers

import (
	"context"
	beego "github.com/beego/beego/v2/server/web"
	pagination2 "github.com/beego/beego/v2/server/web/pagination"
	"github.com/beego/beego/v2/server/web/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"reflect"
	"phonebook/models"
)

var (
	queryResult []models.PhoneBookRecord
	prevQuery   = make(map[string]string)
	recPerPage  = 18
	page        = 0
	bo          models.BaseOptions
	keys        []string
)

type ExtendedController struct {
	beego.Controller
	Session session.Store
}

func IsEmpty(bo models.BaseOptions) bool {
	return reflect.DeepEqual(bo, models.BaseOptions{})
}

func (c *ExtendedController) Prepare() {
	c.Session = c.StartSession()
	if IsEmpty(bo) { getBaseOptions() }
	c.Data["RequestUrl"] = c.Ctx.Input.URL()
	c.Data["pq"] = prevQuery
	c.Data["FirstLevel"] = bo.FirstLevel
	c.Data["SecondLevel"] = keys
	c.Data["ThirdLevel"] = bo.SecondLevel[prevQuery["second_level"]]
}

type MainController struct {
	ExtendedController
}

func getBaseOptions() {
	var err error
	bo, err = models.GetBaseOptions()
	if err != nil {
		log.Println(err)
	}
	keys = make([]string, 0, len(bo.SecondLevel))
	for key, _ := range bo.SecondLevel {
		keys = append(keys, key)
	}
}

func (c *MainController) ListPostsByOffsetAndLimit(offset int, postperpage int) []models.PhoneBookRecord{
	if offset +recPerPage < len(queryResult) {
		return queryResult[offset:offset+recPerPage]
	}	else if offset < len(queryResult){
		return queryResult[offset:len(queryResult)]
	}	else {
		return queryResult
	}
}

func (c *MainController) Get() {
	err := c.Ctx.Input.Bind(&page, "p")
	if err != nil {
		log.Println(err)
	}
	if len(queryResult)>0 {
		paginator := pagination2.SetPaginator(c.Ctx, recPerPage, int64(len(queryResult)))
		c.Data["PhoneRecords"] = c.ListPostsByOffsetAndLimit(paginator.Offset(), recPerPage)
	}
	c.TplName = "base.html"
}

func (c *MainController) SetToken(key string, value string) {
	c.SetSecureCookie(c.Session.SessionID(context.TODO()), key, value)
}

func (c *MainController) GetToken(key string) string {
	cookieValue, success := c.GetSecureCookie(c.Session.SessionID(context.TODO()), key)
	if !success {
		log.Println("Can't get value from cookie...")
	}
	return cookieValue
}

func (c *MainController) Reset() {
	queryResult = []models.PhoneBookRecord{}
	prevQuery =	make(map[string]string)
	page = 0
	c.SetToken("department.first_level", "")
	c.SetToken("department.second_level", "")
	c.SetToken("department.third_level", "")
	c.SetToken("first_name", "")
	c.SetToken("last_name", "")
	c.SetToken("middle_name", "")
	c.Redirect("/", 303)
}

func (c *MainController) Post() {
	qf := models.QueryForm{}
	if err := c.ParseForm(&qf); err != nil {
		log.Println(err)
	}
	prevQuery["first_level"] = qf.FirstLevel
	prevQuery["second_level"] = qf.SecondLevel
	prevQuery["third_level"] = qf.ThirdLevel
	prevQuery["first_name"] = qf.FirstName
	prevQuery["last_name"] = qf.LastName
	prevQuery["middle_name"] = qf.MiddleName
	c.SetToken("department.first_level", qf.FirstLevel)
	c.SetToken("department.second_level", qf.SecondLevel)
	c.SetToken("department.third_level", qf.ThirdLevel)
	c.SetToken("first_name", qf.FirstName)
	c.SetToken("last_name", qf.LastName)
	c.SetToken("middle_name", qf.MiddleName)
	filterString := bson.D{}
	if qf.FirstLevel != "" {
		filterString = append(filterString, primitive.E{Key: "department.first_level", Value: qf.FirstLevel})
	}
	if qf.SecondLevel != "" {
		filterString = append(filterString, primitive.E{Key: "department.second_level", Value: qf.SecondLevel})
	}
	if qf.ThirdLevel != "" {
		filterString = append(filterString, primitive.E{Key: "department.third_level", Value: qf.ThirdLevel})
	}
	if qf.FirstName != "" {
		filterString = append(filterString, primitive.E{Key: "first_name", Value: qf.FirstName})
	}
	if qf.LastName != "" {
		filterString = append(filterString, primitive.E{Key: "last_name", Value: qf.LastName})
	}
	if qf.MiddleName != "" {
		filterString = append(filterString, primitive.E{Key: "middle_name", Value: qf.MiddleName})
	}
	queryResult, _ = models.GetManyByFilter(filterString)
	c.Redirect("/", 303)
}