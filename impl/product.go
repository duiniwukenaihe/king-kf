package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/duiniwukenaihe/king-kf/resource"
	"github.com/duiniwukenaihe/king-utils/common/handle"
)

func GetProduct(c *gin.Context) {
	r := resource.ProductResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleGet(&r)
	c.JSON(responseData.Code, responseData)
}

func ListProduct(c *gin.Context) {
	r := resource.ProductResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleList(&r)
	c.JSON(responseData.Code, responseData)
}

func CreateProduct(c *gin.Context) {
	r := resource.ProductResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleCreate(&r, c)
	c.JSON(responseData.Code, responseData)
}

func DeleteProduct(c *gin.Context) {
	r := resource.ProductResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleDelete(&r)
	c.JSON(responseData.Code, responseData)
}

func UpdateProduct(c *gin.Context) {
	r := resource.ProductResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleUpdate(&r, c)
	c.JSON(responseData.Code, responseData)
}

func CascadeCluster(c *gin.Context) {
	r := resource.ProductResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := handle.HandlerResponse(r.CascadeCluster())
	c.JSON(responseData.Code, responseData)
}

func CascadeAll(c *gin.Context) {
	r := resource.ProductResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := handle.HandlerResponse(r.CascadeAll())
	c.JSON(responseData.Code, responseData)
}

func TreeClusterNamespace(c *gin.Context) {
	r := resource.ProductResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := handle.HandlerResponse(r.TreeClusterNamespace())
	c.JSON(responseData.Code, responseData)
}
