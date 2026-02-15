.PHONY: deploy-live deploy-staging format check lint

deploy-live: format check lint
	gcloud app deploy --project=jamvote

deploy-testing: format check lint
	gcloud app deploy --project=jamvote-testing

format:
	goimports -w .

check:
	go vet ./...

lint:
	staticcheck ./...
