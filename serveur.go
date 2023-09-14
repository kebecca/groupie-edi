package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Structure qui gère les artistes
type Artists struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
}

// Structure qui gère les locations
type Locations struct {
	Id        int
	Locations []string
	Dates     string
}

// Structure qui gère les relations
type Relations struct {
	Id             int
	DatesLocations map[string][]string
}

// Structure qui gère les dates
type Dates struct {
	Id    int
	Dates []string
}

// Structure qui contient toutes les données nécessaires pour la page de détail
type DetailData struct {
	Artist   Artists
	Location Locations
	Date     Dates
	Relation Relations
}

var (
	template1 = template.Must(template.ParseFiles("all-artists.html"))
	template2 = template.Must(template.ParseFiles("personnal-artist.html"))
	templates = template.Must(template.ParseFiles("error.html"))
)

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/accueil", all)
	http.HandleFunc("/artist", unique)
	fmt.Println("Cliquez sur le lien suivant : http://localhost:9999/accueil")
	http.ListenAndServe(":9999", nil)
}

func all(w http.ResponseWriter, r *http.Request) {
	// Récupération des données des artistes depuis une API
	responseArtiste, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	//Cette ligne lit le corps (body) de la réponse HTTP (qui contient les données des artistes) en utilisant ioutil.ReadAll. Le résultat est stocké dans responseDataArtist
	//ioutil.ReadAll lit tout le contenu d'un io.Reader (dans ce cas, le corps de la réponse HTTP) et le retourne sous forme d'un tableau de bytes ([]byte
	responseDataArtiste, err := ioutil.ReadAll(responseArtiste.Body)
	if err != nil {
		log.Fatal(err)
	}
	//Cette ligne convertit les données JSON récupérées en une chaîne de caractères (string) pour faciliter le traitement ultérieur.
	//artisteJSON := string(responseDataArtiste)

	//Cette ligne déclare une variable artistes qui est une slice de Artists. Cette slice sera utilisée pour stocker les données des artistes une fois qu'elles auront été décodées à partir du JSON
	var artistes []Artists
	// Cette ligne utilise json.Unmarshal pour décoder la chaîne JSON en une liste de structures Artists et stocke le résultat dans la variable artistes
	json.Unmarshal(responseDataArtiste, &artistes)

	err1 := template1.ExecuteTemplate(w, "all-artists.html", artistes)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}
func unique(w http.ResponseWriter, r *http.Request) {
	index := r.URL.Query().Get("id")

	var artists []Artists
	var locations Locations
	var dates Dates
	var relations Relations

	// Récupération des données de l'artiste
	responseArtist, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Println("Erreur lors de la requête pour les artistes :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseDataArtist, err := ioutil.ReadAll(responseArtist.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture de la réponse artiste :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Conversion de la réponse JSON en une liste d'artistes
	if err := json.Unmarshal(responseDataArtist, &artists); err != nil {
		log.Println("Erreur lors de la conversion de la réponse artiste JSON en struct :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Recherche de l'artiste avec l'ID spécifié
	var artist Artists
	indexArtist, _ := strconv.Atoi(index)
	for _, a := range artists {
		if a.Id == indexArtist {
			artist = a
			break
		}
	}

	if artist.Id < 1 || artist.Id > 52 {
		err = templates.ExecuteTemplate(w, "error.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Récupération des données des lieux
	responseLocation, err := http.Get("https://groupietrackers.herokuapp.com/api/locations/" + index)
	if err != nil {
		log.Println("Erreur lors de la requête pour les lieux :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseDataLocation, err := ioutil.ReadAll(responseLocation.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture de la réponse lieu :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	//locationJSON := string(responseDataLocation)

	// Conversion des données JSON en une liste de lieux
	if err := json.Unmarshal(responseDataLocation, &locations); err != nil {
		log.Println("Erreur lors de la conversion de la réponse lieu JSON en struct :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Récupération des données des dates
	responseDates, err := http.Get("https://groupietrackers.herokuapp.com/api/dates/" + index)
	if err != nil {
		log.Println("Erreur lors de la requête pour les dates :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseDataDates, err := ioutil.ReadAll(responseDates.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture de la réponse dates :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	//datesJSON := string(responseDataDates)

	// Conversion des données JSON en une liste de detas
	if err := json.Unmarshal(responseDataDates, &dates); err != nil {
		log.Println("Erreur lors de la conversion de la réponse dates JSON en struct :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Récupération des données de la relation
	responseRelation, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + index)
	if err != nil {
		log.Println("Erreur lors de la requête pour les relations :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseDataRelation, err := ioutil.ReadAll(responseRelation.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture de la réponse relation :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	//relationJSON := string(responseDataRelation)

	// Conversion des données JSON en une liste de relation
	if err := json.Unmarshal(responseDataRelation, &relations); err != nil {
		log.Println("Erreur lors de la conversion de la réponse relation JSON en struct :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Création de la structure de données pour la page de détail
	detailData := DetailData{
		Artist:   artist,
		Location: locations,
		Date:     dates,
		Relation: relations,
	}

	// Exécution du modèle et affichage de la page de détail
	err = template2.ExecuteTemplate(w, "personnal-artist.html", detailData)
	if err != nil {
		log.Println("Erreur lors de l'exécution du modèle :", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}
