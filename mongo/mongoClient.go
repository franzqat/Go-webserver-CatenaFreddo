package mongo

import (
    "context"
    "fmt"
    "log"

    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/mongo"
  //  "github.com/mongodb/mongo-go-driver/mongo/options"
)

type Trainer struct {
    Name string
    Age  int
    City string
}

func ConnectToMongo() {

    // Rest of the code will go here
    client, err := mongo.Connect(context.TODO(), "mongodb+srv://utente:unict@progettoapl-zkgjt.mongodb.net/test?retryWrites=true")

if err != nil {
    log.Fatal(err)
}

// Check the connection
err = client.Ping(context.TODO(), nil)

if err != nil {
    log.Fatal(err)
}

fmt.Println("Connected to MongoDB!")


/*
//ash := Trainer{"Ash", 10, "Pallet Town"}
misty := Trainer{"Misty", 10, "Cerulean City"}
brock := Trainer{"Brock", 15, "Pewter City"}

*/
collection := client.Database("test").Collection("trainers")

/* 

trainers := []interface{}{misty, brock}
//POST al database
insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
*/

filter := bson.D{{"name", "Misty"}}


/**
* GET
*/
// create a value into which the result can be decoded
var result Trainer

err = collection.FindOne(context.TODO(), filter).Decode(&result)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found a single document: %+v\n", result)



err = client.Disconnect(context.TODO())

if err != nil {
    log.Fatal(err)
}
fmt.Println("Connection to MongoDB closed.")
}