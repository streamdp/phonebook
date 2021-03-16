package controllers

import (
	"context"
	"github.com/mxmCherry/translit/uknational"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/transform"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"phonebook/models"
	"time"
	"github.com/disintegration/imaging"
)

type RecordController struct {
	ExtendedController
}

func (c *RecordController) GoBack(recoverQuery bool, id string) {
	if recoverQuery {
		queryResult, _ = models.GetManyByFilter(c.GetFilterStringFromCookie())
		c.Redirect("/?p=" + strconv.Itoa(page),302)
	} else {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println(err)
		}
		queryResult, err = models.GetManyByFilter(bson.D{{"_id", objectId}})
		if err != nil {
			log.Println(err)
		}
		c.Redirect("/" ,302)
	}
}

func (c *RecordController) DeleteRecord() {
	var val string
	err := c.Ctx.Input.Bind(&val, "id")
	if err != nil {
		log.Println(err)
	}
	err = models.DeleteOne(val)
	if err != nil {
		log.Println(err)
	}
	c.GoBack(true,"")
}

func (c *RecordController)DeleteManyRecords() {
	var err error
	recordIds := models.RecordIds{}
	if err = c.ParseForm(&recordIds); err != nil {
		log.Println(err)
	}
	for i := range recordIds.ID {
		err = models.DeleteOne(recordIds.ID[i])
		if err != nil {
			log.Println(err)
		}
	}
	c.GoBack(true,"")
}

func (c *RecordController) GetRecord() {
	var val string
	err := c.Ctx.Input.Bind(&val, "id")
	if err != nil {
		log.Println(err)
	}
	c.Data["json"], _ = models.GetOneById(val)
	err = c.ServeJSON()
	if err != nil {
		log.Println(err)
	}
}

func (c *RecordController) SetToken(key string, value string) {
	c.SetSecureCookie(c.Session.SessionID(context.TODO()), key, value)
}

func (c *RecordController) GetToken(key string) string {
	cookieValue, success := c.GetSecureCookie(c.Session.SessionID(context.TODO()), key)
	if !success {
		log.Println("Can't get value from cookie...")
	}
	return cookieValue
}

func (c *RecordController) GetFilterStringFromCookie()  bson.D {
	filterString := bson.D{}
	prevQuery["first_level"] = c.GetToken("department.first_level")
	if prevQuery["first_level"] != "" {
		filterString = append(filterString, primitive.E{Key: "department.first_level", Value: prevQuery["first_level"]})
	}
	prevQuery["second_level"] = c.GetToken("department.second_level")
	if  prevQuery["second_level"] != "" {
		filterString = append(filterString, primitive.E{Key: "department.second_level", Value: prevQuery["second_level"]})
	}
	prevQuery["third_level"] = c.GetToken("department.third_level")
	if prevQuery["third_level"] != "" {
		filterString = append(filterString, primitive.E{Key: "department.third_level", Value: prevQuery["third_level"]})
	}
	prevQuery["first_name"] = c.GetToken("first_name")
	if prevQuery["first_name"] != "" {
		filterString = append(filterString, primitive.E{Key: "first_name", Value: prevQuery["first_name"]})
	}
	prevQuery["last_name"] = c.GetToken("last_name")
	if prevQuery["last_name"] != "" {
		filterString = append(filterString, primitive.E{Key: "last_name", Value: prevQuery["last_name"]})
	}
	prevQuery["middle_name"] = c.GetToken("middle_name")
	if prevQuery["middle_name"] != "" {
		filterString = append(filterString, primitive.E{Key: "middle_name", Value: prevQuery["middle_name"]})
	}
	return filterString
}

func (c *RecordController) Post() {
	var (
		ur models.PhoneBookRecord
		err error
		)
	uf := models.UpdateForm{}
	if err = c.ParseForm(&uf); err != nil {
		log.Println(err)
	}
	if uf.ID != "" {
		ur, err = models.GetOneByIdWithoutIdField(uf.ID)
		if err != nil {
			log.Println(err)
		}
	} else {
		ur = models.PhoneBookRecord{}
		ur.ID = primitive.NewObjectID()
	}
	dep := models.Department{}
	dep.FirstLevel = uf.FirstLevel
	dep.SecondLevel = uf.SecondLevel
	dep.ThirdLevel = uf.ThirdLevel
	ur.FirstName = uf.FirstName
	ur.LastName = uf.LastName
	ur.MiddleName = uf.MiddleName
	ur.Department = dep
	ur.Position = uf.Position
	ur.ServiceNumber = uf.ServiceNumber
	ur.PersonalNumber = uf.PersonalNumber
	ur.ServiceMobileNumber = uf.ServiceMobileNumber
	if uf.ID != "" {
		ur.UpdatedAt = time.Now()
		ur.IsUpdated = true
	} else {
		ur.CreatedAt = time.Now()
		ur.IsUpdated = false
	}
	file, head, err := c.GetFile("photo")
	if err != nil {
		log.Println(err)
	}
	if file != nil {
		defer file.Close()
		recordId := ""

		if uf.ID != "" {
			recordId = uf.ID
		} else {
			recordId = ur.ID.Hex()
		}
		uk := uknational.ToLatin()
		photoFilePath := strings.ReplaceAll(strings.ToLower("static/img/"+uf.FirstLevel+"/"+uf.SecondLevel+"/"+uf.ThirdLevel+"/"), " ", "_")
		photoFilePath, _, err = transform.String(uk.Transformer(), photoFilePath)
		if err != nil {
			log.Println(err)
		}
		re := regexp.MustCompile("[[:^ascii:]]")
		photoFilePath = re.ReplaceAllLiteralString(photoFilePath, "")
		filenameSplited := strings.Split(head.Filename, ".")
		filename:=recordId+"."+filenameSplited[len(filenameSplited) - 1]
		err = os.MkdirAll(photoFilePath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		err = c.SaveToFile("photo", photoFilePath+filename)
		if err != nil {
			log.Println(err)
		}
		src, err := imaging.Open(photoFilePath+filename)
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}
		src = imaging.Resize(src, 415, 0, imaging.Lanczos)
		src = imaging.CropAnchor(src, 415, 415, imaging.Center)
		err = imaging.Save(src, photoFilePath+filename)
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}
		ur.PhotoUrl = photoFilePath+filename
	} else {
		ur.PhotoUrl = "static/img/incognito.jpg"
	}
	if uf.ID != "" {
		err = models.UpdateOne(uf.ID, &ur)
		if err != nil {
			log.Println(err)
		}
		c.GoBack(true, "")
	} else {
		err = models.CreateOne(&ur)
		if err != nil {
			log.Println(err)
		}
		c.GoBack(false, ur.ID.Hex())
	}
}