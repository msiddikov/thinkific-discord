package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"thinkific-discord/internal/email"
	"thinkific-discord/internal/types"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	svc           *sheets.Service
	spreadsheetId = ""
	CodeChan      = make(chan string)
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	email.SendSheetsConsent(os.Getenv("ADMIN_EMAIL"), authURL, "Bot Administrator")
	fmt.Println("Need to authorize this bot to your google drive account. Please check your email")
	//fmt.Println(authURL)

	authCode := <-CodeChan

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func InitService() {
	godotenv.Load(".env")
	ctx := context.Background()
	spreadsheetId = os.Getenv("SHEETS_ID")

	config := oauth2.Config{
		ClientID:     os.Getenv("SHEETS_API_ID"),
		ClientSecret: os.Getenv("SHEETS_API_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		RedirectURL: os.Getenv("SERVER_DOMAIN") + "/sheets/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}
	client := getClient(&config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	svc = srv
}

func AddCourseToUser(userId, courseId int, expire time.Time) ([]types.CurrentRole, error) {
	roleId := GetCourseRole(courseId)
	if roleId == "" {
		return nil, fmt.Errorf("This course is not set")
	}
	if expire.Before(time.Now()) {
		expire = time.Now().AddDate(10, 0, 0)
	}

	currentRoles := GetUserRoles(userId)

	found := false
	for k, v := range currentRoles {
		if v.RoleId == roleId {
			currentRoles[k].Expire = expire
			found = true
		}
	}
	if !found {
		currentRoles = append(currentRoles, types.CurrentRole{
			RoleId: roleId,
			Expire: expire,
		})
	}
	SetUserRoles(userId, currentRoles)

	return currentRoles, nil

}
