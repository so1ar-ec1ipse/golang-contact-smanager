package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	con "github.com/sajjad3k/contactsmanager/controllers"
)

//func main() {
//}

func ServerRoutes() *gin.Engine {

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/contacts", con.ShowallContacts)
		api.GET("/contacts/:name", con.GetcontactbyName)
		//api.GET("/contacts/:number", con.GetcontactbyNumber)
		api.POST("/contacts", con.CreatenewContact)
		api.PUT("/contacts/:name", con.Updatecontact)
		api.DELETE("contacts/:name", con.DeleteContact)
		api.POST("/contacts/upload", con.Uploadlist)
	}

	r.NoRoute(func(c *gin.Context) { c.AbortWithStatus(http.StatusNotFound) })

	return r
}
