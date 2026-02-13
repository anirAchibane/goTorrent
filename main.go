package main

import (
	"crypto/rand"
	"fmt"
	download "goTorrent/Download"
	peer "goTorrent/Peer"
	torrent "goTorrent/Torrent"
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
	
	fmt.Println("Successfully Parsed Torrent File:")
	fmt.Println("Announce URL:", torrent_file.Announce)
	fmt.Println("Name:", torrent_file.Name)
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

	torrent_task := download.Torrent_task{
		Self_ID: peer_id,
		Name: torrent_file.Name,
		Length: torrent_file.Length,
		Piece_length: torrent_file.PieceLength,
		Piece_hashes: torrent_file.PiecesHashes,
		Info_hash: torrent_file.InfoHash,
		Peers: peers,
	}

	buffer, err := torrent_task.Download()
	if err != nil{
		fmt.Println("Couldn't complete download: ", err)
		os.Exit(1)
	}

	output_path := os.Args[2]
	outfile , err := os.Create(output_path)
	if err != nil{
		fmt.Println("Couldn't create output path:", err)
		os.Exit(1)
	}

	defer outfile.Close()
	_, err = outfile.Write(buffer)
	if err != nil{
		fmt.Println("Couldn't complete installation: ", err)
		os.Exit(1)
	}

}
