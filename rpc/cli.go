package rpc

import (
	"flag"
	"time"
)

type Options struct {
	address            string
	dataPath           string
	partitionDuration  time.Duration
	retention          time.Duration
	timestampPrecision string
	walBufferSize      int
}

func DefaultOptions() *Options {
	return &Options{
		address:           ":3000",
		dataPath:          "./data",
		partitionDuration: 24 * time.Hour,
		// 1 year
		retention:          8760 * time.Hour,
		timestampPrecision: "s",
		walBufferSize:      4096,
	}
}

func ParseArgs() (opts *Options, err error) {

	address := flag.String(
		"address",
		"",
		"Address of the tstore server. Ex 127.0.0.1:3000. Default is :3000")

	dataPath := flag.String(
		"dataPath",
		"",
		"Path to the folder where all of the data is stored. Defaults to ./data")

	partitionDuration := flag.Int(
		"partitionDuration",
		0,
		"The length in time of one partition of the database (In hours). Default is 24 hour.")

	retention := flag.Int(
		"retention",
		0,
		"Max age of the stored data (In hours). Default is 8760 or 1 year.")

	timestampPrecision := flag.String(
		"timestampPrecision",
		"",
		"Precision of the passed timestamps. Default is 's' for seconds.")

	walBufferSize := flag.Int(
		"walBufferSize",
		0,
		"The max size of the wal buffer. Default is 4096.")

	flag.Parse()

	opts = DefaultOptions()

	if *address != "" {
		opts.address = *address
	}
	if *dataPath != "" {
		opts.dataPath = *dataPath
	}
	if *partitionDuration != 0 {
		opts.partitionDuration = time.Duration(*partitionDuration) * time.Hour
	}
	if *retention != 0 {
		opts.retention = time.Duration(*retention) * time.Hour
	}
	if *timestampPrecision != "" {
		opts.timestampPrecision = *timestampPrecision
	}
	if *walBufferSize != 0 {
		opts.walBufferSize = *walBufferSize
	}
	return
}
