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

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")

}

func main() {
	go webServer.Listen()
	email.InitServer()
	discordBot.SetGuildId()
	sheets.InitService()
	sheets.UpdateCourses()
	sheets.UpdateRoles()

	// Waiting for exit command from os
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shuting down Server ...")
}
