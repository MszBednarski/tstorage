package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"
	"tstorage"
)

const (
	DAY       = 24 * time.Hour
	TWO_YEARS = 24 * 365 * 2 * time.Hour
)

type tstoragePoint struct {
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

type insertRowsRequest struct {
	Rows []struct {
		Metric string `json:"metric"`
		Labels []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"labels,omitempty"`
		DataPoint tstoragePoint `json:"point"`
	} `json:"rows"`
}

type selectRowsRequest struct {
	Metric string `json:"metric"`
	Labels []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"labels,omitempty"`
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type tstorageSelect struct {
	Metric string
	Labels []tstorage.Label
	Start  int64
	End    int64
}

type tstorageResult struct {
	Points []tstoragePoint `json:"points"`
}

func (req *insertRowsRequest) toTstorageRows() *[]tstorage.Row {
	rows := []tstorage.Row{}
	for _, r := range req.Rows {
		labels := []tstorage.Label{}
		for _, l := range r.Labels {
			labels = append(labels, tstorage.Label{Name: l.Name, Value: l.Value})
		}
		dp := tstorage.DataPoint{Value: r.DataPoint.Value, Timestamp: r.DataPoint.Timestamp}
		trow := tstorage.Row{Metric: r.Metric, Labels: labels, DataPoint: dp}
		rows = append(rows, trow)
	}
	return &rows
}

func (req *selectRowsRequest) toTstorageSelect() *tstorageSelect {
	labels := []tstorage.Label{}
	for _, l := range req.Labels {
		labels = append(labels, tstorage.Label{Name: l.Name, Value: l.Value})
	}
	return &tstorageSelect{Metric: req.Metric, Labels: labels, Start: req.Start, End: req.End}
}

func getPutHandler(s tstorage.Storage) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		requestBody := &insertRowsRequest{}
		err = json.Unmarshal(body, requestBody)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		err = s.InsertRows(*requestBody.toTstorageRows())

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(200)
	}
	return http.HandlerFunc(fn)
}

func tstoragePointsToJson(tpoints []*tstorage.DataPoint) ([]byte, error) {
	points := []tstoragePoint{}
	for _, p := range tpoints {
		points = append(points, tstoragePoint{Value: p.Value, Timestamp: p.Timestamp})
	}
	result := tstorageResult{Points: points}
	return json.Marshal(result)
}

func getGetHandler(s tstorage.Storage) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		requestBody := &selectRowsRequest{}
		err = json.Unmarshal(body, requestBody)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		r := requestBody.toTstorageSelect()
		points, err := s.Select(r.Metric, r.Labels, r.Start, r.End)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		result, err := tstoragePointsToJson(points)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(result)
	}
	return http.HandlerFunc(fn)
}

func main() {
	fmt.Printf("Starting the tstore.")
	storage, err := tstorage.NewStorage(
		tstorage.WithDataPath("./data"),
		tstorage.WithTimestampPrecision("s"),
		tstorage.WithRetention(TWO_YEARS),
		tstorage.WithPartitionDuration(DAY),
	)
	if err != nil {
		panic(err)
	}
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		fmt.Printf("Shutting down.")
		storage.Close()
		os.Exit(0)
	}()
	http.Handle("/put", getPutHandler(storage))
	http.Handle("/get", getGetHandler(storage))
	http.ListenAndServe(fmt.Sprintf(":%d", 3000), nil)
}
