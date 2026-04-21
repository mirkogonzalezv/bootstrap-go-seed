package version

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func BuildRoutes(router *gin.Engine, basePrefix string) {
	registry := GetRegistry()
	for versionNum, modules := range registry {
		versionGroup := router.Group(fmt.Sprintf("/%s/v%d", basePrefix, versionNum))
		for _, module := range modules {
			moduleGroup := versionGroup.Group("/" + module.Name)
			module.RegisterFunc(moduleGroup, module.Controller)
		}
	}
}
