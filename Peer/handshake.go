package peer

import (
	"errors"
	"io"
)

type Handshake struct{
	Pstr [19]byte
	Info_hash [20]byte
	Peer_id [20]byte
}

// create handshake for own machine
func Create_handshake(info_hash [20]byte,peer_id [20]byte) (*Handshake){
	var pstr [19]byte
	_ = copy(pstr[:], "BitTorrent protocol")

	return &Handshake{
		Pstr: pstr,
		Info_hash: info_hash,
		Peer_id: peer_id,
	}
}

// Serialize handshake
func Serialize_handshake(handshake *Handshake) ([]byte) {	
	str := make([]byte, 68) // 1 + 19 + 8 + 20 + 20

	str[0] = byte(len(handshake.Pstr))			// adding pstr len
	i := 1	
	i += copy(str[i:],handshake.Pstr[:])  		// adding pstr
	i += copy(str[i:],make([]byte,8))  			// adding 8 reserved bytes
	i += copy(str[i:], handshake.Info_hash[:])  // adding info hash
	i += copy(str[i:], handshake.Peer_id[:]) 	// adding peer id

	return str
}

// parse handshake
func Parse_handshake(r io.Reader) (*Handshake, error){
	var Pstr_val [19]byte
	var Info_hash_val [20]byte
	var Peer_id_val [20]byte

	buff := make([]byte, 68)

	_ , err := io.ReadFull(r, buff)
	if err != nil{
		return nil,err
	}

	pstr_len := int(buff[0])
	if pstr_len != 19{
		err = errors.New("Incompatible Protocol")
		return nil, err
	}
	
	_ = copy(Pstr_val[:], buff[1:20])
	_ = copy(Info_hash_val[:], buff[28:48])
	_ = copy(Peer_id_val[:], buff[48:68])

	h := Handshake{
		Pstr: Pstr_val,
		Info_hash: Info_hash_val,
		Peer_id: Peer_id_val,
	}

	return &h, nil
}