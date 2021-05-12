package grab

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"

type MyArtistFull struct {
	ID             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Locations      []string            `json:"locations"`
	ConcertDates   []string            `json:"concertDates"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type MyArtist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type MyLocation struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type MyRelation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type MyDate struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type MyDates struct {
	Index []MyDate `json:"index"`
}

type MyLocations struct {
	Index []MyLocation `json:"index"`
}
type MyRelations struct {
	Index []MyRelation `json:"index"`
}

var ArtistsFull []MyArtistFull
var Artists []MyArtist
var Dates MyDates
var Locations MyLocations
var Relations MyRelations

func GetArtistsData() error {
	resp, err := http.Get(baseURL + "/artists")
	if err != nil {
		return errors.New("Error by get")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Artists)
	return nil
}

func GetDatesData() error {
	resp, err := http.Get(baseURL + "/dates")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Dates)
	return nil
}

func GetLocationsData() error {
	resp, err := http.Get(baseURL + "/locations")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Locations)
	return nil
}

func GetRelationsData() error {
	resp, err := http.Get(baseURL + "/relation")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Relations)
	return nil
}

func GetData() error {
	if len(ArtistsFull) != 0 {
		return nil
	}
	err1 := GetArtistsData()
	err2 := GetLocationsData()
	err3 := GetDatesData()
	err4 := GetRelationsData()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return errors.New("Error by get data artists, locations, dates")
	}
	for i := range Artists {
		var tmpl MyArtistFull
		tmpl.ID = i + 1
		tmpl.Image = Artists[i].Image
		tmpl.Name = Artists[i].Name
		tmpl.Members = Artists[i].Members
		tmpl.CreationDate = Artists[i].CreationDate
		tmpl.FirstAlbum = Artists[i].FirstAlbum
		tmpl.Locations = Locations.Index[i].Locations
		tmpl.ConcertDates = Dates.Index[i].Dates
		tmpl.DatesLocations = Relations.Index[i].DatesLocations
		ArtistsFull = append(ArtistsFull, tmpl)
	}
	return nil
}
