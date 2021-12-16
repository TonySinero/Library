package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/library/docs"
	"github.com/library/src/config"
	"github.com/library/src/controllers"
	"github.com/library/src/database"
	"github.com/library/src/middlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func main() {
	r := gin.Default()

	database.ConnectDataBase()

	// routes
	r.GET("/library/status", controllers.Status)
	r.POST("library/signIn", controllers.SignIn)

	apiAdmin := r.Group("/library/admin")
	{
		apiAdmin.POST("/", controllers.CreateAdmin)
	}
  {
	apiUser := r.Group("/library/user")
	apiUser.Use(middlewares.AuthorizeJWT())
	{
		apiUser.POST("/", controllers.CreateUsers)
		apiUser.GET("/", controllers.GetUsers)
		apiUser.GET("/:id", controllers.DetailUser)
		apiUser.PATCH("/:id", controllers.UpdateUser)
	}

	apiLending := r.Group("/library/order")
	apiLending.Use(middlewares.AuthorizeJWT())
	{
		apiLending.POST("/", controllers.CreateForm)
		apiLending.GET("/", controllers.GetAllForm)
		apiLending.GET("/:id", controllers.GetForm)
		apiLending.DELETE("/:id", controllers.DeleteForm)
	}

	apiAcceptance := r.Group("/library/accept")
	apiAcceptance.Use(middlewares.AuthorizeJWT())
	{
		apiAcceptance.POST("/", controllers.CreateAccept)
		apiAcceptance.GET("/", controllers.GetAllAccept)
		apiAcceptance.GET("/:id", controllers.GetAccept)
		apiAcceptance.DELETE("/:id", controllers.DeleteAccept)
	}


	apiBook := r.Group("/library/book")
	apiBook.Use(middlewares.AuthorizeJWT())
	{
		apiBook.POST("/", controllers.CreateBook)
		apiBook.GET("/", controllers.FindBooks)
		apiBook.GET("/:id", controllers.DetailBooks)
		apiBook.DELETE("/:id", controllers.DeleteBook)
		apiBook.PATCH("/:id", controllers.UpdateBook)
	}
}
	apiCategory := r.Group("/library/category")
	{
		apiCategory.POST("/", controllers.CreateCategory)
		apiCategory.GET("/", controllers.FindCategory)
	}

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Library API"
	docs.SwaggerInfo.Description = "This is a library API golang."
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = "localhost:8000"
	docs.SwaggerInfo.BasePath = "/library"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(config.Port)
}
