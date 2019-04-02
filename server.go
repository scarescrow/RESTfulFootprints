package main

import (
    "database/sql"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
    "fmt"
    "strconv"
    "encoding/json"
    // "reflect"
    // "strings"
)

type CountResponse struct {
    Count string `json:"count"`
}

func sayHello(w http.ResponseWriter, r *http.Request) {
    
}

func getBuildingsByYear(w http.ResponseWriter, r *http.Request) {

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
    http.HandleFunc("/", sayHello)
    http.HandleFunc("/getBuildingsByYear", getBuildingsByYear)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    } else {
        fmt.Println("Server started")
    }
}