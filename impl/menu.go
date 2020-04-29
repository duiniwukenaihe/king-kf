package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-kf/resource"
	"github.com/open-kingfisher/king-utils/common/handle"
)

func GetMenu(c *gin.Context) {
	r := resource.MenuResource{Params: handle.GenerateCommonParams(c, nil)}
	responseData := HandleGet(&r)
	c.JSON(responseData.Code, responseData)
}
