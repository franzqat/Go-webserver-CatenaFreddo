package mongo

import (
    "context"
    "fmt"
    "log"

    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/mongo"
    "github.com/mongodb/mongo-go-driver/mongo/options"
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


func GetWarnings(sensorID string, Client *mongo.Client) string {
    // Pass these options to the Find method
  findOptions := options.Find()
  findOptions.SetLimit(50)

  collection := Client.Database("test").Collection(sensorID)

  filter := bson.D{{"warning", "1"}}

  // Here's an array in which you can store the decoded documents
  var results []*Messaggio

  // Passing nil as the filter matches all documents in the collection
  cur, err := collection.Find(context.TODO(), filter, findOptions)

  if err != nil {
      log.Fatal(err)
  }

  // Finding multiple documents returns a cursor
  // Iterating through the cursor allows us to decode documents one at a time
  for cur.Next(context.TODO()) {

      // create a value into which the single document can be decoded
      var elem Messaggio
      err := cur.Decode(&elem)
      if err != nil {
          log.Fatal(err)
      }
      results = append(results, &elem)
  }

  if err := cur.Err(); err != nil {
      log.Fatal(err)
  }

  // Close the cursor once finished
  cur.Close(context.TODO())

  return fmt.Sprintf("%d", len(results))  
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