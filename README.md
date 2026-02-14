# GoTorrent - Single-file BitTorrent Client

A concurrent BitTorrent client implementation in Go that downloads single files from .torrent metadata files using the BitTorrent protocol.

## Overview

GoTorrent is a minimalist BitTorrent client built from scratch in Go. It demonstrates core BitTorrent protocol mechanics including:

- **Torrent file parsing** - Decodes bencoded .torrent files to extract metadata
- **Tracker communication** - Requests peer lists from HTTP trackers
- **Peer protocol** - Implements BitTorrent peer wire protocol (handshake, messages, bitfields)
- **Concurrent downloading** - Uses goroutines to download pieces from multiple peers simultaneously

## Installation

### Dependencies

```bash
go get github.com/jackpal/bencode-go
```

### Usage

Try downloading [debian.iso](https://cdimage.debian.org/debian-cd/current/amd64/bt-cd/#indexlist)!

```bash
# Clone the repository
git clone 
cd goTorrent

# Run directly
go run main.go debian-10.2.0-amd64-netinst.iso.torrent debian.iso

# Or using the compiled binary
./goTorrent debian-10.2.0-amd64-netinst.iso.torrent debian.iso
```

## Architecture

```
goTorrent/
├── main.go                  # Entry point
├── Torrent/
│   ├── torrent.go          # Torrent file parsing and hashing
│   └── tracker.go          # Tracker communication
├── Peer/
│   ├── message.go          # BitTorrent message protocol
│   ├── peers.go            # Peer parsing from tracker response
│   ├── handshake.go        # Handshake protocol implementation
│   └── bitfield.go         # Bitfield management
└── Download/
    ├── client.go           # Peer connection management
    └── download.go         # Download orchestration and worker pool
```

## Future Areas for Improvement

- UDP tracker support
- DHT (Distributed Hash Table) for trackerless operation
- Multi-file torrent support
- Magnet link support

## Resources

- [BitTorrent Protocol Specification](http://www.bittorrent.org/beps/bep_0003.html)
- [bencode-go](https://github.com/jackpal/bencode-go) for bencoding support
- [Jesse Li](https://blog.jse.li/posts/torrent/)'s blog for building a single-file client.