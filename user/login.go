package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/log"
	"github.com/open-kingfisher/king-utils/db"
	"github.com/open-kingfisher/king-utils/kit"
	jwtAuth "github.com/open-kingfisher/king-utils/middleware/jwt"
	"net/http"
	"time"
)

type LoginPost struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	responseData = &common.ResponseData{}
)

func Login(c *gin.Context) {
	config := common.ConfigDB{}
	if err := db.GetById(common.ConfigTable, common.ConfigID, &config); err != nil {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, responseData)
		return
	}
	// 获取登录信息
	login := LoginPost{}
	c.BindJSON(&login)
	// 校验用户是否存在
	userDB := common.User{}
	if err := db.Get(common.UserTable, map[string]interface{}{"$.name": login.Username, "$.authMode": config.LDAPDB.Mode}, &userDB); err != nil {
		responseData.Msg = "The user does not exist"
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, responseData)
		return
	}

	user := make(map[string]string)
	// 认证模式是ldap的
	if config.LDAPDB.Mode == "ldap" {
		client := *kit.LdapLookup(config.LDAPDB.URL,
			config.LDAPDB.SearchDN,
			config.LDAPDB.SearchPassword,
			config.LDAPDB.BaseDN,
			config.LDAPDB.UserFilter,
			config.LDAPDB.TLS)
		// ldap进行验证,获取用户信息
		var err error
		var ok bool
		ok, user, err = client.Authenticate(login.Username, login.Password)
		defer client.Close()
		if err != nil || !ok {
			responseData.Msg = err.Error()
			responseData.Data = ""
			responseData.Code = http.StatusInternalServerError
			c.JSON(http.StatusInternalServerError, responseData)
			log.Errorf("ldap authenticate error: %s", err)
			return
		}
	} else {
		match := CheckPasswordHash(login.Password, userDB.Password)
		if !match {
			responseData.Msg = "password error"
			responseData.Data = ""
			responseData.Code = http.StatusInternalServerError
			c.JSON(http.StatusInternalServerError, responseData)
			log.Error("password error")
			return
		}
	}

	//生成token
	j := &jwtAuth.JWT{
		SigningKey: []byte(common.Signing),
	}
	claims := jwtAuth.CustomClaims{
		ID:    userDB.Id,
		Name:  login.Username,
		Email: user["mail"],
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(common.ExpiresAt).Unix(),
		},
	}
	if token, err := j.CreateToken(claims); err == nil {
		responseData.Msg = ""
		responseData.Data = token
		responseData.Code = http.StatusOK
	} else {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
	}
	c.JSON(responseData.Code, responseData)
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/home", common.KingfisherDomain, false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func GetUserLookup(c *gin.Context) {
	username := c.Query("username")
	config := common.ConfigDB{}
	if err := db.GetById(common.ConfigTable, common.ConfigID, &config); err != nil {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, responseData)
		return
	}
	client := *kit.LdapLookup(config.LDAPDB.URL,
		config.LDAPDB.SearchDN,
		config.LDAPDB.SearchPassword,
		config.LDAPDB.BaseDN,
		config.LDAPDB.UserFilter,
		config.LDAPDB.TLS)
	user, err := client.GetUser(username)
	defer client.Close()
	if err != nil {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
	} else {
		responseData.Msg = ""
		responseData.Data = user
		responseData.Code = http.StatusOK
	}
	c.JSON(responseData.Code, responseData)
}
