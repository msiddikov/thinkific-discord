package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"thinkific-discord/internal/discordBot"
	"thinkific-discord/internal/email"
	"thinkific-discord/internal/sheets"
	"thinkific-discord/internal/webServer"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Days().At("00:00").Do(discordBot.AdjustRoles)
}

func main() {
	go webServer.Listen()
	email.InitServer()
	discordBot.SetGuildId()
	sheets.InitService()
	sheets.UpdateCourses()
	discordBot.UpdateRoles()

	// Waiting for exit command from os
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shuting down Server ...")
}
