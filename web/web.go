package web

import (
	"github.com/gin-gonic/gin"
	"github.com/nange/gospider/web/router"
	"github.com/sirupsen/logrus"
)

func Run() {
	engine := gin.Default()
	router.Route(engine)

	logrus.Fatal(engine.Run())
}