package torrent

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

// Structures:

type bencode_inf struct {
	Pieces       string `bencode:"pieces"`
	PiecesLength int	`bencode:"piece length"`
	Length       int	`bencode:"length"`
	Name         string	`bencode:"name"`
}
type bencode_torrent struct {
	Announce string;
	Info     bencode_inf
}

type Torrent_file struct {
	Announce     string
	InfoHash     [20]byte
	PiecesHashes [][20]byte
	PieceLength  int
	Length       int
	Name         string
}

// Functions:

func Open(path string) (*bencode_torrent, error) { // return bencoded struct from .torrent path
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	bto := bencode_torrent{}

	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return nil, err
	}

	return &bto, nil
}

func To_torrent(bto *bencode_torrent) (Torrent_file, error) {
	infoHash, err := hash(&bto.Info)
	if err != nil {
		fmt.Println("Error converting to Torrent:", err)
		os.Exit(1)
	}

	piecesHashes, err := pieces_hash(&bto.Info)
	if err != nil {
		fmt.Println("Error converting to Torrent:", err)
		os.Exit(1)
	}
	return Torrent_file{
		Announce:     bto.Announce,
		PieceLength:  bto.Info.PiecesLength,
		InfoHash:     infoHash,
		PiecesHashes: piecesHashes,
		Length:       bto.Info.Length,
		Name:         bto.Info.Name,
	}, nil
}

// hashing functions: hashes info and pieces using SHA1

func hash(i *bencode_inf) ([20]byte, error) {

	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)	// bencoding info struct into buf

	if err != nil {
		return [20]byte{}, err
	}

	hash := sha1.Sum(buf.Bytes())		// hashing the buf
	return hash, nil
}

func pieces_hash(i *bencode_inf) ([][20]byte, error) {
	hashLen := 20

	buf := []byte(i.Pieces)

	if len(buf)%hashLen != 0 { // if the pieces are malformed
		err := errors.New("Malformed Piece lengths")
		return nil, err
	}

	numHashes := len(buf) / hashLen			
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])	// copying each piece hash into the hashes slice
	}

	return hashes, nil
}