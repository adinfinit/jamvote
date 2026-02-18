.PHONY: deploy-production deploy-staging format check lint run emulator

deploy-production: format check lint
	gcloud app deploy --project=apt-vote app.yaml

deploy-testing: format check lint
	gcloud app deploy --project=apt-vote app.test.yaml

format:
	goimports -w .

check:
	go vet ./...

lint:
	staticcheck ./...

emulator:
	@bash -c ' \
		trap "echo \"\nCaught Ctrl-C. Killing Datastore emulator...\"; lsof -ti:8081 | xargs kill -9 2>/dev/null; exit 0" SIGINT SIGTERM EXIT; \
		gcloud beta emulators datastore start --project=local-dev --host-port=localhost:8081 --no-store-on-disk & \
		wait $$! \
	'
run:
	DEVELOPMENT=1 GOOGLE_CLOUD_PROJECT=local-dev DATASTORE_EMULATOR_HOST=localhost:8081 go run main.go
