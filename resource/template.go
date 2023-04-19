package resource

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/duiniwukenaihe/king-utils/common"
	"github.com/duiniwukenaihe/king-utils/common/handle"
	"github.com/duiniwukenaihe/king-utils/common/log"
	"github.com/duiniwukenaihe/king-utils/db"
	"github.com/duiniwukenaihe/king-utils/kit"
	"time"
)

type TemplateResource struct {
	Params   *handle.Resources
	PostData *common.TemplateDB
}

func (r *TemplateResource) Get() (interface{}, error) {
	template := common.TemplateDB{}
	if err := db.GetById(common.TemplateTable, r.Params.Name, &template); err != nil {
		return nil, err
	}
	return template, nil
}

func (r *TemplateResource) List() (interface{}, error) {
	template := make([]common.TemplateDB, 0)
	if err := db.List(common.DataField, common.TemplateTable, &template, ""); err != nil {
		return nil, err
	}
	return template, nil
}

func (r *TemplateResource) Create(c *gin.Context) (err error) {
	template := common.TemplateDB{}
	if err = c.BindJSON(&template); err != nil {
		return err
	}
	r.PostData = &template
	// 对提交的数据进行校验
	if err = c.ShouldBindWith(r.PostData, binding.Query); err != nil {
		return err
	}
	TemplateList := make([]*common.TemplateDB, 0)
	if err = db.List(common.DataField, common.TemplateTable, &TemplateList, "WHERE data-> '$.name'=?", r.PostData.Name); err == nil {
		if len(TemplateList) > 0 {
			return errors.New("the Template name already exists")
		}
	} else {
		return
	}
	r.PostData.Id = kit.UUID("p")
	r.PostData.CreateTime = time.Now().Unix()
	r.PostData.ModifyTime = time.Now().Unix()
	if err = db.Insert(common.TemplateTable, r.PostData); err != nil {
		log.Errorf("Template create error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return err
	}
	auditLog := handle.AuditLog{
		Kind:       common.Template,
		ActionType: common.Create,
		Resources:  r.Params,
		PostData:   r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *TemplateResource) Delete() (err error) {
	if err = db.Delete(common.TemplateTable, r.Params.Name); err != nil {
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Template,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *TemplateResource) Update(c *gin.Context) (err error) {
	template := common.TemplateDB{}
	if err = c.BindJSON(&template); err != nil {
		return err
	}
	r.PostData = &template
	// 对提交的数据进行校验
	if err = c.ShouldBindWith(r.PostData, binding.Query); err != nil {
		return err
	}
	TemplateList := make([]*common.TemplateDB, 0)
	if err = db.List(common.DataField, common.TemplateTable, &TemplateList, "WHERE data-> '$.name'=?", r.PostData.Name); err == nil {
		if len(TemplateList) > 0 {
			for _, v := range TemplateList {
				if v.Id != r.PostData.Id {
					return errors.New("the Template name already exists")
				}
			}
		}
	} else {
		return
	}
	Templates := common.TemplateDB{}
	if err = db.GetById(common.TemplateTable, r.PostData.Id, &Templates); err != nil {
		log.Errorf("Template update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	r.PostData.CreateTime = Templates.CreateTime
	r.PostData.ModifyTime = time.Now().Unix()
	if err = db.Update(common.TemplateTable, r.PostData.Id, r.PostData); err != nil {
		log.Errorf("Template update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Template,
		ActionType: common.Update,
		Resources:  r.Params,
		PostData:   r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}
