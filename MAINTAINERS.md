# How to release

Login with gcloud, if you already haven't:

```
gcloud auth login
```

Set your project with (for APT it's `apt-vote`):

```
gcloud config set project PROJECT_ID
```

Then deploy:

```
gcloud app deploy
```