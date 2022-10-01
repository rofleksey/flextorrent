package main

import (
	"context"
	"flag"
	"flextorrent/flextorrent"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func printUsageAndExit() {
	fmt.Println("Usage: go flextorrent.go --file <path to torrent file>")
	flag.PrintDefaults()
	os.Exit(1)
}

func handleInterrupt(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()
}

func main() {
	filePath := flag.String("file", "", "torrent file path")
	fileIndicesStr := flag.String("indices", "", "file indices to download (separated with ',' or '-' for ranges, e.g. 5,6-8,10)")
	metadataOnly := flag.Bool("metadata", false, "metadata only mode")
	downloadDir := flag.String("dir", "", "download directory path")
	flag.Parse()

	if len(*filePath) == 0 {
		printUsageAndExit()
	}

	var fileIndices flextorrent.FileIndices
	var err error

	if len(*fileIndicesStr) > 0 {
		fileIndices, err = flextorrent.ParseFileIndices(*fileIndicesStr)
		if err != nil {
			printUsageAndExit()
		}
	} else {
		fileIndices, _ = flextorrent.ParseFileIndices("*")
	}

	ctx, cancel := context.WithCancel(context.Background())
	handleInterrupt(cancel)
	if *metadataOnly {
		err = flextorrent.GetMetadata(ctx, *filePath)
	} else {
		err = flextorrent.DownloadTorrent(ctx, *filePath, *downloadDir, fileIndices)
	}
	if err != nil {
		fmt.Printf("error,%v\n", err)
		os.Exit(1)
	}
}
