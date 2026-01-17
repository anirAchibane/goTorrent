package main

import (
	"fmt"
	torrent "goTorrent/Torrent"
	peer "goTorrent/Peer"
	"crypto/rand"
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
	torrent_file, err := torrent.To_torrent(bencoded_torrent)
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
	
	// retreiving peers:
	var peer_id [20]byte	// generating random peer id
	_, err = rand.Read(peer_id[:])
	if err != nil {
		fmt.Println("failed to read random bytes: %w", err)
		os.Exit(1)
	}
	
	tracker_repsonse, err := torrent.Request_peers(&torrent_file, peer_id, uint16(9090))
	if err != nil{
		fmt.Println("Connection to tracker failed: ", err)
		os.Exit(1)
	}

	peers, err := peer.Parse_peers(&tracker_repsonse)
	if err != nil{
		fmt.Println("Parsing peers failed: ", err)
	}
	for i := range peers{
		fmt.Println(i,"/ ip: ",peers[i].IP,", port: ", peers[i].Port) 
	}
}
