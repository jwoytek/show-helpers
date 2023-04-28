package main

import (
	"encoding/json"
	"errors"
	"flag"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hypebeast/go-osc/osc"
	"github.com/spf13/viper"
)

var Timers map[string]*Timer

type TimerValue struct {
	HMS          string `json:"hms"`
	HMSIndicator string `json:"hms_indicator"`
	Seconds      int    `json:"seconds"`
	Over         bool   `json:"over,omitempty"`
	Type         int    `json:"type"`
	Running      bool   `json:"running"`
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
	tv.HMSIndicator = t.HMSIndicator()
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

func oscHandleTimerStart(msg *osc.Message) {
	if msg.CountArguments() != 1 {
		log.Printf("Bad OSC /timer/start message: %s", msg)
		return
	}
	name := msg.Arguments[0].(string)
	t, ok := Timers[name]
	if !ok {
		log.Printf("Asked to start unknown timer: %s", name)
	}
	log.Printf("OSC timer start for %s", name)
	t.Start()
}

func oscHandleTimerStop(msg *osc.Message) {
	if msg.CountArguments() != 1 {
		log.Printf("Bad OSC /timer/stop message: %s", msg)
		return
	}
	name := msg.Arguments[0].(string)
	t, ok := Timers[name]
	if !ok {
		log.Printf("Asked to stop unknown timer: %s", name)
	}
	log.Printf("OSC timer stop for %s", name)
	t.Stop()
}

func oscHandleTimerReset(msg *osc.Message) {
	if msg.CountArguments() != 1 {
		log.Printf("Bad OSC /timer/reset message: %s", msg)
		return
	}
	name := msg.Arguments[0].(string)
	t, ok := Timers[name]
	if !ok {
		log.Printf("Asked to reset unknown timer: %s", name)
	}
	log.Printf("OSC timer reset for %s", name)
	t.Reset()
}

func main() {
	//ticker := time.NewTicker(1 * time.Second)
	var configFile = flag.String("config", "showtimer.yaml", "name of configuration file to read")
	flag.Parse()

	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	timerList := viper.Get("timers")
	log.Printf("Timers: %v", timerList)

	for _, timerDef := range timerList.([]interface{}) {
		//log.Printf("timerList[%v] = %v", timerDef, timerDef) //timerList[timerDef])
		log.Printf("timer id: %s\ttimer name: %s", timerDef[0])
	}
	panic("die")

	ips, err := findMyIPs()
	if err != nil {
		log.Fatalf("Error getting IP address: %s", err)
	}
	for _, ip := range ips {
		log.Printf("Listening for web traffic on %s:%d", ip, 8080)
	}

	Timers = make(map[string]*Timer)

	//var err error
	//duration, _ := time.ParseDuration("10s")
	Timers["act1"], err = NewTimer(TimerCountUp, "Act 1", time.Duration(0))
	if err != nil {
		log.Fatalf("Error creating new timer: %s", err)
	}
	// log.Printf("act 1 val: %s", Timers["act1"].HMS())
	//Timers["act1"].Start()

	duration, _ := time.ParseDuration("1m")
	Timers["intermission"], err = NewTimer(TimerCountDown, "Intermission", duration)
	if err != nil {
		log.Fatalf("Error creating new timer: %s", err)
	}
	// log.Printf("act 1 val: %s", Timers["act1"].HMS())
	Timers["intermission"].Start()

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

	oscDispatcher := osc.NewStandardDispatcher()
	oscDispatcher.AddMsgHandler("/timer/start", oscHandleTimerStart)
	oscDispatcher.AddMsgHandler("/timer/stop", oscHandleTimerStop)
	oscDispatcher.AddMsgHandler("/timer/reset", oscHandleTimerReset)
	oscServer := &osc.Server{
		Addr:       "0.0.0.0:8000",
		Dispatcher: oscDispatcher,
	}

	go func() {
		oscServer.ListenAndServe()
	}()

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
}

func findMyIPs() ([]string, error) {
	ifaces, err := net.Interfaces()
	ips := make([]string, 0, 10)
	if err != nil {
		return ips, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		//if iface.Flags&net.FlagLoopback != 0 {
		//	continue // loopback interface
		//}
		addrs, err := iface.Addrs()
		if err != nil {
			return ips, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil { //|| ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ips = append(ips, ip.String())
			//return ip.String(), nil
		}
	}
	if len(ips) == 0 {
		return ips, errors.New("are you connected to the network?")
	}
	return ips, nil
}
