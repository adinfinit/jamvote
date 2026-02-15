.PHONY: deploy-production deploy-staging format check lint

deploy-production: format check lint
	gcloud app deploy --project=jamvote app.yaml

deploy-testing: format check lint
	gcloud app deploy --project=jamvote-testing app-testing.yaml

format:
	goimports -w .

check:
	go vet ./...

lint:
	staticcheck ./...
