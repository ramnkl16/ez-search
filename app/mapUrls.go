package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ramnkl16/ez-search/controllers"
)

func mapUrls(rout *gin.Engine) {
	//controllers.BrandController.RegisterRouter(rout)

	controllers.EventQueueController.RegisterRouter(rout)
	controllers.EventQueueHistoryController.RegisterRouter(rout)
	controllers.SearchController.RegisterRouter(rout)
	controllers.UserController.RegisterRouter(rout)
	controllers.UserMenuController.RegisterRouter(rout)
	controllers.NamespaceController.RegisterRouter(rout)
	controllers.UserGroupController.RegisterRouter(rout)
	controllers.MenuController.RegisterRouter(rout)
	controllers.WidgetMetaController.RegisterRouter(rout)

}
