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
	Timestamp   string
	Temperatura string
	Warning     string
}

//si connette al database e restituisce il client connesso
func ConnectToMongo() *mongo.Client {

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

//esegue una query al database e restituisce il numero di warnings 
func GetWarnings(sensorID string, Client *mongo.Client) string {

	findOptions := options.Find()
	collection := Client.Database("test").Collection(sensorID)

	filter := bson.D{{"warning", "1"}} //filtra tutti i messaggi con warning

	var results []*Messaggio

	// Passare filter come filtro consente di matchare tutti i valori che matchano il filtro nella collezione
	cur, err := collection.Find(context.TODO(), filter, findOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Itera il cursore su tutti i risultati trovati
	for cur.Next(context.TODO()) {

		var elem Messaggio
		err := cur.Decode(&elem) //decodifica l'elemento puntato dal cursore
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem) //appende l'elemento al vettore di risultati
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Chiude il cursore una volta finito
	cur.Close(context.TODO())

	return fmt.Sprintf("%d", len(results)) //restiuisce la dimensione del vettore di risultati sotto forma di stringa
}

//esegue la post al database mongo con i valori ricevuti
func PostTemperature(sensorID string, timestamp string, temperature string, warning string, Client *mongo.Client) {

	collection := Client.Database("test").Collection(sensorID)

	msg := Messaggio{timestamp, temperature, warning}
	//POST al database
	_, err := collection.InsertOne(context.TODO(), msg)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Inserito elemento del sensore " + sensorID + " su mongoDB")
	}
}

func Disconnect(Client *mongo.Client) {
	err := Client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
