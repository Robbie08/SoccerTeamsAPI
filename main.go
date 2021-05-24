package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"

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

	json.NewDecoder(request.Body).Decode(&newTeam)     // Reads our json data and stores it in our newTeam obj
	newTeam.ID = fmt.Sprint(GenerateUID(newTeam.Name)) // creates and sets our newTeam's id

	teams = append(teams, newTeam)        // add the value to our "DB" of teams
	json.NewEncoder(resp).Encode(newTeam) // this sends the response to the client with our newTeam obj
}

// This Function will allow us to update a team based on the given
// id provided by the request.
func updateTeam(resp http.ResponseWriter, request *http.Request) {
	log.Printf("You got hit in the updateTeam")
	resp.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request) // get our paramters from the request

	isFound := false // flag for if the value is found

	// iterate through teams and return the one with the same id
	for i, item := range teams {
		if item.ID == params["id"] {
			var team Team                               // will contain the team we wish to update
			teams = append(teams[:i], teams[i+1:]...)   // this removes the team from our slice
			json.NewDecoder(request.Body).Decode(&team) // modified the team with the passed in params from request
			team.ID = params["id"]                      // we want to keep the same id for this updated team
			teams = append(teams, team)                 // add new team to the array
			json.NewEncoder(resp).Encode(team)
			isFound = true
			break
		}
	}

	// if the team wasn't found let the client and console know
	if isFound == false {
		log.Printf("Update Error: The team with that give id is not in our database.")
		resp.Write([]byte("Could not find the team with the id: " + params["id"]))
	}

}

// This function will remove an item from our
func deleteTeam(resp http.ResponseWriter, request *http.Request) {
	log.Printf("You got hit in the deleteRoll")
	resp.Header().Set("Content-Type", "application/json")

	params := mux.Vars(request)
	isFound := false
	for i, item := range teams {
		if item.ID == params["id"] {
			teams = append(teams[:i], teams[i+1:]...) // delete team from slice
			isFound = true
			break
		}
	}

	// if the item was removed then we can respond with the updated teams list
	if isFound == true {
		json.NewEncoder(resp).Encode(teams)
	} else {
		log.Printf("Update Error: The team with that give id is not in our database.")
		resp.Write([]byte("Could not find the team with the id: " + params["id"]))
	}
}

// This function will generate a Unique ID by using fnv lib.
// @params str string: Name of the team we wish to create a UID for
// @return uint64 : will be the UID for the given team as a uint64
func GenerateUID(str string) uint64 {
	uid := fnv.New64a()
	uid.Write([]byte(str))
	return uid.Sum64()
}

func main() {
	router := mux.NewRouter() // create a new router with mux

	// simulate having something in our "DB"
	teams = append(teams, Team{ID: fmt.Sprint(GenerateUID("Real Madrid")), ImageNumber: "898989", Name: "Real Madrid", League: "La Liga"})

	router.HandleFunc("/soccer/teams", getTeams).Methods("GET") // handle for our defualt entry point to server
	router.HandleFunc("/soccer/teams/{id}", getTeam).Methods("GET")
	router.HandleFunc("/soccer/teams", createTeam).Methods("POST")
	router.HandleFunc("/soccer/teams/{id}", updateTeam).Methods("POST")
	router.HandleFunc("/soccer/teams/{id}", deleteTeam).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router)) // if something goes wrong it will output an error
}
