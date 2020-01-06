package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"fmt"
)

//DBURL ...
type DBURL struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	LongURL string `bson:"longURL" json:"longURL"`
	ShortURL string `bson:"shortURL" json:"shortURL"`
}


func (_dbURL *DBURL) findInDB(m interface{})(error){
	fmt.Printf("%v sent map",m)
	ctx := context.TODO()
	coll := env.db.Collection(env.urlCollection)
	findResult := coll.FindOne(ctx,m)
	if err := findResult.Err(); err != nil{ // not found in db
		fmt.Printf("find result for %v, %v\n",m,err)
		return err
	}
	err := findResult.Decode(_dbURL)
	if err != nil{
		fmt.Printf("decode find result %v\n",err)
		return err
	}
	return nil
}

func (_dbURL *DBURL) insertInDB() (*mongo.InsertOneResult,error) {
	ctx := context.TODO()
	coll := env.db.Collection(env.urlCollection)
	_dbURL.ID = primitive.NewObjectID()
	result, err := coll.InsertOne(ctx,_dbURL)
		if err != nil{
			fmt.Printf("insert result %v\n",err)
			return nil ,err
		}
		return result ,nil
}


func (_dbURL *DBURL) addShortURLToDB() error{
	fmt.Printf("sent obj %v\n",_dbURL)
	doc := bson.M{"longURL":_dbURL.LongURL}
	err := _dbURL.findInDB(doc)
	if err != nil{ // not found in db
		fmt.Printf("find result %v\n",err)
		result ,err := _dbURL.insertInDB()
		if err != nil{
			return err
		}
		_id := result.InsertedID.(primitive.ObjectID)
		err = _dbURL.findInDB(bson.D{primitive.E{Key :"_id", Value :
	_id}})
		if err != nil{
			return err
		}
	}
	return nil
}
