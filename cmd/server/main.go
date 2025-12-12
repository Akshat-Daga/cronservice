package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spillbox/cronservice/pkg/crontab"
)

var scheduler = crontab.NewScheduler()

func main() {
	http.HandleFunc("/start", startScheduler)
	http.HandleFunc("/stop", stopScheduler)
	http.HandleFunc("/jobs", listJobs)
	http.HandleFunc("/schedule/cron", scheduleCronJob)
	http.HandleFunc("/schedule/every", scheduleIntervalJob)

	fmt.Println("API server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startScheduler(w http.ResponseWriter, r *http.Request) {
	scheduler.StartAsync()
	fmt.Fprintln(w, "Scheduler started")
}

func stopScheduler(w http.ResponseWriter, r *http.Request) {
	scheduler.Stop()
	fmt.Fprintln(w, "Scheduler stopped")
}

func listJobs(w http.ResponseWriter, r *http.Request) {
	jobs := scheduler.Jobs()
	json.NewEncoder(w).Encode(jobs)
}

func scheduleCronJob(w http.ResponseWriter, r *http.Request) {
	cronExpr := r.URL.Query().Get("cron")
	if cronExpr == "" {
		http.Error(w, "missing cron expression", 400)
		return
	}

	err := scheduler.Cron(cronExpr, func() {
		fmt.Println("Running cron job:", cronExpr)
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintln(w, "Cron job scheduled:", cronExpr)
}

func scheduleIntervalJob(w http.ResponseWriter, r *http.Request) {
	interval := r.URL.Query().Get("interval")
	unit := r.URL.Query().Get("unit")

	if interval == "" || unit == "" {
		http.Error(w, "missing interval or unit", 400)
		return
	}

	var val uint64
	fmt.Sscanf(interval, "%d", &val)

	err := scheduler.EveryInterval(val, unit, func() {
		fmt.Println("Running interval job:", interval, unit)
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintln(w, "Interval job scheduled:", interval, unit)
}
