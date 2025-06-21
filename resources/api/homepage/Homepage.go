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

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Credentials struct {
	ProjectID string `json:"project_id"`
}

// Struct for billing record
type BillingRecord struct {
	ResourceName   string  `json:"resource_name,omitempty"`
	Service        string  `json:"service"`
	SKU            string  `json:"sku"`
	CostUSD        float64 `json:"cost_usd"`
	UsageStartTime string  `json:"usage_start_time"`
	UsageEndTime   string  `json:"usage_end_time"`
}

func FetchCost(c *gin.Context) {
	ctx := context.Background()

	// Decode base64 credentials
	CredsFile, _ := base64.StdEncoding.DecodeString(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if len(CredsFile) == 0 {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS environment variable not set or empty")
	}

	// Extract project ID from credentials
	var creds Credentials
	if err := json.NewDecoder(bytes.NewReader(CredsFile)).Decode(&creds); err != nil {
		log.Fatalf("Failed to decode credentials JSON: %v", err)
	}

	// Connect to BigQuery using credentials
	bqClient, err := bigquery.NewClient(ctx, creds.ProjectID, option.WithCredentialsJSON(CredsFile))
	if err != nil {
		log.Fatalf("Failed to create BigQuery client: %v", err)
	}
	defer bqClient.Close()

	// Replace with your billing dataset name
	dataset := "your_billing_dataset" // TODO: Replace with actual dataset name

	// BigQuery query
	query := fmt.Sprintf(`
		SELECT
			resource.name AS resource_name,
			service.description AS service,
			sku.description AS sku,
			ROUND(SUM(cost), 2) AS cost_usd,
			usage_start_time,
			usage_end_time
		FROM
			`+"`%s.%s.gcp_billing_export_v1_*`"+`
		WHERE
			usage_start_time >= TIMESTAMP_TRUNC(CURRENT_TIMESTAMP(), DAY)
		GROUP BY
			resource_name, service, sku, usage_start_time, usage_end_time
		ORDER BY
			cost_usd DESC
		LIMIT 100
	`, creds.ProjectID, dataset)

	it, err := bqClient.Query(query).Read(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute BigQuery query"})
		return
	}

	var records []BillingRecord

	for {
		var row struct {
			ResourceName   bigquery.NullString    `bigquery:"resource_name"`
			Service        bigquery.NullString    `bigquery:"service"`
			SKU            bigquery.NullString    `bigquery:"sku"`
			CostUSD        float64                `bigquery:"cost_usd"`
			UsageStartTime bigquery.NullTimestamp `bigquery:"usage_start_time"`
			UsageEndTime   bigquery.NullTimestamp `bigquery:"usage_end_time"`
		}

		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read query result"})
			return
		}

		records = append(records, BillingRecord{
			ResourceName:   row.ResourceName.StringVal,
			Service:        row.Service.StringVal,
			SKU:            row.SKU.StringVal,
			CostUSD:        row.CostUSD,
			UsageStartTime: row.UsageStartTime.Timestamp.String(),
			UsageEndTime:   row.UsageEndTime.Timestamp.String(),
		})
	}

	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"project_id": creds.ProjectID,
		"records":    records,
	})
}
