package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"encoding/json"
	"errors"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"encoding/hex"
	"encoding/base64"
	"crypto/md5"
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"log"
)



//ENV ...
type ENV struct{
	mongoClient *mongo.Client
	db *mongo.Database
	urlCollection string
}

var env  = &ENV{urlCollection : "urls"}


func shortURLHandler(rw http.ResponseWriter, r *http.Request){
	v := mux.Vars(r);
	shortURL := v["shortURL"]
	_dbURL := &DBURL{}
	idDoc := bson.D{primitive.E{Key :"shortURL", Value :shortURL}}
	err := _dbURL.findInDB(idDoc)
	if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(rw,r,_dbURL.LongURL,http.StatusMovedPermanently)
}

func getMD5Hash(text string) string{
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getBase64(text string) string{
	md5Hash := getMD5Hash(text)
	b64 := base64.URLEncoding.EncodeToString([]byte(md5Hash))
	shortURL := []byte(b64)
	return string(shortURL[0:6])
}




func getShortURLHandler(rw http.ResponseWriter, r *http.Request) {
	var _dbURL DBURL
	var mr *malformedRequest
	err := decodeJSONBody(rw,r,&_dbURL)
	if err != nil{
		if errors.As(err,&mr){
			http.Error(rw,mr.msg,mr.status)
			return
		}
		log.Println(err.Error())
		http.Error(rw,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}
	base64String := getBase64(_dbURL.LongURL)
	_dbURL.ShortURL = base64String
	err = (&_dbURL).addShortURLToDB()
	js, err := json.Marshal(_dbURL)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)	
}

func connectDB() {
	// Base context.
	ctx := context.TODO()
	// Options to the database.
	clientOpts := options.Client().ApplyURI("mongodb://mongo_user_1:Password1@ds255403.mlab.com:55403/url-shortner?retryWrites=false")
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
			fmt.Println(err)
			return
	}
	if err != nil {
    log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
			log.Fatal(err)
	}

	env.mongoClient = client
	db := client.Database("url-shortner")
	env.db = db
	fmt.Println(db.Name())

	fmt.Println("Connected to MongoDB!")
}

func main(){
	connectDB()
	r := mux.NewRouter()
	r.HandleFunc("/{shortURL}",shortURLHandler).Methods("GET")
	r.HandleFunc("/getShortURL",getShortURLHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080",r))
}
