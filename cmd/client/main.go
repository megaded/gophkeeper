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
	if _, err := tea.NewProgram(ui.InitialMainModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
	/* client := proto.NewKeeperClient()
	err := client.Register(context.TODO(), "1", "2")
	if err != nil {
		logger.Log.Error(err.Error())
	} */
}
