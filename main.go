package main

import (
	Homepage "GoogleCloudCostAnalysis/resources/api/homepage"
	k8sclustergo "GoogleCloudCostAnalysis/resources/api/k8sCluster.go"
	"GoogleCloudCostAnalysis/resources/api/login"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Getting Values")

	router := gin.Default()
	login.Login()

	// Trusted Proxy
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Homepage
	router.GET("/", Homepage.FetchCost)
	router.GET("/k8sclusters", k8sclustergo.ListCluster)
	// router.GET("/login", login.Login)

	// Server
	err := router.Run(":3000")
	if err != nil {
		panic(err)
	}
}
