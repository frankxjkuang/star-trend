package main

import (
	"context"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"star-trend/config"
	"star-trend/controller"
	"star-trend/github"
	"time"
)

func main() {
	log.Infof("Cfg: %+v, Gh: %+v", config.Cfg, github.Gh)
	router := gin.Default()
	router.LoadHTMLFiles("./static/index.html")

	router.GET("/", controller.RenderIndex)
	router.GET("/:owner/:repo", controller.RenderGraph)
	router.POST("/", controller.RenderRepoDetail)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Cfg.Port),
		Handler: router,
	}
	e := log.WithField("listen port", config.Cfg.Port)
	e.Info("starting up...")
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.WithError(err).Error("Failed to start up server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown err: %v", err)
	}
	log.Info("Server exiting")
}
