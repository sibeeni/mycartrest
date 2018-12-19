package main

import (
	Controller "./controllers"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"io"
	"os"
	"log"
)

func SetupRouter() *gin.Engine {
   router := gin.Default()
	 router.Use(cors.Default())

	 //Router URL
	 router.GET("/cart/getProducts", Controller.GetProductFromCart)
	 router.POST("/cart/add", Controller.AddProductToCart)
	
   return router
}

func main() {
  //Access request URL log file
  f, _ := os.Create("gin.log")
  gin.DefaultWriter = io.MultiWriter(f)

	//Error log file
	errlogfile, _ := os.Create("error.log")
	gin.DefaultErrorWriter = io.MultiWriter(errlogfile)

	log.Printf("Initing a new routine...")

	router := SetupRouter()
	router.Run(":3001")
}
