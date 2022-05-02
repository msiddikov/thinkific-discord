package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"thinkific-discord/internal/discord"
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
	go webServer.Listen()
	email.InitServer()
	discordBot.SetGuildId()
	sheets.InitService()

	s := gocron.NewScheduler(time.UTC)
	interval := os.Getenv("SERVER_UPDATE_INTERVAL")
	if interval == "" {
		interval = "24h"
	}
	s.Every(interval).Do(discordBot.UpdateRoles)
	s.Every(interval).Do(sheets.UpdateCourses)
	s.Every("24h").Do(discordBot.AdjustRoles)
	s.StartAsync()
	discord.GenerateLink("")
}

func main() {
	// sheets.UpdateCourses()
	// discordBot.UpdateRoles()

	// Waiting for exit command from os
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shuting down Server ...")
}
