package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)


type AppStatus struct {
	Id             string `json:"app_id"`
	Name           string `json:"app_name"`
	Status         string `json:"status"`
	UpdatedAt      string `json:"updated_at"`
}

type BuildUpdate struct {
	CreatedAt   string `json:"created_at"`
	Action      string `json:"action"`    
	Data         struct {
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		PublishedAt    string `json:"published_at"`
		Status       string `json:"status"`
		App           struct {
			Id             string `json:"id"`
			Name           string `json:"name"`
		}
	}
}

func getAppStatus(appName string) (AppStatus, error) {
	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	var status string
	var lastUpdate string
	var appId string
	dbErr := db.QueryRow(`SELECT app_id, status, last_update FROM status WHERE app_name=$1;`, appName).Scan(&appId, &status, &lastUpdate)
	defer db.Close()
	if dbErr != nil {
		return AppStatus{}, dbErr
	}
	return AppStatus{ appId, appName, status, lastUpdate }, nil
}


func setupDb() {
  db, dbConnectErr := sql.Open("postgres", os.Getenv("DATABASE_URL"))
  if dbConnectErr != nil {
    log.Fatal(dbConnectErr)
  }
	// Status table
	// Only tracking the last status to minimize data storage
	db.Exec(`CREATE TABLE IF NOT EXISTS status (
		app_name VARCHAR (50) NOT NULL,
		app_id VARCHAR (50) NOT NULL UNIQUE,
		status VARCHAR (20) NOT NULL,
		last_update TIMESTAMP WITHOUT TIME ZONE NOT NULL
	)`)
	defer db.Close()
}

func main() {
	setupDb();

	/**
	 	* Return status info for the app. Requires a specific app name.
		* 
		* Query Params: 
		* 	app_name=<string> : name of the app deployed in heroku
		*/
	getStatusHandler := func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		appName := req.FormValue("app_name")
		status, dbErr := getAppStatus(appName)
		if dbErr != nil {
			log.Println(dbErr)
			// If there was no match above then it is an unknown word
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "None Found", http.StatusBadRequest)
			return
		}
		jsonResponse, jsonError := json.Marshal(status)
		if jsonError != nil {
			log.Println(jsonError)
		  fmt.Println("Unable to encode JSON")
		}
    w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
		return
	}

	getBadgeHandler := func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		appName := req.FormValue("app_name")
		appStatus, dbErr := getAppStatus(appName)
		status := appStatus.Status
		color := "inactive"
		if (status == "succeeded") {
			color = "success"
		}
		if (status == "pending") {
			color = "yellow"
		}
		if (status == "failed") {
			color = "critical"
		}
		log.Println(appStatus)
		if dbErr != nil {
			status = "unknown"
			log.Println(dbErr)
			// If there was no match above then it is an unknown word
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "None Found", http.StatusBadRequest)
			return
		}
		log.Println(status)
		log.Println(color)
		badgeRes, _ := http.Get(fmt.Sprintf("https://img.shields.io/badge/Build-%s-%s", status, color))
		// badgeRes, _ := http.Get("https://img.shields.io/badge/test-foo-red")
		badge, _ := ioutil.ReadAll(badgeRes.Body)
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(badge)
		return
	}

	/**
	 	* Webhook Listener for heroku build updates.
		* 
		*/
		buildUpdateHandler := func(w http.ResponseWriter, req *http.Request) {
			log.Println("Received Build Update")
			db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
			var postBody BuildUpdate
			decoder := json.NewDecoder(req.Body)
			decodePostErr := decoder.Decode(&postBody)
			if decodePostErr != nil {
				log.Println(decodePostErr)
				panic(decodePostErr)
			}
			data := postBody.Data
			log.Println(data.CreatedAt)
			log.Println(data.App.Id)
			log.Println(data.App.Name)
			log.Println(data.Status)
			// Separting Insert and Update. Both are optional and only 1 should execute
			// depending on if its a brand new app or an update to an existing one
			_, errInsert := db.Exec(`
			INSERT INTO status (app_id, app_name, status, last_update)
       	VALUES ($1, $2, $3, $4)
       	ON CONFLICT DO NOTHING;`, data.App.Id, data.App.Name, data.Status, data.CreatedAt )
			_, errUpdate := db.Exec(`
			UPDATE status SET status=$2, last_update=$3 WHERE app_id=$1 AND last_update<=$3;
			`, data.App.Id, data.Status, data.CreatedAt )
			defer db.Close()
			if errInsert != nil {
				log.Println(errInsert)
			}
			if errUpdate != nil {
				log.Println(errUpdate)
			}
			w.Write([]byte("Success"))
			return
		}
	
	// Router setup
	router := mux.NewRouter().StrictSlash(true)

	port := getPort()

	router.HandleFunc("/", getBadgeHandler).Methods("GET","OPTIONS")
	log.Println(fmt.Sprintf("Listening for requests at GET http://localhost%s/", port))

	router.HandleFunc("/status", getStatusHandler).Methods("GET","OPTIONS")
	log.Println(fmt.Sprintf("Listening for requests at GET http://localhost%s/status", port))

	router.HandleFunc("/build-update", buildUpdateHandler).Methods("POST","OPTIONS")
	log.Println(fmt.Sprintf("Listening for requests at POST http://localhost%s/build-update", port))


	// TODO: Return a API doc page w/ examples like type ahead
	http.Handle("/", router)
	http.ListenAndServe(port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization" }), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),handlers.AllowedOrigins([]string{"*"}))(router))
}

func getPort() string {
  p := os.Getenv("PORT")
  if p != "" {
    return ":" + p
  }
  return ":9000"
}


func CheckError(err error) {
	if err != nil {
			panic(err)
	}
}