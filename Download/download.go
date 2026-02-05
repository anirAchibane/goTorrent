package download

import (
	peer "goTorrent/Peer"
	torrent "goTorrent/Torrent"
)

type Piece_work struct{
	Index int
	Hash [20]byte
	Length int
}

type Piece_result struct{
	Index int
	Buf []byte
}

type Torrent_task struct{
	Name string
	Length int
	Piece_length int
	Piece_hashes [][20]byte
	Info_hash [20]byte
	Peers []peer.Peer
}

// main download function:
func (t *Torrent_task) Download() ([]byte, error){

}

// helper functions:

func (t *Torrent_task) startWorker(p peer.Peer, work_queue chan *Piece_work, results chan *Piece_result){

}

func (t *Torrent_task) calculate_bounds_for_piece(index int) (begin int, end int){

}

func (t *Torrent_task) calculate_piece_size(index int) int{

}

func check_intergrity(piece_work *Piece_work, buf []byte) (error){

}