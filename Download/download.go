package download

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	peer "goTorrent/Peer"
	"log"
	"runtime"
	"time"
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
	Self_ID [20]byte
	Name string
	Length int
	Piece_length int
	Piece_hashes [][20]byte
	Info_hash [20]byte
	Peers []peer.Peer
}

// main download function:
func (t *Torrent_task) Download() ([]byte, error){
	fmt.Println("Starting Download for", t.Name+"...")

	work_queue := make(chan *Piece_work, len(t.Piece_hashes))
	results := make(chan *Piece_result, len(t.Piece_hashes))

	// fill work queue
	for index, hash := range t.Piece_hashes{
		length := t.calculate_piece_size(index)
		work_queue <- &Piece_work{
			Index: index,
			Hash: hash,
			Length: length,
		}
	}

	// download from peers
	fmt.Println("Connecting With Peers...")
	for _, peer := range t.Peers{
		go t.start_worker(peer, work_queue, results)
	}

	// receive results
	buffer := make([]byte, t.Length)
	done_piece := 0
	for done_piece < len(t.Piece_hashes){
		res := <- results
		begin, end := t.calculate_bounds_for_piece(res.Index)
		copy(buffer[begin:end],res.Buf)

		done_piece++

		percentage := (float64(done_piece) / float64(len(t.Piece_hashes)))*100
		number_of_workers := runtime.NumGoroutine() - 1
		log.Println("(",int(percentage),"%) ---- Downloaded piece ", res.Index, " from ",number_of_workers, "peers")
	}
	close(work_queue)

	return buffer,nil
}

// helper functions:

func (t *Torrent_task) start_worker(p peer.Peer, work_queue chan *Piece_work, results chan *Piece_result) {
	client, err := Connect_to_peer(&p,t.Info_hash,t.Self_ID)
	if err != nil{
		return
	}

	defer client.Connection.Close()

	client.Send_unchoke()
	client.Send_interested()

	for client.IsChoked {
		msg, err := client.Read()
		if err != nil{
			log.Println("Error waiting for unchoke:", err.Error())
			return
		}

		if msg == nil{
			continue
		}

		if msg.ID == peer.MsgUnchoke {
			client.IsChoked = false
			break
		}
	}
	
	for pw := range work_queue{
		if !client.Bitfield.Has_piece(pw.Index){
			work_queue <- pw
			continue
		}

		result_buf, err := t.attempt_download_piece(client, pw)
		if err != nil{
			work_queue <- pw
			return
		}

		err = check_intergrity(pw,result_buf)
		if err != nil{
			log.Printf("Invalid Piece: %s", err.Error())
			work_queue <- pw
			continue
		}

		res := Piece_result{
			Index: pw.Index,
			Buf: result_buf,
		}

		results <- &res

		err = client.Send_have(pw.Index)
		if err != nil{
			log.Println("Exiting worker:", err.Error())
			return
		}
	}
}

func (t *Torrent_task) attempt_download_piece(c *Peer_client, pw *Piece_work) ([]byte, error){
	c.Connection.SetDeadline(time.Now().Add(time.Second*30))
	defer c.Connection.SetDeadline(time.Time{})
	
	begin := 0
	piece_content := make([]byte, pw.Length)
	
	for begin < pw.Length{
		length := 16384
		if begin + length >= pw.Length{
			length = pw.Length - begin
		}
		
		err := c.Send_request(pw.Index, begin, length)
		if err != nil{
			return nil, err
		}
		
		for {

			message, err := c.Read()
			if err != nil{
				return nil, err
			}

			if message == nil{
				continue
			}

			switch message.ID{
			case peer.MsgChoke:
				c.IsChoked = true
				return nil, fmt.Errorf("Choked by Peer During Download")
			
			case peer.MsgUnchoke:
				c.IsChoked = false
			
			case peer.MsgHave:
				index := binary.BigEndian.Uint32(message.Payload)
				c.Bitfield.Have_update(int(index))
			
			case peer.MsgPiece:
				if len(message.Payload) < 8{
					return nil, fmt.Errorf("Invalid Piece Message!")
				}

				message_index := binary.BigEndian.Uint32(message.Payload[:4])
				message_begin := binary.BigEndian.Uint32(message.Payload[4:8])

				if int(message_index) != pw.Index {
					continue
				}
				
				if int(message_begin) != begin {
					continue
				}

				begin += copy(piece_content[begin:], message.Payload[8:])

				goto next_block
			}
		}

		next_block:
		
	}

	return piece_content, nil
}

func (t *Torrent_task) calculate_bounds_for_piece(index int) (begin int, end int){
	begin = index * t.Piece_length
	end = begin + t.Piece_length
	
	if end > t.Length{
		end = t.Length
	}

	return begin, end
}

func (t *Torrent_task) calculate_piece_size(index int) int{
	begin, end := t.calculate_bounds_for_piece(index)
	return end - begin
}

func check_intergrity(piece_work *Piece_work, buf []byte) (error){
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], piece_work.Hash[:]){
		return fmt.Errorf("Index %d failed integrity!", piece_work.Index)
	}
	return nil
}