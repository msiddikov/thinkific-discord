package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lavina-Tech-LLC/lavinagopackage/llog"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {

}

func startApp() {

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	//router.LoadHTMLGlob("lvn-tools/internal/templates/*")

	// Define the route for the index page and display the index.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	router.GET("/hello", Default)

	// Start serving the application
	//router.Run()

	srv := &http.Server{
		Addr:    "192.168.1.133:8085",
		Handler: router,
	}

	llog.Info(fmt.Sprintf("Server is listening to %s", srv.Addr))

	go func() {
		var err error
		err = srv.ListenAndServe()
		// service connections
		if err != nil && err != http.ErrServerClosed {
			llog.Error(fmt.Sprintf("listen: %s\n", err))
		} else {
			llog.Info(fmt.Sprintf("Server is listening to %s", srv.Addr))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	llog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		llog.Error(fmt.Sprintf("Server Shutdown: %s", err))
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		llog.Info("timeout of 5 seconds.")
	}
	llog.Info("Server exiting")
}

func Default(c *gin.Context) {

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title": "Home Page",
		},
	)
}
