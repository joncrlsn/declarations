## Development Tasks
- [ ] Test that the application can parse the declarations.txt file properly
- [ ] Create the Golang API
- [ ] Create the web client
- [ ] Test that the web client can connect to the Golang API

## Google Cloud Deployment Tasks
- [ ] Build a Dockerfile that creates a container that runs on Intel
- [ ] Add Google Cloud Storage integration to replace local file storage
- [ ] Add environment variable support for GCS bucket name and credentials
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
- [ ] Set up monitoring and logging