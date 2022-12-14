package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type id struct {
	Id string
}

type points struct {
	Total int
}

type receipt_id struct {
	Uuid id
	Path string
}

type item struct {
	Description string `json:"shortDescription"`
	Price       string `json:"price"`
}

type receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []item `json:"items"`
}

var id_for_list_points receipt_id

// Creates a JSON with a UUID in it for it's only field
func generateId(path string) id {
	id1 := uuid.New()
	uuid := id{id1.String()}
	return uuid
}

// Needs to create and return a struct (to be turned into a JSON)
func listPoints(passed_id receipt_id) int {
	string_path := passed_id.Path
	data, err := ioutil.ReadFile(string_path)
	if err != nil {
		fmt.Print(err)
	}

	var s_receipt receipt
	error := json.Unmarshal([]byte(data), &s_receipt)

	if error != nil {
		fmt.Println("JSON decode error!")
	}

	val := 0

	// One point for every alphanumeric character in the retailer name.
	fmt.Println()
	for i := 0; i < len(s_receipt.Retailer); i++ {
		if ('a' <= s_receipt.Retailer[i] && s_receipt.Retailer[i] <= 'z') ||
			('A' <= s_receipt.Retailer[i] && s_receipt.Retailer[i] <= 'Z') ||
			(s_receipt.Retailer[i] >= '0' && s_receipt.Retailer[i] <= '9') {
			val++
		}
	}

	// 6 points if the day in the purchase date is odd.
	// Only have to check and see if the last character is 1,3,5,7,9
	last_char := s_receipt.PurchaseDate[len(s_receipt.PurchaseDate)-1]
	if last_char == '1' || last_char == '3' || last_char == '5' || last_char == '7' || last_char == '9' {
		val = val + 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	first_two := s_receipt.PurchaseTime[0:2]
	hour, err := strconv.Atoi(first_two)
	if err != nil {
		log.Fatal(err)
	}
	if hour >= 14 && hour <= 16 {
		val = val + 10
	}

	total_cents, err := strconv.ParseFloat(s_receipt.Total, 64)
	if err != nil {
		log.Fatal(err)
	}
	// 50 points if the total is a round dollar amount with no cents.
	if math.Mod(total_cents, 1.00) == 0 {
		val = val + 50
	}
	// 25 points if the total is a multiple of 0.25.
	if math.Mod(total_cents, 0.25) == 0 {
		val = val + 25
	}

	loopCount := 0
	//There will be a variable number of Item objects in this part of the code, need a loop
	for i := 0; i < len(s_receipt.Items); i++ {
		// If the trimmed length of the item description is a multiple of 3,
		// multiply the price by 0.2 and round up to the nearest integer.
		// The result is the number of points earned.
		// Not sure what kinds of characters to trim?
		trimmed := strings.Trim(s_receipt.Items[i].Description, "!@#$%^&*()")
		if len(trimmed)%3 == 0 {
			item_total, err := strconv.ParseFloat(s_receipt.Items[i].Price, 64)
			if err != nil {
				log.Fatal(err)
			}
			pts_to_add := math.Round(item_total * 0.2)
			val = val + int(pts_to_add)
		}

		// 5 points for every two items on the receipt.
		loopCount++
		if loopCount%2 == 0 {
			val = val + 5
		}
	}

	return val
}

func enter_json(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/receipts/process" {
		http.Error(w, "404 PAGE NOT FOUND", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		path := r.FormValue("path")
		// Print out the id json object for the user to enter in later
		// path is the JSON name to be used for generateId(), be sure to include .json
		// http://localhost:8080/receipts/process
		fmt.Fprintf(w, "JSON name = %s\n", path)
		id_to_pass_on := generateId(path)
		json_data, err := json.Marshal(id_to_pass_on)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", string(json_data))

		//Declares the global variable to be used later on
		path_for_list_points := "examples/" + path
		id_for_list_points = receipt_id{id_to_pass_on, path_for_list_points}

	default:
		fmt.Fprintf(w, "Only get and post are allowed")
	}
}

func main() {

	http.HandleFunc("/receipts/process", enter_json)

	http.HandleFunc("/receipts/points/", func(w http.ResponseWriter, r *http.Request) {
		//the path will be /receipts/points/(id generated), get the last element from the URL
		//last element of the URL is the input for listPoints()
		//All that needs to be done is print out the points json
		//Be sure that JSON processing works with bigger examples first
		//Also be sure to make sure that docker works
		// http://localhost:8080/receipts/points/
		fmt.Fprintf(w, "Points for id: %s\n", id_for_list_points.Uuid.Id)
		points_display := points{listPoints(id_for_list_points)}
		json_data, err := json.Marshal(points_display)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", string(json_data))
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}

}
