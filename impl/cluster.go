package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-kf/resource"
	"github.com/open-kingfisher/king-utils/common/handle"
)

func GetCluster(c *gin.Context) {
	r := resource.ClusterResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleGet(&r)
	c.JSON(responseData.Code, responseData)
}

func ListCluster(c *gin.Context) {
	r := resource.ClusterResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleList(&r)
	c.JSON(responseData.Code, responseData)
}

func DeleteCluster(c *gin.Context) {
	r := resource.ClusterResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleDelete(&r)
	c.JSON(responseData.Code, responseData)
}

func UpdateCluster(c *gin.Context) {
	r := resource.ClusterResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleUpdate(&r, c)
	c.JSON(responseData.Code, responseData)
}

func CreateCluster(c *gin.Context) {
	r := resource.ClusterResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleCreate(&r, c)
	c.JSON(responseData.Code, responseData)
}
