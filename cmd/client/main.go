package main

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/ui"
	"gophkeeper/internal/logger"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	Version string
	Time    string
)

func main() {
	logger.SetupLogger("info")
	/* client := proto.NewKeeperClient() */

	/* _, err := client.Login(context.Background(), "1234", "1234")
	if err != nil {
		logger.Log.Error(err.Error())
	} */
	/* client.AddCreditCard(context.Background(), dto.Card{Number: "111111", Description: "Хуй", CVV: "1345", Exp: "2020"})
	if err != nil {
		fmt.Println(err)
	}
	cards, err := client.GetCreditCards(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		fmt.Println("Карты")
		fmt.Println(cards)
	}

	err = client.AddCredentials(context.Background(), dto.Credentials{Login: "1111", Password: "gfg", Description: "454545"})
	if err != nil {
		fmt.Println(err)
	}

	creds, err := client.GetCredentials(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		fmt.Println("Креды")
		fmt.Println(creds)
	} */

	/* 	var filePath = "I:/ChromeDownload/vv-2.zip"
	   	file, err := os.Open(filePath)
	   	if err != nil {
	   		panic(err)
	   	}
	   	defer file.Close()
	   	s, err := file.Stat()
	   	if err != nil {
	   		logger.Log.Error(err.Error())
	   	}
	   	err = client.UploadBinaryFile(file, s.Name(), "описание")
	   	if err != nil {
	   		logger.Log.Error(err.Error())
	   	}
	   	files, err := client.GetBinaryFileList(context.Background())
	   	if err != nil {
	   		logger.Log.Error(err.Error())
	   	}
	   	fmt.Println(files)
	   	err = client.DownloadBinaryFile(context.Background(), files[0].Id)
	   	if err != nil {
	   		logger.Log.Error(err.Error())
	   	} */

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigChan
		cancel()
	}()
	if _, err := tea.NewProgram(ui.InitialMainModel(ctx)).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

}
