package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/duiniwukenaihe/king-kf/resource"
	"github.com/duiniwukenaihe/king-utils/common/handle"
)

func GetPlugin(c *gin.Context) {
	r := resource.PluginResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleGet(&r)
	c.JSON(responseData.Code, responseData)
}

func ListPlugin(c *gin.Context) {
	r := resource.PluginResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleList(&r)
	c.JSON(responseData.Code, responseData)
}

func CreatePlugin(c *gin.Context) {
	r := resource.PluginResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleCreate(&r, c)
	c.JSON(responseData.Code, responseData)
}

func DeletePlugin(c *gin.Context) {
	r := resource.PluginResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleDelete(&r)
	c.JSON(responseData.Code, responseData)
}

func UpdatePlugin(c *gin.Context) {
	r := resource.PluginResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleUpdate(&r, c)
	c.JSON(responseData.Code, responseData)
}
