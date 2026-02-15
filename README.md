# JamVote

JamVote is a web application for managing game jams and voting.

## Local Development

### Prerequisites

- [Go](https://go.dev/dl/)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
- Java JRE 11+ (required by the Datastore emulator)

Install the Datastore emulator component:

```
gcloud components install cloud-datastore-emulator
```

On macOS you can install Java with:

```
brew install openjdk
```

On Debian/Ubuntu:

```
sudo apt install default-jre
```

### Running locally

1. Start the Datastore emulator in one terminal:

```
make emulator
```

2. Start the app in another terminal:

```
make run
```

3. Open http://localhost:8080

The app starts in development mode with seed data: 50 users, 12 events across all lifecycle stages (registration, voting, closed, revealed), each with 10 teams.

Use the "Development Login" form to log in. Enter "Admin" to log in as an admin user, or any of the seeded user names.

## Deployment

```
make deploy-production
make deploy-testing
```
