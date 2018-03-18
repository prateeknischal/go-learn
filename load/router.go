package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ServiceResponse - default response
type ServiceResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Load - load response struct
type Load struct {
	ServiceResponse ServiceResponse
	Load            int32 `json:"load"`
}

// Router - Main router that returns the configured router
func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/status", GetStatus).Methods("GET")
	router.HandleFunc("/load", GetLoadStats).Methods("GET")
	return router
}

// GetStatus - returns the status
func GetStatus(w http.ResponseWriter, r *http.Request) {
	Add()
	defer Del()
	// simulate random load on the server - greater the load greater the wait
	// sleep time range 2000ms to 11000ms
	sleepTime := rand.Int31n(Get() * 100)
	fmt.Println("Sleeping for : ", sleepTime, "ms, Current Load : ", Get())
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	resp, _ := json.Marshal(ServiceResponse{Status: 200, Message: "OK"})
	w.Write(resp)
}

// GetLoadStats - returns the load statistics
func GetLoadStats(w http.ResponseWriter, r *http.Request) {
	Add()
	defer Del()
	resp, _ := json.Marshal(Load{ServiceResponse: ServiceResponse{Status: 200, Message: "OK"}, Load: Get()})
	w.Write(resp)
}
