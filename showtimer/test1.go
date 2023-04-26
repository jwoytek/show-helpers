package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var Timers map[string]*Timer

func webHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("timers.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, Timers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func timerValueHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name, present := query["name"]
	if !present || len(name) == 0 {
		log.Println("timer name not present")
		http.Error(w, "invalid parameters; name not specified", http.StatusBadRequest)
		return
	}
	if present && len(name) != 1 {
		log.Println("timer name specified more than once")
		http.Error(w, "invalid parameters; name can only be specified once", http.StatusBadRequest)
		return
	}
	t, ok := Timers[name[0]]
	if !ok {
		log.Println("timer name not found")
		http.Error(w, "timer name not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, t.HMS())
}

func main() {
	//ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	Timers = make(map[string]*Timer)

	var err error
	duration, _ := time.ParseDuration("10s")
	Timers["act1"], err = NewTimer(TimerCountDown, "Act 1", duration)
	if err != nil {
		log.Fatalf("Error creating new timer: %s", err)
	}
	//intermission := Timers["intermission"]
	//act2 := Timers["act2"]

	log.Printf("act 1 val: %s", Timers["act1"].HMS())
	Timers["act1"].Start()

	/*
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
	*/

	http.HandleFunc("/timer/", timerValueHandler)
	http.HandleFunc("/", webHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	//time.Sleep(10 * time.Second)
	//act1.Reset()
	//act1.Start()
	//time.Sleep(10 * time.Second)
	//ticker.Stop()
	done <- true
}
