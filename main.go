package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"

	grab "./data"
)

var data []grab.MyArtistFull

func mainPage(w http.ResponseWriter, r *http.Request) {
	err := grab.GetData()
	if err != nil {
		errors.New("Error by get data")
	}
	if !(len(data) != 0) {
		data = Datass()
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}

func Datass() []grab.MyArtistFull {
	return grab.ArtistsFull
}

func main() {

	file, err := os.OpenFile("Errors.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	http.HandleFunc("/", mainPage)

	fs := http.FileServer(http.Dir("CSS")) //Allow the link of the CSS to the html file
	http.Handle("/CSS/", http.StripPrefix("/CSS/", fs))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}
}
