package controllers

import (
	"log"
	"phonebook/models"
)

type OptionsController struct {
	ExtendedController
}

func (c *OptionsController) Get() {
	var (
		val string
		err error
	)
	err = c.Ctx.Input.Bind(&val, "id")
	if err != nil {
		log.Println(err)
	}
	a, err := models.GetBaseOptions()
	if err != nil {
		log.Println(err)
	}
	c.Data["json"] = a.SecondLevel[val]
	err = c.ServeJSON()
	if err != nil {
		log.Println(err)
	}
}
