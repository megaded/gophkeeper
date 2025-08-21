package main

import (
	"fmt"
	"gophkeeper/internal/client/proto"
	"gophkeeper/internal/logger"
	"os"
)

func main() {
	logger.SetupLogger("info")
	/*if _, err := tea.NewProgram(ui.InitialMainModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	} */

	client := proto.NewKeeperClient()
	/* /*err := client.Register(context.TODO(), "1", "2")
	if err != nil {
		logger.Log.Error(err.Error())
	} */

	var filePath = "I:/Torrents/God Is a Bullet (2023)WEB-DLRip-AVC.mkv"
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	s, err := file.Stat()
	if err != nil {
		logger.Log.Error(err.Error())
	}
	size := s.Size()
	name := s.Name()
	logger.Log.Info(name)
	logger.Log.Info(fmt.Sprint(size))
	/* cfg := config.GetConfig()
	minio, err := minio.NewStorage(cfg)
	err = minio.UploadFile(context.Background(), "keeperrr", name, file, size)
	if err != nil {
		logger.Log.Error(err.Error())
	} */
	err = client.UploadBinaryFile(file, name, "описание", size)
	if err != nil {
		logger.Log.Error(err.Error())
	}

}
