package controllers

import (
	"context"
	beego "github.com/beego/beego/v2/server/web"
	pagination2 "github.com/beego/beego/v2/server/web/pagination"
	"github.com/beego/beego/v2/server/web/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"phonebook/models"
	"reflect"
	"strings"
)

var (
	recPerPage = 18
	bo         models.BaseOptions
	keys       []string
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
	if IsEmpty(bo) {
		getBaseOptions()
	}
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

func (c *MainController) ListPostsByOffsetAndLimit(querySet []models.PhoneBookRecord, offset int, postperpage int) []models.PhoneBookRecord {
	if offset+postperpage < len(querySet) {
		return querySet[offset : offset+recPerPage]
	} else if offset < len(querySet) {
		return querySet[offset:len(querySet)]
	} else {
		return querySet
	}
}

func (c *MainController) Get() {
	var page int
	err := c.Ctx.Input.Bind(&page, "p")
	if err != nil {
		log.Println(err)
	}
	var queryResultGet []models.PhoneBookRecord
	pqtkn, success := c.GetToken("prevQuery")
	if success && len(pqtkn) > 5 {
		queryResultGet, _ = models.GetManyByFilter(c.GetFilterStringFromCookie())
	}
	if len(queryResultGet) > 0 {
		paginator := pagination2.SetPaginator(c.Ctx, recPerPage, int64(len(queryResultGet)))
		c.Data["PhoneRecords"] = c.ListPostsByOffsetAndLimit(queryResultGet, paginator.Offset(), recPerPage)
	}
	c.BindTemplateData()
	c.Data["CurrentPage"] = page
	c.TplName = "base.html"
}

func (c *MainController) BindTemplateData() {
	pqtkn, success := c.GetToken("prevQuery")
	prevQuery := make(map[string]string)
	if success {
		pqstr := strings.Split(pqtkn, ";")
		prevQuery["first_level"] = strings.TrimSpace(pqstr[0])
		prevQuery["second_level"] = strings.TrimSpace(pqstr[1])
		prevQuery["third_level"] = strings.TrimSpace(pqstr[2])
		prevQuery["first_name"] = strings.TrimSpace(pqstr[3])
		prevQuery["last_name"] = strings.TrimSpace(pqstr[4])
		prevQuery["middle_name"] = strings.TrimSpace(pqstr[5])
	}
	c.Data["RequestUrl"] = c.Ctx.Input.URL()
	c.Data["pq"] = prevQuery
	c.Data["FirstLevel"] = bo.FirstLevel
	c.Data["SecondLevel"] = keys
	c.Data["ThirdLevel"] = bo.SecondLevel[prevQuery["second_level"]]
}

func (c *MainController) ShowOne() {
	var id string
	err := c.Ctx.Input.Bind(&id, "id")
	if err != nil {
		log.Println(err)
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
	}
	queryResultGet, err := models.GetManyByFilter(bson.D{{"_id", objectId}})
	if err != nil {
		log.Println(err)
	}
	if len(queryResultGet) > 0 {
		paginator := pagination2.SetPaginator(c.Ctx, recPerPage, int64(len(queryResultGet)))
		c.Data["PhoneRecords"] = c.ListPostsByOffsetAndLimit(queryResultGet, paginator.Offset(), recPerPage)
	}
	c.BindTemplateData()
	c.TplName = "base.html"
}

func (c *MainController) SetToken(key string, value string) {
	c.SetSecureCookie(c.Session.SessionID(context.TODO()), key, value)
}

func (c *MainController) GetToken(key string) (string, bool) {
	cookieValue, success := c.GetSecureCookie(c.Session.SessionID(context.TODO()), key)
	if !success {
		log.Println("Can't get value from cookie...")
	}
	return cookieValue, success
}

func (c *MainController) Reset() {
	c.SetToken("prevQuery", ";;;;;")
	c.Redirect("/", 303)
}

func (c *MainController) Post() {
	qf := models.QueryForm{}
	if err := c.ParseForm(&qf); err != nil {
		log.Println(err)
	}
	pqstring := qf.FirstLevel + ";" + qf.SecondLevel + ";" + qf.ThirdLevel + ";" + qf.FirstName + ";" + qf.LastName + ";" + qf.MiddleName
	c.SetToken("prevQuery", pqstring)
	c.Redirect("/", 303)
}

func (c *MainController) GetFilterStringFromCookie() bson.D {
	filterString := bson.D{}
	pqtkn, success := c.GetToken("prevQuery")
	if success && len(pqtkn) > 0 {
		pqstr := strings.Split(pqtkn, ";")
		if strings.TrimSpace(pqstr[0]) != "" {
			filterString = append(filterString, primitive.E{Key: "department.first_level", Value: strings.TrimSpace(pqstr[0])})
		}
		if strings.TrimSpace(pqstr[1]) != "" {
			filterString = append(filterString, primitive.E{Key: "department.second_level", Value: strings.TrimSpace(pqstr[1])})
		}
		if strings.TrimSpace(pqstr[2]) != "" {
			filterString = append(filterString, primitive.E{Key: "department.third_level", Value: strings.TrimSpace(pqstr[2])})
		}
		if strings.TrimSpace(pqstr[3]) != "" {
			filterString = append(filterString, primitive.E{Key: "first_name", Value: strings.TrimSpace(pqstr[3])})
		}
		if strings.TrimSpace(pqstr[4]) != "" {
			filterString = append(filterString, primitive.E{Key: "last_name", Value: strings.TrimSpace(pqstr[4])})
		}
		if strings.TrimSpace(pqstr[5]) != "" {
			filterString = append(filterString, primitive.E{Key: "middle_name", Value: strings.TrimSpace(pqstr[5])})
		}
	}
	return filterString
}
