package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/log"
	"github.com/open-kingfisher/king-utils/db"
	"github.com/open-kingfisher/king-utils/kit"
	"github.com/open-kingfisher/king-utils/middleware/jwt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sync"
	"time"
)

func getProduct(user *common.User) ([]map[string]string, error) {
	productResponse := make([]map[string]string, 0)
	productList := make([]common.ProductTree, 0)
	if err := db.List(common.DataField, common.ProductTreeTable, &productList, ""); err != nil {
		log.Errorf("product list error: %s", err)
		return []map[string]string{}, err
	}
	for _, p := range user.Product {
		for _, pp := range productList {
			if p == pp.Id {
				productResponse = append(productResponse, map[string]string{
					"id":   pp.Id,
					"name": pp.Name,
				})
			}
		}
	}
	return productResponse, nil
}

func getCluster(user *common.User) ([]map[string]string, error) {
	clusterResponse := make([]map[string]string, 0)
	clusterList := make([]common.ClusterDB, 0)
	if err := db.List(common.DataField, common.ClusterTable, &clusterList, ""); err != nil {
		log.Errorf("cluster list error: %s", err)
		return []map[string]string{}, err
	}
	for _, c := range user.Cluster {
		for _, cc := range clusterList {
			if c == cc.Id {
				clusterResponse = append(clusterResponse, map[string]string{
					"id":   cc.Id,
					"name": cc.Name,
				})
			}
		}
	}
	return clusterResponse, nil
}

func getNamespace(user *common.User) ([]map[string]string, error) {
	namespaceResponse := make([]map[string]string, 0)
	namespaceList := make([]common.NamespaceDB, 0)
	if err := db.List(common.DataField, common.NamespaceTable, &namespaceList, ""); err != nil {
		log.Errorf("Namespace get error: %s", err)
		return []map[string]string{}, err
	}
	for _, n := range user.Namespace {
		for _, nn := range namespaceList {
			if n == nn.Id {
				namespaceResponse = append(namespaceResponse, map[string]string{
					"id":   nn.Id,
					"name": nn.Name,
				})
			}
		}
	}
	return namespaceResponse, nil
}

func getPlatformRole(user *common.User) ([]map[string]string, []string, error) {
	platformRoleResponse := make([]map[string]string, 0)
	platformRole := common.PlatformRoleDB{}
	if err := db.GetById(common.PlatformRoleTable, user.Role, &platformRole); err != nil {
		log.Errorf("Platform get error: %s", err)
		return []map[string]string{}, []string{}, err
	}
	platformRoleResponse = append(platformRoleResponse, map[string]string{
		"id":   platformRole.Id,
		"name": platformRole.Name,
	})
	return platformRoleResponse, platformRole.Access, nil
}

func ListUser(c *gin.Context) {
	userResponse := make([]map[string]interface{}, 0)
	responseData := common.ResponseData{}
	users := make([]*common.User, 0)
	err := db.List(common.DataField, common.UserTable, &users, "")
	if err == nil {
		var gg errgroup.Group
		var mu sync.Mutex
		for _, user := range users {
			user := user
			gg.Go(func() error {
				var g errgroup.Group
				product := make([]map[string]string, 0)
				cluster := make([]map[string]string, 0)
				namespace := make([]map[string]string, 0)
				platformRole := make([]map[string]string, 0)
				// 获取产品线相关信息
				g.Go(func() error {
					if product, err = getProduct(user); err != nil {
						return err
					}
					return nil
				})
				// 获取集群相关信息
				g.Go(func() error {
					if cluster, err = getCluster(user); err != nil {
						return err
					}
					return nil
				})
				// 获取Namespace相关信息
				g.Go(func() error {
					if namespace, err = getNamespace(user); err != nil {
						return err
					}
					return nil
				})
				// 获取用户角色信息
				g.Go(func() error {
					if platformRole, _, err = getPlatformRole(user); err != nil {
						return err
					}
					return nil
				})
				if err := g.Wait(); err != nil {
					log.Errorf("user g.Wait: %s", err)
				}
				namespace, _ = getNamespace(user)
				mu.Lock()
				userResponse = append(userResponse, map[string]interface{}{
					"id":         user.Id,
					"name":       user.Name,
					"realName":   user.RealName,
					"mobile":     user.Mobile,
					"mail":       user.Mail,
					"product":    product,
					"cluster":    cluster,
					"namespace":  namespace,
					"role":       platformRole,
					"authMode":   user.AuthMode,
					"modifyTime": user.ModifyTime,
					"createTime": user.CreateTime,
				})
				mu.Unlock()
				return nil
			})
		}
		if err := gg.Wait(); err != nil {
			log.Errorf("user gg.Wait: %s", err)
		}
		responseData.Msg = ""
		responseData.Data = userResponse
		responseData.Code = http.StatusOK
	} else {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		log.Errorf("User list: %s", err)
	}
	c.JSON(responseData.Code, responseData)
}

func GetUser(c *gin.Context) {
	userResponse := map[string]interface{}{}
	id := c.Param("id")
	responseData := common.ResponseData{}
	user := &common.User{}
	err := db.GetById(common.UserTable, id, user)
	if err == nil {
		var g errgroup.Group
		var mu sync.Mutex
		product := make([]map[string]string, 0)
		cluster := make([]map[string]string, 0)
		namespace := make([]map[string]string, 0)
		platformRole := make([]map[string]string, 0)
		platformRoleAccess := make([]string, 0)
		// 获取产品线相关信息
		g.Go(func() error {
			if product, err = getProduct(user); err != nil {
				return err
			}
			return nil
		})
		// 获取集群相关信息
		g.Go(func() error {
			if cluster, err = getCluster(user); err != nil {
				return err
			}
			return nil
		})
		// 获取Namespace相关信息
		g.Go(func() error {
			if namespace, err = getNamespace(user); err != nil {
				return err
			}
			return nil
		})
		// 获取用户角色信息
		g.Go(func() error {
			if platformRole, platformRoleAccess, err = getPlatformRole(user); err != nil {
				return err
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			log.Errorf("user g.Wait: %s", err)
		}
		mu.Lock()
		userResponse = map[string]interface{}{
			"id":         user.Id,
			"name":       user.Name,
			"realName":   user.RealName,
			"mobile":     user.Mobile,
			"mail":       user.Mail,
			"product":    product,
			"cluster":    cluster,
			"namespace":  namespace,
			"access":     platformRoleAccess,
			"role":       platformRole,
			"authMode":   user.AuthMode,
			"modifyTime": user.ModifyTime,
			"createTime": user.CreateTime,
		}
		mu.Unlock()
		responseData.Msg = ""
		responseData.Data = userResponse
		responseData.Code = http.StatusOK
	} else {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		log.Errorf("User get: %s", err)
	}
	c.JSON(responseData.Code, responseData)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	user := c.MustGet("user").(*jwt.CustomClaims)
	users := common.User{}
	if err := db.GetById(common.UserTable, id, &users); err != nil {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		log.Errorf("User get: %s", err)
		return
	}
	responseData := common.ResponseData{}
	auditLog := common.AuditLog{
		Type:       common.UserTable,
		Id:         id,
		Name:       users.Name,
		User:       user.Name,
		ProductId:  "",
		Cluster:    "",
		Json:       "",
		ActionTime: time.Now().Unix(),
		ActionType: common.Delete,
		Namespace:  "",
		Result:     true,
		Msg:        "",
	}
	err := db.Delete(common.UserTable, id)
	if err == nil {
		responseData.Msg = ""
		responseData.Data = ""
		responseData.Code = http.StatusOK
	} else {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		auditLog.Result = false
		auditLog.Msg = err.Error()
		log.Errorf("User delete: %s", err)
	}
	if err := db.Insert(common.AuditLogTable, auditLog); err != nil {
		log.Errorf("User insert audit log: %s", err)
	}
	c.JSON(responseData.Code, responseData)
}

func CreateUser(c *gin.Context) {
	user := c.MustGet("user").(*jwt.CustomClaims)
	responseData := common.ResponseData{}
	uuid := kit.UUID("u")
	users := common.User{
		Id:         uuid,
		CreateTime: time.Now().Unix(),
		ModifyTime: time.Now().Unix(),
	}
	c.BindJSON(&users)
	users.Password, _ = HashPassword(users.Password)
	// 对提交的数据进行校验
	if err := c.ShouldBindWith(&users, binding.Query); err != nil {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusBadRequest
		c.JSON(http.StatusBadRequest, responseData)
		return
	}
	// 判断用户名是否已经存在
	userList := make([]*common.User, 0)
	if err := db.List(common.DataField, common.UserTable, &userList, "WHERE data-> '$.name'='"+users.Name+"'"); err == nil {
		if len(userList) > 0 {
			responseData.Msg = "The user name already exists"
			responseData.Data = ""
			responseData.Code = http.StatusBadRequest
			c.JSON(http.StatusBadRequest, responseData)
			return
		}
	} else {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, responseData)
		return
	}
	auditLog := common.AuditLog{
		Type:       common.UserTable,
		Id:         uuid,
		Name:       c.PostForm("name"),
		User:       user.Name,
		ProductId:  "",
		Cluster:    "",
		Json:       users,
		ActionTime: time.Now().Unix(),
		ActionType: common.Create,
		Namespace:  "",
		Result:     true,
		Msg:        "",
	}
	if err := db.Insert(common.UserTable, users); err == nil {
		responseData.Msg = ""
		responseData.Data = ""
		responseData.Code = http.StatusOK
	} else {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		auditLog.Msg = err.Error()
		log.Errorf("User create error:%s; Json:%+v; Name:%s", err, users, users.Name)
	}
	if err := db.Insert(common.AuditLogTable, auditLog); err != nil {
		log.Errorf("User insert audit log: %s", err)
	}
	c.JSON(responseData.Code, responseData)
}

func UpdateUser(c *gin.Context) {
	id := c.Query("id")
	user := c.MustGet("user").(*jwt.CustomClaims)
	responseData := common.ResponseData{}
	users := common.User{
		Id: id,
	}
	if err := db.GetById(common.UserTable, id, &users); err != nil {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusBadRequest
		c.JSON(http.StatusBadRequest, responseData)
		return
	}
	password := users.Password
	c.BindJSON(&users)
	users.ModifyTime = time.Now().Unix()
	// 编辑不能修改密码
	users.Password = password
	// 对提交的数据进行校验
	if err := c.ShouldBindWith(&users, binding.Query); err != nil {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusBadRequest
		c.JSON(http.StatusBadRequest, responseData)
		return
	}
	// 判断用户名是否已经存在
	userList := make([]*common.ClusterDB, 0)
	if err := db.List(common.DataField, common.UserTable, &userList, "WHERE data-> '$.name'='"+users.Name+"'"); err == nil {
		if len(userList) > 0 {
			for _, v := range userList {
				if v.Id != id {
					responseData.Msg = "The User name already exists"
					responseData.Data = ""
					responseData.Code = http.StatusBadRequest
					c.JSON(http.StatusBadRequest, responseData)
					return
				}
			}
		}
	} else {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, responseData)
		return
	}
	auditLog := common.AuditLog{
		Type:       common.UserTable,
		Id:         id,
		Name:       users.Name,
		User:       user.Name,
		ProductId:  "",
		Cluster:    "",
		Json:       users,
		ActionTime: time.Now().Unix(),
		ActionType: common.Update,
		Namespace:  "",
		Result:     true,
		Msg:        "",
	}
	if err := db.Update(common.UserTable, id, users); err == nil {
		responseData.Msg = ""
		responseData.Data = ""
		responseData.Code = http.StatusOK
	} else {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		auditLog.Msg = err.Error()
		log.Errorf("User create error:%s; Json:%+v; Name:%s", err, users, users.Name)
	}
	if err := db.Insert(common.AuditLogTable, auditLog); err != nil {
		log.Errorf("User insert audit log: %s", err)
	}
	c.JSON(responseData.Code, responseData)
}

func GetUserByToken(c *gin.Context) {
	token := c.GetHeader(common.HeaderSigning)
	jsonWebToken := jwt.JWT{SigningKey: []byte(common.Signing)}
	j, err := jsonWebToken.ParseToken(token)
	if err != nil {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, responseData)
		return
	}
	responseData := common.ResponseData{}
	user := &common.User{}
	err = db.GetById(common.UserTable, j.ID, user)
	if err == nil {
		var g errgroup.Group
		product := make([]map[string]string, 0)
		platformRoleAccess := make([]string, 0)
		// 获取产品线相关信息
		g.Go(func() error {
			if product, err = getProduct(user); err != nil {
				return err
			}
			return nil
		})
		// 获取用户角色信息
		g.Go(func() error {
			if _, platformRoleAccess, err = getPlatformRole(user); err != nil {
				return err
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			log.Errorf("user g.Wait: %s", err)
		}
		userResponse := map[string]interface{}{
			"id":       user.Id,
			"name":     user.Name,
			"realName": user.RealName,
			"mobile":   user.Mobile,
			"mail":     user.Mail,
			"product":  product,
			"access":   platformRoleAccess,
		}
		responseData.Msg = ""
		responseData.Data = userResponse
		responseData.Code = http.StatusOK
	} else {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		log.Errorf("Get user by token: %s", err)
	}
	c.JSON(responseData.Code, responseData)
}

func GetUserAuthMode(c *gin.Context) {
	config := common.ConfigDB{}
	if err := db.GetById(common.ConfigTable, common.ConfigID, &config); err != nil {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, responseData)
		return
	} else {
		if config.LDAPDB.Mode == "ldap" {
			responseData.Data = map[string]string{"mode": "ldap"}
		} else {
			responseData.Data = map[string]string{"mode": "local"}
		}
		responseData.Msg = ""
		responseData.Code = http.StatusOK
	}
	c.JSON(responseData.Code, responseData)
}

func ChangePassword(c *gin.Context) {
	responseData := common.ResponseData{}
	var changePassword struct {
		UserId   string `json:"userId"`
		Password string `json:"password"`
	}
	c.BindJSON(&changePassword)
	users := common.User{}
	if err := db.GetById(common.UserTable, changePassword.UserId, &users); err != nil {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusBadRequest
		c.JSON(http.StatusBadRequest, responseData)
		return
	}
	users.ModifyTime = time.Now().Unix()
	users.Password, _ = HashPassword(changePassword.Password)
	// 对提交的数据进行校验
	if err := c.ShouldBindWith(&users, binding.Query); err != nil {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusBadRequest
		c.JSON(http.StatusBadRequest, responseData)
		return
	}
	auditLog := common.AuditLog{
		Type:       common.UserTable,
		Id:         users.Id,
		Name:       users.Name,
		User:       users.Name,
		ProductId:  "",
		Cluster:    "",
		Json:       users,
		ActionTime: time.Now().Unix(),
		ActionType: common.Update,
		Namespace:  "",
		Result:     true,
		Msg:        "",
	}
	if err := db.Update(common.UserTable, changePassword.UserId, users); err == nil {
		responseData.Msg = ""
		responseData.Data = ""
		responseData.Code = http.StatusOK
	} else {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		auditLog.Msg = err.Error()
		log.Errorf("User create error:%s; Json:%+v; Name:%s", err, users, users.Name)
	}
	if err := db.Insert(common.AuditLogTable, auditLog); err != nil {
		log.Errorf("User insert audit log: %s", err)
	}
	c.JSON(responseData.Code, responseData)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
