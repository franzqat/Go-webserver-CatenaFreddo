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
)  

//Create a struct that holds information to be displayed in our HTML file
type Welcome struct {
   Name string
   Time string
}

type Page struct {
    Title string
    Body  []byte //https://blog.golang.org/go-slices-usage-and-internals
}
var Client = mongo.ConnectToMongo()


//Go application entrypoint
func main() {
   var root = flag.String("root", "./sensori" , "file system path")

   fmt.Println("Listening")
   http.Handle("/", http.FileServer(http.Dir(*root)))

 http.HandleFunc("/save/", saveHandler)

   http.ListenAndServe(":8080", nil)
}


func saveHandler(w http.ResponseWriter, r *http.Request,) {

    body := r.FormValue("body")
    r.ParseForm()

    println(r.Form.Get("Device Id"))
    mongo.PostTemperature(r.Form.Get("Device Id"), r.Form.Get("timestamp"),r.Form.Get("temperatura") , Client)

    p := &Page{Title: r.Form.Get("Device Id"), Body: []byte(body)}
    err := p.save()
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    http.Redirect(w, r, "/sensori/"+r.Form.Get("Device Id"), http.StatusFound)
}

func (p *Page) save() error {
    filenameJpg := p.Title + ".jpg"
    index := "index.html"

    var percorso = "sensori/" 

    //The octal integer literal 0600, passed as the third parameter to WriteFile, indicates that the file should be created with read-write permissions for the current user only
    os.MkdirAll(percorso+p.Title, os.FileMode(0522))

    //TODO: controllare se esiste il jpg, in caso contrario crearlo
    //creare index se non esiste
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
      // path/to/whatever does *not* exist
      ioutil.WriteFile(percorso+p.Title+"/"+index, p.Body, 0600)
    } else {
      return err
    }
    return nil
}
