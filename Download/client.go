package download

import (
	"errors"
	"fmt"
	peer "goTorrent/Peer"
	"io"
	"net"
	"time"
)

type Peer_client struct {
	Bit_field    peer.Bit_field
	Peer         peer.Peer
	isChoked     bool
	Connection   net.Conn
}

func Connect_to_peer(p *peer.Peer,info_hash [20]byte, self_peer_id [20]byte) (*Peer_client,error) {
	// establish connection with peer
	conn, err := establish_connection(*p)
	if err != nil{
		return nil, err
	}

	conn.SetDeadline(time.Now().Add(time.Second*5))
	defer conn.SetDeadline(time.Time{})

	// make handshake
	handshake, err := make_handshake(conn, info_hash,self_peer_id)
	if err != nil{
		return nil, err
	}
	if handshake.Pstr[:] == nil{
		return nil, errors.New("Incompatible peer")
	}

	return &Peer_client{
		Bit_field: nil,
		Peer: *p,
		isChoked: true,
		Connection: conn,
	}, nil
}

func (client *Peer_client) Run() ([][]byte, error){
	for{
		client.Connection.SetDeadline(time.Now().Add(time.Second*5))
		defer client.Connection.SetDeadline(time.Time{})

		message, err := peer.Parse_message(client.Connection)
		if err != nil{
			return nil, err
		}
		switch message.ID{
			case 0: // Choke
				client.isChoked = true
			case 1: // unChoke
				client.isChoked = false				
			case 4: // have
				
			case 5: // bitfield
				client.Bit_field = message.Payload
			case 6: // request
				
			case 7: // piece
				
			case 8: // cancel

			default:
				// for other cases, do nothing, since we're just leeching
		}
	}
}
// Helper functions for connecting to peer
func make_handshake(conn net.Conn, info_hash [20]byte, self_peer_id [20]byte) (peer.Handshake, error) {
	// send handshake
	self_handshake := peer.Create_handshake(info_hash,self_peer_id)
	self_message := peer.Serialize_handshake(self_handshake)

	_ ,err := conn.Write(self_message)
	if err != nil{
		return peer.Handshake{}, err
	}
	
	// receive handshake
	received_message := make([]byte, 68)
	received_len, err := io.ReadFull(conn, received_message)
	if received_len != 68{
		return peer.Handshake{}, errors.New("Wrong Handshake form!")
	}

	received_handshake, err := peer.Parse_handshake(received_message)
	if err != nil{
		return peer.Handshake{}, err
	}

	// check received handshake
	if string(received_handshake.Pstr[:]) != "BitTorrent protocol" || received_handshake.Info_hash != info_hash{
		return peer.Handshake{}, nil
	}
	return *received_handshake, nil
}

func establish_connection(p peer.Peer) (net.Conn, error) {

	address := fmt.Sprintf("%s:%d", p.IP.String(),p.Port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil{
		return nil, err
	}

	return conn, nil
}