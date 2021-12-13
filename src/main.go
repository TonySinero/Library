package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/library/docs"
	"github.com/library/src/config"
	"github.com/library/src/controllers"
	"github.com/library/src/models"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func main() {
	r := gin.Default()

	models.ConnectDataBase()

	// routes
	r.GET("/library/status", controllers.Status)
	r.POST("library/signIn", controllers.SignIn)

	apiUser := r.Group("/library/user")
	{
		apiUser.POST("/", controllers.CreateUser)
		apiUser.GET("/", controllers.GetUsers)

	}

	apiBook := r.Group("/library/book")
	//apiBook.Use(middlewares.AuthorizeJWT())
	{
		apiBook.POST("/", controllers.CreateBook)
		apiBook.GET("/", controllers.FindBooks)
		apiBook.GET("/:id", controllers.DetailBooks)
		apiBook.DELETE("/:id", controllers.DeleteBook)
		apiBook.PATCH("/:id", controllers.UpdateBook)
	}

	apiCategory := r.Group("/library/category")
	{
		apiCategory.POST("/", controllers.CreateCategory)
		apiCategory.GET("/", controllers.FindCategory)
	}

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Library API"
	docs.SwaggerInfo.Description = "This is a sample server golang."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8000"
	docs.SwaggerInfo.BasePath = "/library"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(config.Port)
}
