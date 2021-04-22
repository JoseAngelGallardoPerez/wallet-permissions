package routes

import (
	"net/http"

	"github.com/Confialink/wallet-permissions/internal/app/di"
	"github.com/Confialink/wallet-permissions/internal/authentication"
	"github.com/Confialink/wallet-permissions/internal/http/controller"
	"github.com/Confialink/wallet-permissions/internal/http/middleware"
	"github.com/Confialink/wallet-permissions/internal/version"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine
var c = di.Container

func initRoutes() {

	r = gin.New()
	mwAuth := authentication.Middleware(c.ServiceLogger().New("middleware", "authentication"))
	mwCors := middleware.Cors(c.Config().Cors)

	r.GET("/permissions/health-check", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/permissions/build", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.BuildInfo)
	})

	r.OPTIONS("/*cors", mwCors)

	apiGroup := r.Group("permissions")
	apiGroup.Use(
		mwCors,
		gin.Recovery(),
		gin.Logger(),
		errors.ErrorHandler(c.ServiceLogger().New("middleware", "errors")),
	)

	privateGroup := apiGroup.Group("/private", mwAuth)
	{
		v1Group := privateGroup.Group("/v1")
		{
			groupController := controller.Factory.GroupController()
			permissionController := controller.Factory.PermissionController()

			adminGroup := v1Group.Group("/admin", middleware.AdminOnly)
			{
				adminGroup.GET("/permission", permissionController.ListAction)
				adminGroup.GET("/permission/:key", permissionController.GetAction)

				adminGroup.GET("/group", groupController.ListAction)
				adminGroup.GET("/group/:id", groupController.GetAction)
				adminGroup.POST("/group", groupController.CreateAction)
				adminGroup.POST("/group/:id", groupController.UpdateAction)
				adminGroup.DELETE("/group/:id", groupController.DeleteAction)
				adminGroup.DELETE("/group/:id/action/:action", groupController.RemoveActionAction)

				actionsController := controller.Factory.ActionController()
				adminGroup.GET("/action/:id", actionsController.GetAction)
				adminGroup.GET("/action", actionsController.ListAction)

				categoryController := controller.Factory.CategoryController()
				adminGroup.GET("/category", categoryController.List)
			}

			clientGroup := v1Group.Group("/client", middleware.MainUserOnly)
			{
				clientGroup.GET("/group", groupController.ListAction)
				clientGroup.GET("/permission", permissionController.ListAction)
				clientGroup.GET("/permission/:key", permissionController.GetAction)
			}
		}
	}
}

func GetRouter() *gin.Engine {
	if nil == r {
		initRoutes()
	}
	return r
}
