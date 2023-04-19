package resource

import (
	"github.com/gin-gonic/gin"
	"github.com/duiniwukenaihe/king-utils/common"
	"github.com/duiniwukenaihe/king-utils/common/handle"
	"github.com/duiniwukenaihe/king-utils/common/log"
	"github.com/duiniwukenaihe/king-utils/db"
	"github.com/duiniwukenaihe/king-utils/kit"
	"time"
)

type ConfigResource struct {
	Params     *handle.Resources
	ConfigType string
	PostData   *common.ConfigDB
}

func (r *ConfigResource) Get() (interface{}, error) {
	Config := common.ConfigDB{}
	if err := db.GetById(common.ConfigTable, common.ConfigID, &Config); err != nil {
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
		r.PostData.Id = common.ConfigID
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
	configs := common.ConfigDB{}
	if err = db.GetById(common.ConfigTable, common.ConfigID, &configs); err != nil {
		log.Errorf("Config update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	configs.ModifyTime = time.Now().Unix()
	switch r.ConfigType {
	case "ldap":
		configs.LDAPDB = r.PostData.LDAPDB
	}
	if err = db.Update(common.ConfigTable, common.ConfigID, configs); err != nil {
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

func (r *ConfigResource) LDAPTest(c *gin.Context) (err error) {
	config := common.ConfigDB{}
	if err = c.BindJSON(&config); err != nil {
		return err
	}
	client := *kit.LdapLookup(config.LDAPDB.URL,
		config.LDAPDB.SearchDN,
		config.LDAPDB.SearchPassword,
		config.LDAPDB.BaseDN,
		config.LDAPDB.UserFilter,
		config.LDAPDB.TLS)
	err = client.Connect()
	if err != nil {
		return err
	}
	return
}
