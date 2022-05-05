package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"thinkific-discord/internal/discord"
	"thinkific-discord/internal/discordBot"
	"thinkific-discord/internal/email"
	"thinkific-discord/internal/sheets"
	"thinkific-discord/internal/tgbot"
	"thinkific-discord/internal/webServer"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func init() {

	godotenv.Load(".env")
	go webServer.Listen()

}

func main() {
	var wg sync.WaitGroup
	for {
		wg.Add(1)
		go runApp(&wg)
		wg.Wait()
		fmt.Println("Recovering in 10 sec...")
		time.Sleep(10 * time.Second)
		// fmt.Println("Recoverered...")
	}
}

func runApp(wg *sync.WaitGroup) {
	defer func() {
		err := recover()
		if err != nil {
			tgbot.SendString(fmt.Sprint(err))
		}
		fmt.Println("Panic error...")
		wg.Done()
		return
	}()
	tgbot.Start()
	email.InitServer()
	sheets.InitService()
	discordBot.SetGuildId()

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

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shuting down Server ...")
}
