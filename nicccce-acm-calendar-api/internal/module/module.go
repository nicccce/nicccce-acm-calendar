package module

import (
	"github.com/gin-gonic/gin"
	"nicccce-acm-calendar-api/internal/module/crawler"
	"nicccce-acm-calendar-api/internal/module/ping"
)

type Module interface {
	GetName() string
	Init()
	InitRouter(r *gin.RouterGroup)
}

var Modules []Module

func registerModule(m []Module) {
	Modules = append(Modules, m...)
}

func init() {
	// Register your module here
	registerModule([]Module{
		&ping.ModulePing{},
		&crawler.ModuleCrawler{},
	})
}
