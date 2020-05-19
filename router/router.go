package router

import (
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-kf/impl"
	"github.com/open-kingfisher/king-kf/user"
	"github.com/open-kingfisher/king-utils/common"
	jwtAuth "github.com/open-kingfisher/king-utils/middleware/jwt"
	"net/http"
)

func SetupRouter(r *gin.Engine) *gin.Engine {
	// token
	//r.GET(common.KingfisherPath+"token", user.GetToken)

	//重新定义404
	r.NoRoute(NoRoute)
	// login
	r.POST(common.KingfisherPath+"login", user.Login)
	r.GET(common.KingfisherPath+"logout", user.Logout)
	r.GET(common.KingfisherPath+"token", user.Index)

	authorize := r.Group("/", jwtAuth.JWTAuth())
	{

		// user lookup
		authorize.GET(common.KingfisherPath+"userlookup", user.GetUserLookup)

		// product
		authorize.GET(common.KingfisherPath+"product", impl.ListProduct)
		authorize.GET(common.KingfisherPath+"product/:name", impl.GetProduct)
		authorize.POST(common.KingfisherPath+"product", impl.CreateProduct)
		authorize.DELETE(common.KingfisherPath+"product/:name", impl.DeleteProduct)
		authorize.PUT(common.KingfisherPath+"product", impl.UpdateProduct)

		// cluster
		authorize.GET(common.KingfisherPath+"cluster", impl.ListCluster)
		authorize.GET(common.KingfisherPath+"cluster/:name", impl.GetCluster)
		authorize.DELETE(common.KingfisherPath+"cluster/:name", impl.DeleteCluster)
		authorize.POST(common.KingfisherPath+"cluster", impl.CreateCluster)
		authorize.PUT(common.KingfisherPath+"cluster", impl.UpdateCluster)

		// user
		authorize.GET(common.KingfisherPath+"user", user.ListUser)
		authorize.GET(common.KingfisherPath+"user/:id", user.GetUser)
		authorize.DELETE(common.KingfisherPath+"user/:id", user.DeleteUser)
		authorize.POST(common.KingfisherPath+"user", user.CreateUser)
		authorize.PUT(common.KingfisherPath+"user", user.UpdateUser)
		authorize.GET(common.KingfisherPath+"userByToken", user.GetUserByToken)

		// audit log
		authorize.GET(common.KingfisherPath+"audit", impl.ListAuditLog)

		// cascade
		authorize.GET(common.KingfisherPath+"cascade", impl.CascadeCluster)
		authorize.GET(common.KingfisherPath+"cascadeAll", impl.CascadeAll)
		authorize.GET(common.KingfisherPath+"treeClusterNamespace", impl.TreeClusterNamespace)

		// plugin
		authorize.GET(common.KingfisherPath+"plugin", impl.ListPlugin)
		authorize.GET(common.KingfisherPath+"plugin/:name", impl.GetPlugin)
		authorize.POST(common.KingfisherPath+"plugin", impl.CreatePlugin)
		authorize.DELETE(common.KingfisherPath+"plugin/:name", impl.DeletePlugin)
		authorize.PUT(common.KingfisherPath+"plugin", impl.UpdatePlugin)

		// cluster plugin
		authorize.GET(common.KingfisherPath+"istio/install", impl.Install)
		authorize.GET(common.KingfisherPath+"istio/uninstall", impl.UnInstall)
		authorize.GET(common.KingfisherPath+"podDebug/install", impl.Install)
		authorize.GET(common.KingfisherPath+"podDebug/uninstall", impl.UnInstall)
		authorize.GET(common.KingfisherPath+"inspect/install", impl.Install)
		authorize.GET(common.KingfisherPath+"inspect/uninstall", impl.UnInstall)
		authorize.GET(common.KingfisherPath+"clusterplugin", impl.ListClusterPlugin)

		// platform role
		authorize.GET(common.KingfisherPath+"platformRole", impl.ListPlatformRole)
		authorize.GET(common.KingfisherPath+"platformRole/:name", impl.GetPlatformRole)
		authorize.POST(common.KingfisherPath+"platformRole", impl.CreatePlatformRole)
		authorize.DELETE(common.KingfisherPath+"platformRole/:name", impl.DeletePlatformRole)
		authorize.PUT(common.KingfisherPath+"platformRole", impl.UpdatePlatformRole)
		authorize.GET(common.KingfisherPath+"platformRoleTree", impl.ListPlatformRoleTree)

		// menu
		authorize.GET(common.KingfisherPath+"menu", impl.GetMenu)

		// template
		authorize.GET(common.KingfisherPath+"template", impl.ListTemplate)
		authorize.GET(common.KingfisherPath+"template/:name", impl.GetTemplate)
		authorize.POST(common.KingfisherPath+"template", impl.CreateTemplate)
		authorize.DELETE(common.KingfisherPath+"template/:name", impl.DeleteTemplate)
		authorize.PUT(common.KingfisherPath+"template", impl.UpdateTemplate)

		// config
		authorize.GET(common.KingfisherPath+"config", impl.ListConfig)
		authorize.GET(common.KingfisherPath+"config/:name", impl.GetConfig)
		authorize.POST(common.KingfisherPath+"config", impl.CreateConfig)
		authorize.DELETE(common.KingfisherPath+"config/:name", impl.DeleteConfig)
		authorize.PUT(common.KingfisherPath+"config", impl.UpdateConfig)

	}
	return r
}

// 重新定义404错误
func NoRoute(c *gin.Context) {
	responseData := common.ResponseData{Code: http.StatusNotFound, Msg: "404 Not Found"}
	c.JSON(http.StatusNotFound, responseData)
}
