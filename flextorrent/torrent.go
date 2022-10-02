package flextorrent

import (
	"context"
	"encoding/json"
	"fmt"
	torrentLib "github.com/anacrolix/torrent"
)

func GetMetadata(ctx context.Context, filePath string) error {
	fc, err := newFlexClient(ctx, filePath, "")
	if err != nil {
		return err
	}
	name := fc.info.BestName()
	metadata := TorrentMetadata{
		Name:  name,
		Files: fc.files,
	}
	bytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	fc.close()
	return nil
}

func DownloadTorrent(ctx context.Context, filePath, downloadPath string, fileIndices FileIndices) error {
	fc, err := newFlexClient(ctx, filePath, downloadPath)
	if err != nil {
		return err
	}
	for i := range fc.files {
		cFile := fc.torrent.Files()[i]
		cFile.SetPriority(torrentLib.PiecePriorityNone)
	}
	for i := range fc.files {
		if fileIndices.Contains(i) {
			file := fc.torrent.Files()[i]
			file.SetPriority(torrentLib.PiecePriorityNormal)
			fc.downloadLength += uint(file.Length())
			fc.selected[i] = true
		}
	}
	err = fc.join(ctx)
	fc.close()
	return err
}
