package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Model for our Team
type Team struct {
	ID          string `json:"id"`
	ImageNumber string `json:"imagenumber"`
	Name        string `json:"name"`
	League      string `json:"league"`
}

// this slice will act as our temporary database
var teams []Team // will store an array of Rolls'

// This function will convert our slice into json and send it
// to the response stream
func getTeams(resp http.ResponseWriter, request *http.Request) {

	log.Printf("You got hit!\n")
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(teams) // function encodes our slice obj into json and sends
}

// This function will return the Team when found. It must iterate
// through our slice and find the value that contains the same 'id'
// as the one the client is searching for.
func getTeam(resp http.ResponseWriter, request *http.Request) {
	log.Printf("You got hit in the getTeam!\n")
	resp.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	isFound := false // flag to determine if team was found or not
	var team Team    // contains the team

	// iterate through teams and return the one with the same id
	for _, item := range teams {
		if item.ID == params["id"] {
			team = item // store the team found
			isFound = true
			break
		}
	}

	// if we could not find the team then we must let the user know
	if isFound == true {
		json.NewEncoder(resp).Encode(team)
		log.Printf("Team found!")
	} else {
		log.Printf("Team Not found -- Please try again")
		resp.Write([]byte("Could not find the team with the id: " + params["id"]))
	}
}

// This function will allow us to create a team
func createTeam(resp http.ResponseWriter, request *http.Request) {
	log.Printf("You got hit in the createTeam")
	resp.Header().Set("Content-Type", "application/json")

	var newTeam Team // this will contain the temp obj of the new team to add

	json.NewDecoder(request.Body).Decode(&newTeam) // Reads our json data and stores it in our newTeam obj
	newTeam.ID = strconv.Itoa(len(teams) + 1)      // creates and sets our newTeam's id

	teams = append(teams, newTeam)        // add the value to our "DB" of teams
	json.NewEncoder(resp).Encode(newTeam) // this sends the response to the client with our newTeam obj
}

func updateTeam(resp http.ResponseWriter, request *http.Request) {
	log.Printf("You got hit in the updateTeam")
}

func deleteTeam(resp http.ResponseWriter, request *http.Request) {
	log.Printf("You got hit in the deleteRoll")
}

func main() {
	router := mux.NewRouter() // create a new router with mux

	// simulate having something in our "DB"
	teams = append(teams, Team{ID: "1", ImageNumber: "898989", Name: "Real Madrid", League: "La Liga"})

	router.HandleFunc("/soccer/teams", getTeams).Methods("GET") // handle for our defualt entry point to server
	router.HandleFunc("/soccer/teams/{id}", getTeam).Methods("GET")
	router.HandleFunc("/soccer/teams", createTeam).Methods("POST")
	router.HandleFunc("/soccer/teams/{id}", updateTeam).Methods("POST")
	router.HandleFunc("/soccer/teams/{id}", deleteTeam).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router)) // if something goes wrong it will output an error
}
