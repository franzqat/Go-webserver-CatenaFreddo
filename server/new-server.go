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
   var root = flag.String("root", "./sensori" , "file system path")
   
 // templates := template.Must(template.ParseFiles("templates/body.html"))

  http.Handle("/front", http.FileServer(http.Dir("sensori"))) 

/*
  http.HandleFunc("/sensori/" , func(w http.ResponseWriter, r *http.Request, deviceid string) {
    sensore := Sensore{"deviceid"}
    if err := templates.ExecuteTemplate(w, "body.html", sensore); err != nil {
         http.Error(w, err.Error(), http.StatusInternalServerError)
      }
   })
*/

   fmt.Println("Listening")
   http.Handle("/", http.FileServer(http.Dir(*root)))
   http.HandleFunc("/save/", saveHandler)

   // fs := http.FileServer(http.Dir("./sensori"))

  //  http.HandleFunc("/sensori",handler)

  //  http.Handle("/front", adaptFileServer(fs))


   http.ListenAndServe(":8080", nil)
}


func saveHandler(w http.ResponseWriter, r *http.Request,) {

    body := r.FormValue("body") // al momento è vuoto
    r.ParseForm()

    println(r.Form.Get("Device Id"))
    mongo.PostTemperature(r.Form.Get("Device Id"), r.Form.Get("timestamp"),r.Form.Get("temperatura") , Client)
    
    

    p := &Page{Title: r.Form.Get("Device Id"), Body: []byte(body)}
    err := p.save()
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    aggiornaTabellaR(r.Form.Get("Device Id"))

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
    } else {
      return err
    }
    
    if _, err := os.Stat(percorso+p.Title+"/"+index); err == nil {
      //il file esiste
    } else if os.IsNotExist(err) {
    //creare index se non esiste
      ioutil.WriteFile(percorso+p.Title+"/"+index, nil, 0600)
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
