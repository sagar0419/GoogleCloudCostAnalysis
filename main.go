package main

import (
	homepage "resources/api/homepage"

	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	homepage "github.com/sagar0419/GoogleCloudCostAnalysis/resources/api/Homepage"
)

func main() {
	fmt.Println("Getting Values")

	// Getting values from OS ENV variables.
	region := os.Getenv("GOOGLE_CLOUD_REGION")
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if region == "" || projectID == "" {
		log.Fatal("Value of variable region or Project ID has not been passed")
	}

	// Setting up router
	router := gin.Default()
	// Trusted Proxy
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.Get("/", homepage.HomePage())
}
