package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"phonebook/connectionhelper"
	"time"
)

type BaseOptions struct {
	FirstLevel			[]string			`json:"first_level" bson:"first_level"`
	SecondLevel			map[string][]string	`json:"second_level" bson:"second_level"`
}

type Department struct {
	FirstLevel			string				`json:"first_level" bson:"first_level"`
	SecondLevel			string				`json:"second_level" bson:"second_level"`
	ThirdLevel			string				`json:"third_level" bson:"third_level"`
}

type QueryForm struct {
	FirstLevel			string				`form:"first_level"`
	SecondLevel			string				`form:"second_level"`
	ThirdLevel			string				`form:"third_level"`
	FirstName   		string			    `form:"first_name"`
	LastName			string			    `form:"last_name"`
	MiddleName			string			    `form:"middle_name"`
}

type UpdateForm struct {
	ID					string				`form:"modal_id_user"`
	FirstLevel			string				`form:"modal_first_level"`
	SecondLevel			string				`form:"modal_second_level"`
	ThirdLevel			string				`form:"modal_third_level"`
	FirstName   		string			    `form:"modal_first_name"`
	LastName			string			    `form:"modal_last_name"`
	MiddleName			string			    `form:"modal_middle_name"`
	Position			string				`form:"modal_position"`
	ServiceNumber		[]string			`form:"service_num"`
	PersonalNumber		[]string			`form:"personal_num"`
	ServiceMobileNumber	[]string			`form:"service_mobile_num"`
}

type RecordIds struct {
	ID					[]string 			`json:"id" form:"should_be_removed"`
}

type PhoneBookRecord struct {
	ID         	 		primitive.ObjectID	`json:"id" bson:"_id,omitempty"`
	FirstName   		string			    `json:"first_name" bson:"first_name"`
	LastName			string			    `json:"last_name" bson:"last_name"`
	MiddleName			string			    `json:"middle_name" bson:"middle_name"`
	Department			Department			`json:"department" bson:"department"`
	Position 			string			    `json:"position" bson:"position"`
	ServiceNumber		[]string		    `json:"service_number" bson:"service_number"`
	PersonalNumber  	[]string			`json:"personal_number" bson:"personal_number"`
	ServiceMobileNumber	[]string			`json:"service_mobile_number" bson:"service_mobile_number"`
	CreatedAt   		time.Time     	    `json:"created_at" bson:"created_at"`
	UpdatedAt   		time.Time       	`json:"updated_at" bson:"updated_at"`
	IsUpdated			bool				`json:"is_updated" bson:"is_updated"`
	PhotoUrl			string				`json:"photo_url" bson:"photo_url"`
}

func CreateOne(record *PhoneBookRecord) error {
	//Get MongoDB connection using connectionhelper.
	client, err := connectionhelper.GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionhelper.DB).Collection(connectionhelper.ISSUES)
	//Perform InsertOne operation & validate against the error.
	res, err := collection.InsertOne(context.TODO(), record)
	if err != nil {
		return err
	}
	//Return success without any error.
	fmt.Println("from method", res.InsertedID)
	return nil
}

func UpdateOne(id string, record *PhoneBookRecord) error {
	client, err := connectionhelper.GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(connectionhelper.DB).Collection(connectionhelper.ISSUES)
	objectId, _ := primitive.ObjectIDFromHex(id)
	opts := options.Update().SetUpsert(true)
	update := bson.M{
		"$set": &record,
	}
	_, err = collection.UpdateOne(context.TODO(), bson.D{{"_id", objectId}} , &update, opts)
	if err != nil {
		return err
	}
	return nil
}

func GetBaseOptions() (BaseOptions, error) {
	result := BaseOptions{}
	//Define filter query for fetching specific document from collection
	filter := bson.D{{Key: "code", Value: "options"}}
	//Get MongoDB connection using connectionhelper.
	client, err := connectionhelper.GetMongoClient()
	if err != nil {
		return result, err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionhelper.DB).Collection("options")
	//Perform FindOne operation & validate against the error.
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	//Return result without any error.
	return result, nil
}

func GetManyByFilter(filter bson.D) ([]PhoneBookRecord, error) {
	var result []PhoneBookRecord
	//Get MongoDB connection using connectionhelper.
	client, err := connectionhelper.GetMongoClient()
	if err != nil {
		return result, err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionhelper.DB).Collection(connectionhelper.ISSUES)
	//Perform FindOne operation & validate against the error.
	cur, finErr := collection.Find(context.TODO(), filter)
	if finErr != nil {
		return result, err
	}
	for cur.Next(context.TODO()) {
		pbr := PhoneBookRecord{}
		err := cur.Decode(&pbr)
		if err != nil {
			return result, err
		}
		result = append(result, pbr)
	}
	cur.Close(context.TODO())
	if len(result) == 0 {
		return result, mongo.ErrNoDocuments
	}
	//Return result without any error.
	return result, nil
}

func GetOneByIdWithoutIdField(id string) (PhoneBookRecord, error) {
	result := PhoneBookRecord{}
	objectId, _ := primitive.ObjectIDFromHex(id)
	client, err := connectionhelper.GetMongoClient()
	if err != nil {
		return result, err
	}
	collection := client.Database(connectionhelper.DB).Collection(connectionhelper.ISSUES)
	err = collection.FindOne(context.TODO(), bson.D{{"_id", objectId}}, options.FindOne().SetProjection(bson.M{"_id": 0})).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, err
}

func GetOneById(id string) (PhoneBookRecord, error) {
	result := PhoneBookRecord{}
	objectId, _ := primitive.ObjectIDFromHex(id)
	client, err := connectionhelper.GetMongoClient()
	if err != nil {
		return result, err
	}
	collection := client.Database(connectionhelper.DB).Collection(connectionhelper.ISSUES)
	err = collection.FindOne(context.TODO(), bson.D{{"_id", objectId}}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, err
}

func DeleteOne(id string) error {
	//Define filter query for fetching specific document from collection
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{primitive.E{Key: "_id", Value: objectId}}
	//Get MongoDB connection using connectionhelper.
	client, err := connectionhelper.GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionhelper.DB).Collection(connectionhelper.ISSUES)
	//Perform DeleteOne operation & validate against the error.
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}
