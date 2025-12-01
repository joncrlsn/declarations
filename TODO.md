## Development Tasks
- [x] Test that the application can parse the declarations.txt file properly
- [x] Create the Golang API
- [x] Create the web client
- [x] Test that the web client can connect to the Golang API

## Google Cloud Deployment Tasks
- [x] Build a Dockerfile that creates a container that runs on Intel
- [x] Store ESV api token in Google Secret Manager
- [x] During gcloud run deploy, set ESV_API_TOKEN to the secret value from Google Secret Manager (GSM)
- [ ] When running in gcloud, retrieve declarations from Google Cloud Storage, otherwise use local file
- [ ] Create a script to ...
      - Create Google Cloud Storage bucket for declarations file
      - Upload initial declarations.txt file to GCS bucket
- [ ] Create a script to ...
      -  Set up Google Cloud Run service configuration
      - Configure IAM roles for Cloud Run to access Cloud Storage
      - Create cloudbuild.yaml for automated deployment
- [ ] Set up environment variables in Cloud Run (GCS_BUCKET_NAME, etc.)
- [ ] Test deployment in Google Cloud Run
- [ ] Configure custom domain (optional)