package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var Timers map[string]*Timer

type TimerValue struct {
	HMS     string `json:"hms"`
	Seconds int    `json:"seconds"`
	Over    bool   `json:"over,omitempty"`
	Type    int    `json:"type"`
	Running bool   `json:"running"`
}

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
	//query := r.URL.Query()
	//log.Println(r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	out := json.NewEncoder(w)
	path := strings.SplitN(r.URL.Path[1:], "/", -1)
	//log.Println(path)

	if len(path) != 2 {
		log.Println("invalid path in timerValueHandler")
		http.Error(w, "invalid parameters; name not specified", http.StatusBadRequest)
		return
	}
	t, ok := Timers[path[1]]
	if !ok {
		log.Println("timer name not found")
		http.Error(w, "timer name not found", http.StatusNotFound)
		return
	}
	var tv TimerValue
	tv.HMS = t.HMS()
	tv.Seconds = t.Seconds()
	tv.Over = t.Over()
	tv.Type = t.Type()
	tv.Running = t.Running()
	err := out.Encode(tv)
	if err != nil {
		log.Fatalf("Unable to encode response: %s", err)
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
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

	staticServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticServer))
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
