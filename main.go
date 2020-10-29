package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type FlightDetails []struct {
	FlightID    int       `json:"flightID"`
	Destination string    `json:"Destination"`
	DepartFrom  string    `json:"DepartFrom"`
	DepartsAt   time.Time `json:"DepartsAt"`
	Seats       []struct {
		Number   int  `json:"Number"`
		IsLocked bool `json:"IsLocked"`
	} `json:"Seats"`
}
type FlightDetail struct {
	FlightID    int       `json:"flightID"`
	Destination string    `json:"Destination"`
	DepartFrom  string    `json:"DepartFrom"`
	DepartsAt   time.Time `json:"DepartsAt"`
	Seats       []struct {
		Number   int  `json:"Number"`
		IsLocked bool `json:"IsLocked"`
	} `json:"Seats"`
}

var (
	flightDetails = FlightDetails{}
)

func getConfig() (FlightDetails, error) {
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

func GetFlightDetails(w http.ResponseWriter, r *http.Request) {

	data, _ := getConfig()
	sendData, _ := json.Marshal(data)
	w.Write([]byte(sendData))
}

func PostFlightDetails(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var Data FlightDetails
	err = json.Unmarshal(b, &Data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	flightDetail, _ := getConfig()
	flightDetail = append(flightDetail, Data...)
	//update file
	jaonData, _ := json.Marshal(flightDetail)

	ioutil.WriteFile("flight.json", jaonData, os.ModePerm)
}

type FlightID struct {
	ID int
}

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
	fmt.Println("++++", flilghtID)

	flightDetail, _ := getConfig()

	for i, v := range flightDetail {
		if v.FlightID == flilghtID.ID {
			flightDetail = append(flightDetail[:i], flightDetail[i+1:]...)

		}

	}
	jaonData, _ := json.Marshal(flightDetail)

	ioutil.WriteFile("flight.json", jaonData, os.ModePerm)
}

func UpdateFlightDetails(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var flilghtID FlightDetail
	err = json.Unmarshal(b, &flilghtID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("++++", flilghtID)

	flightDetail, _ := getConfig()
	flightDetail1 := flightDetail[:0]

	for _, v := range flightDetail {
		if v.FlightID == flilghtID.FlightID {
			flightDetail1 = append(flightDetail1, flilghtID)
		}
	}
	jaonData, _ := json.Marshal(flightDetail1)
	ioutil.WriteFile("flight.json", jaonData, os.ModePerm)
}
func main() {
	fmt.Println("Listening...")
	http.HandleFunc("/GetFlightDetails", GetFlightDetails)
	http.HandleFunc("/PostFlightDetails", PostFlightDetails)
	http.HandleFunc("/UpdateFlightDetails", UpdateFlightDetails)
	http.HandleFunc("/DeleteFlight", DeleteFlightDetails)
	http.HandleFunc("/BookFlight", DeleteFlightDetails)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
