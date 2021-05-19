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

type relationValue struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
type artists struct {
	Id              int      `json:"id"`
	Image           string   `json:"image"`
	Name            string   `json:"name"`
	Members         []string `json:"members"`
	CreationDate    int      `json:"creationDate"`
	Firstalbum      string   `json:"firstAlbum"`
	Locations       string   `json:"locations"`
	ConcertDates    string   `json:"concertDates"`
	Relations       string   `json:"relations"`
	RelationsValues relationValue
}

func LogError(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
func init() {
	tpl = template.Must(template.ParseGlob("assets/html/*.html"))
}
func GetAPI(lien string, adresseArtists *[]artists) {
	// Prend la donnée
	response, err := http.Get(lien)
	LogError(err)
	tmp, err := ioutil.ReadAll(response.Body)
	LogError(err)
	err = json.Unmarshal(tmp, &adresseArtists)
	LogError(err)
}
func getAPIValue(lien string) []byte {
	response, err := http.Get(lien)
	LogError(err)
	tmp, err := ioutil.ReadAll(response.Body)
	LogError(err)
	return tmp
}
func GetValues(data *[]artists) {
	for i := range *data {
		err := json.Unmarshal(getAPIValue((*data)[i].Relations), &(*data)[i].RelationsValues) //marche
		LogError(err)
	}
}

// Error Code à ajouter
func main() {
	// Créer/charge le fichier des erreurs
	file, err := os.OpenFile("Errors.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	LogError(err)
	defer file.Close()
	log.SetOutput(file)
	//	Charge les données de l'API
	var data []artists
	GetAPI("https://groupietrackers.herokuapp.com/api/artists", &data)
	GetValues(&data)
	// Partie serveur
	handleur := func(w http.ResponseWriter, r *http.Request) {
		err := tpl.ExecuteTemplate(w, "index.html", data)
		if LogError(err) {
			http.Error(w, "404 Not Found", 404)
		}
	}
	map_handleur := func(w http.ResponseWriter, r *http.Request) {
		err := tpl.ExecuteTemplate(w, "map.html", nil)
		if LogError(err) {
			http.Error(w, "404 Not Found", 404)
		}

	}
	http.HandleFunc("/", handleur)
	http.HandleFunc("/map", map_handleur)
	//Allow the link of the CSS to the html file
	fs := http.FileServer(http.Dir("CSS"))
	http.Handle("/CSS/", http.StripPrefix("/CSS/", fs))
	err = http.ListenAndServe(":8080", nil)
	LogError(err)
}
