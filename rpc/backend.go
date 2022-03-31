package main

import (
	"tstorage"
)

type Backend struct {
	storage tstorage.Storage
}

type SelectArgs struct {
	Metric string
	Labels []struct {
		Name  string
		Value string
	}
	Start int64
	End   int64
}

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

func (b *Backend) Select(args *SelectArgs, reply *[]*tstorage.DataPoint) error {
	points, err := b.storage.Select(args.toArgs())
	if err != nil {
		return err
	}
	*reply = points
	return nil
}

type InsertRowsArgs []struct {
	Metric string
	Labels []struct {
		Name  string
		Value string
	}
	DataPoint struct {
		Value     float64
		Timestamp int64
	}
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

func (b *Backend) InsertRows(args *InsertRowsArgs, reply *bool) error {
	*reply = false
	err := b.storage.InsertRows(*args.toRows())
	if err != nil {
		return err
	}
	*reply = true
	return nil
}
