GOOGLE_PROJECT_ID = declarations-34xd
REGION = us-central1
GAR_IMAGE = $(REGION)-docker.pkg.dev/$(GOOGLE_PROJECT_ID)/docker/declarations:latest

run:
	@go run .

docker-build:
	echo "GAR_IMAGE=$(GAR_IMAGE)"
	@docker buildx build --platform linux/amd64 -t $(GAR_IMAGE) .

#docker-push: docker-build
docker-push:
	echo "GAR_IMAGE=$(GAR_IMAGE)"
	@docker push $(GAR_IMAGE)

docker-run: 
	docker run --rm -p 8080 $(GAR_IMAGE)

deploy:
	@gcloud run deploy declarations \
	  --image=$(GAR_IMAGE) \
	  --region=$(REGION) \
	  --platform=managed \
	  --allow-unauthenticated \
	  --update-secrets=ESV_API_TOKEN=esv-api-token:latest \
	  --update-env-vars=DECLARATIONS_BUCKET_NAME=jon-storage-34xd \
	  --port=8080 \
	  --max-instances=10 \
	  --cpu=.5 \
	  --memory=128Mi
	  # --set-env-vars=ESV_API_TOKEN_SECRET=esv-api-token \

copy: 
	cp ~/declarations ./declarations.txt
	#scp *.go 192.168.1.181:/Users/joncarlson/work/declarations
	#scp *.md 192.168.1.181:/Users/joncarlson/work/declarations
	#scp go.* 192.168.1.181:/Users/joncarlson/work/declarations
	#scp Makefile 192.168.1.181:/Users/joncarlson/work/declarations
	#scp Dockerfile 192.168.1.181:/Users/joncarlson/work/declarations

# gsutil versioning set on gs://jon-storage-34xd
# gcloud storage cp ~/declarations  gs://jon-storage-34xd/declarations
