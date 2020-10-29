package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type FlightDetails struct {
	FlightID    int       `json:"flightID"`
	Destination string    `json:"Destination"`
	DepartFrom  string    `json:"DepartFrom"`
	DepartsAt   time.Time `json:"DepartsAt"`
	Seats       []struct {
		Number   int  `json:"Number"`
		IsLocked bool `json:"IsLocked"`
	} `json:"Seats"`
}

type FlightID struct {
	ID int
}

func getConfig() ([]FlightDetails, error) {
	var flightDetails = []FlightDetails{}
	file, err := os.Open("flight.json")
	if err != nil {
		fmt.Println("Error during opening configuration file: ", err)
		return nil, errors.New("Unable to readfile")
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&flightDetails); err != nil {
		fmt.Println("Error during decoding configuration file: ", err)
		return nil, errors.New("Unable to Decode")
	}
	return flightDetails, nil
}

//read data from json file
func GetFlightDetails(w http.ResponseWriter, r *http.Request) {
	data, err := getConfig()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return

	}
	sendData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return

	}
	w.Write([]byte(sendData))
}

func Equal(a, b []FlightDetails) bool {

	for i, v := range a {
		if v.FlightID != b[i].FlightID {
			return false
		}
	}
	return true
}

//Add data to Json file
func PostFlightDetails(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var Data []FlightDetails
	err = json.Unmarshal(b, &Data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	flightDetail, _ := getConfig()
	if Equal(flightDetail, Data) {

		http.Error(w, "Flight Details already Exist", 500)
		return
	}
	//update file
	flightDetail = append(flightDetail, Data...)
	jsonData, _ := json.Marshal(flightDetail)
	ioutil.WriteFile("flight.json", jsonData, os.ModePerm)
}

//DeleteFlightDetails
func DeleteFlightDetails(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var flilghtID FlightID
	err = json.Unmarshal(b, &flilghtID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	flightDetail, _ := getConfig()

	for i, v := range flightDetail {
		if v.FlightID == flilghtID.ID {
			flightDetail = append(flightDetail[:i], flightDetail[i+1:]...)
		}
	}
	jsonData, _ := json.Marshal(flightDetail)
	ioutil.WriteFile("flight.json", jsonData, os.ModePerm)
}

//UpadateFlightDetials
func UpdateFlightDetails(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var flilghtID FlightDetails
	err = json.Unmarshal(b, &flilghtID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	flightDetail, _ := getConfig()

	for i, _ := range flightDetail {
		attr := &flightDetail[i]
		if attr.FlightID == flilghtID.FlightID {
			attr.Destination = flilghtID.Destination
			attr.DepartFrom = flilghtID.DepartFrom
			attr.DepartsAt = flilghtID.DepartsAt
			attr.Seats = flilghtID.Seats
		}
	}
	jsonData, _ := json.Marshal(flightDetail)
	ioutil.WriteFile("flight.json", jsonData, os.ModePerm)
}

type BookFlight struct {
	FlightID   int
	SeatNumber int
}

//UpadateFlightDetials
func BookFlightFunction(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var bookseat BookFlight
	err = json.Unmarshal(b, &bookseat)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var mutex = &sync.Mutex{}
	mutex.Lock()
	flightDetail, _ := getConfig()
	for i, _ := range flightDetail {
		attr := &flightDetail[i]
		if attr.FlightID == bookseat.FlightID {
			fmt.Println("Check the availability")
			for k, seat := range attr.Seats {
				if seat.Number == bookseat.SeatNumber {

					if seat.IsLocked {
						http.Error(w, "Already Booked", 500)
						return
					}

					attr.Seats[k].IsLocked = true
				}
			}
		}
	}
	jsonData, _ := json.Marshal(flightDetail)
	ioutil.WriteFile("flight.json", jsonData, os.ModePerm)
	mutex.Unlock()

}
func main() {

	router := mux.NewRouter()
	fmt.Println("Listening...")
	router.HandleFunc("/GetFlightDetails", GetFlightDetails)
	router.HandleFunc("/PostFlightDetails", PostFlightDetails)
	router.HandleFunc("/UpdateFlightDetails", UpdateFlightDetails)
	router.HandleFunc("/DeleteFlight", DeleteFlightDetails)
	router.HandleFunc("/BookFlight", BookFlightFunction)

	http.ListenAndServe(":8080", router)
}
