package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"within.website/ln"
)

var (
	readTimes = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "blogpage_read_times",
		Help: "This tracks how much time people spend reading articles on my blog",
	}, []string{"path"})
)

func init() {
	_ = prometheus.Register(readTimes)
}

func handlePageViewTimer(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("DNT") == "1" {
		http.NotFound(w, r)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ln.Error(r.Context(), err, ln.Info("while reading data"))
		http.Error(w, "oopsie whoopsie uwu", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	type metricsData struct {
		Path      string    `json:"path"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	}
	var md metricsData
	err = json.Unmarshal(data, &md)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	diff := md.EndTime.Sub(md.StartTime).Seconds()

	readTimes.WithLabelValues(md.Path).Observe(float64(diff))
}
