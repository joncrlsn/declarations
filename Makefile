GOOGLE_PROJECT_ID = declarations-34xd
REGION = us-central1
GAR_IMAGE = $(REGION)-docker.pkg.dev/$(GOOGLE_PROJECT_ID)/docker/declarations:latest 

build:
	@go build -o declarations .

docker-build:
	echo "GAR_IMAGE=$(GAR_IMAGE)"
	@docker build -t $(GAR_IMAGE) .

docker-push: docker-build
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
	  --port=8080 \
	  --max-instances=10 \
	  --cpu=.5 \
	  --memory=256Mi
