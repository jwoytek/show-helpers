package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

const (
	TimerCountUp   = iota
	TimerCountDown = iota
)

type Timer struct {
	Name            string
	timerType       int
	totalSecs       float64
	initialDuration time.Duration
	set             bool
	running         bool
	timerStop       chan bool
}

func NewTimer(timerType int, name string, initialDuration time.Duration) (t *Timer, err error) {
	if timerType < TimerCountUp || timerType > TimerCountDown {
		err = errors.New("Timer type not one of 'TimerCountUp' or 'TimerCountDown'")
	}
	log.Printf("Creating new timer '%s' with initial duration of %fs", name, initialDuration.Seconds())
	t = new(Timer)
	t.timerType = timerType
	t.Name = name
	t.initialDuration = initialDuration
	t.totalSecs = initialDuration.Seconds()

	return t, nil
}

func (t *Timer) update(duration time.Duration) {
	t.set = true
	t.totalSecs = duration.Seconds()
	//e.hours = int(e.totalSecs/(60*60)) % 24
	//e.minutes = int(e.totalSecs/60) % 60
	//e.seconds = int(e.totalSecs) % 60
}

func (t *Timer) Start() {
	if t.timerStop == nil {
		t.timerStop = make(chan bool)
	}

	go func() {
		log.Println("timer started")
		t.running = true
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		start := time.Now()
		end := start.Add(time.Duration(t.initialDuration))
		for {
			select {
			case <-t.timerStop:
				log.Println("timer told to stop")
				return
			case tick := <-ticker.C:
				switch t.timerType {
				case TimerCountUp:
					t.update(tick.Sub(start))
					//log.Println("Elapsed:", t.HMS())
				case TimerCountDown:
					t.update(end.Sub(tick))
					//log.Println("Remaining:", t.HMS())
				}
			}
		}
	}()
}

func (t *Timer) Stop() {
	if t.running {
		log.Println("stopping timer")
		t.timerStop <- true
		t.running = false
	}
}

func (t *Timer) Reset() {
	t.Stop()
	t.set = false
	t.totalSecs = t.initialDuration.Seconds()
}

func (t Timer) HMS() string {
	if !t.set {
		return fmt.Sprint("--:--:--")
	}
	log.Printf("totalSecs = %f", t.totalSecs)
	secs := t.totalSecs
	prefix := ""
	if t.totalSecs < 0 {
		secs = math.Abs(t.totalSecs)
		prefix = "- "
	}
	hours := int(secs/(60*60)) % 24
	minutes := int(secs/60) % 60
	seconds := int(secs) % 60
	return fmt.Sprintf("%s%02d:%02d:%02d", prefix, hours, minutes, seconds)
}
