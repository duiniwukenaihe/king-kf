package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/log"
	"github.com/open-kingfisher/king-utils/db"
	"net/http"
)

func ListAuditLog(c *gin.Context) {
	responseData := common.ResponseData{}
	audit := make([]*common.AuditLog, 0)
	err := db.List(common.DataField, common.AuditLogTable, &audit, "order by data -> '$.action_time' desc")
	if err == nil {
		responseData.Msg = ""
		responseData.Data = audit
		responseData.Code = http.StatusOK
	} else {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		log.Errorf("Audit log list: %s", err)
	}
	c.JSON(responseData.Code, responseData)
}
