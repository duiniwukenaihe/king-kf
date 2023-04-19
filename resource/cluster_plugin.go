package resource

import (
	"github.com/duiniwukenaihe/king-utils/common"
	"github.com/duiniwukenaihe/king-utils/common/handle"
	"github.com/duiniwukenaihe/king-utils/common/log"
	"github.com/duiniwukenaihe/king-utils/db"
	"github.com/duiniwukenaihe/king-utils/kit"
	"time"
)

type ClusterPluginResource struct {
	Params   *handle.Resources
	PostData *common.PluginDB
	Plugin   string `json:"plugin "`
}

func (r *ClusterPluginResource) Install() error {
	pluginList := common.ClusterPluginDB{}
	query := map[string]interface{}{
		"$.plugin":  r.Plugin,
		"$.cluster": r.Params.Cluster,
	}
	if err := db.Get(common.ClusterPluginTable, query, &pluginList); err == nil {
		pluginList.Status = 1
		pluginList.Timestamp = time.Now().Unix()
		if err = db.Update(common.ClusterPluginTable, pluginList.Id, pluginList); err != nil {
			return err
		}
	} else {
		clusterPlugin := common.ClusterPluginDB{
			Id:        kit.UUID("p"),
			Plugin:    r.Plugin,
			Cluster:   r.Params.Cluster,
			Status:    1,
			Timestamp: time.Now().Unix(),
		}
		if err = db.Insert(common.ClusterPluginTable, clusterPlugin); err != nil {
			log.Errorf("Cluster Plugin add error:%s; Json:%+v;", err, r.PostData)
			return err
		}
	}
	return nil
}

func (r *ClusterPluginResource) UnInstall() error {
	clusterPlugin := common.ClusterPluginDB{}
	query := map[string]interface{}{
		"$.plugin":  r.Plugin,
		"$.cluster": r.Params.Cluster,
	}
	if err := db.Get(common.ClusterPluginTable, query, &clusterPlugin); err == nil {
		clusterPlugin.Status = 0
		clusterPlugin.Timestamp = time.Now().Unix()
		if err = db.Update(common.ClusterPluginTable, clusterPlugin.Id, clusterPlugin); err != nil {
			return err
		} else {
			return err
		}
	}
	return nil
}

func (r *ClusterPluginResource) List() (interface{}, error) {
	plugin := make([]common.PluginDB, 0)
	if err := db.List(common.DataField, common.PluginTable, &plugin, ""); err != nil {
		return nil, err
	}
	clusterPlugin := make([]common.ClusterPluginDB, 0)
	if err := db.List(common.DataField, common.ClusterPluginTable, &clusterPlugin, "WHERE data-> '$.cluster'=?", r.Params.Cluster); err != nil {
		return nil, err
	}
	plugins := make([]map[string]interface{}, 0)
	status := 0
	for _, p := range plugin {
		for _, c := range clusterPlugin {
			if p.Id == c.Plugin && c.Status == 1 {
				status = 1
				break
			} else {
				status = 0
			}
		}
		plugins = append(plugins, map[string]interface{}{
			"id":       p.Id,
			"name":     p.Name,
			"describe": p.Describe,
			"status":   status,
		})
	}
	return plugins, nil
}
