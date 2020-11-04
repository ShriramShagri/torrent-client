package main

import(
	"os"
	"log"

	"torrentfile"
)

func main(){
	// Get path for torrent file and path to save the downloads
	torrentFilePath := os.Args[1]
	downloadPath := os.Args[2]

	tf, err = torrentfile.Open(torrentFilePath)
	if err != nil {
		log.fatal(err)
	}
}