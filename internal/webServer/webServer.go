package webServer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"thinkific-discord/internal/discord"
	"thinkific-discord/internal/discordBot"
	"thinkific-discord/internal/email"
	"thinkific-discord/internal/sheets"
	"thinkific-discord/internal/tgbot"
	"thinkific-discord/internal/types"
	"time"

	"github.com/gin-gonic/gin"
)

func Listen() {

	router := gin.Default()
	router.Use(ErrorHandler)

	router.GET("/", Default)
	router.GET("/discord/auth", discordAuth)
	router.GET("/discord/sendinvite", sendInviteLink)
	router.GET("/sheets/auth", sheetsAuth)
	router.POST("/thinkific/order", newOrder)
	router.POST("/thinkific/course", newCourse)
	router.GET("/thinkific/updateAll", updateAll)
	router.Static("/assets", "./internal/email/resources")

	srv := &http.Server{
		Addr:    os.Getenv("SERVER_ADDRESS"),
		Handler: router,
	}

	fmt.Println(fmt.Sprintf("Server is listening to %s", srv.Addr))

	cert := os.Getenv("SERVER_CERT")
	key := os.Getenv("SERVER_KEY")

	var err error
	if cert == "" && key == "" {
		err = srv.ListenAndServe()
	} else {
		err = srv.ListenAndServeTLS(cert, key)
	}
	// service connections
	if err != nil && err != http.ErrServerClosed {
		fmt.Println(fmt.Sprintf("listen: %s\n", err))
	} else {
		fmt.Println(fmt.Sprintf("Server is listening to %s", srv.Addr))
	}

}

func Listen80() {

	router := gin.Default()

	address := strings.Split(os.Getenv("SERVER_ADDRESS"), ":")[0] + ":80"

	router.Static("/", "./internal/email/resources")

	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	fmt.Println(fmt.Sprintf("Server is listening to %s", srv.Addr))

	var err error
	err = srv.ListenAndServe()
	// service connections
	if err != nil && err != http.ErrServerClosed {
		fmt.Println(fmt.Sprintf("listen: %s\n", err))
	} else {
		fmt.Println(fmt.Sprintf("Server is listening to %s", srv.Addr))
	}

}

func Default(c *gin.Context) {
	c.Writer.WriteHeader(204)
}

func updateAll(c *gin.Context) {
	// defer func() {
	// 	err := recover()
	// 	if err != nil {
	// 		tgbot.SendString(fmt.Sprint(err))
	// 		c.Writer.WriteHeader(500)
	// 		c.Writer.WriteString(fmt.Sprint(err))
	// 	}

	// }()

	updateAllMembers()
}

func sendInviteLink(c *gin.Context) {

	thinkificId := c.Request.URL.Query()["id"][0]
	emailAddress := c.Request.URL.Query()["mail"][0]
	name := c.Request.URL.Query()["name"][0]

	link := discord.GenerateLink(thinkificId)
	err := email.SendInviteLink(emailAddress, link, name)
	if err != nil {
		c.Writer.WriteHeader(500)
		c.Writer.WriteString(err.Error())
		return
	}
	c.Writer.WriteHeader(200)
}

func discordAuth(c *gin.Context) {
	code, _ := c.Request.URL.Query()["code"]
	state, _ := c.Request.URL.Query()["state"]

	id, _ := strconv.Atoi(state[0])
	discord.AddToGroup(code[0], id)

	c.Redirect(307, "https://discord.gg/"+discordBot.GetInviteLink())
}

func sheetsAuth(c *gin.Context) {
	code, _ := c.Request.URL.Query()["code"]

	sheets.CodeChan <- code[0]

	c.Writer.WriteHeader(200)
	c.Writer.WriteString("Cool! You are good to go!")
}

func newOrder(c *gin.Context) {

	defer func() {
		err := recover()
		if err != nil {
			tgbot.SendString(fmt.Sprint(err))
			c.Writer.WriteHeader(500)
			c.Writer.WriteString(fmt.Sprint(err))
		}

	}()
	t1 := time.Now()
	order := types.WebhookOrder{}
	c.ShouldBindJSON(&order)
	handleNewOrder(order, true)
	fmt.Println(time.Now().Sub(t1))
}

func newCourse(c *gin.Context) {
	sheets.UpdateCourses()
	c.Writer.WriteHeader(200)
}

func ErrorHandler(c *gin.Context) {
	//c.Next()

	if len(c.Errors) == 0 {
		return
	}
	msg, _ := json.Marshal(c.Errors)
	c.JSON(http.StatusInternalServerError, msg)
}
