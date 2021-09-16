// Package to test rate limiter
package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

const apiURL = "http://localhost:8080/api"

func main() {
	//log.Println("Checking ipLimiter (max 11 events in 1 minute). Should discard 4 requests")
	//for i := 0; i < 13; i++ {
	//	resp, err := http.Get(apiURL)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	if resp.StatusCode != http.StatusOK {
	//		log.Print("Request discarded with status code: ", resp.StatusCode)
	//	}
	//	time.Sleep(1 * time.Second)
	//}

	//log.Println("Checking userLimiter (max 2 events in 1 second). Should discard 3 requests")
	//for i := 0; i < 5; i++ {
	//	client := &http.Client{}
	//	req, err := http.NewRequest("GET", apiURL, nil)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	req.Header.Set("Authorization", "User1")
	//	resp, err := client.Do(req)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	if resp.StatusCode != http.StatusOK {
	//		log.Print("Request discarded with status code: ", resp.StatusCode)
	//	}
	//	time.Sleep(100 * time.Millisecond)
	//}

	log.Println("Checking uriLimiter (max 10 events in 1 second). Should discard 1 request")
	for i := 0; i < 11; i++ {
		client := &http.Client{}
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", "User" + strconv.Itoa(i))
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Print("Request discarded with status code: ", resp.StatusCode)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
