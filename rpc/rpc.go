package rpc

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"tstorage"
)

type JSONRPC struct {
	opts    *Options
	storage tstorage.Storage
}

func NewJSONRPC(opts *Options) (*JSONRPC, error) {
	storage, err := tstorage.NewStorage(tstorage.WithDataPath(opts.dataPath),
		tstorage.WithPartitionDuration(opts.partitionDuration),
		tstorage.WithRetention(opts.retention),
		tstorage.WithTimestampPrecision(tstorage.TimestampPrecision(opts.timestampPrecision)),
		tstorage.WithWALBufferedSize(opts.walBufferSize),
	)
	if err != nil {
		return nil, err
	}
	log.Printf("Starting on %s", opts.address)
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		log.Printf("Shutting down gracefully.")
		storage.Close()
		os.Exit(0)
	}()
	return &JSONRPC{opts: opts, storage: storage}, nil
}

func (jrpc *JSONRPC) Main() error {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(jrpc, "ts")

	r := mux.NewRouter()
	r.Handle("/", s)
	return http.ListenAndServe(jrpc.opts.address, r)
}

type SelectArgs struct {
	Metric string      `json:"metric"`
	Labels []jsonLabel `json:"labels,omitempty"`
	Start  int64       `json:"start"`
	End    int64       `json:"end"`
}

type SelectRes []jsonPoint

func (args *SelectArgs) toArgs() (metric string, labels []tstorage.Label, start int64, end int64) {
	metric = args.Metric
	start = args.Start
	end = args.End
	labels = []tstorage.Label{}
	for _, l := range args.Labels {
		labels = append(labels, tstorage.Label{Name: l.Name, Value: l.Value})
	}
	return
}

func toSelectRes(points []*tstorage.DataPoint) *SelectRes {
	res := &SelectRes{}
	for _, p := range points {
		*res = append(*res, jsonPoint{Value: p.Value, Timestamp: p.Timestamp})
	}
	return res
}

func (b *JSONRPC) Select(r *http.Request, args *SelectArgs, reply *SelectRes) error {
	points, err := b.storage.Select(args.toArgs())
	if err != nil {
		return err
	}
	*reply = *toSelectRes(points)
	return nil
}

type jsonPoint struct {
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

type jsonLabel struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type InsertRowsArgs []struct {
	Metric    string      `json:"metric"`
	Labels    []jsonLabel `json:"labels,omitempty"`
	DataPoint jsonPoint   `json:"point"`
}

type InsertRowsRes struct {
	Result bool `json:"result"`
}

func (args *InsertRowsArgs) toRows() *[]tstorage.Row {
	rows := []tstorage.Row{}
	for _, r := range *args {
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

func (b *JSONRPC) InsertRows(r *http.Request, args *InsertRowsArgs, reply *InsertRowsRes) error {
	reply.Result = false
	err := b.storage.InsertRows(*args.toRows())
	if err != nil {
		return err
	}
	reply.Result = true
	return nil
}
