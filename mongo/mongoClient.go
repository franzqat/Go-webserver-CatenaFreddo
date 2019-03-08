package mongo

import (
    "context"
    "fmt"
    "log"

    _ "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/mongo"
)

type Messaggio struct {
    Timestamp  string
    Temperatura string
    Warning string
}


func ConnectToMongo() (*mongo.Client) {
    //"mongodb+srv://utente:unict@progettoapl-zkgjt.mongodb.net/test?retryWrites=true"
    Client, err := mongo.Connect(context.TODO(), "mongodb://127.0.0.1:27017/admin")

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

func PostTemperature(sensorID string, timestamp string, temperature string, warning string, Client *mongo.Client) {

    collection := Client.Database("test").Collection(sensorID)


    msg := Messaggio{timestamp, temperature,warning}
    //POST al database
    insertResult, err := collection.InsertOne(context.TODO(), msg)
    if err != nil {
        log.Fatal(err)
    } else {
    fmt.Println("Inserted a single document: ",  insertResult.InsertedID, msg.Warning)
    }
}


func Disconnect(Client *mongo.Client){
    err := Client.Disconnect(context.TODO())


    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connection to MongoDB closed.")
}