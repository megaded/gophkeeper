package main

import (
	"context"
	"fmt"
	"gophkeeper/internal/client/proto"
	"gophkeeper/internal/client/ui"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/dto"
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
	TestClient()
}

func main_client() {
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

func TestClient() {
	client := proto.NewKeeperClient()
	ctx := context.Background()

	_, err := client.Login(context.Background(), "1234", "1234")
	if err != nil {
		logger.Log.Error(err.Error())
	}
	// Карты
	/* client.AddCreditCard(context.Background(), dto.Card{Number: "111111", Description: "Хуй", CVE: "1345", Exp: "2020"})
	if err != nil {
		fmt.Println(err)
	}
	cards, err := client.GetCreditCards(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		fmt.Println("Карты :")
		fmt.Println(cards)
	}
	card := cards[0]
	err = client.UpdateCreditCard(ctx, dto.Card{Id: card.Id, Number: "1244545", Exp: "56565", CVE: "564565", Description: "Новое описание"})
	if err != nil {
		fmt.Println(err)
	}

	err = client.DeleteCreditCard(ctx, card.Id)
	if err != nil {
		fmt.Println(err)
	} */

	//Креды
	err = client.AddCredentials(context.Background(), dto.Credentials{Login: "1111", Password: "gfg", Description: "454545"})
	if err != nil {
		logger.Log.Error(err.Error())
	}

	creds, err := client.GetCredentials(context.Background())
	if err != nil {
		logger.Log.Error(err.Error())
	}

	cred := creds[0]
	cred.Description = "Поменяли"
	err = client.UpdateCreditial(ctx, cred)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	err = client.DeleteCreditial(ctx, cred.Id)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	//Текст
	err = client.UploadText(ctx, "Я контент", "Описание")
	if err != nil {
		logger.Log.Error(err.Error())
	}

	texts, err := client.GetTextList(ctx)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	text := texts[0]

	err = client.UpdateText(ctx, dto.Text{Content: "new", Description: "new", Id: text.Id})
	if err != nil {
		logger.Log.Error(err.Error())
	}

	err = client.DeleteText(ctx, text.Id)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	//Файлы
	var filePath = "I:/ChromeDownload/vv-2.zip"
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err != nil {
		logger.Log.Error(err.Error())
	}
	err = client.UploadBinaryFile(ctx, filePath, "описание")
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
	}
	updFile := "I:/ChromeDownload/go-1.24.5-win32-x64.zip"
	err = client.UpdateBinaryFile(ctx, updFile, "новое описание", files[0].Id)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	err = client.DownloadBinaryFile(context.Background(), files[0].Id)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	err = client.DeleteBinaryFile(ctx, files[0].Id)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}
