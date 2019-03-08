package main


import (
    "net/http"
    "fmt"
    _ "time"
    _ "html/template"
    "flag"
    "webserver/mongo"
    "io/ioutil"
    "os"
    "log"
    "os/exec"
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
   var root = flag.String("root", "./" , "file system path")
   
 // templates := template.Must(template.ParseFiles("templates/body.html"))

/*
  http.HandleFunc("/sensori/" , func(w http.ResponseWriter, r *http.Request, deviceid string) {
    sensore := Sensore{"deviceid"}
    if err := templates.ExecuteTemplate(w, "body.html", sensore); err != nil {
         http.Error(w, err.Error(), http.StatusInternalServerError)
      }
   })
*/
   generaFrontIndex()
   fmt.Println("Listening")

   http.Handle("/", http.FileServer(http.Dir(*root)))
   http.HandleFunc("/save/", saveHandler)



   http.ListenAndServe(":8080", nil)
}


func saveHandler(w http.ResponseWriter, r *http.Request,) {

    body := r.FormValue("body") // al momento è vuoto
    r.ParseForm()

    println("Ricevuto dato da: "+r.Form.Get("Device Id") + " Warning:" +  r.Form.Get("warning"))
    mongo.PostTemperature(r.Form.Get("Device Id"), r.Form.Get("timestamp"),r.Form.Get("temperatura"), r.Form.Get("warning") , Client)
    
    //if warning -> aggiorna indexhtml
  


    p := &Page{Title: r.Form.Get("Device Id"), Body: []byte(body)}
    err := p.save()
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    aggiornaTabellaR(r.Form.Get("Device Id"))

    updateIndex(r.Form.Get("Device Id"))

    http.Redirect(w, r, "/sensori/"+r.Form.Get("Device Id"), http.StatusFound)

}




func (p *Page) save() error {
    filenameJpg := p.Title + ".jpg"
    index := "index.html"

    var percorso = "sensori/" 

    //The octal integer literal 0600, passed as the third parameter to WriteFile, indicates that the file should be created with read-write permissions for the current user only
    os.MkdirAll(percorso+p.Title, os.FileMode(0522))

    //controlla se esiste il jpg, in caso contrario crearlo
    
    if _, err := os.Stat(percorso+p.Title+"/"+filenameJpg); err == nil {
      //il file esiste
      
    } else if os.IsNotExist(err) {
      // path/to/whatever does *not* exist

      ioutil.WriteFile(percorso+p.Title+"/"+filenameJpg, p.Body, 0600)
      
      generaFrontIndex()

    } else {
      return err
    }
    
    if _, err := os.Stat(percorso+p.Title+"/"+index); err == nil {
      //il file esiste
    } else if os.IsNotExist(err) {
    //creare index se non esiste
      bodyindex := `
      <!DOCTYPE html>
      <head>
      <link rel="stylesheet" href="/static/stylesheets/template.css">
      </head>
      <body>
      <div class="center"> <p><a href="#" onclick="history.go(-1)"> Torna Indietro</a></p></div>

      <div class="welcome center">Sensore` +  p.Title + `</div> 
      <div>
      <img class="center" src="` +  p.Title + `.jpg" width="600" height="600" />
      </div>      
      </body>`
      ioutil.WriteFile(percorso+p.Title+"/"+index, []byte(bodyindex), 0600)
    } else {
      return err
    }
    return nil
}

func aggiornaTabellaR(id string) {
  _, err := exec.Command("c://PROGRA~1/R/R-3.5.2/bin/x64/Rscript.exe","--vanilla C:/Users/franz/go/src/webserver/R-Handler.R " + id).Output()
  if err != nil {
    log.Fatal(err)
  }

}

func generaFrontIndex(){
    println("Front index generato")
    files, err := ioutil.ReadDir("./sensori/")
    if err != nil {
        log.Fatal(err)
    }
    indirizzi :=""
    for _, f := range files {
     indirizzi += `<p> <a href="http://localhost:8080/sensori/`+ f.Name() + `/" > Sensore #`+ f.Name()  +`</a>  </p>
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
        <div class="welcome center">Frontpage</div>`+ indirizzi  +`
        </body>`
        ioutil.WriteFile("./index.html", []byte(bodyindex), 0600)
      } else {
        log.Fatal(err)
      }

  }
    func updateIndex(deviceID string) {
    
    
    var numeroWarnings = mongo.GetWarnings(deviceID,Client)
    var percorso = "sensori/"+deviceID

      if _, err := os.Stat(percorso+"/index.html"); err == nil {
        //il file esiste

         bodyindex := `
      <!DOCTYPE html>
      <head>
      <link rel="stylesheet" href="/static/stylesheets/template.css">
      </head>
      <body>
      <p><a href="#" onclick="history.go(-1)"> Torna Indietro</a></p>

      <div class="welcome center">Sensore ` +  deviceID + `</div> `

      if (numeroWarnings != "0") {
        bodyindex+=`<h2>  <div class="center"> <font color="red"> Numero di warnings ` + numeroWarnings + `</font></div>  </h2> `
      }
      bodyindex+=`<div>
                  <img class="center" src="` + deviceID + `.jpg" width="600" height="600" />
                  </div>
                  </body>`
      println("Warning del sensore " + deviceID + " aggiornato")
      ioutil.WriteFile(percorso+"/index.html", []byte(bodyindex), 0600)
      } else {
        log.Fatal(err)
      }
    }


