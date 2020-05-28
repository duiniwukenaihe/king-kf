package resource

import (
	"encoding/json"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/handle"
	"github.com/open-kingfisher/king-utils/db"
)

const menuData = `[
  {
    "path": "/login",
    "name": "login",
    "meta": {
      "title": "Login - 登录",
      "hideInMenu": true
    },
    "component": "login"
  },
  {
    "path": "/",
    "name": "_home",
    "redirect": "/home",
    "component": "main",
    "meta": {
      "hideInBread": true,
      "notCache": true
    },
    "children": [
      {
        "path": "/home",
        "name": "home",
        "meta": {
          "title": "首页",
          "notCache": true,
          "icon": "md-home"
        },
        "component": "home"
      }
    ]
  },
  {
    "path": "/cluster",
    "name": "cluster",
    "meta": {
      "icon": "logo-buffer",
      "title": "集群",
      "access": ["cluster"]
    },
    "component": "main",
    "children": [
      {
        "path": "node",
        "name": "node",
        "meta": {
          "icon": "md-git-network",
          "title": "节点",
          "cache": false,
          "access": ["node"]
        },
        "component": "node"
      },
      {
        "path": "node/:nodename",
        "name": "nodedetail",
        "meta": {
          "title": "节点信息",
          "icon": "ios-document",
          "hideInMenu": true,
          "cache": false,
          "access": ["node"]
        },
        "component": "nodeDetail"

      },
      {
        "path": "/storage",
        "component": "parentView",
        "name": "storage",
        "meta": {
          "title": "存储",
          "hideInMenu": false,
          "icon": "ios-filing-outline",
          "access": ["admin1"]
        },
        "children": [
          {
            "path": "pvc",
            "name": "pvc",
            "meta": {
              "title": "持久卷申请",
              "icon": "ios-create"
            },
            "component": "pvc"
          },
          {
            "path": "pvc/:pvcname",
            "name": "pvcdetail",
            "meta": {
              "title": "持久卷申请信息",
              "icon": "ios-document",
              "hideInMenu": true,
              "cache": false
            },
            "component": "pvcDetail"
          },
          {
            "path": "pv",
            "name": "pv",
            "meta": {
              "title": "持久卷",
              "icon": "ios-create-outline"
            },
            "component": "pv"
          },
          {
            "path": "pv/:pvname",
            "name": "pvdetail",
            "meta": {
              "title": "持久卷信息",
              "icon": "ios-document",
              "hideInMenu": true,
              "cache": false
            },
            "component": "pvDetail"
          },
          {
            "path": "sclass",
            "name": "sclass",
            "meta": {
              "title": "存储类",
              "icon": "ios-link"
            },
            "component": "sclass"
          },
          {
            "path": "sclass/:scname",
            "name": "scdetail",
            "meta": {
              "title": "存储类信息",
              "icon": "ios-document",
              "hideInMenu": true,
              "cache": false
            },
            "component": "sclassDetail"
          }
        ]
      },
      {
        "path": "/rbac",
        "name": "rbac",
        "meta": {
          "icon": "logo-buffer",
          "title": "角色控制",
          "access": ["admin1"]
        },
        "component": "parentView",
        "children": [
          {
            "path": "srvaccount",
            "name": "srvaccount",
            "meta": {
              "title": "服务账户",
              "icon": "ios-key-outline"
            },
            "component": "srvaccount"
          },
          {
            "path": "role",
            "name": "role",
            "meta": {
              "title": "角色",
              "icon": "ios-ionic-outline"
            },
            "component": "role"
          },
          {
            "path": "rolebind",
            "name": "rolebind",
            "meta": {
              "title": "角色绑定",
              "icon": "ios-ionic-outline"
            },
            "component": "rolebind"
          },
          {
            "path": "clusterrole",
            "name": "clusterrole",
            "meta": {
              "title": "集群角色",
              "icon": "md-ionic"
            },
            "component": "clusterrole"
          },
          {
            "path": "clusterrolebind",
            "name": "clusterrolebind",
            "meta": {
              "title": "集群角色绑定",
              "icon": "md-ionic"
            },
            "component": "clusterrolebind"
          }
        ]
      }
    ]
  },
  {
    "path": "/application",
    "component": "main",
    "name": "application",
    "meta": {
      "title": "应用管理",
      "icon": "ios-apps",
      "access": ["application"]
    },
    "children": [
      {
        "path": "apply",
        "name": "apply",
        "meta": {
          "title": "应用",
          "cache": false,
          "icon": "ios-document",
          "access": ["apply"]
        },
        "component": "apply"
      },
      {
        "path": "apply/:setName/:ctrl",
        "name": "applyinfo",
        "meta": {
          "title": "应用信息",
          "icon": "ios-document",
          "hideInMenu": true,
          "cache": false,
          "access": ["apply"]
        },
        "component": "applyInfo"
      },{
        "path": "template",
        "name": "template",
        "meta": {
          "title": "模板",
          "cache": false,
          "icon": "ios-browsers-outline",
          "access": ["template"]
        },
        "component": "template"
      }
    ]
  },
  {
    "path": "/workload",
    "component": "main",
    "name": "workload",
    "meta": {
      "title": "工作区",
      "icon": "ios-stats",
      "access": ["workload"]
    },
    "children": [
      {
        "path": "deployment",
        "name": "deployment",
        "meta": {
          "title": "部署",
          "cache": false,
          "icon": "ios-open-outline",
          "access": ["deployment"]
        },
        "component": "deployment"
      },
      {
        "path": "deployment/:setName/:ctrl",
        "name": "deploymentinfo",
        "meta": {
          "title": "部署信息",
          "icon": "ios-document",
          "hideInMenu": true,
          "cache": false,
          "access": ["deployment"]
        },
        "component": "deploymentInfo"
      },
      {
        "path": "pod",
        "name": "pod",
        "meta": {
          "title": "Pod",
          "icon": "ios-cube-outline",
          "cache": false,
          "access": ["pod"]
        },
        "component": "pod"
      },
      {
        "path": "podexec/:podName/:containerName",
        "name": "podexec",
        "meta": {
          "title": "Pod终端",
          "icon": "ios-cube-outline",
          "cache": false,
          "hideInMenu": true,
          "access": ["pod"]
        },
        "component": "podexec"
      },
      {
        "path": "kubectl/:podName/:containerName",
        "name": "kubectl",
        "meta": {
          "title": "kubectl",
          "icon": "ios-cube-outline",
          "cache": false,
          "hideInMenu": true,
          "access": ["pod"]
        },
        "component": "podexeckubectl"
      },
      {
        "path": "poddebug/:podName/:containerName",
        "name": "poddebug",
        "meta": {
          "title": "Pod 调试",
          "icon": "ios-cube-outline",
          "cache": false,
          "hideInMenu": true,
          "access": ["pod"]
        },
        "component": "poddebug"
      },
      {
        "path": "pod/:podName",
        "name": "podinfo",
        "meta": {
          "title": "Pod 信息",
          "icon": "ios-document",
          "hideInMenu": true,
          "cache": false,
          "access": ["pod"]
        },
        "component": "podInfo"
      }
    ]
  },
  {
    "path": "/discovery",
    "component": "main",
    "name": "discovery",
    "meta": {
      "title":"服务及入口",
      "hideInMenu": false,
      "icon": "md-options",
      "access": ["discovery"]
    },
    "children": [
      {
        "path": "service",
        "name": "service",
        "meta": {
          "title": "服务",
          "icon": "md-compass",
          "access": ["service"]
        },
        "component": "service"
      },
      {
        "path": "service/:servicename",
        "name": "serviceinfo",
        "meta": {
          "title": "服务信息",
          "icon": "md-compass",
          "hideInMenu": true,
          "cache": false,
          "access": ["service"]
        },
        "component": "serviceInfo"
      },
      {
        "path": "ingress",
        "name": "ingress",
        "meta": {
          "title": "Ingress",
          "icon": "md-contract",
          "access": ["ingress"]
        },
        "component": "ingress"
      },
      {
        "path": "ingress/:ingressname",
        "name": "ingressinfo",
        "meta": {
          "title": "入口信息",
          "icon": "logo-contract",
          "hideInMenu": true,
          "cache": false,
          "access": ["ingress"]
        },
        "component": "ingressInfo"
      }
    ]
  },
  {
    "path": "/configuration",
    "component": "main",
    "name": "configuration",
    "meta": {
      "title": "配置与存储",
      "hideInMenu": false,
      "icon": "md-cog",
      "access": ["configuration"]
    },
    "children": [
      {
        "path": "map",
        "name": "map",
        "meta": {
          "title": "配置字典",
          "icon": "ios-build-outline",
          "access": ["config_map"]
        },
        "component": "map"
      },
      {
        "path": "secrets",
        "name": "secrets",
        "meta": {
          "title": "保密字典",
          "icon": "ios-eye-outline",
          "access": ["secret_map"]
        },
        "component": "secrets"
      }
    ]
  },
  {
    "path": "/cplugin",
    "component": "main",
    "name": "cplugin",
    "meta": {
      "title":"集群插件",
      "hideInMenu": false,
      "icon": "md-barcode",
      "access": ["cluster_plugin"]
    },
    "children": [
      {
        "path": "clusterPlugin",
        "name": "clusterPlugin",
        "meta": {
          "title": "集群插件",
          "icon": "md-barcode",
          "access": ["cluster_plugin"]
        },
        "component": "clusterPlugin"
      }
    ]
  },
  {
    "path": "/manage",
    "component": "main",
    "name": "manage",
    "meta": {
      "title":"系统管理",
      "hideInMenu": false,
      "icon": "ios-construct",
      "access": ["manage"]
    },
    "children": [
      {
        "path": "cluster",
        "name": "cluster_manage",
        "meta": {
          "title": "集群管理",
          "icon": "md-grid",
          "access": ["cluster_manage"]
        },
        "component": "cluster"
      },
      {
        "path": "namespace",
        "name": "namespace",
        "meta": {
          "title": "命名空间",
          "icon": "md-filing",
          "access": ["namespace"]
        },
        "component": "namespace"
      },
      {
        "path": "product",
        "name": "product",
        "meta": {
          "title": "产品线",
          "icon": "md-list",
          "access": ["product"]
        },
        "component": "product"
      },
      {
        "path": "platform_role",
        "name": "platform_role",
        "meta": {
          "title": "平台角色",
          "icon": "md-color-filter",
          "access": ["platform_role"]
        },
        "component": "platformRole"
      },
      {
        "path": "user",
        "name": "user",
        "meta": {
          "title": "用户管理",
          "icon": "md-person",
          "access": ["user"]
        },
        "component": "user"
      },
      {
        "path": "plugin",
        "name": "plugin",
        "meta": {
          "title": "插件管理",
          "icon": "md-barcode",
          "access": ["plugin"]
        },
        "component": "plugin"
      },
      {
        "path": "config",
        "name": "config",
        "meta": {
          "title": "配置管理",
          "icon": "md-create",
          "access": ["config"]
        },
        "component": "audit"
      },
      {
        "path": "audit",
        "name": "audit",
        "meta": {
          "title": "审计日志",
          "icon": "ios-analytics",
          "access": ["audit"]
        },
        "component": "audit"
      }
    ]
  },
  {
    "path": "/_search",
    "component": "main",
    "name": "_search",
    "meta": {
      "hideInMenu": true
    },
    "children": [
      {
        "path": "/search",
        "name": "search",
        "meta": {
          "title": "搜索",
          "hideInMenu": true
        },
        "component": "search"
      }
    ]
  },
  {
    "path": "/_profile",
    "component": "main",
    "name": "_profile",
    "meta": {
      "hideInMenu": true
    },
    "children": [
      {
        "path": "/profile",
        "name": "profile",
        "meta": {
          "hideInMenu": true
        },
        "component": "search"
      }
    ]
  },
  {
    "path": "/401",
    "name": "error_401",
    "meta": {
      "hideInMenu": true
    },
    "component": "errorPage401"
  },
  {
    "path": "/500",
    "name": "error_500",
    "meta": {
      "hideInMenu": true
    },
    "component": "errorPage500"
  },
  {
    "path": "*",
    "name": "error_404",
    "meta": {
      "hideInMenu": true
    },
    "component": "errorPage404"
  }
]`

const istioMenu = `{
	"path": "/istio",
	"component": "main",
	"name": "istio",
	"meta": {
		"title": "微服务治理",
		"hideInMenu": false,
		"icon": "md-expand",
		"access": ["istio"]
	},
	"children": [{
			"path": "gateway",
			"name": "gateway",
			"meta": {
				"title": "网关管理",
				"icon": "ios-sync",
				"access": ["gateway"]
			},
			"component": "gateway"
		},
		{
			"path": "gateway/:resourceName",
			"name": "gatewayinfo",
			"meta": {
				"title": "网关信息",
				"icon": "md-compass",
				"hideInMenu": true,
				"cache": false,
				"access": ["gateway"]
			},
			"component": "gatewayInfo"
		},
		{
			"path": "destinationrule",
			"name": "destinationrule",
			"meta": {
				"title": "目标规则",
				"icon": "ios-shuffle",
				"access": ["dr"]
			},
			"component": "destinationrule"
		},
		{
			"path": "destinationrule/:resourceName",
			"name": "destinationruleinfo",
			"meta": {
				"title": "目标规则信息",
				"icon": "md-compass",
				"hideInMenu": true,
				"cache": false,
				"access": ["dr"]
			},
			"component": "destinationruleInfo"
		},
		{
			"path": "virtualservice",
			"name": "virtualservice",
			"meta": {
				"title": "虚拟服务",
				"icon": "md-copy",
				"access": ["vs"]
			},
			"component": "virtualservice"
		},
		{
			"path": "virtualservice/:resourceName",
			"name": "virtualserviceinfo",
			"meta": {
				"title": "虚拟服务信息",
				"icon": "md-podium",
				"hideInMenu": true,
				"cache": false,
				"access": ["vs"]
			},
			"component": "virtualserviceInfo"
		},
		{
			"path": "serviceentry",
			"name": "serviceentry",
			"meta": {
				"title": "服务入口",
				"icon": "ios-arrow-dropup",
				"access": ["se"]
			},
			"component": "serviceentry"
		},
		{
			"path": "serviceentry/:resourceName",
			"name": "serviceentryinfo",
			"meta": {
				"title": "服务入口信息",
				"icon": "ios-arrow-dropup",
				"hideInMenu": true,
				"cache": false,
				"access": ["se"]
			},
			"component": "serviceentryInfo"
		}
	]
}`

const inspectMenu = `{
    "path": "/cinspect",
    "component": "main",
    "name": "cinspect",
    "meta": {
      "title":"健康检查",
      "hideInMenu": false,
      "icon": "md-recording",
      "access": ["inspect"]
    },
    "children": [
      {
        "path": "inspect",
        "name": "inspect",
        "meta": {
          "title": "健康检查",
          "icon": "md-recording",
          "access": ["inspect"]
        },
        "component": "inspect"
      },
      {
        "path": "inspect/:resourceName",
        "name": "inspectinfo",
        "meta": {
          "title": "健康检查信息",
          "icon": "md-compass",
          "hideInMenu": true,
          "cache": false,
          "access": ["inspect"]
        },
        "component": "inspectInfo"
      }
    ]
  }`

type MenuResource struct {
	Params   *handle.Resources
	PostData *common.PlatformRoleDB
}

type menu struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Redirect string `json:"redirect,omitempty"`
	Meta     struct {
		Title       string   `json:"title,omitempty"`
		HideInMenu  bool     `json:"hideInMenu,omitempty"`
		HideInBread bool     `json:"hideInBread,omitempty"`
		NotCache    bool     `json:"notCache,omitempty"`
		Cache       bool     `json:"cache,omitempty"`
		Icon        string   `json:"icon,omitempty"`
		Access      []string `json:"access,omitempty"`
	} `json:"meta"`
	Component string `json:"component"`
	Children  []menu `json:"children,omitempty"`
}

func menuDataUnSerialize(data string) ([]menu, error) {
	menuList := make([]menu, 0)
	if err := json.Unmarshal([]byte(data), &menuList); err != nil {
		return nil, err
	}
	return menuList, nil
}

func menuUnSerialize(data string) (menu, error) {
	menuList := menu{}
	if err := json.Unmarshal([]byte(data), &menuList); err != nil {
		return menuList, err
	}
	return menuList, nil
}

func SliceInsert(s []menu, index int, value menu) []menu {
	rear := append([]menu{}, s[index:]...)
	return append(append(s[:index], value), rear...)
}

func (r *MenuResource) Get() (interface{}, error) {
	menuData, err := menuDataUnSerialize(menuData)
	if err != nil {
		return nil, err
	}
	plugin := make([]common.PluginDB, 0)
	if err := db.List(common.DataField, common.PluginTable, &plugin, ""); err != nil {
		return nil, err
	}
	clusterPlugin := make([]common.ClusterPluginDB, 0)
	if err := db.List(common.DataField, common.ClusterPluginTable, &clusterPlugin, "WHERE data-> '$.cluster'=?", r.Params.Cluster); err != nil {
		return nil, err
	}
	menuLen := len(menuData)
	for _, c := range clusterPlugin {
		if c.Status == 1 {
			for _, p := range plugin {
				if c.Plugin == p.Id && p.Name == "Istio" {
					if isitoMenuData, err := menuUnSerialize(istioMenu); err != nil {
						return nil, err
					} else {
						menuData = SliceInsert(menuData, 8, isitoMenuData)
					}
				}
				if c.Plugin == p.Id && p.Name == "Inspect" {
					if isitoMenuData, err := menuUnSerialize(inspectMenu); err != nil {
						return nil, err
					} else {
						menuData = SliceInsert(menuData, 8+(len(menuData)-menuLen), isitoMenuData)
					}
				}
			}
		}
	}
	return menuData, nil
}
