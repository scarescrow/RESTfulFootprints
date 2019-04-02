package main

import (
    "os"
    "fmt"
    "encoding/csv"
    "io"
    "bufio"
    "log"
    "database/sql"
    "time"
    "strconv"
    _ "github.com/mattn/go-sqlite3"
)

type Building struct {
    shape, status, source string
    height, length, area float64
    bin, doitt_id, feat_code, construct_year int64
    ground_elev, base_bbl, mpluto_bbl, date int64
}

func getBuilding(line []string) *Building {
    
    time_layout := "01/02/2006 15:04:05 PM -0700"
    
    building := new(Building)
    parsedTime, _ := time.Parse(time_layout, line[4])

    building.shape = line[0]
    building.bin, _ = strconv.ParseInt(line[1], 10, 64)
    building.construct_year, _ = strconv.ParseInt(line[2], 10, 64)
    building.date = parsedTime.Unix()
    building.status = line[5]
    building.doitt_id, _ = strconv.ParseInt(line[6], 10, 64)
    building.height, _ = strconv.ParseFloat(line[7], 64)
    building.feat_code, _ = strconv.ParseInt(line[8], 10, 64)
    building.ground_elev, _ = strconv.ParseInt(line[9], 10, 64)
    building.area, _ = strconv.ParseFloat(line[10], 64)
    building.length, _ = strconv.ParseFloat(line[11], 64)
    building.base_bbl, _ = strconv.ParseInt(line[12], 10, 64)
    building.mpluto_bbl, _ = strconv.ParseInt(line[13], 10, 64)
    building.source = line[14]

    return building
}

func main() {
    args := os.Args[1:]
    header := false
    if len(args) < 1 {
        fmt.Println("Please provide a CSV file to perform ETL on")
    } else {
        csvFile, err := os.Open(args[0])
        if err != nil {
            fmt.Println(err)
        }

        database, err := sql.Open("sqlite3", "./assets/buildings.db")
        if err != nil {
            fmt.Println(err)
        }

        statement, err := database.Prepare("DROP TABLE IF EXISTS buildings")
        if err != nil {
            fmt.Println(err)
        }
        statement.Exec()

        statement, err = database.Prepare("CREATE TABLE buildings (id INTEGER PRIMARY KEY AUTOINCREMENT, shape TEXT, bin INTEGER, construct_year INTEGER, date INTEGER, status TEXT, doitt_id INTEGER, height REAL, feat_code INTEGER, ground_elev INTEGER, area REAL, length REAL, base_bbl INTEGER, mpluto_bbl INTEGER, source TEXT)")
        if err != nil {
            fmt.Println(err)
        }

        _, err = statement.Exec()
        if err != nil {
            fmt.Println(err)
        }

        reader := csv.NewReader(bufio.NewReader(csvFile))

        tx, _ := database.Begin()

        i := 0
        for {
            line, error := reader.Read()

            if !header {
                header = true
                continue
            }

            i += 1
            
            var building = getBuilding(line)
        
            if error == io.EOF {
                break
            } else if error!= nil {
                log.Fatal(error)
            }

            statement, err = database.Prepare("INSERT INTO buildings VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
            if err != nil {
                fmt.Println(err)
            }

            _, err = statement.Exec(i, 
                building.shape, 
                building.bin, 
                building.construct_year,
                building.date,
                building.status,
                building.height,
                building.doitt_id,
                building.feat_code,
                building.ground_elev,
                building.area,
                building.length,
                building.base_bbl,
                building.mpluto_bbl,
                building.source)

            if err != nil {
                fmt.Println(err)
                break
            }

            if i % 20000 == 0 {
                fmt.Println(i, "rows have been pushed to the DB")
                tx.Commit()
            }
        }

        tx.Commit()
    }
}