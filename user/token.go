package user

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-utils/common"
	jwtAuth "github.com/open-kingfisher/king-utils/middleware/jwt"
	"net/http"
	"time"
)

func GetToken(c *gin.Context) {
	responseData := common.ResponseData{}
	j := &jwtAuth.JWT{
		SigningKey: []byte(common.Signing),
	}
	claims := jwtAuth.CustomClaims{
		ID:    "1",
		Name:  "Tom",
		Email: "tom@kingfisher.com",
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
