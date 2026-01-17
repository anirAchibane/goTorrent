package peer

import (
	"encoding/binary"
	"errors"
	"goTorrent/Torrent"
	"net"
)

type Peer struct{
	IP net.IP;
	Port uint16
}

func Parse_peers(t *torrent.Tracker_response) ([]Peer, error){
	const peer_size = 6
	
	// if Peers string length is incompatible
	if len(t.Peers) % peer_size != 0{
		err := errors.New("Malformed Peers string")
		return nil, err
	}

	number_of_peers := len(t.Peers) / peer_size
	peers_list := make([]Peer,number_of_peers)
	
	// adding peers one by one
	for i := range number_of_peers{
		step := i * peer_size
		peers_list[i] = Peer{
			IP: net.IP([]byte(t.Peers[step:step+4])),
			Port: binary.BigEndian.Uint16([]byte(t.Peers[step+4:step+6])),
		}
	}

	return peers_list, nil
}