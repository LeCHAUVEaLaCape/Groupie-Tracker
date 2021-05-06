package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var tpl *template.Template

type Artists struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	Firstalbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

func init() {
	tpl = template.Must(template.ParseGlob("*.html"))
}

// Error Code à ajouter
func main() {
	// Créer le fichier des erreurs
	file, err := os.OpenFile("Errors.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	// Prend la donnée
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Println(err)
	}
	tmp, _ := ioutil.ReadAll(response.Body)
	var articles []Artists // la donnée
	err = json.Unmarshal(tmp, &articles)
	if err != nil {
		log.Println(err)
	}

	// Partie serveur
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = tpl.ExecuteTemplate(w, "index.html", articles)
		if err != nil {
			log.Println(err)
		}
	})
	fs := http.FileServer(http.Dir("CSS")) //Allow the link of the CSS to the html file
	http.Handle("/CSS/", http.StripPrefix("/CSS/", fs))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}
}
