package torrentfile

import(
	"os"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func Open(path string) (TorrentFile, error) {
	// Open torrent File
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()
	// To be continued
}