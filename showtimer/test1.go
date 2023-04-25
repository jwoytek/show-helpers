package main

import (
	//"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
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
	done := make(chan bool)

	Timers = make(map[string]*Timer)

	var err error
	duration, _ := time.ParseDuration("10s")
	Timers["act1"], err = NewTimer(TimerCountUp, "Act 1", duration)
	if err != nil {
		log.Fatalf("Error creating new timer: %s", err)
	}
	//intermission := Timers["intermission"]
	//act2 := Timers["act2"]

	Timers["act1"].Start()

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				break
			case <-ticker.C:
				log.Println("Timer:", Timers["act1"].HMS())
			}
		}
	}()

	http.HandleFunc("/", webHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	//time.Sleep(10 * time.Second)
	//act1.Reset()
	//act1.Start()
	//time.Sleep(10 * time.Second)
	//ticker.Stop()
	done <- true
}
