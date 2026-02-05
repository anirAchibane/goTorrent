package download

import (
	"encoding/binary"
	"errors"
	"fmt"
	peer "goTorrent/Peer"
	"io"
	"net"
	"time"
)

type Peer_client struct {
	Bitfield  peer.Bit_field
	Peer       peer.Peer
	isChoked   bool
	Connection net.Conn
}

func Connect_to_peer(p *peer.Peer, info_hash [20]byte, self_peer_id [20]byte) (*Peer_client, error) {
	// establish connection with peer
	conn, err := establish_connection(*p)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer conn.SetDeadline(time.Time{})

	// make handshake
	handshake, err := make_handshake(conn, info_hash, self_peer_id)
	if err != nil {
		return nil, err
	}
	if string(handshake.Pstr[:]) != "BitTorrent protocol" {
		return nil, errors.New("Incompatible peer")
	}

	// receive bitfield
	bitfield, err := receive_bitfield(conn)
	if err != nil{
		conn.Close()
		return nil, err
	}

	return &Peer_client{
		Peer:       *p,
		isChoked:   true,
		Bitfield: bitfield,
		Connection: conn,
	}, nil
}

// Helper functions for connecting to peer
func make_handshake(conn net.Conn, info_hash [20]byte, self_peer_id [20]byte) (peer.Handshake, error) {
	conn.SetDeadline(time.Now().Add(time.Second*5))
	defer conn.SetDeadline(time.Time{})

	// send handshake
	self_handshake := peer.Create_handshake(info_hash, self_peer_id)
	self_message := peer.Serialize_handshake(self_handshake)

	_, err := conn.Write(self_message)
	if err != nil {
		return peer.Handshake{}, err
	}

	// receive handshake
	received_message := make([]byte, 68)
	received_len, err := io.ReadFull(conn, received_message)
	if received_len != 68 {
		return peer.Handshake{}, errors.New("Wrong Handshake form!")
	}

	received_handshake, err := peer.Parse_handshake(received_message)

	if err != nil {
		return peer.Handshake{}, err
	}

	// check received handshake
	if string(received_handshake.Pstr[:]) != "BitTorrent protocol" || received_handshake.Info_hash != info_hash {
		return peer.Handshake{}, nil
	}
	return *received_handshake, nil
}

func establish_connection(p peer.Peer) (net.Conn, error) {

	address := fmt.Sprintf("%s:%d", p.IP.String(), p.Port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func receive_bitfield(conn net.Conn) ([]byte, error){
	conn.SetDeadline(time.Now().Add(time.Second*5))
	defer conn.SetDeadline(time.Time{})

	message, err := peer.Parse_message(conn)
	if err != nil{
		return nil,err
	}

	if message == nil{
		return nil, errors.New("Expected bitfield but got empty message!")
	}

	if message.ID != peer.MsgBitfield{
		return nil, fmt.Errorf("Expected Bitfield but go ID %d",message.ID)
	}

	return message.Payload, nil
}
 
// messaging methods:

func (client *Peer_client) Read() (*peer.Message,error){
	message, err := peer.Parse_message(client.Connection)
	return message, err
}

func (client *Peer_client) Send_interested() (error){
	message := peer.Message{
		Length: 1,
		ID: peer.MsgInterested,
		Payload:nil,
	}

	serialized_message, err := peer.Serialize_message(&message)
	if err != nil{
		return err
	}

	_, err = client.Connection.Write(serialized_message)
	if err != nil{
		return err
	}

	return nil
}

func (client *Peer_client) Send_not_interested() (error){
	message := peer.Message{
		Length: 1,
		ID: peer.MsgNotInterested,
		Payload: nil,
	}

	serialized_message, err := peer.Serialize_message(&message)
	if err != nil{
		return err
	}

	_, err = client.Connection.Write(serialized_message)
	if err != nil{
		return err
	}

	return nil
}


func (client *Peer_client) Send_request(index, begin, length int) (error){
	request_payload := make([]byte, 12)
	binary.BigEndian.PutUint32(request_payload[0:4],uint32(index))
	binary.BigEndian.PutUint32(request_payload[4:8],uint32(begin))
	binary.BigEndian.PutUint32(request_payload[8:12],uint32(length))

	request_message := peer.Message{
		Length: 13,
		ID: peer.MsgRequest,
		Payload: request_payload,
	}

	serialized_message, err := peer.Serialize_message(&request_message) 
	if err != nil{
		return err
	}

	_, err = client.Connection.Write(serialized_message)
	if err != nil{
		return err
	}

	return nil
}

func (client *Peer_client) Send_unchoke() (error){
	message := peer.Message{
		Length: 1,
		ID: peer.MsgUnchoke,
		Payload: nil,
	}

	serialized_message, err := peer.Serialize_message(&message)
	if err != nil{
		return err
	}
	
	_,err = client.Connection.Write(serialized_message)
	if err != nil{
		return err
	}

	return nil
}

func (client *Peer_client) Send_have(index int) (error){
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload,uint32(index))

	message := peer.Message{
		Length: 5,
		ID: peer.MsgHave,
		Payload: payload,
	}

	serialized_message, err := peer.Serialize_message(&message)
	if err != nil{
		return err
	}

	_, err = client.Connection.Write(serialized_message)
	if err != nil{
		return err
	}

	return nil
}
