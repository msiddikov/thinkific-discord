package discordBot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"thinkific-discord/internal/email"
	"thinkific-discord/internal/types"
	"time"
)

var (
	GuildId string
	host    = "https://discord.com/api"
)

type (
	guildResp []struct {
		Id   string
		Name string
	}

	chanResp []struct {
		Id string
	}

	inviteResp struct {
		Code string
	}

	setRolesBody struct {
		Roles []string `json:"roles"`
	}
)

func discordReq(req *http.Request) []byte {
	client := &http.Client{}

	req.Header.Add("Authorization", "Bot "+os.Getenv("DISCORD_BOT_SECRET"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error when sending request to the server")
	}

	if resp.StatusCode > 299 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Error when sending request to Discord " + req.URL.Path + " " + string(responseBody))
		return nil
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseBody
}

func SetGuildId() {
	req, _ := http.NewRequest(http.MethodGet, host+"/users/@me/guilds", nil)
	responseBody := discordReq(req)

	guilds := guildResp{}
	json.Unmarshal(responseBody, &guilds)
	if len(guilds) > 0 {
		GuildId = guilds[0].Id
		fmt.Println("Bot is managing " + guilds[0].Name)
	} else {
		botLink := "https://discord.com/api/oauth2/authorize?client_id=" + os.Getenv("DISCORD_CLIENT_ID") + "&permissions=268435511&scope=bot applications.commands"
		email.SendBotAddLink(os.Getenv("ADMIN_EMAIL"), botLink, "Bot Administrator")
		fmt.Println("Need to add this bot to your server. Please check your email")
	}
}

func GetRoles() types.RolesResp {
	req, _ := http.NewRequest(http.MethodGet, host+"/guilds/"+GuildId+"/roles", nil)
	responseBody := discordReq(req)

	roles := types.RolesResp{}
	json.Unmarshal(responseBody, &roles)
	return roles
}

func BanUser(userId string) {

}

func SetRoles(discordUserId string, roles []types.CurrentRole) []types.CurrentRole {
	rolesToSet := []string{}
	i := 0
	for i < len(roles) {
		if roles[i].Expire.Before(time.Now()) || roles[i].RoleId == "" {
			roles[i] = roles[len(roles)-1]
			roles = roles[:len(roles)-1]
			continue
		}

		rolesToSet = append(rolesToSet, roles[i].RoleId)
		i++
	}
	body := setRolesBody{rolesToSet}
	bodybytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPatch, host+"/guilds/"+GuildId+"/members/"+discordUserId, bytes.NewReader(bodybytes))
	discordReq(req)
	return roles
}

func GetInviteLink() string {

	body := struct {
		MaxAge int `json:"max_age"`
	}{
		MaxAge: 60,
	}
	bodybytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, host+"/channels/"+os.Getenv("DISCORD_INVITE_CHAN_ID")+"/invites", bytes.NewReader(bodybytes))
	responseBody := discordReq(req)
	fmt.Println(responseBody)

	invite := inviteResp{}
	json.Unmarshal(responseBody, &invite)

	return invite.Code
}
