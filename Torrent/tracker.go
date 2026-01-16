package torrent

import (
	"io"
	"net/http"
	"encoding/json"
)

type TrackerResponse struct{
	Interval int "json:Interval"
	Peers string "json:Peers"
}

// Building tracker URL with params
func buildTrackerURL(t *TorrentFile, peerID string) (string,error) {
	return "",nil
}

// Making GET request to tracker 
func RequestPeers(peerID [20]byte, port uint16) (TrackerResponse, error) {
	return TrackerResponse{},nil
}

// HELPER FUNC
// Parse Peers from compact format to list of IP:Port
func parsePeers(){

}