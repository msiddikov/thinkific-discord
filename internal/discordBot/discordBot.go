package discordBot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"thinkific-discord/internal/email"
	"thinkific-discord/internal/sheets"
	"thinkific-discord/internal/tgbot"
	"thinkific-discord/internal/thinkific"
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

func discordReq(req *http.Request) ([]byte, error) {
	client := &http.Client{}

	req.Header.Add("Authorization", "Bot "+os.Getenv("DISCORD_BOT_SECRET"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error when sending request to the server")
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("Error when sending request to Discord " + req.URL.Path + " " + string(responseBody))
	}
	return responseBody, nil
}

func SetGuildId() error {
	req, _ := http.NewRequest(http.MethodGet, host+"/users/@me/guilds", nil)
	responseBody, err := discordReq(req)
	if err != nil {
		return err
	}
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
	return nil
}

func GetRoles() (types.RolesResp, error) {
	req, _ := http.NewRequest(http.MethodGet, host+"/guilds/"+GuildId+"/roles", nil)
	responseBody, err := discordReq(req)
	if err != nil {
		return nil, err
	}
	roles := types.RolesResp{}
	json.Unmarshal(responseBody, &roles)
	return roles, nil
}

func BanUser(userId string) {

}

func UpdateRoles() error {
	defer func() {
		err := recover()
		if err != nil {
			tgbot.SendString(fmt.Sprint(err))
		}
		return
	}()
	roles, err := GetRoles()
	if err != nil {
		return err
	}
	sheets.UpdateRoles(roles)
	return nil
}

func SetRoles(userId int, roles []types.CurrentRole) error {
	defer func() {
		err := recover()
		if err != nil {
			tgbot.SendString(fmt.Sprint(err))
		}
	}()
	discordUserId := sheets.GetDiscordIdByUserId(userId)

	rolesToSet := []string{}
	i := 0
	set := false
	for i < len(roles) {

		if roles[i].RoleId == "" {
			roles[i] = roles[len(roles)-1]
			roles = roles[:len(roles)-1]
			continue
		}

		if roles[i].Expire.Before(time.Now()) {
			set = true
			courses, _ := sheets.GetRoleCourses(roles[i].RoleId)
			lastExpire := time.Now().Add(-1 * time.Hour)
			for _, courseId := range courses {
				expire, _ := thinkific.GetExpiryDate(userId, courseId)
				if lastExpire.Before(expire) {
					lastExpire = expire
				}
			}

			if time.Now().Before(lastExpire) {
				roles[i].Expire = lastExpire
			} else {
				roles[i] = roles[len(roles)-1]
				roles = roles[:len(roles)-1]
				continue
			}
		}

		rolesToSet = append(rolesToSet, roles[i].RoleId)
		i++
	}

	if len(rolesToSet) == 0 {
		role, err := sheets.GetCourseRole(0)
		if err != nil {
			return err
		}
		rolesToSet = append(rolesToSet, role)
	}

	if discordUserId != "" {
		body := setRolesBody{rolesToSet}
		bodybytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPatch, host+"/guilds/"+GuildId+"/members/"+discordUserId, bytes.NewReader(bodybytes))
		bodyBytes, err := discordReq(req)
		if err != nil {
			return fmt.Errorf("Discord error: %s: %s", err, string(bodyBytes))
		}
	}

	if set {
		err := sheets.SetUserRoles(userId, roles)
		if err != nil {
			return err
		}
	}

	return nil
}

func AdjustRoles() {
	defer func() {
		err := recover()
		if err != nil {
			tgbot.SendString(fmt.Sprint(err))
		}
		return
	}()
	users := sheets.GetUsersRoles()
	i := 0
	for _, v := range users {
		SetRoles(v.Id, v.Roles)
		i++
		if i%20 == 0 {
			time.Sleep(60 * time.Second)
		}
	}
}

func GetInviteLink() string {

	// body := struct {
	// 	MaxAge int `json:"max_age"`
	// }{
	// 	MaxAge: 60,
	// }
	// bodybytes, _ := json.Marshal(body)
	// req, _ := http.NewRequest(http.MethodPost, host+"/channels/"+os.Getenv("DISCORD_INVITE_CHAN_ID")+"/invites", bytes.NewReader(bodybytes))

	// responseBody := discordReq(req)

	// invite := inviteResp{}
	// json.Unmarshal(responseBody, &invite)

	// return invite.Code
	return os.Getenv("DISCORD_INVITE_LINK")
}
