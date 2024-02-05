package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Movie struct represents a movie with its attributes
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// Director struct represents a director with their name
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// movies is a global slice to store movie data (for simplicity)
var movies []Movie

// getMovies retrieves all movies and encodes them as JSON
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		log.Println("Error encoding movies:", err)
		http.Error(w, "Error retrieving movies", http.StatusInternalServerError)
		return
	}
}

// deleteMovie deletes a movie with the specified ID
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	foundIndex := -1
	for i, movie := range movies {
		if movie.ID == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	movies = append(movies[:foundIndex], movies[foundIndex+1:]...)
}

// getMovie retrieves a movie with the specified ID
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	for _, movie := range movies {
		if movie.ID == id {
			err := json.NewEncoder(w).Encode(movie)
			if err != nil {
				log.Println("Error encoding movie:", err)
				http.Error(w, "Error retrieving movie", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "Movie not found", http.StatusNotFound)
}

// createMovie creates a new movie and adds it to the movies slice
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		log.Println("Error decoding movie:", err)
		http.Error(w, "Error creating movie", http.StatusBadRequest)
		return
	}

	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)

	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		log.Println("Error encoding movie:", err)
		http.Error(w, "Error creating movie", http.StatusInternalServerError)
		return
	}
}

// updateMovie updates a movie with the specified ID
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	updatedMovie := Movie{}
	err := json.NewDecoder(r.Body).Decode(&updatedMovie)
	if err != nil {
		log.Println("Error decoding updated movie:", err)
		http.Error(w, "Error updating movie", http.StatusBadRequest)
		return
	}

	for i, movie := range movies {
		if movie.ID == id {
			movies = append(movies[:i], movies[i+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = id
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie Allok", Director: &Director{Firstname: "Allok", Lastname: "Raj"}})
	movies = append(movies, Movie{ID: "2", Isbn: "438228", Title: "Movie Apurve", Director: &Director{Firstname: "Apurv", Lastname: "Raj"}})
	movies = append(movies, Movie{ID: "3", Isbn: "438229", Title: "Movie Raj", Director: &Director{Firstname: "Jigyasha", Lastname: "Raj"}})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
