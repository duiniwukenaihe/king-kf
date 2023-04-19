package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/duiniwukenaihe/king-kf/resource"
	"github.com/duiniwukenaihe/king-utils/common/handle"
)

func GetMenu(c *gin.Context) {
	r := resource.MenuResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleGet(&r)
	c.JSON(responseData.Code, responseData)
}
