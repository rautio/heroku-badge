package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)


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


func main() {
	// Connect to DB
  db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
  if err != nil {
    log.Fatal(err)
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

	/**
	 	* Return a badge image.
		* 
		* Query Params: 
		* 	app_name=<string> : name of the app deployed in heroku
		*/
	getBadgeHandler := func(w http.ResponseWriter, req *http.Request) {
		db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		req.ParseForm()
		appName, hasAppName := req.Form["app_name"]
		log.Println("Get Badge Request!")
		log.Println(req.Form)
		log.Println(hasAppName)
		log.Println(appName)
		var status string
		var lastUpdate string
		var appId string
		err := db.QueryRow(`SELECT app_id, status, last_update FROM status WHERE app_name=$1;`, appName).Scan(&appId, &status, &lastUpdate)
		defer db.Close()
		if err != nil {
			log.Println(err)
			// If there was no match above then it is an unknown word
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "None Found", http.StatusBadRequest)
			return
		}
		result := map[string]interface{}{ "status": status, "app_id": appId, "last_update": lastUpdate }
		jsonResponse, jsonError := json.Marshal(result)
		if jsonError != nil {
			log.Println(jsonError)
		  fmt.Println("Unable to encode JSON")
		}
    w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
		return
	}

	/**
	 	* Return a badge image.
		* 
		* Query Params: 
		* 	app_name=<string> : name of the app deployed in heroku
		*/
		buildUpdateHandler := func(w http.ResponseWriter, req *http.Request) {
			db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
			log.Println("Received Build Update")
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
			// Update status info
			// _, err = db.Exec(`
			// UPDATE status SET status=$2, last_update=$3 WHERE app_id=$1 AND last_update<=$3;`, data.App.Id, data.Status, data.CreatedAt)
			// _, err = db.Exec(`
			// UPDATE status SET status=$3, last_update=$4 WHERE app_id=$1 AND last_update<=$4;
			// INSERT INTO status (app_id, app_name, status, last_update)
      //  	VALUES ($1, $2, $3, $4)
      //  	WHERE (NOT) EXISTS (SELECT 1 FROM status WHERE app_id=$1);`, data.App.Id, data.App.Name, data.Status, data.CreatedAt )
			_, err = db.Exec(`
			UPDATE status SET status=$3, last_update=$4 WHERE app_id=$1 AND last_update<=$4;
			INSERT INTO status (app_id, app_name, status, last_update)
       	VALUES ($1, $2, $3, $4)
       	ON CONFLICT DO NOTHING);`, data.App.Id, data.App.Name, data.Status, data.CreatedAt )
			// _, err := db.Exec(`
			// INSERT INTO status (app_id, app_name, status, last_update)
			// VALUES ($1, $2, $3, $4)`, data.App.Id, data.App.Name, data.Status, data.CreatedAt )
			defer db.Close()
			if err != nil {
				log.Println(err)
			}
			w.Write([]byte("Success"))
			// Status table
			return
		}
	
	// Router setup
	router := mux.NewRouter().StrictSlash(true)

	port := getPort()

	router.HandleFunc("/status", getBadgeHandler).Methods("GET","OPTIONS")
	log.Println(fmt.Sprintf("Listening for requests at GET http://localhost%s/", port))


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