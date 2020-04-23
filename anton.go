package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"anton/helper/appsettings"
	"anton/helper/mongodbhelper"
	"anton/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// this is the
var settings struct {
	Log struct {
		MinFilter string `envconfig:"optional"`
	}
	Mongo struct {
		URL        string `envconfig:"optional"`
		DB         string `envconfig:"optional"`
		Collection string `envconfig:"optional"`
	}
}

var playerParameters []string

func main() {
	log.Printf("%+v", settings)

	playerParameters = []string{"firstname", "lastname", "company", "status", "linkedin", "email", "telnumber"}

	// filling the variable with the settings file and env vars
	if err := appsettings.ReadFromFileAndEnv(&settings); err != nil {
		panic(err)
	}

	// do something with the settings
	log.Print("anton config settings:")
	log.Printf("%+v", settings)
	log.Print("service Anton started")

	r := mux.NewRouter()
	r.HandleFunc("/player", createPlayerHandler).Methods("POST")
	r.HandleFunc("/players", getPlayerHandler).Methods("GET")
	r.HandleFunc("/player/{id}", getPlayerHandlerId).Methods("GET")
	r.HandleFunc("/player/{id}", upDatePlayerHandler).Methods("PUT")
	log.Fatal(http.ListenAndServe(":3000", r))
}

//Get a list of players .
//if search params are defined we will perform a search on the database
func getPlayerHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Request received to get a Player list")
	w.Header().Set("Content-Type", "application/json")

	//expirimental
	if err := r.ParseForm(); err != nil {
		// Handle error
	}

	//Compose the filter
	// our searchables are firstname , lastname, company and status for now
	q := make(bson.M)

	for _, p := range playerParameters {
		if r.Form.Get(p) != "" {
			q[p] = bson.M{"$eq": r.Form.Get(p)}
			log.Print("done")
		}
	}

	var ps []models.Player

	//Connection mongoDB with mongodbhelper class
	collection := mongodbhelper.ConnectDB(settings.Mongo.URL, settings.Mongo.DB, settings.Mongo.Collection)

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), q)

	if err != nil {
		mongodbhelper.GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var p models.Player
		// & character returns the memory address of the following variable.
		err := cur.Decode(&p) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		ps = append(ps, p)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(ps) // encode similar to serialize process.

}

//Retreive a Player by its mongo id
func getPlayerHandlerId(w http.ResponseWriter, r *http.Request) {
	log.Print("Request received to get a Player")
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var p models.Player
	// we get params with mux.
	var params = mux.Vars(r)
	log.Print("params", params)
	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := mongodbhelper.ConnectDB(settings.Mongo.URL, settings.Mongo.DB, settings.Mongo.Collection)

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}

	err := collection.FindOne(context.TODO(), filter).Decode(&p)

	if err != nil {
		mongodbhelper.GetError(err, w)

		return
	}

	json.NewEncoder(w).Encode(p)

}

//update a Player
func upDatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var p models.Player

	collection := mongodbhelper.ConnectDB(settings.Mongo.URL, settings.Mongo.DB, settings.Mongo.Collection)

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&p)

	// u := make(bson.D)

	// for _, p := range playerParameters {
	// 	if r.Form.Get(p) != "" {
	// 		q[p] = bson.M{"$eq": r.Form.Get(p)}
	// 		log.Print("done")
	// 	}
	// }

	partial, err := toDoc(p)
	if err != nil {
		log.Print(err)
	}

	// prepare update model.
	u := bson.D{
		{"$set", partial},
	}

	log.Print("%+v", u)

	err = collection.FindOneAndUpdate(context.TODO(), filter, u).Decode(&p)

	if err != nil {
		mongodbhelper.GetError(err, w)
		return
	}

	p.ID = id

	json.NewEncoder(w).Encode(p)
}

//create a Player
func createPlayerHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var player = models.PlayerDefaults()
	_ = json.NewDecoder(r.Body).Decode(&player)

	log.Print("succesfully received request to create Player", player)

	// connect db
	collection := mongodbhelper.ConnectDB(settings.Mongo.URL, settings.Mongo.DB, settings.Mongo.Collection)

	// insert our player model.
	result, err := collection.InsertOne(context.TODO(), player)

	if err != nil {
		mongodbhelper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// utilities ... might want to move to mongodbhelper
func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func deletePlayerHandler(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])

	collection := mongodbhelper.ConnectDB(settings.Mongo.URL, settings.Mongo.DB, settings.Mongo.Collection)

	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		mongodbhelper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}
