# GoogleCloudCostAnalysis

### To run this program you need to pass the project ID and right credentials to access your GCP account.

To run this code you need to encode your credntails file. To encode the JSON as base64 string and pass it as the environment variable

Run the below command to print the base64 string of the credentials file
`bash> cat /home/me/gcp-creds.json | base64`

Then set a new environment variable
`export GCP_CREDS_JSON_BASE64="paste_base64_output_here"`