package version

import "github.com/gin-gonic/gin"

func ConfigControllerVersion(moduleName string, controller Versioned, registerFunc func(*gin.RouterGroup, any)) {
	version := controller.GetVersion()

	RegisterModule(version, moduleName, registerFunc, controller)
}
