package main


import (
      "net/http"
     "fmt"
     _ "time"
     _ "html/template"
      "flag"
    
)  

//Create a struct that holds information to be displayed in our HTML file
type Welcome struct {
   Name string
   Time string
}



//Go application entrypoint
func main() {
   var root = flag.String("root", "./sensori" , "file system path")

   fmt.Println("Listening")
   http.Handle("/sensori", http.FileServer(http.Dir(*root)))
   http.ListenAndServe(":8080", nil)
}