package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-kf/resource"
	"github.com/open-kingfisher/king-utils/common/handle"
)

func GetConfig(c *gin.Context) {
	r := resource.ConfigResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleGet(&r)
	c.JSON(responseData.Code, responseData)
}

func ListConfig(c *gin.Context) {
	r := resource.ConfigResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleList(&r)
	c.JSON(responseData.Code, responseData)
}

func CreateConfig(c *gin.Context) {
	r := resource.ConfigResource{Params: handle.GenerateCommonParams(c, nil)}
	r.ConfigType = c.Query("configType")
	responseData := HandleCreate(&r, c)
	c.JSON(responseData.Code, responseData)
}

func DeleteConfig(c *gin.Context) {
	r := resource.ConfigResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleDelete(&r)
	c.JSON(responseData.Code, responseData)
}

func UpdateConfig(c *gin.Context) {
	r := resource.ConfigResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleUpdate(&r, c)
	c.JSON(responseData.Code, responseData)
}
