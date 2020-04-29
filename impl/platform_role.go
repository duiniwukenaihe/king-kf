package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-kf/resource"
	"github.com/open-kingfisher/king-utils/common/handle"
)

func GetPlatformRole(c *gin.Context) {
	r := resource.PlatformRoleResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleGet(&r)
	c.JSON(responseData.Code, responseData)
}

func ListPlatformRole(c *gin.Context) {
	r := resource.PlatformRoleResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleList(&r)
	c.JSON(responseData.Code, responseData)
}

func CreatePlatformRole(c *gin.Context) {
	r := resource.PlatformRoleResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleCreate(&r, c)
	c.JSON(responseData.Code, responseData)
}

func DeletePlatformRole(c *gin.Context) {
	r := resource.PlatformRoleResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleDelete(&r)
	c.JSON(responseData.Code, responseData)
}

func UpdatePlatformRole(c *gin.Context) {
	r := resource.PlatformRoleResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleUpdate(&r, c)
	c.JSON(responseData.Code, responseData)
}

func ListPlatformRoleTree(c *gin.Context) {
	r := resource.PlatformRoleResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := handle.HandlerResponse(r.ListPlatformRoleTree())
	c.JSON(responseData.Code, responseData)
}
