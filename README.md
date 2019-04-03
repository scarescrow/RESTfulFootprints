# Data Handling for the Building Footprints Dataset

### Installation and Configuration

* The code is written in [Golang](https://golang.org/), so go needs to be installed on the system first.
* [SQLite](https://www.sqlite.org/index.html) has been used as the database. So to get the server and ETL script running, it is required that the [Golang driver for SQLite](https://github.com/mattn/go-sqlite3) is installed. The instructions are present at [this link](https://github.com/mattn/go-sqlite3#installation).
* If the data stored in the database needs to be viewed, the [DB Browser for SQLite](https://sqlitebrowser.org/) may be used. This is an optional step and is not required for setting up of the server or the data transformation.

### Running the ETL Script

* Clone this repository, `git clone https://github.com/scarescrow/RESTfulFootprints.git`
* Naviate to this folder, `cd RESTfulFootprints`
* Create a folder "assets", `mkdir assets`
* Download the [Building Footprints](https://data.cityofnewyork.us/Housing-Development/Building-Footprints/nqwf-w8eh) dataset as a CSV file.
* Run the `etl.go` file by passing the location of the CSV file (downloaded in the previous step) as an argument. For example, 
`go run etl.go data\building.csv`
This will read the CSV file, and insert rows into the `buildings.db` file, which will be automatically created inside the `assets` folder. The script also shows progress in terms of the number of rows pushed.

### Running the Server

* Simply run the command `go run server.py`. This will start an HTTP server at port 8080 by default.
* There are currently 3 endpoints:
    - `/`
    - `/getBuildingsByYear?year=<year>`
    - `/getAverageHeightByYear?year=<year>`
* The first endpoint is simply a sanity check, while the other two query the DB to get the desired results, which are presented as a JSON response. Sample response: `{"average":"539320.923128"}` (average height of buildings for the year 1939)