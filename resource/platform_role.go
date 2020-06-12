package resource

import (
	"encoding/json"
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

type PlatformRoleResource struct {
	Params   *handle.Resources
	PostData *common.PlatformRoleDB
}

func (r *PlatformRoleResource) Get() (interface{}, error) {
	platformRole := common.PlatformRoleDB{}
	if err := db.GetById(common.PlatformRoleTable, r.Params.Name, &platformRole); err != nil {
		return nil, err
	}
	return platformRole, nil
}

func (r *PlatformRoleResource) List() (interface{}, error) {
	platformRole := make([]common.PlatformRoleDB, 0)
	if err := db.List(common.DataField, common.PlatformRoleTable, &platformRole, ""); err != nil {
		return nil, err
	}
	return platformRole, nil
}

func (r *PlatformRoleResource) Create(c *gin.Context) (err error) {
	platformRole := common.PlatformRoleDB{}
	if err = c.BindJSON(&platformRole); err != nil {
		return err
	}
	r.PostData = &platformRole
	// 对提交的数据进行校验
	if err = c.ShouldBindWith(r.PostData, binding.Query); err != nil {
		return err
	}
	platformRoleList := make([]*common.PlatformRoleDB, 0)
	if err = db.List(common.DataField, common.PlatformRoleTable, &platformRoleList, "WHERE data-> '$.name'=?", r.PostData.Name); err == nil {
		if len(platformRoleList) > 0 {
			return errors.New("the platform role name already exists")
		}
	} else {
		return
	}
	r.PostData.Id = kit.UUID("r")
	r.PostData.CreateTime = time.Now().Unix()
	r.PostData.ModifyTime = time.Now().Unix()
	if err = db.Insert(common.PlatformRoleTable, r.PostData); err != nil {
		log.Errorf("PlatformRole create error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return err
	}
	auditLog := handle.AuditLog{
		Kind:       common.PlatformRole,
		ActionType: common.Create,
		Resources:  r.Params,
		PostData:   r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *PlatformRoleResource) Delete() (err error) {
	if err = db.Delete(common.PlatformRoleTable, r.Params.Name); err != nil {
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.PlatformRole,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *PlatformRoleResource) Update(c *gin.Context) (err error) {
	platformRole := common.PlatformRoleDB{}
	if err = c.BindJSON(&platformRole); err != nil {
		return err
	}
	r.PostData = &platformRole
	// 对提交的数据进行校验
	if err = c.ShouldBindWith(r.PostData, binding.Query); err != nil {
		return err
	}
	platformRoleList := make([]*common.PlatformRoleDB, 0)
	if err = db.List(common.DataField, common.PlatformRoleTable, &platformRoleList, "WHERE data-> '$.name'=?", r.PostData.Name); err == nil {
		if len(platformRoleList) > 0 {
			for _, v := range platformRoleList {
				if v.Id != r.PostData.Id {
					return errors.New("the platform pole name already exists")
				}
			}
		}
	} else {
		return
	}
	platformRoles := common.PlatformRoleDB{}
	if err = db.GetById(common.PlatformRoleTable, r.PostData.Id, &platformRoles); err != nil {
		log.Errorf("PlatformRole update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	r.PostData.CreateTime = platformRoles.CreateTime
	r.PostData.ModifyTime = time.Now().Unix()
	if err = db.Update(common.PlatformRoleTable, r.PostData.Id, r.PostData); err != nil {
		log.Errorf("PlatformRole update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.PlatformRole,
		ActionType: common.Update,
		Resources:  r.Params,
		PostData:   r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *PlatformRoleResource) ListPlatformRoleTree() (interface{}, error) {
	type role struct {
		Title    string `json:"title"`
		Expand   bool   `json:"expand"`
		Checked  bool   `json:"checked"`
		Level    int32  `json:"level"`
		Name     string `json:"name"`
		Children []role `json:"children"`
	}
	tree := make([]role, 0)
	roles := `[
    {
        "title": "Kingfisher",
        "expand": true,
		"checked": false,
        "level": 0,
		"name": "kingfisher",
        "children": [
            {
                "title": "集群管理",
                "expand": true,
                "level": 1,
				"name": "cluster",
                "children": [
                    {
                        "title": "节点",
                        "expand": true,
                		"level": 2,
						"name": "node",
                		"children": [
							{
		                        "title": "查看节点",
								"name": "show_node",
		                		"level": 3
                			},
                			{
		                        "title": "添加节点",
								"name": "add_node",
		                		"level": 3
                			},
                			{
		                        "title": "编辑节点",
								"name": "edit_node",
		                		"level": 3
                			},
                			{
		                        "title": "删除节点",
								"name": "del_node",
		                		"level": 3
                			},
                			{
		                        "title": "调度节点",
								"name": "schedule_node",
		                		"level": 3
                			},
							{
		                        "title": "删除Pod",
								"name": "del_node_pod",
		                		"level": 3
                			}
                		]
                    }
                ]
            },
			{
                "title": "应用管理",
                "expand": true,
                "level": 1,
				"name": "application",
                "children": [
                    {
                        "title": "应用",
                        "expand": true,
                		"level": 2,
						"name": "apply",
                		"children": [
							{
		                        "title": "查看应用",
								"name": "show_apply",
		                		"level": 3
                			},
                			{
		                        "title": "添加应用",
								"name": "add_apply",
		                		"level": 3
                			},
                			{
		                        "title": "编辑应用",
								"name": "edit_apply",
		                		"level": 3
                			},
                			{
		                        "title": "删除应用",
								"name": "del_apply",
		                		"level": 3
                			}
                		]
                    },
                    {
                        "title": "模板",
                        "expand": true,
                		"level": 2,
						"name": "template",
                		"children": [
							{
		                        "title": "查看模板",
								"name": "show_tem",
		                		"level": 3
                			},
                			{
		                        "title": "添加模板",
								"name": "add_tem",
		                		"level": 3
                			},
                			{
		                        "title": "编辑模板",
								"name": "edit_tem",
		                		"level": 3
                			},
                			{
		                        "title": "删除模板",
								"name": "del_tem",
		                		"level": 3
                			}
                		]
                    }
                ]
            },
            {
                "title": "工作区",
                "expand": true,
                "level": 1,
				"name": "workload",
                "children": [
                    {
                        "title": "部署",
                        "expand": true,
						"name": "deployment",
                		"level": 2,
                		"children": [
							{
		                        "title": "查看部署",
								"name": "show_dep",
		                		"level": 3
                			},
                			{
		                        "title": "添加部署",
		                        "name": "add_dep",
		                		"level": 3
                			},
							{
		                        "title": "高级添加部署",
		                        "name": "a_add_dep",
		                		"level": 3
                			},
                			{
		                        "title": "编辑部署",
								"name": "edit_dep",
		                		"level": 3							
                			},
							{
		                        "title": "高级编辑部署",
		                        "name": "a_edit_dep",
		                		"level": 3
                			},
							{
		                        "title": "伸缩部署",
								"name": "scale_dep",
		                		"level": 3							
                			},
							{
		                        "title": "复制部署",
								"name": "copy_dep",
		                		"level": 3							
                			},
                			{
		                        "title": "删除部署",
		                		"level": 3,
								"name": "del_dep"
                			},
							{
		                        "title": "重启部署",
		                		"level": 3,
								"name": "restart_dep"
                			},
							{
		                        "title": "保存为模板",
		                		"level": 3,
								"name": "save_tem_dep"
                			},
                			{
		                        "title": "删除Pod",
								"level": 3,
								"name": "del_pod_dep"
                			},
                			{
		                        "title": "添加HPA",
		                		"level": 3,
								"name": "add_hpa"
                			},
                			{
		                        "title": "删除HPA",
		                		"level": 3,
								"name": "del_hpa"
                			},
                			{
		                        "title": "编辑HPA",
		                		"level": 3,
								"name": "edit_hpa"
                			},
                			{
		                        "title": "删除副本集",
		                		"level": 3,
								"name": "del_rep"
                			}
                		]
                    },
                    {
                        "title": "Pod",
                        "expand": true,
                		"level": 2,
						"name": "pod",
                		"children": [
							{
		                        "title": "查看Pod",
								"name": "show_pod",
		                		"level": 3
                			},
                			{
		                        "title": "添加Pod",
		                		"level": 3,
								"name": "add_pod"
                			},
                			{
		                        "title": "编辑Pod",
		                		"level": 3,
								"name": "edit_pod"
                			},
                			{
		                        "title": "删除Pod",
		                		"level": 3,
								"name": "del_pod"
                			},
                			{
		                        "title": "Pod日志",
								"name": "pod_log",
		                		"level": 3
                			},
                			{
		                        "title": "Pod终端",
		                		"level": 3,
								"name": "pod_terminal"
                			},
                			{
		                        "title": "Pod调试",
		                		"level": 3,
								"name": "debug_pod"
                			},
							{
		                        "title": "Pod救援",
		                		"level": 3,
								"name": "rescue_pod"
                			}
                		]
                    }
                ]
            },
            {
                "title": "配置",
                "expand": true,
                "level": 1,
				"name": "configuration",
                "children": [
                    {
                        "title": "配置字典",
                        "expand": true,
                		"level": 2,
						"name": "config_map",
                		"children": [
							{
		                        "title": "查看配置字典",
								"name": "show_config_map",
		                		"level": 3
                			},
                			{
		                        "title": "添加配置字典",
		                        "name": "add_config_map",
		                		"level": 3
                			},
                			{
		                        "title": "编辑配置字典",
		                        "name": "edit_config_map",
		                		"level": 3
                			},
                			{
		                        "title": "删除配置字典",
		                        "name": "del_config_map",
		                		"level": 3
                			}
                		]
                    },
                    {
                        "title": "保密字典",
                        "expand": true,
                		"level": 2,
                		"name": "secret_map",
                		"children": [
							{
		                        "title": "查看保密字典",
								"name": "show_secret_map",
		                		"level": 3
                			},
                			{
		                        "title": "添加保密字典",
		                        "name": "add_secret_map",
		                		"level": 3
                			},
                			{
		                        "title": "编辑保密字典",
		                        "name": "edit_secret_map",
		                		"level": 3
                			},
                			{
		                        "title": "删除保密字典",
		                        "name": "del_secret_map",
		                		"level": 3
                			}
                		]
                    }
                ]
            },
            {
                "title": "服务与入口",
                "expand": true,
                "level": 1,
                "name": "discovery",
                "children": [
                    {
                        "title": "服务",
                        "expand": true,
                		"level": 2,
                		"name": "service",
                		"children": [
							{
		                        "title": "查看服务",
								"name": "show_service",
		                		"level": 3
                			},
                			{
		                        "title": "添加服务",
		                        "name": "add_service",
		                		"level": 3
                			},
							{
		                        "title": "高级添加服务",
		                        "name": "a_add_service",
		                		"level": 3
                			},
                			{
		                        "title": "编辑服务",
		                        "name": "edit_service",
		                		"level": 3
                			},
							{
		                        "title": "高级编辑服务",
		                        "name": "a_edit_service",
		                		"level": 3
                			},
                			{
		                        "title": "删除服务",
		                        "name": "del_service",
		                		"level": 3
                			},
							{
		                        "title": "复制服务",
								"name": "copy_service",
		                		"level": 3							
                			},
                			{
		                        "title": "删除Pod",
		                        "name": "del_service_pod",
		                		"level": 3
                			}
                		]
                    },
                    {
                        "title": "入口",
                        "expand": true,
                		"level": 2,
                		"name": "ingress",
                		"children": [
							{
		                        "title": "查看入口",
								"name": "show_ingress",
		                		"level": 3
                			},
                			{
		                        "title": "添加入口",
		                        "name": "add_ingress",
		                		"level": 3
                			},
							{
		                        "title": "高级添加入口",
		                        "name": "a_add_ingress",
		                		"level": 3
                			},
                			{
		                        "title": "编辑入口",
		                        "name": "edit_ingress",
		                		"level": 3
                			},
							{
		                        "title": "高级编辑入口",
		                        "name": "a_edit_ingress",
		                		"level": 3
                			},
                			{
		                        "title": "删除入口",
		                        "name": "del_ingress",
		                		"level": 3
                			}
                		]
                    }
                ]
            },
            {
                "title": "服务治理",
                "expand": true,
                "level": 1,
                "name": "istio",
                "children": [
                    {
                        "title": "虚拟服务",
                        "expand": true,
                		"level": 2,
                		"name": "vs",
                		"children": [
							{
		                        "title": "查看虚拟服务",
								"name": "show_vs",
		                		"level": 3
                			},
                			{
		                        "title": "添加虚拟服务",
		                        "name": "add_vs",
		                		"level": 3
                			},
                			{
		                        "title": "高级添加虚拟服务",
		                        "name": "a_add_vs",
		                		"level": 3
                			},
                			{
		                        "title": "编辑虚拟服务",
		                        "name": "edit_vs",
		                		"level": 3
                			},
                			{
		                        "title": "高级编辑虚拟服务",
		                        "name": "a_edit_vs",
		                		"level": 3
                			},
                			{
		                        "title": "删除虚拟服务",
		                        "name": "del_vs",
		                		"level": 3
                			}
                		]
                    },
                    {
                        "title": "网关",
                        "expand": true,
                		"level": 2,
                		"name": "gateway",
                		"children": [
							{
		                        "title": "查看网关",
								"name": "show_gateway",
		                		"level": 3
                			},
                			{
		                        "title": "添加网关",
		                        "name": "add_gateway",
		                		"level": 3
                			},
							{
		                        "title": "高级添加网关",
		                        "name": "a_add_gateway",
		                		"level": 3
                			},
                			{
		                        "title": "编辑网关",
		                        "name": "edit_gateway",
		                		"level": 3
                			},
							{
		                        "title": "高级编辑网关",
		                        "name": "a_edit_gateway",
		                		"level": 3
                			},
                			{
		                        "title": "删除网关",
		                        "name": "del_gateway",
		                		"level": 3
                			}
                		]
                    },
                    {
                        "title": "目标规则",
                        "expand": true,
                		"level": 2,
                		"name": "dr",
                		"children": [
							{
		                        "title": "查看目标规则",
								"name": "show_dr",
		                		"level": 3
                			},
                			{
		                        "title": "添加目标规则",
		                        "name": "add_dr",
		                		"level": 3
                			},
                			{
		                        "title": "高级添加目标规则",
		                        "name": "a_add_dr",
		                		"level": 3
                			},
                			{
		                        "title": "编辑目标规则",
		                        "name": "edit_dr",
		                		"level": 3
                			},
                			{
		                        "title": "高级编辑目标规则",
		                        "name": "a_edit_dr",
		                		"level": 3
                			},
                			{
		                        "title": "删除目标规则",
		                        "name": "del_dr",
		                		"level": 3
                			}
                		]
                    },
					{
                        "title": "服务入口",
                        "expand": true,
                		"level": 2,
                		"name": "se",
                		"children": [
							{
		                        "title": "查看服务入口",
								"name": "show_se",
		                		"level": 3
                			},
                			{
		                        "title": "添加服务入口",
		                        "name": "add_se",
		                		"level": 3
                			},
							{
		                        "title": "高级添加服务入口",
		                        "name": "a_add_se",
		                		"level": 3
                			},
                			{
		                        "title": "编辑服务入口",
		                        "name": "edit_se",
		                		"level": 3
                			},
							{
		                        "title": "高级编辑服务入口",
		                        "name": "a_edit_se",
		                		"level": 3
                			},
                			{
		                        "title": "删除服务入口",
		                        "name": "del_se",
		                		"level": 3
                			}
                		]
                    }
                ]
            },
            {
                "title": "监控告警",
                "expand": true,
                "level": 1,
                "name": "alert",
                "children": [
                    {
                        "title": "规则管理",
                        "expand": true,
                		"level": 2,
                		"name": "rule",
                		"children": [
							{
		                        "title": "查看规则",
								"name": "show_rule",
		                		"level": 3
                			},
                			{
		                        "title": "添加规则",
		                        "name": "add_rule",
		                		"level": 3
                			},
                			{
		                        "title": "编辑规则",
		                        "name": "edit_rule",
		                		"level": 3
                			},
                			{
		                        "title": "删除规则",
		                        "name": "del_rule",
		                		"level": 3
                			}
                		]
                    },
                    {
                        "title": "Scrape",
                        "expand": true,
                		"level": 2,
                		"name": "scrape",
                		"children": [
							{
		                        "title": "查看Scrape",
								"name": "show_scrape",
		                		"level": 3
                			},
                			{
		                        "title": "添加Scrape",
		                        "name": "add_scrape",
		                		"level": 3
                			},
                			{
		                        "title": "编辑Scrape",
		                        "name": "edit_scrape",
		                		"level": 3
                			},
                			{
		                        "title": "删除Scrape",
		                        "name": "del_scrape",
		                		"level": 3
                			}
                		]
                    }
                ]
            },
            {
                "title": "集群插件",
                "expand": true,
                "level": 1,
                "name": "cluster_plugin",
                "children": [
					{
						"title": "查看插件",
						"name": "show_cluster_plugin",
						"level": 3
					},
                    {
                        "title": "安装插件",
                        "expand": true,
                        "name": "install_cluster_plugin",
                		"level": 3
                    },
                    {
                        "title": "卸载插件",
                        "expand": true,
                        "name": "uninstall_cluster_plugin",
                		"level": 3
                    },
					{
                        "title": "kubectl",
                        "expand": true,
                        "name": "kubectl",
                		"level": 3
                    }
                ]
            },
            {
                "title": "系统管理",
                "expand": true,
                "level": 1,
                "name": "manage",
                "children": [
                    {
                        "title": "集群管理",
                        "expand": true,
                		"level": 2,
                		"name": "cluster_manage",
                		"children": [
							{
		                        "title": "查看集群",
								"name": "show_cluster_manage",
		                		"level": 3
                			},
                			{
		                        "title": "添加集群",
		                        "name": "add_cluster_manage",
		                		"level": 3
                			},
                			{
		                        "title": "编辑集群",
		                        "name": "edit_cluster_manage",
		                		"level": 3
                			},
                			{
		                        "title": "删除集群",
		                        "name": "del_cluster_manage",
		                		"level": 3
                			}
                		]
                    },
                    {
                        "title": "命名空间",
                        "expand": true,
                		"level": 2,
                		"name": "namespace",
                		"children": [
							{
		                        "title": "查看命名空间",
								"name": "show_namespace",
		                		"level": 3
                			},
                			{
		                        "title": "添加命名空间",
		                        "name": "add_namespace",
		                		"level": 3
                			},
							{
		                        "title": "编辑命名空间",
		                        "name": "edit_namespace",
		                		"level": 3
                			},
                			{
		                        "title": "删除命名空间",
		                        "name": "del_namespace",
		                		"level": 3
                			}
                		]
                    },
                	{
                        "title": "产品线",
                        "expand": true,
                		"level": 2,
                		"name": "product",
                		"children": [
							{
		                        "title": "查看产品线",
								"name": "show_product",
		                		"level": 3
                			},
                			{
		                        "title": "添加产品线",
		                        "name": "add_product",
		                		"level": 3
                			},
                			{
		                        "title": "编辑产品线",
		                        "name": "edit_product",
		                		"level": 3
                			},
                			{
		                        "title": "删除产品线",
		                        "name": "del_product",
		                		"level": 3
                			}
                		]
                    },
                    {
                        "title": "用户管理",
                        "expand": true,
                		"level": 2,
                		"name": "user",
                		"children": [
							{
		                        "title": "查看用户",
								"name": "show_user",
		                		"level": 3
                			},
                			{
		                        "title": "添加用户",
		                        "name": "add_user",
		                		"level": 3
                			},
                			{
		                        "title": "编辑用户",
		                        "name": "edit_user",
		                		"level": 3
                			},
                			{
		                        "title": "删除用户",
		                        "name": "del_user",
		                		"level": 3
                			},
                			{
		                        "title": "修改密码",
		                        "name": "change_pw",
		                		"level": 3
                			}
                		]
                    },
                    {
                        "title": "平台角色",
                        "expand": true,
                		"level": 2,
                		"name": "platform_role",
                		"children": [
							{
		                        "title": "查看平台角色",
								"name": "show_platform_role",
		                		"level": 3
                			},
                			{
		                        "title": "添加角色",
		                        "name": "add_platform_role",
		                		"level": 3
                			},
                			{
		                        "title": "编辑角色",
		                        "name": "edit_platform_role",
		                		"level": 3
                			},
                			{
		                        "title": "删除角色",
		                        "name": "del_platform_role",
		                		"level": 3
                			}
                		]
                    },
                    {
                        "title": "插件管理",
                        "expand": true,
                		"level": 2,
                		"name": "plugin",
                		"children": [
							{
		                        "title": "查看插件",
								"name": "show_plugin",
		                		"level": 3
                			},
                			{
		                        "title": "添加插件",
		                        "name": "add_plugin",
		                		"level": 3
                			},
                			{
		                        "title": "编辑插件",
		                        "name": "edit_plugin",
		                		"level": 3
                			},
                			{
		                        "title": "删除插件",
		                        "name": "del_plugin",
		                		"level": 3
                			}
                		]
                    },
                    {
                        "title": "健康检查",
                        "expand": true,
                		"level": 2,
                		"name": "inspect",
                		"children": [
							{
		                        "title": "查看健康检查",
								"name": "show_inspect",
		                		"level": 3
                			},
                			{
		                        "title": "添加健康检查",
		                        "name": "add_inspect",
		                		"level": 3
                			},
                			{
		                        "title": "编辑健康检查",
		                        "name": "edit_inspect",
		                		"level": 3
                			},
                			{
		                        "title": "删除健康检查",
		                        "name": "del_inspect",
		                		"level": 3
                			},
                			{
		                        "title": "执行健康检查",
		                        "name": "action_inspect",
		                		"level": 3
                			}
                		]
                    },
					{
                        "title": "配置管理",
                        "expand": true,
                        "name": "config",
                		"level": 3
                    },
                    {
                        "title": "审计日志",
                        "expand": true,
                        "name": "audit",
                		"level": 3
                    }
                ]
            }
        ]
    }
]`
	if err := json.Unmarshal([]byte(roles), &tree); err != nil {
		return nil, err
	}
	return tree, nil
}
