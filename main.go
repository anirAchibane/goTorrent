package main

import (
	"fmt"
	torrent "goTorrent/Torrent"

	//torrent "goTorrent/Torrent"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <torrentFile Path> <Download Path>")
		os.Exit(1)
	}

	// Parsing raw torrent file to bencoded structure:
	torrent_path := os.Args[1]

	bencoded_torrent, err := torrent.Open(torrent_path)
	if err != nil {
		fmt.Println("Error parsing torrent: ", err)
		os.Exit(1)
	}
	
	// Converting bencoded structure to TorrentFile structure:
	torrent_file, err := torrent.ToTorrent(bencoded_torrent)
	if err != nil {
		fmt.Println("Error converting to TorrentFile: ", err)
		os.Exit(1)
	}
	
	fmt.Println("Successfully parsed torrent file:")
	fmt.Println("Announce URL:", torrent_file.Announce)
	fmt.Println("Name:", torrent_file.Name)
	fmt.Println("Length:", torrent_file.Length)
	fmt.Println("Piece Length:", torrent_file.PieceLength)
	fmt.Printf("Info Hash: %x\n", torrent_file.InfoHash)
	fmt.Println("Number of Pieces:", len(torrent_file.PiecesHashes))

}
