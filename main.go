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

type Locations struct {
	Locations []string `json:"locations"`
}
type ConcertDates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}
type Artists struct {
	Id                int          `json:"id"`
	Image             string       `json:"image"`
	Name              string       `json:"name"`
	Members           []string     `json:"members"`
	CreationDate      int          `json:"creationDate"`
	Firstalbum        string       `json:"firstAlbum"`
	Locations         string       `json:"locations"`
	ConcertDates      string       `json:"concertDates"`
	Relations         string       `json:"relations"`
	LocationsValue    Locations    `json:"locationsValues"`
	ConcertDatesValue ConcertDates `json:"concertDatesValue"`
}

func init() {
	tpl = template.Must(template.ParseGlob("*.html"))
}
func getAPI(lien string, adresseArtists *[]Artists) {
	// Prend la donnée
	response, err := http.Get(lien)
	if err != nil {
		log.Println(err)
	}
	tmp, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(tmp, &adresseArtists)
	if err != nil {
		log.Println(err)
	}
}
func getAPIValue(lien string) []byte {
	response, err := http.Get(lien)
	if err != nil {
		log.Println(err)
	}
	tmp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	return tmp
}
func getValues(data *[]Artists) {
	for i := range *data {
		err := json.Unmarshal(getAPIValue((*data)[i].Locations), &(*data)[i].LocationsValue) //marche
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(getAPIValue((*data)[i].ConcertDates), &(*data)[i].ConcertDatesValue)
		if err != nil {
			log.Println(err)
		}
	}
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

	var data []Artists
	getAPI("https://groupietrackers.herokuapp.com/api/artists", &data)
	getValues(&data)

	// Partie serveur
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = tpl.ExecuteTemplate(w, "index.html", data)
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
