.PHONY: deploy-production deploy-staging format check lint run emulator

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

emulator:
	exec gcloud beta emulators datastore start --project=local-dev --host-port=localhost:8081 --no-store-on-disk

run:
	DEVELOPMENT=1 GOOGLE_CLOUD_PROJECT=local-dev DATASTORE_EMULATOR_HOST=localhost:8081 go run main.go
