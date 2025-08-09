package server

import (
	"os"
)

func main() {
	var filePath = "I:/Torrents/God Is a Bullet (2023)WEB-DLRip-AVC.mkv"
	os.ReadFile()
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	file.Read()
}
