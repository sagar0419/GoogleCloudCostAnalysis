package Homepage

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	billing "cloud.google.com/go/billing/apiv1"
	"cloud.google.com/go/billing/apiv1/billingpb"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type Credentials struct {
	ProjectID string `json:"project_id"`
}

func HomePage(c *gin.Context) {
	ctx := context.Background()

	// Loading creds
	CredsFile, _ := base64.StdEncoding.DecodeString(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if len(CredsFile) == 0 {
		log.Fatal("GOOGLE_CLOUD_PROJECTID or GOOGLE_APPLICATION_CREDENTIALS environment variable not set")
	}

	// Billing client
	client, err := billing.NewCloudBillingClient(ctx, option.WithCredentialsJSON(CredsFile))
	if err != nil {
		log.Fatalf("Failed to create billing client: %v\n", err)
	}
	defer client.Close()

	// Loading Project ID
	var creds Credentials

	if err := json.NewDecoder(bytes.NewReader(CredsFile)).Decode(&creds); err != nil {
		log.Fatalf("Failed to decode credentials from JSON: %v\n", err)
	}
	projectName := fmt.Sprintf("projects/%s", creds.ProjectID)

	// calling API
	billingInfo, err := client.GetProjectBillingInfo(ctx, &billingpb.GetProjectBillingInfoRequest{
		Name: projectName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch billing info"})
		log.Fatalf("Failed to get billing info: %v\n", err)
	}

	// Sending the billing info as JSON response
	c.JSON(http.StatusOK, gin.H{
		"billing_account_name": billingInfo.BillingAccountName,
		"billing_enabled":      billingInfo.BillingEnabled,
		"project_id":           creds.ProjectID,
	})
}
