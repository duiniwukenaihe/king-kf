package resource

import (
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/handle"
	"github.com/open-kingfisher/king-utils/common/log"
	"github.com/open-kingfisher/king-utils/db"
	"time"
)

const ID = "c_000000000001"

type ConfigResource struct {
	Params     *handle.Resources
	ConfigType string
	PostData   *common.ConfigDB
}

func (r *ConfigResource) Get() (interface{}, error) {
	Config := common.ConfigDB{}
	if err := db.GetById(common.ConfigTable, ID, &Config); err != nil {
		return nil, err
	}
	return Config, nil
}

func (r *ConfigResource) List() (interface{}, error) {
	Config := make([]common.ConfigDB, 0)
	if err := db.List(common.DataField, common.ConfigTable, &Config, ""); err != nil {
		return nil, err
	}
	return Config, nil
}

func (r *ConfigResource) Create(c *gin.Context) (err error) {
	Config := common.ConfigDB{}
	if err = c.BindJSON(&Config); err != nil {
		return err
	}
	r.PostData = &Config
	configList := make([]common.ConfigDB, 0)
	if err := db.List(common.DataField, common.ConfigTable, &configList, ""); err != nil {
		return err
	}
	// 如果存在说明是更新
	if len(configList) > 0 {
		return r.Update(c)
	} else {
		r.PostData.Id = ID
		r.PostData.CreateTime = time.Now().Unix()
		r.PostData.ModifyTime = time.Now().Unix()
		if err = db.Insert(common.ConfigTable, r.PostData); err != nil {
			log.Errorf("Config create error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
			return err
		}
		auditLog := handle.AuditLog{
			Kind:       common.Config,
			ActionType: common.Create,
			Resources:  r.Params,
			PostData:   r.PostData,
		}
		if err = auditLog.InsertAuditLog(); err != nil {
			return
		}
	}
	return
}

func (r *ConfigResource) Delete() (err error) {
	if err = db.Delete(common.ConfigTable, r.Params.Name); err != nil {
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Config,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *ConfigResource) Update(c *gin.Context) (err error) {
	Configs := common.ConfigDB{}
	if err = db.GetById(common.ConfigTable, ID, &Configs); err != nil {
		log.Errorf("Config update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	Configs.ModifyTime = time.Now().Unix()
	switch r.ConfigType {
	case "ldap":
		Configs.LDAPDB = r.PostData.LDAPDB
	}
	if err = db.Update(common.ConfigTable, ID, Configs); err != nil {
		log.Errorf("Config update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Config,
		ActionType: common.Update,
		Resources:  r.Params,
		PostData:   r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}
