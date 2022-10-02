FlexTorrent
=============
Simple torrent cli tool for single-use torrent jobs build on top of [anacrolix/torrent](https://github.com/anacrolix/torrent).
```text
Usage: flextorrent -f <path to torrent file>
  -d string
        download directory path
  -f string
        torrent file path
  -i string
        indices of files to download (separated with ',' or '-' for ranges, e.g. 0,5,6-8,10)
  -m    print metadata in JSON format (e.g. file list) and exit
```
Periodically outputs progress to stdout:
```text
progress,<downloaded bytes count>,<total download length in bytes>
```
Prints message and exits with code 1 on error:
```text
error,<error message>
```
Finishes with exit code 0 on success

Can resume download on restart

**Each process creates .torrent.db, .torrent.db-shm and .torrent.db-wal files in the working directory, so you MUST NOT use the same working directory for multiple simultaneously running processes**