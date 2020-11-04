package torrentfile

import(
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

// Structure of Torrent File
type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type bencodeInfo struct{
	Pieces 		string	`bencode:"pieces"`
	PieceLength	int		`bencode:"piece length"`
	Length 		int		`bencode:"length"`
	Name 		string	`bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

// Temporary
func (i *TorrentFile) print(){
	fmt.Println(i.Name)
}

func Open(path string) (TorrentFile, error) {

	// Open the file from specified path
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}

	// Close file when return from the function
	defer file.Close()

	// Data of .torrent file
	fileData := bencodeTorrent{}

	// imported function to extract data from bencode
	err = bencode.Unmarshal(file, &fileData)
	if err != nil {
		return TorrentFile{}, err
	}

	// Writedata to file data object and return
	return fileData.toTorrentFile()
}

func (i *bencodeInfo) hash() ([20]byte, error) {

	// Declare bytes Buffer
	var buf bytes.Buffer

	// Unmarshal a bencode stream into an bytes Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}

	// 
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {

	// Length of SHA-1 hash
	hashLen := 20 

	// Create buffer with all pieces
	buf := []byte(i.Pieces)

	// If length ofhash of pieces is not a multiple of hashlength of sha1, then the hash is corrupted
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}

	// Number of hashes/Pieces of hash
	numHashes := len(buf) / hashLen

	// Create slice of arrays to return back each separate hash
	hashes := make([][20]byte, numHashes)

	// Separate hash of each data piece0
	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (fileData *bencodeTorrent) toTorrentFile() (TorrentFile, error) {

	// Hash value of the torrentfile itself
	infoHash, err := fileData.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}

	// Hashes of each piece of data
	pieceHashes, err := fileData.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}

	// Write to the TorrentFile object and return
	finalData := TorrentFile{
		Announce:    fileData.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: fileData.Info.PieceLength,
		Length:      fileData.Info.Length,
		Name:        fileData.Info.Name,
	}
	return finalData, nil
}