package main

import (
	beego "github.com/beego/beego/v2/server/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "phonebook/routers"
)

func isDivisibleByThree(in int)(out int){
	if (in + 1) % 3 == 0 {
		out = 1
	} else {
		out = 0
	}
	return
}

func getHexFromObjectID(in primitive.ObjectID)(out string){
	return in.Hex()
}

func main() {
	_ = beego.AddFuncMap("getHexFromObjectID", getHexFromObjectID)
	_ = beego.AddFuncMap("isDivisibleByThree", isDivisibleByThree)
	beego.Run()
}

