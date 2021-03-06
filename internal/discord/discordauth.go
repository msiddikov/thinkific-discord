package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"thinkific-discord/internal/discordBot"
	"thinkific-discord/internal/sheets"

	"golang.org/x/oauth2"
)

var (
	Conf   *oauth2.Config
	Client *http.Client
)

type (
	userResp struct {
		Id            string
		Username      string
		Email         string
		Discriminator string
	}
	bodyStruct struct {
		Access_token string `json:"access_token"`
	}
)

func GenerateLink(userId string) string {

	Conf = &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		Scopes:       []string{"guilds.join", "identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		RedirectURL: os.Getenv("SERVER_DOMAIN") + "/discord/auth",
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := Conf.AuthCodeURL(userId, oauth2.AccessTypeOffline)
	return url
}

func AddToGroup(code string, thinkificId int) error {
	ctx := context.Background()
	tok, err := Conf.Exchange(ctx, code)
	if err != nil {
		return err
	}

	Client = Conf.Client(ctx, tok)
	resp, err := Client.Get("https://discord.com/api/users/@me")

	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%s: %s", err, string(bodyBytes))
		}

		user := userResp{}
		json.Unmarshal(bodyBytes, &user)

		err = sheets.SetDiscordIdByUserId(thinkificId, user.Id, user.Username+"#"+user.Discriminator)
		if err != nil {
			return err
		}

		clientBot := &http.Client{}
		bodyObj := bodyStruct{tok.AccessToken}
		bodyBytes1, _ := json.Marshal(bodyObj)
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("https://discord.com/api/guilds/%s/members/%s", discordBot.GuildId, user.Id), bytes.NewReader(bodyBytes1))
		req.Header.Add("Authorization", "Bot "+os.Getenv("DISCORD_BOT_SECRET"))
		req.Header.Add("Content-Type", "application/json")
		resp, err = clientBot.Do(req)
		bodyBytes, _ = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = discordBot.SetRoles(thinkificId, sheets.GetUserRoles(thinkificId))
		if err != nil {
			return err
		}
	}
	return nil
}
