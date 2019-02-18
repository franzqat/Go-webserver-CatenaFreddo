package mongo

import (
    "context"
    "fmt"
    "log"

    _ "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/mongo"
  //  "github.com/mongodb/mongo-go-driver/mongo/options"
)

type Messaggio struct {
    Timestamp  string
    Temperatura string
}


func ConnectToMongo() (*mongo.Client) {

    // Rest of the code will go here
    Client, err := mongo.Connect(context.TODO(), "mongodb+srv://utente:unict@progettoapl-zkgjt.mongodb.net/test?retryWrites=true")

    if err != nil {
        log.Fatal(err)
    }

    // Check the connection
    err = Client.Ping(context.TODO(), nil)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")
    return Client
}

func PostTemperature(sensorID string, timestamp string, temperature string, Client *mongo.Client) {


    collection := Client.Database("test").Collection(sensorID)


    msg := Messaggio{timestamp, temperature}
    //POST al database
    insertResult, err := collection.InsertOne(context.TODO(), msg)
    if err != nil {
        log.Fatal(err)
    } else {
    fmt.Println("Inserted a single document: ", insertResult.InsertedID)
    }
}

/*
func Get(sensorID string){
    filter := bson.D{{"name", "Misty"}}
    collection := client.Database("test").Collection(sensorID)
    /**
    * GET
    */
    // create a value into which the result can be decoded
 /*   var result Messaggio

    err = collection.FindOne(context.TODO(), filter).Decode(&result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found a single document: %+v\n", result)

*/

func Disconnect(Client *mongo.Client){
    err := Client.Disconnect(context.TODO())


    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connection to MongoDB closed.")
}