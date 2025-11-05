package module

import (
	"github.com/gin-gonic/gin"
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
		&user.ModuleUser{},
		&ping.ModulePing{},
		&crawler.ModuleCrawler{},
	})
}
