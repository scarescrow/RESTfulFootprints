package main

import (
    "database/sql"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
    "fmt"
    "strconv"
    "encoding/json"
)

// This is the struct to store JSON response of the count endpoint

type CountResponse struct {
    Count string `json:"count"`
}

// This is the struct to store JSON response of the average height endpoint

type AvgResponse struct {
    Avg string `json:"average"`
}

// This function is just a sanity check to see if the server is running
// as expected

func sayHello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("I am alive"))
}

// This is the function which handles the average height by year endpoint

func getAverageHeightByYear(w http.ResponseWriter, r *http.Request) {

    var avg float64

    // Get the year parameter from the URL
    keys, ok := r.URL.Query()["year"]
    message := ""

    // Read the database. For scalability purposes, it is better to open
    // the database at every endpoint, because SQLite is a small SQL application
    // which is not designed for multiple concurrent reads and writes. While
    // it is theoretically possible, it is better for the app to be a bit slow,
    // but robust.

    database, err := sql.Open("sqlite3", "./assets/buildings.db")
    if err != nil || database == nil {
        fmt.Println(err)
    }
    
    if !ok || len(keys[0]) < 1 {
        
        // If the year is not specified, prompt the client to provide one

        message = "Please specify a year"

    } else {
        
        // If the database has been loaded successfully, and year is present,
        // we can query the DB to get the desired result

        year := keys[0]
        err := database.QueryRow("SELECT AVG(height) AS h FROM buildings WHERE construct_year = ?", year).Scan(&avg)
        defer database.Close()
        if err != nil {
        
            message = err.Error()
        
        } else {

            // Convert the float to a string to display as a response

            message = fmt.Sprintf("%f", avg)

        }

    }

    // We can now create a JSON object out of the result using the structs
    // we created earlier. JSON is being used because it is a highly popular 
    // and well supported means of data transfer among REST APIs

    response := &AvgResponse{Avg: message}
    js, err := json.Marshal(response)
    if err != nil {
        fmt.Println(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(js)
}

// This is the function which handles the average height by year endpoint

func getBuildingsByYear(w http.ResponseWriter, r *http.Request) {

    // This function is quite similar to the previous one, with some
    // minor tweaks

    var count int
    keys, ok := r.URL.Query()["year"]
    message := ""

    database, err := sql.Open("sqlite3", "./assets/buildings.db")
    if err != nil || database == nil {
        fmt.Println(err)
    }
    
    if !ok || len(keys[0]) < 1 {
        
        message = "Please specify a year"

    } else {
        
        year := keys[0]
        err := database.QueryRow("SELECT COUNT(*) FROM buildings WHERE construct_year = ?", year).Scan(&count)
        defer database.Close()
        if err != nil {
        
            message = err.Error()
        
        } else {

            // Convert int to string in this case

            message = strconv.Itoa(count)

        }

    }
    response := &CountResponse{Count: message}
    js, err := json.Marshal(response)
    if err != nil {
        fmt.Println(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(js)
}

func main() {

    // Define the routes for the different endpoints

    http.HandleFunc("/", sayHello)
    http.HandleFunc("/getBuildingsByYear", getBuildingsByYear)
    http.HandleFunc("/getAverageHeightByYear", getAverageHeightByYear)

    // Start the server at port 8080

    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    } else {
        fmt.Println("Server started")
    }
}