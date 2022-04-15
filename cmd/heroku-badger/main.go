package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Word struct {
	Word             string `json:"word"`
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
			log.Println("=====START=====")
			log.Println("Build Update!")
			log.Println(req)
			var postBody map[string]interface{}
			decoder := json.NewDecoder(req.Body)
			decodePostErr := decoder.Decode(&postBody)
			if decodePostErr != nil {
				log.Println(decodePostErr)
				panic(decodePostErr)
			}
			log.Println(postBody)
  		reqBody, _ := ioutil.ReadAll(req.Body)
			var marshalledData map[string]interface{}
		  json.Unmarshal(reqBody, &marshalledData)
			data := marshalledData
			log.Println(data)
			log.Println("=====DATA=====")
			log.Println(postBody["data"])
			log.Println("=====CREATED=====")
			log.Println(postBody["created_at"])
			log.Println(data["created_at"])
			log.Println("=====APP=====")
			log.Println(data["app"])
			log.Println("=====STATUS=====")
			log.Println(data["status"])
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