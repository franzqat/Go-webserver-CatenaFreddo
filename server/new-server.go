package main

import (
	"fmt" 
	"io/ioutil" //per scrittura su file
	"log"
	"net/http" // funzionalità da server http
	"os" //permette l'esecuzione come da riga di comando
	"os/exec" //permette l'esecuzione come da riga di comando
	"webserver/mongo" //permette query a mongo
)

//Create a struct that holds information to be displayed in our HTML file
type Page struct {
	Title string
	Body  []byte
}

var Client = mongo.ConnectToMongo()

//Main
func main() {

	creaFrontIndex() //creazione dell'index della frontpage all'avvio del server
	
  fmt.Println("Listening")

	http.Handle("/", http.FileServer(http.Dir("./"))) //inizializza un fileserver nella root
	http.HandleFunc("/save/", saveHandler) //handler delle post al server

	http.ListenAndServe(":8080", nil) // sta in ascolto di chiamate http pronto a servirle
}

//Gestore delle post
func saveHandler(w http.ResponseWriter, r *http.Request) {

	body := r.FormValue("body") // al momento è vuoto
	r.ParseForm() //ParseForm popola r.Form con il contenuto della Request

	println("Ricevuto dato da: " + r.Form.Get("Device Id") + "con Warning:" + r.Form.Get("warning"))
  //post a mongoDB
	mongo.PostTemperature(r.Form.Get("Device Id"), r.Form.Get("timestamp"), r.Form.Get("temperatura"), r.Form.Get("warning"), Client)

	p := &Page{Title: r.Form.Get("Device Id"), Body: []byte(body)}
	err := p.save() //salva la pagina del sensore
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
  //esegue lo script in R
	aggiornaTabellaR(r.Form.Get("Device Id"))

	//se c'è un warning aggiorna index indicando il numero corretto di warning
	if r.Form.Get("warning") != "0" {
		updateIndex(r.Form.Get("Device Id"))
	}

	http.Redirect(w, r, "/sensori/"+r.Form.Get("Device Id"), http.StatusFound)
}

//crea la pagina del sensore se non esiste
//(p *Page) è un ricevitore (receiver) e si sta dichiarando il metodo save associato al ricevitore
func (p *Page) save() error {

	filenameJpg := p.Title + ".jpg"
	index := "index.html"

	var percorso = "sensori/"

	//creazione di tutte le directory fino al percorso finale
	os.MkdirAll(percorso+p.Title, os.FileMode(0522))

	//Controllo esistenza dei files se il sensore è nuovo
	if _, err := os.Stat(percorso + p.Title + "/" + filenameJpg); err == nil {
		//il file esiste
		//non faccio nulla
	} else if os.IsNotExist(err) {
		//il file jpg non esiste
		//crea il jpg vuoto
		ioutil.WriteFile(percorso+p.Title+"/"+filenameJpg, p.Body, 0600)
		//crea l'index
		creaFrontIndex()

	} else {
		return err
	}

	//se l'index del sensore non esiste
	if _, err := os.Stat(percorso + p.Title + "/" + index); err == nil {
		//il file esiste
	} else if os.IsNotExist(err) {
		//creare index se non esiste
    println("Creo l'index del sensore")
		scriviIndexSensore(p.Title,percorso+p.Title+ "/" +index, "0")

	} else {
		return err
	}
	return nil
}

func scriviIndexSensore(deviceID string, path string, warnings string) {
	bodyindex := `
  <!DOCTYPE html>
  <head>
  <link rel="stylesheet" href="/static/stylesheets/template.css">
  </head>

  <body> 
  <div class="welcome center">Sensore ` + deviceID + `</div>
	<h2>  <div class="warning center">  Numero di warnings ` + warnings + `</font></div>  </h2> 
  <div><img class="center" src="` + deviceID + `.jpg" width="600" height="600" /> </div>      
  <p><a href="#" onclick="history.go(-1)"> Torna Indietro</a></p>
  </body>`

	ioutil.WriteFile(path, []byte(bodyindex), 0600)
}

func aggiornaTabellaR(id string) {
	_, err := exec.Command("c://PROGRA~1/R/R-3.5.2/bin/x64/Rscript.exe", "--vanilla C:/Users/franz/go/src/webserver/R-Handler.R "+id).Output()
	if err != nil {
		log.Fatal(err)
	}

}

func creaFrontIndex() {
  os.MkdirAll("./sensori/", os.FileMode(0522))
	//legge il path relativo ./sensori/
	files, err := ioutil.ReadDir("./sensori/")
	if err != nil {
		log.Fatal(err)
	}
	indirizzi := ""
	for _, f := range files {
		indirizzi += `<li><a href="http://localhost:8080/sensori/` + f.Name() + `/" > Sensore #` + f.Name() + `</a></li>
     `
	}

	if _, err := os.Stat("./index.html"); err == nil {
		//il file esiste
		bodyindex := `
        <!DOCTYPE html>
        <head>
      <link rel="stylesheet" href="/static/stylesheets/template.css">
      </head>
        <body>
        <div class="welcome center titolo">Lista dei sensori di temperatura</div>
        <ul>` + indirizzi + `
        </ul>
        </body>`
		ioutil.WriteFile("./index.html", []byte(bodyindex), 0600)
		println("Front index creato")
	} else {
		log.Fatal(err)
	}

}
func updateIndex(deviceID string) {

	var numeroWarnings = mongo.GetWarnings(deviceID, Client)
	var path = "sensori/" + deviceID + "/index.html"

	if _, err := os.Stat(path); err == nil {
		//il file esiste
		scriviIndexSensore(deviceID, path, numeroWarnings)
	} else {
		log.Fatal(err)
	}
}
