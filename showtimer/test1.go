package main

import (
	//"fmt"
	"html/template"
	"log"
	"net/http"
)

func webHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("timers.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, "00:00:00")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var Timers map[string]*Timer

func main() {
	//ticker := time.NewTicker(1 * time.Second)
	//done := make(chan bool)

	Timers = make(map[string]*Timer)

	Timers["act1"] = &Timer{}
	intermission := Timers["intermission"]
	act2 := Timers["act2"]

	Timers["act1"].Start()

	act1.Start()

	http.HandleFunc("/", webHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	//time.Sleep(10 * time.Second)
	//act1.Reset()
	//act1.Start()
	//time.Sleep(10 * time.Second)
	//ticker.Stop()
	//done <- true
}
