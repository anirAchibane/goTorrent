package peer

type Handshake struct{
	Pstr string
	InfoHash [20]byte
	PeerID [20]byte
}