package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// https://stackoverflow.com/questions/36719525/how-to-log-messages-to-the-console-and-a-file-both-in-golang
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	//Load env
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		fmt.Println("Error loading .env file")
	}
	fmt.Println("Starting Server")
	// mux handler
	r := mux.NewRouter()
	r.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(200)
		fmt.Fprintf(writer, "Server is up and running")
		log.Println("Server is up and running")
	})
	r.HandleFunc("/registerWebhook", registerWebhook)
	// displays the json data to the user in the browser in json format
	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./data.json")
	})
	// Opens the log.txt file in browser
	r.HandleFunc("/systemlogs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./log.txt")
	})
	//r.PathPrefix("/data/").Handler(http.StripPrefix("/data/", http.FileServer(http.Dir("data.json"))))
	// handles the callback
	//r.HandleFunc("/callback", callbackhandler)
	// checks for crc everywhour
	r.HandleFunc("/webhook/twitter", CrcCheck).Methods("GET")
	//Listen to webhook event and  handle
	r.HandleFunc("/webhook/twitter", twitterwebhookdump).Methods("POST")
	//Start Server
	server := &http.Server{
		Handler: r,
	}
	server.Addr = os.Getenv("SERVER_PORT")
	server.ListenAndServe()
}
