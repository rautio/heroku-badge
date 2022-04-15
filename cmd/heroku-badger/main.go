package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)


type BuildUpdate struct {
	CreatedAt   string `json:"created_at"`
	data         struct {
		CreatedAt    string `json:"created_at"`
		app           struct {
			Status        string `json:"status"`
		}
	}
}

func main() {

	/**
	 	* Return a badge image.
		* 
		* Query Params: 
		* 	app_name=<string> : name of the app deployed in heroku
		*/
	getBadgeHandler := func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		_, hasAppName := req.Form["app_name"]
		appName, _ := strconv.Atoi(req.URL.Query().Get("app_name"))
		log.Println("Get Badge Request!")
		log.Println(req.Form)
		log.Println(hasAppName)
		log.Println(appName)
		w.Write([]byte("Hello World"))
		return
	}

	/**
	 	* Return a badge image.
		* 
		* Query Params: 
		* 	app_name=<string> : name of the app deployed in heroku
		*/
		buildUpdateHandler := func(w http.ResponseWriter, req *http.Request) {
			log.Println("Build Update!")
			log.Println("=====START=====")
			log.Println(req)
			log.Println("=====BODY=====")
			log.Println(req.Body)
			var postBody BuildUpdate
			decoder := json.NewDecoder(req.Body)
			decodePostErr := decoder.Decode(&postBody)
			if decodePostErr != nil {
				log.Println(decodePostErr)
				panic(decodePostErr)
			}
			log.Println(postBody)
			data := postBody.data
			log.Println(data)
			log.Println("=====DATA=====")
			log.Println(data)
			log.Println("=====CREATED=====")
			log.Println(postBody.CreatedAt)
			log.Println(data.CreatedAt)
			log.Println("=====APP=====")
			log.Println(data.app)
			log.Println("=====STATUS=====")
			log.Println(data.app.Status)
			log.Println("=====END=====")
			w.Write([]byte("Success"))
			return
		}
	
	// Router setup
	router := mux.NewRouter().StrictSlash(true)

	port := getPort()

	router.HandleFunc("/", getBadgeHandler).Methods("GET","OPTIONS")
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