package resource

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/handle"
	"github.com/open-kingfisher/king-utils/common/log"
	"github.com/open-kingfisher/king-utils/db"
	"github.com/open-kingfisher/king-utils/kit"
	"time"
)

type ProductResponse struct {
	common.Info
	Cluster    []common.Info `json:"cluster"`
	Namespace  []common.Info `json:"namespace"`
	CreateTime int64         `json:"createTime"`
	ModifyTime int64         `json:"modifyTime"`
}

type ProductResource struct {
	Params   *handle.Resources
	PostData *common.ProductDB
}

func (r *ProductResource) Get() (interface{}, error) {
	product := common.ProductDB{}
	if err := db.GetById(common.ProductTable, r.Params.Name, &product); err != nil {
		return nil, err
	}
	clusterList := make([]common.Info, 0)
	namespaceList := make([]common.Info, 0)
	productResponse := ProductResponse{}
	for _, v := range product.Cluster {
		clusterDB := common.ClusterDB{}
		if err := db.GetById(common.ClusterTable, v, &clusterDB); err != nil {
			log.Errorf("get cluster failed, %s ", err)
		}
		clusterList = append(clusterList, common.Info{
			Id:   clusterDB.Id,
			Name: clusterDB.Name,
		})
	}
	for _, v := range product.Namespace {
		namespaceDB := common.NamespaceDB{}
		if err := db.GetById(common.NamespaceTable, v, &namespaceDB); err != nil {
			log.Errorf("get namespace failed, %s ", err)
		}
		namespaceList = append(namespaceList, common.Info{
			Id:   namespaceDB.Id,
			Name: namespaceDB.Name,
		})
	}
	productTreeDB := common.ProductTree{}
	if err := db.GetById(common.ProductTreeTable, product.Id, &productTreeDB); err != nil {
		log.Errorf("get product failed, %s ", err)
	}
	productResponse = ProductResponse{
		Info: common.Info{
			Id:   product.Id,
			Name: productTreeDB.Name,
		},
		Cluster:   clusterList,
		Namespace: namespaceList,
	}
	return productResponse, nil
}

func (r *ProductResource) List() (interface{}, error) {
	product := make([]common.ProductDB, 0)
	// 获取产品信息
	if err := db.List(common.DataField, common.ProductTable, &product, ""); err != nil {
		return nil, err
	}
	productResponse := make([]*ProductResponse, 0)
	for _, p := range product {
		clusterList := make([]common.Info, 0)
		namespaceList := make([]common.Info, 0)
		for _, v := range p.Cluster {
			clusterDB := common.ClusterDB{}
			if err := db.GetById(common.ClusterTable, v, &clusterDB); err != nil {
				log.Errorf("get cluster failed, %s ", err)
			}
			clusterList = append(clusterList, common.Info{
				Id:   clusterDB.Id,
				Name: clusterDB.Name,
			})
		}

		for _, v := range p.Namespace {
			namespaceDB := common.NamespaceDB{}
			if err := db.GetById(common.NamespaceTable, v, &namespaceDB); err != nil {
				log.Errorf("get namespace failed, %s ", err)
			}
			clusterDB := common.ClusterDB{}
			if err := db.GetById(common.Cluster, namespaceDB.Cluster, &clusterDB); err != nil {
				log.Errorf("get cluster failed, %s ", err)
			}
			namespaceList = append(namespaceList, common.Info{
				Id:   namespaceDB.Id,
				Name: fmt.Sprintf("%s |__ %s", clusterDB.Name, namespaceDB.Name),
			})
		}
		productResponse = append(productResponse, &ProductResponse{
			Info: common.Info{
				Id:   p.Id,
				Name: p.Name,
			},
			Cluster:    clusterList,
			Namespace:  namespaceList,
			CreateTime: p.CreateTime,
			ModifyTime: p.ModifyTime,
		})
	}
	return productResponse, nil
}

func (r *ProductResource) Create(c *gin.Context) (err error) {
	if err = c.BindJSON(&r.PostData); err != nil {
		return err
	}
	uuid := kit.UUID("p")
	// 判断集产品线否已经存在
	productList := make([]*common.ProductDB, 0)
	if err = db.List(common.DataField, common.ProductTable, &productList, "WHERE data-> '$.name'=?", r.PostData.Name); err == nil {
		if len(productList) > 0 {
			return errors.New("the product name already exists")
		}
	} else {
		return
	}
	cluster := make([]string, 0)
	for _, v := range r.PostData.Namespace {
		namespaceDB := common.NamespaceDB{}
		if err := db.GetById(common.NamespaceTable, v, &namespaceDB); err != nil {
			log.Errorf("get namespace failed, %s ", err)
		}
		cluster = append(cluster, namespaceDB.Cluster)
	}
	r.PostData.Cluster = UniqueList(cluster)
	r.PostData.Id = uuid
	r.PostData.CreateTime = time.Now().Unix()
	r.PostData.ModifyTime = time.Now().Unix()
	if err = db.Insert(common.ProductTable, r.PostData); err != nil {
		log.Errorf("Product create error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return err
	}
	auditLog := handle.AuditLog{
		Kind:        common.Product,
		ActionType:  common.Create,
		Resources:   r.Params,
		ProductData: r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *ProductResource) Delete() (err error) {
	if err = db.Delete(common.ProductTable, r.Params.Name); err != nil {
		return
	}
	// 从user表中删除product
	user := make([]*common.User, 0)
	if err := db.List(common.DataField, common.UserTable, &user, ""); err == nil {
		for _, v := range user {
			v.Product = kit.DeleteItemForList(r.Params.Name, v.Product)
			if err := db.Update(common.UserTable, v.Id, v); err != nil {
				log.Errorf("User table delete product :%s error:%s", v.Id, err)
			}
		}
	}
	auditLog := handle.AuditLog{
		Kind:       common.Product,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *ProductResource) Update(c *gin.Context) (err error) {
	product := common.ProductDB{}
	if err = c.BindJSON(&product); err != nil {
		return err
	}
	r.PostData = &product
	// 判断集产品线否已经存在
	var createTime int64
	productList := make([]*common.ProductDB, 0)
	if err = db.List(common.DataField, common.ProductTable, &productList, "WHERE data-> '$.name'=?", r.PostData.Name); err == nil {
		if len(productList) > 0 {
			for _, v := range productList {
				if v.Id != r.PostData.Id {
					return errors.New("the product name already exists")
				} else {
					createTime = v.CreateTime
				}
			}
		}
	} else {
		return
	}
	cluster := make([]string, 0)
	for _, v := range r.PostData.Namespace {
		namespaceDB := common.NamespaceDB{}
		if err := db.GetById(common.NamespaceTable, v, &namespaceDB); err != nil {
			log.Errorf("get namespace failed, %s ", err)
		}
		cluster = append(cluster, namespaceDB.Cluster)
	}
	r.PostData.Cluster = UniqueList(cluster)
	r.PostData.CreateTime = createTime
	r.PostData.ModifyTime = time.Now().Unix()
	if err = db.Update(common.ProductTable, r.PostData.Id, r.PostData); err != nil {
		log.Errorf("Product update error:%s; Json:%+v; Name:%s", err, r.PostData, r.PostData.Id)
		return
	}
	auditLog := handle.AuditLog{
		Kind:        common.Product,
		ActionType:  common.Update,
		Resources:   r.Params,
		ProductData: r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func UniqueList(l []string) []string {
	ll := make([]string, 0)
	m := map[string]string{}
	for _, v := range l {
		m[v] = ""
	}
	for k, _ := range m {
		ll = append(ll, k)
	}
	return ll
}

func In(s string, l []string) bool {
	for _, v := range l {
		if s == v {
			return true
		}
	}
	return false
}

type Cascade struct {
	common.InfoType
	Children []*Cascade `json:"children"`
}

// 根据用户权限的级联列出产品线--集群--Namespace
func (r *ProductResource) CascadeCluster() (interface{}, error) {
	user := common.User{}
	// 获取用户产品线信息
	if err := db.GetById(common.UserTable, r.Params.User.ID, &user); err != nil {
		log.Errorf("User get error: %s", err)
		return nil, err
	}
	product := make([]common.ProductDB, 0)
	// 获取产品信息
	if err := db.List(common.DataField, common.ProductTable, &product, ""); err != nil {
		log.Errorf("Cascade list error: %s", err)
		return nil, err
	}
	productResponse := make([]*Cascade, 0)
	for _, p := range product {
		// 判断产品线是否在用户的产品线中
		if In(p.Id, user.Product) {
			cluster := make([]*Cascade, 0)
			for _, cc := range p.Cluster {
				if In(cc, user.Cluster) {
					clusterDB := common.ClusterDB{}
					if err := db.GetById(common.ClusterTable, cc, &clusterDB); err != nil {
						log.Errorf("Cluster get error: %s", err)
					}
					ns := make([]*Cascade, 0)
					namespace := make([]*common.NamespaceDB, 0)
					if err := db.List(common.DataField, common.NamespaceTable, &namespace, "WHERE data-> '$.cluster'=?", cc); err != nil {
						log.Errorf("Namespace list error: %s", err)
					}
					for _, n := range p.Namespace {
						if In(n, user.Namespace) {
							for _, nn := range namespace {
								if n == nn.Id {
									ns = append(ns, &Cascade{
										InfoType: common.InfoType{
											Value: nn.Id,
											Label: nn.Name,
										},
										Children: []*Cascade{},
									})
								}
							}
						}
					}
					cluster = append(cluster, &Cascade{
						InfoType: common.InfoType{
							Value: clusterDB.Id,
							Label: clusterDB.Name,
						},
						Children: ns,
					})
				}
			}
			productResponse = append(productResponse, &Cascade{
				InfoType: common.InfoType{
					Value: p.Id,
					Label: p.Name,
				},
				Children: cluster,
			})
		}
	}
	return productResponse, nil
}

// 级联查询所有产品线--集群--Namespace
func (r *ProductResource) CascadeAll() (interface{}, error) {
	product := make([]common.ProductDB, 0)
	// 获取产品信息
	if err := db.List(common.DataField, common.ProductTable, &product, ""); err != nil {
		log.Errorf("Cascade list error: %s", err)
		return nil, err
	}
	productResponse := make([]*Cascade, 0)
	for _, p := range product {
		cluster := make([]*Cascade, 0)
		for _, cc := range p.Cluster {
			clusterDB := common.ClusterDB{}
			if err := db.GetById(common.ClusterTable, cc, &clusterDB); err != nil {
				log.Errorf("Cluster get error: %s", err)
			}
			ns := make([]*Cascade, 0)
			namespace := make([]*common.NamespaceDB, 0)
			if err := db.List(common.DataField, common.NamespaceTable, &namespace, "WHERE data-> '$.cluster'=?", cc); err != nil {
				log.Errorf("Namespace list error: %s", err)
			}
			for _, n := range p.Namespace {
				for _, nn := range namespace {
					if n == nn.Id {
						ns = append(ns, &Cascade{
							InfoType: common.InfoType{
								Value: nn.Id,
								//Label: nn.Name,
								Level:  2,
								Title:  nn.Name,
								Expand: true,
							},
							Children: []*Cascade{},
						})
					}
				}
			}
			cluster = append(cluster, &Cascade{
				InfoType: common.InfoType{
					Value: clusterDB.Id,
					//Label: clusterDB.Name,
					Level:  1,
					Title:  clusterDB.Name,
					Expand: true,
				},
				Children: ns,
			})
		}
		productResponse = append(productResponse, &Cascade{
			InfoType: common.InfoType{
				Value:  p.Id,
				Level:  0,
				Title:  p.Name,
				Expand: true,
			},
			Children: cluster,
		})
	}
	return productResponse, nil
}

type TreeCascade struct {
	Id       string         `json:"id"`
	Label    string         `json:"label"`
	Children []*TreeCascade `json:"children,omitempty"`
}

// 级联列出集群--Namespace
func (r *ProductResource) TreeClusterNamespace() (interface{}, error) {
	productResponse := make([]*Cascade, 0)
	clusterDB := make([]common.ClusterDB, 0)
	if err := db.List(common.DataField, common.ClusterTable, &clusterDB, ""); err != nil {
		log.Errorf("Cluster list error: %s", err)
	} else {
		for _, v := range clusterDB {
			cluster := &Cascade{
				InfoType: common.InfoType{
					Value:  v.Id,
					Level:  0,
					Title:  v.Name,
					Expand: true,
				},
			}
			namespace := make([]*common.NamespaceDB, 0)
			if err := db.List(common.DataField, common.NamespaceTable, &namespace, "WHERE data-> '$.cluster'=?", v.Id); err != nil {
				log.Errorf("Namespace list error: %s", err)
			} else {
				for _, n := range namespace {
					cluster.Children = append(cluster.Children, &Cascade{
						InfoType: common.InfoType{
							Value:  n.Id,
							Level:  1,
							Title:  n.Name,
							Expand: true,
						},
					})
				}
			}
			productResponse = append(productResponse, cluster)
		}
	}
	return productResponse, nil
}
