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
	client := *kit.LdapLookup(config.LDAPDB.URL,
		config.LDAPDB.SearchDN,
		config.LDAPDB.SearchPassword,
		config.LDAPDB.BaseDN,
		config.LDAPDB.UserFilter,
		config.LDAPDB.TLS)
	login := LoginPost{}
	c.BindJSON(&login)
	// ldap进行验证,获取用户信息
	ok, user, err := client.Authenticate(login.Username, login.Password)
	defer client.Close()
	if err != nil || !ok {
		responseData.Msg = err.Error()
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, responseData)
		log.Errorf("ldap authenticate error: %s", err)
		return
	}
	users := common.User{}
	if err := db.Get(common.UserTable, map[string]interface{}{"$.name": login.Username}, &users); err != nil {
		responseData.Msg = "The user does not exist"
		responseData.Data = ""
		responseData.Code = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, responseData)
		return
	}
	//生成token
	j := &jwtAuth.JWT{
		SigningKey: []byte(common.Signing),
	}
	claims := jwtAuth.CustomClaims{
		ID:    users.Id,
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
