package controllers

import (
	"context"
	"github.com/disintegration/imaging"
	"github.com/mxmCherry/translit/uknational"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/transform"
	"log"
	"os"
	"phonebook/models"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RecordController struct {
	ExtendedController
}

func (c *RecordController) GoBack(recoverQuery bool, page int, id string) {
	if recoverQuery {
		c.Redirect("/?p="+strconv.Itoa(page), 302)
	} else {
		c.Redirect("/showOne?id="+id, 302)
	}
}

func (c *RecordController) DeleteRecord() {
	var (
		val  string
		page int
	)
	err := c.Ctx.Input.Bind(&val, "id")
	if err != nil {
		log.Println(err)
	}
	err = c.Ctx.Input.Bind(&page, "p")
	if err != nil {
		log.Println(err)
	}
	err = models.DeleteOne(val)
	if err != nil {
		log.Println(err)
	}
	c.GoBack(true, page, "")
}

func (c *RecordController) DeleteManyRecords() {
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
	c.GoBack(true, recordIds.PageNum, "")
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

func (c *RecordController) GetToken(key string) (string, bool) {
	cookieValue, success := c.GetSecureCookie(c.Session.SessionID(context.TODO()), key)
	if !success {
		log.Println("Can't get value from cookie...")
	}
	return cookieValue, success
}

func (c *RecordController) GetFilterStringFromCookie() bson.D {
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

func (c *RecordController) Post() {
	var (
		ur  models.PhoneBookRecord
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
		filename := recordId + "." + filenameSplited[len(filenameSplited)-1]
		err = os.MkdirAll(photoFilePath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		err = c.SaveToFile("photo", photoFilePath+filename)
		if err != nil {
			log.Println(err)
		}
		src, err := imaging.Open(photoFilePath + filename)
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}
		src = imaging.Resize(src, 415, 0, imaging.Lanczos)
		src = imaging.CropAnchor(src, 415, 415, imaging.Center)
		err = imaging.Save(src, photoFilePath+filename)
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}
		ur.PhotoUrl = photoFilePath + filename
	} else {
		ur.PhotoUrl = "static/img/incognito.jpg"
	}
	if uf.ID != "" {
		err = models.UpdateOne(uf.ID, &ur)
		if err != nil {
			log.Println(err)
		}
		c.GoBack(true, uf.PageNum, "")
	} else {
		err = models.CreateOne(&ur)
		if err != nil {
			log.Println(err)
		}
		c.GoBack(false, uf.PageNum, ur.ID.Hex())
	}
}
