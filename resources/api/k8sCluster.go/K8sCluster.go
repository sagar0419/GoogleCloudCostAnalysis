package k8scluster

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

// Simulated cluster struct
type Cluster struct {
	Name           string   `json:"name"`
	Location       string   `json:"location"`
	Status         string   `json:"status"`
	MasterVersion  string   `json:"master_version"`
	NodeVersion    string   `json:"node_version"`
	NodeCount      int64    `json:"node_count"`
	Network        string   `json:"network"`
	Subnetwork     string   `json:"subnetwork"`
	Logging        string   `json:"logging_service"`
	Monitoring     string   `json:"monitoring_service"`
	PrivateCluster bool     `json:"private_cluster"`
	NodePools      []string `json:"node_pools"`
}

type Credentials struct {
	ProjectID string `json:"project_id"`
}

func ListCluster(c *gin.Context) {
	ctx := context.Background()

	CredsFile, _ := base64.StdEncoding.DecodeString(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if len(CredsFile) == 0 {
		log.Fatal("GOOGLE_CLOUD_PROJECTID or GOOGLE_APPLICATION_CREDENTIALS environment variable not set")
	}
	var creds Credentials

	if err := json.NewDecoder(bytes.NewReader(CredsFile)).Decode(&creds); err != nil {
		log.Fatalf("Failed to decode credentials from JSON: %v\n", err)
	}

	client, err := container.NewService(ctx, option.WithCredentialsJSON(CredsFile))
	if err != nil {
		log.Fatalf("Failed to create kubernetes client %v", err)
	}

	resp, err := client.Projects.Zones.Clusters.List(creds.ProjectID, "-").Context(ctx).Do()
	if err != nil {
		log.Fatalf("Failed to list all cluster: %v ", err)
	}

	if len(resp.Clusters) == 0 {
		log.Fatalln("No Clusters found in project", creds.ProjectID)
		return
	}

	var output []Cluster
	for _, cl := range resp.Clusters {
		var nodePools []string
		for _, np := range cl.NodePools {
			nodePools = append(nodePools, fmt.Sprintf("%s (initial: %d)", np.Name, np.InitialNodeCount))
		}

		output = append(output, Cluster{
			Name:           cl.Name,
			Location:       cl.Location,
			Status:         cl.Status,
			MasterVersion:  cl.CurrentMasterVersion,
			NodeVersion:    cl.CurrentNodeVersion,
			NodeCount:      cl.CurrentNodeCount,
			Network:        cl.Network,
			Subnetwork:     cl.Subnetwork,
			Logging:        cl.LoggingService,
			Monitoring:     cl.MonitoringService,
			PrivateCluster: cl.PrivateClusterConfig != nil,
			NodePools:      nodePools,
		})
	}
	c.JSON(200, output)
}
