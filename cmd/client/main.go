package main

import (
	"fmt"
	"gophkeeper/internal/client/ui"
	"gophkeeper/internal/logger"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logger.SetupLogger("info")
	/* client := proto.NewKeeperClient()

	_, err := client.Login(context.Background(), "1234", "1234")
	if err != nil {
		logger.Log.Error(err.Error())
	}
	client.AddCreditCard(context.Background(), dto.Card{Number: "111111", Description: "Хуй", CVV: "1345", Exp: "2020"})
	r, err := client.GetCreditCards(context.Background())
	fmt.Println(r) */
	if _, err := tea.NewProgram(ui.InitialMainModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	//client := proto.NewKeeperClient()

	/* var filePath = "I:/Torrents/God Is a Bullet (2023)WEB-DLRip-AVC.mkv"
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
	logger.Log.Info(fmt.Sprint(size)) */
	/* cfg := config.GetConfig()
	minio, err := minio.NewStorage(cfg)
	err = minio.UploadFile(context.Background(), "keeperrr", name, file, size)
	if err != nil {
		logger.Log.Error(err.Error())
	} */
	/* err = client.UploadBinaryFile(file, name, "описание", size)
	if err != nil {
		logger.Log.Error(err.Error())
	} */

}
