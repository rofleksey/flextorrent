package flextorrent

import (
	"context"
	"fmt"
	torrentLib "github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"math"
	"time"
)

type FlexClient struct {
	client         *torrentLib.Client
	torrent        *torrentLib.Torrent
	info           *metainfo.Info
	files          []TorrentFile
	selected       []bool
	downloadLength uint
}

type TorrentFile struct {
	Path   string `json:"path"`
	Length uint   `json:"length"`
}

type TorrentMetadata struct {
	Name  string        `json:"name"`
	Files []TorrentFile `json:"files"`
}

func newFlexClient(ctx context.Context, filePath, downloadPath string) (*FlexClient, error) {
	clientConfig := torrentLib.NewDefaultClientConfig()
	if len(downloadPath) > 0 {
		clientConfig.DataDir = downloadPath
	}
	client, err := torrentLib.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}
	// drop any pre-saved torrents
	retainedTorrents := client.Torrents()
	for _, t := range retainedTorrents {
		t.Drop()
	}
	torrent, err := client.AddTorrentFromFile(filePath)
	if err != nil {
		return nil, err
	}
	select {
	case <-torrent.GotInfo():
		break
	case <-ctx.Done():
		return nil, ErrInterrupted
	}
	info := torrent.Info()
	files := make([]TorrentFile, 0, len(info.Files))
	selected := make([]bool, 0, len(info.Files))
	for _, f := range torrent.Files() {
		file := TorrentFile{
			Path:   f.DisplayPath(),
			Length: uint(f.Length()),
		}
		files = append(files, file)
		selected = append(selected, false)
	}
	return &FlexClient{
		client:   client,
		torrent:  torrent,
		info:     info,
		files:    files,
		selected: selected,
	}, nil
}

func (fc *FlexClient) join(ctx context.Context) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	cFiles := fc.torrent.Files()
	var lastBytesRead uint
	lastBytesRead = math.MaxUint
	for {
		select {
		case <-ctx.Done():
			return ErrInterrupted
		case <-fc.torrent.Closed():
			return ErrClientError
		case <-ticker.C:
			bytesRead := uint(0)
			for i := range fc.files {
				if fc.selected[i] {
					bytesRead += uint(cFiles[i].BytesCompleted())
				}
			}
			if bytesRead != lastBytesRead {
				fmt.Printf("progress,%d,%d\n", bytesRead, fc.downloadLength)
				lastBytesRead = bytesRead
			}
			if bytesRead < fc.downloadLength {
				continue
			}
			return nil
		}
	}
}

func (fc *FlexClient) close() {
	fc.torrent.Drop()
	<-fc.torrent.Closed()
	fc.client.Close()
	<-fc.client.Closed()
}
