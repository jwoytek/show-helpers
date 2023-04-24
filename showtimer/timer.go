package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	CountUp   = iota
	CountDown = iota
)

func MakeTimer(timerType int) (t *Timer, err error) {
	if timerType < CountUp || timerType > CountDown {
		err = errors.New("Timer type not one of 'CountUp' or 'CountDown'")
	}
	return &Timer{timerType: timerType}, nil
}

type Timer struct {
	timerType int
	totalSecs float64
	set       bool
	running   bool
	timerStop chan bool
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
		for {
			select {
			case <-t.timerStop:
				log.Println("timer told to stop")
				return
			case tick := <-ticker.C:
				switch t.timerType {
				case CountUp:
					t.update(tick.Sub(start))
					log.Println("Elapsed:", t.HMS())
				case CountDown:
					t.update(start.Sub(tick))
					log.Println("Remaining:", t.HMS())
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
	t.totalSecs = 0
}

func (t Timer) HMS() string {
	if !t.set {
		return fmt.Sprint("--:--:--")
	}
	hours := int(t.totalSecs/(60*60)) % 24
	minutes := int(t.totalSecs/60) % 60
	seconds := int(t.totalSecs) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
