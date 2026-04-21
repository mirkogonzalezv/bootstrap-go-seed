package version

import "github.com/gin-gonic/gin"

type ModuleInfo struct {
	Name         string
	RegisterFunc func(*gin.RouterGroup, any)
	Controller   any
}

var Registry = make(map[int][]ModuleInfo)

func RegisterModule(version int, moduleName string, registerFunc func(*gin.RouterGroup, any), controller any) {
	info := ModuleInfo{
		Name:         moduleName,
		RegisterFunc: registerFunc,
		Controller:   controller,
	}
	Registry[version] = append(Registry[version], info)
}

func GetRegistry() map[int][]ModuleInfo {
	return Registry
}
