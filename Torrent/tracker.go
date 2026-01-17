package torrent

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
)

type Tracker_response struct{
	Interval int `bencode:"interval"`
	Peers string `bencode:"peers"`
}

// Building tracker URL with params (helper func)
func build_tracker_URL(t *Torrent_file, peerID [20]byte, port uint16) (string,error) { 
	// parse to url struct
	requ_url, err := url.Parse(t.Announce)
	
	if err != nil{
		return "", err
	}
	// add params
	params := url.Values{
		"info_hash": []string{string(t.InfoHash[:])},
		"peer_id" : []string{string(peerID[:])},
		"port": []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
        "downloaded": []string{"0"},
        "compact":    []string{"1"},
        "left":       []string{strconv.Itoa(t.Length)},	
	}

	requ_url.RawQuery = params.Encode()
	return requ_url.String(), nil
}

// Making GET request to tracker 
func Request_peers(t *Torrent_file,peerID [20]byte, port uint16) (Tracker_response, error) {
	
	requ_url, err := build_tracker_URL(t,peerID,port)
	if err != nil{
		return Tracker_response{}, err
	}

	resp, err := http.Get(requ_url)
	if err != nil{
		return Tracker_response{}, err
	}

	defer resp.Body.Close()

	trackerResponse := Tracker_response{}
	err = bencode.Unmarshal(resp.Body, &trackerResponse)
	if err != nil{
		return Tracker_response{} ,err
	}

	return trackerResponse,nil
}