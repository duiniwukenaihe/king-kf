package resource

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/handle"
	"github.com/open-kingfisher/king-utils/common/log"
	"github.com/open-kingfisher/king-utils/db"
	"github.com/open-kingfisher/king-utils/kit"
	"time"
)

type PluginResource struct {
	Params   *handle.Resources
	PostData *common.PluginDB
}

func (r *PluginResource) Get() (interface{}, error) {
	plugin := common.PluginDB{}
	if err := db.GetById(common.PluginTable, r.Params.Name, &plugin); err != nil {
		return nil, err
	}
	return plugin, nil
}

func (r *PluginResource) List() (interface{}, error) {
	plugin := make([]common.PluginDB, 0)
	if err := db.List(common.DataField, common.PluginTable, &plugin, ""); err != nil {
		return nil, err
	}
	return plugin, nil
}

func (r *PluginResource) Create(c *gin.Context) (err error) {
	plugin := common.PluginDB{}
	if err = c.BindJSON(&plugin); err != nil {
		return err
	}
	r.PostData = &plugin
	// 对提交的数据进行校验
	if err = c.ShouldBindWith(r.PostData, binding.Query); err != nil {
		return err
	}
	pluginList := make([]*common.PluginDB, 0)
	if err = db.List(common.DataField, common.PluginTable, &pluginList, "WHERE data-> '$.name'=?", r.PostData.Name); err == nil {
		if len(pluginList) > 0 {
			return errors.New("the plugin name already exists")
		}
	} else {
		return
	}
	r.PostData.Id = kit.UUID("p")
	r.PostData.CreateTime = time.Now().Unix()
	r.PostData.ModifyTime = time.Now().Unix()
	if err = db.Insert(common.PluginTable, r.PostData); err != nil {
		log.Errorf("Plugin create error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return err
	}
	auditLog := handle.AuditLog{
		Kind:       common.Plugin,
		ActionType: common.Create,
		Resources:  r.Params,
		PostData:   r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *PluginResource) Delete() (err error) {
	if err = db.Delete(common.PluginTable, r.Params.Name); err != nil {
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Plugin,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *PluginResource) Update(c *gin.Context) (err error) {
	plugin := common.PluginDB{}
	if err = c.BindJSON(&plugin); err != nil {
		return err
	}
	r.PostData = &plugin
	// 对提交的数据进行校验
	if err = c.ShouldBindWith(r.PostData, binding.Query); err != nil {
		return err
	}
	pluginList := make([]*common.PluginDB, 0)
	if err = db.List(common.DataField, common.PluginTable, &pluginList, "WHERE data-> '$.name'=?", r.PostData.Name); err == nil {
		if len(pluginList) > 0 {
			for _, v := range pluginList {
				if v.Id != r.PostData.Id {
					return errors.New("the plugin name already exists")
				}
			}
		}
	} else {
		return
	}
	plugins := common.PluginDB{}
	if err = db.GetById(common.PluginTable, r.PostData.Id, &plugins); err != nil {
		log.Errorf("Plugin update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	r.PostData.CreateTime = plugins.CreateTime
	r.PostData.ModifyTime = time.Now().Unix()
	if err = db.Update(common.PluginTable, r.PostData.Id, r.PostData); err != nil {
		log.Errorf("Plugin update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Plugin,
		ActionType: common.Update,
		Resources:  r.Params,
		PostData:   r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}
