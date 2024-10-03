package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Reservation struct {
	ID           int    `json:"id"`
	GuestName    string `json:"guest_name"`
	RoomNumber   string `json:"room_number"`
	CheckInDate  string `json:"check_in_date"`
	CheckOutDate string `json:"check_out_date"`
	Status       string `json:"status"` // "pending", "confirmed", or "canceled"
}

var reservations []Reservation
var currentID = 1

func createReservation(w http.ResponseWriter, r *http.Request) {
	var reservation Reservation
	json.NewDecoder(r.Body).Decode(&reservation)
	reservation.ID = currentID
	currentID++
	reservation.Status = "pending" // Default status
	reservations = append(reservations, reservation)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}
func getReservations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}
func getReservation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for _, reservation := range reservations {
		if reservation.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(reservation)
			return
		}
	}
	http.Error(w, "Reservation not found", http.StatusNotFound)
}
func updateReservation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for index, reservation := range reservations {
		if reservation.ID == id {
			var updatedReservation Reservation
			json.NewDecoder(r.Body).Decode(&updatedReservation)
			updatedReservation.ID = id
			reservations[index] = updatedReservation

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedReservation)
			return
		}
	}
	http.Error(w, "Reservation not found", http.StatusNotFound)
}
func deleteReservation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for index, reservation := range reservations {
		if reservation.ID == id {
			reservations = append(reservations[:index], reservations[index+1:]...)
			fmt.Fprintf(w, "Reservation deleted")
			return
		}
	}
	http.Error(w, "Reservation not found", http.StatusNotFound)
}
func main() {
	r := mux.NewRouter()

	// Route Handlers
	r.HandleFunc("/reservations", createReservation).Methods("POST")
	r.HandleFunc("/reservations", getReservations).Methods("GET")
	r.HandleFunc("/reservations/{id}", getReservation).Methods("GET")
	r.HandleFunc("/reservations/{id}", updateReservation).Methods("PUT")
	r.HandleFunc("/reservations/{id}", deleteReservation).Methods("DELETE")

	// Start the server
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
