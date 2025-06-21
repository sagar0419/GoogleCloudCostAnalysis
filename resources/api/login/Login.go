package login

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
)

type Credentials struct {
	ProjectID string `json:"project_id"`
}

func Login() {

	CredsFile, _ := base64.StdEncoding.DecodeString(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	if len(CredsFile) == 0 {
		log.Fatal("GOOGLE_CLOUD_PROJECTID or GOOGLE_APPLICATION_CREDENTIALS environment variable not set")
	}
	var creds Credentials

	if err := json.NewDecoder(bytes.NewReader(CredsFile)).Decode(&creds); err != nil {
		log.Fatalf("Failed to decode credentials from JSON: %v\n", err)
	}
}
