package main

import (
    "database/sql"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
    // "strings"
)

database, err := sql.Open("sqlite3", "./buildings.db")
    if err != nil {
        fmt.Println(err)
    }

func sayHello(w http.ResponseWriter, r *http.Request) {
    keys, ok := r.URL.Query()["key"]
    message := ""
    if !ok || len(keys[0]) < 1 {
        message = "Key not found"
    } else {
        message = keys[0]
    }
    w.Write([]byte (message))
}

func getBuildingsByYear(w http.ResponseWriter, r *http.Request) {

}

func main() {
    http.HandleFunc("/", sayHello)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}