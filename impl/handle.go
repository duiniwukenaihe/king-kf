package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/handle"
	"github.com/open-kingfisher/king-utils/common/log"
)

type HandleGetInterface interface {
	Get() (interface{}, error)
}

type HandleListInterface interface {
	List() (interface{}, error)
}

type HandleDeleteInterface interface {
	Delete() error
}

type HandleUpdateInterface interface {
	Update(c *gin.Context) error
}

type HandleCreateInterface interface {
	Create(c *gin.Context) error
}

type HandleInstallInterface interface {
	Install() error
}

type HandleUnInstallInterface interface {
	UnInstall() error
}

type HandleLDAPTestInterface interface {
	LDAPTest(c *gin.Context) error
}

func HandleGet(r HandleGetInterface) *common.ResponseData {
	responseData := handle.HandlerResponse(r.Get())
	return responseData
}

func HandleList(r HandleListInterface) *common.ResponseData {
	responseData := handle.HandlerResponse(r.List())
	return responseData
}

func HandleDelete(r HandleDeleteInterface) *common.ResponseData {
	responseData := handle.HandlerResponse(nil, r.Delete())
	log.Error(responseData)
	return responseData
}

func HandleUpdate(r HandleUpdateInterface, c *gin.Context) *common.ResponseData {
	responseData := handle.HandlerResponse(nil, r.Update(c))
	return responseData
}

func HandleCreate(r HandleCreateInterface, c *gin.Context) *common.ResponseData {
	responseData := handle.HandlerResponse(nil, r.Create(c))
	return responseData
}

func HandleInstall(r HandleInstallInterface) *common.ResponseData {
	responseData := handle.HandlerResponse(nil, r.Install())
	return responseData
}

func HandleUnInstall(r HandleUnInstallInterface) *common.ResponseData {
	responseData := handle.HandlerResponse(nil, r.UnInstall())
	return responseData
}

func HandleLDAPTest(r HandleLDAPTestInterface, c *gin.Context) *common.ResponseData {
	responseData := handle.HandlerResponse(nil, r.LDAPTest(c))
	return responseData
}
