package main

import (
	"flag"
	"fmt"
	_ "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	_ "time"
	"webserver/mongo"
)

//Create a struct that holds information to be displayed in our HTML file

type Page struct {
	Title string
	Body  []byte
}

var Client = mongo.ConnectToMongo()

/*
var base = template.Must(template.New("base").Parse("header\n{{template \"content\"}}\nfooter"))
//var content1 = template.Must(template.Must(base.Clone()).Parse(`{{define "content"}}<img src="`+ r.Form.Get("Device Id") + `.jpg" width="600" height="600" alt="My Pic">{{end}}`))

func handler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("body.html")
    t.Execute(w, "Body: Hi this is my body")
}*/

//Go application entrypoint
func main() {
	var root = flag.String("root", "./", "file system path")

	// templates := template.Must(template.ParseFiles("templates/body.html"))

	/*
	   http.HandleFunc("/sensori/" , func(w http.ResponseWriter, r *http.Request, deviceid string) {
	     sensore := Sensore{"deviceid"}
	     if err := templates.ExecuteTemplate(w, "body.html", sensore); err != nil {
	          http.Error(w, err.Error(), http.StatusInternalServerError)
	       }
	    })
	*/
	creaFrontIndex()
	fmt.Println("Listening")

	http.Handle("/", http.FileServer(http.Dir(*root)))
	http.HandleFunc("/save/", saveHandler)

	http.ListenAndServe(":8080", nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

	body := r.FormValue("body") // al momento è vuoto
	r.ParseForm()

	println("Ricevuto dato da: " + r.Form.Get("Device Id") + " Warning:" + r.Form.Get("warning"))
	mongo.PostTemperature(r.Form.Get("Device Id"), r.Form.Get("timestamp"), r.Form.Get("temperatura"), r.Form.Get("warning"), Client)

	p := &Page{Title: r.Form.Get("Device Id"), Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aggiornaTabellaR(r.Form.Get("Device Id"))

	//if warning -> aggiorna index
	if r.Form.Get("warning") != "0" {
		updateIndex(r.Form.Get("Device Id"))
	}

	http.Redirect(w, r, "/sensori/"+r.Form.Get("Device Id"), http.StatusFound)

}

func (p *Page) save() error {
	filenameJpg := p.Title + ".jpg"
	index := "index.html"

	var percorso = "sensori/"

	//The octal integer literal 0600, passed as the third parameter to WriteFile, indicates that the file should be created with read-write permissions for the current user only
	os.MkdirAll(percorso+p.Title, os.FileMode(0522))

	//Controllo esistenza dei files se il sensore è nuovo
	if _, err := os.Stat(percorso + p.Title + "/" + filenameJpg); err == nil {
		//il file esiste
		//non faccio nulla
	} else if os.IsNotExist(err) {
		//il file jpg non esiste
		//crea il jpg
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
  <p><a href="#" onclick="history.go(-1)"> Torna Indietro</a></p><div class="welcome center">Sensore ` + deviceID + `</div>
	<h2>  <div class="center"> <font color="red"> Numero di warnings ` + warnings + `</font></div>  </h2> 
  <div><img class="center" src="` + deviceID + `.jpg" width="600" height="600" /> </div>      
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
	//legge il path relativo ./sensori/
	files, err := ioutil.ReadDir("./sensori/")
	if err != nil {
		log.Fatal(err)
	}
	indirizzi := ""
	for _, f := range files {
		indirizzi += `<p> <a href="http://localhost:8080/sensori/` + f.Name() + `/" > Sensore #` + f.Name() + `</a>  </p>
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
        <div class="welcome center">Frontpage</div>` + indirizzi + `
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
