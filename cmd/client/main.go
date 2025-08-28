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
