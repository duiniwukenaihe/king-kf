package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/duiniwukenaihe/king-kf/resource"
	"github.com/duiniwukenaihe/king-utils/common/handle"
)

func Install(c *gin.Context) {
	r := resource.ClusterPluginResource{Params: handle.GenerateCommonParams(c, nil)}
	r.Plugin = c.Query("plugin")
	responseData := HandleInstall(&r)
	c.JSON(responseData.Code, responseData)
}

func UnInstall(c *gin.Context) {
	r := resource.ClusterPluginResource{Params: handle.GenerateCommonParams(c, nil)}
	r.Plugin = c.Query("plugin")
	responseData := HandleUnInstall(&r)
	c.JSON(responseData.Code, responseData)
}

func ListClusterPlugin(c *gin.Context) {
	r := resource.ClusterPluginResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleList(&r)
	c.JSON(responseData.Code, responseData)
}
