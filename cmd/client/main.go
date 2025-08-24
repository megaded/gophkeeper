package main

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/proto"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/dto"
)

func main() {
	logger.SetupLogger("info")
	client := proto.NewKeeperClient()

	_, err := client.Login(context.Background(), "1234", "1234")
	if err != nil {
		logger.Log.Error(err.Error())
	}
	client.AddCreditCard(context.Background(), dto.Card{Number: "111111", Description: "Хуй", CVV: "1345", Exp: "2020"})
	if err != nil {
		fmt.Println(err)
	}
	cards, err := client.GetCreditCards(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		fmt.Println(cards)
	}

	err = client.AddCredentials(context.Background(), dto.Credentials{Login: "1111", Password: "gfg", Description: "454545"})
	if err != nil {
		fmt.Println(err)
	}

	creds, err := client.GetCreditCards(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		fmt.Println(creds)
	}

	/* if _, err := tea.NewProgram(ui.InitialMainModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	} */

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
