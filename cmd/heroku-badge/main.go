package main

import (
	"fmt"
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
		// Choose a word at random from the most frequent sub-list
		w.WriteHeader([]byte("Hello World"))
		return
	}

	
	// Router setup
	router := mux.NewRouter().StrictSlash(true)

	port := getPort()

	router.HandleFunc("/", getBadgeHandler).Methods("GET","OPTIONS")
	log.Println(fmt.Sprintf("Listening for requests at GET http://localhost%s/", port))


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