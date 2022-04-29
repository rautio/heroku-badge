# Heroku Badger

App that allows you to get a badge image to render in markdown for heroku builds.

Note: The app must be configured with the heroku-badger add on in Heroku
and have had at least 1 build since the add on was added.

## Getting Started

### 1. Configure your Heroku app

Within the Heroku administration for your app navigate tot the webhooks page. From there create a new webhook with:

1. Payload URL: `https://heroku-badger.herokuapp.com/build-update`
2. Select the `api:build` and `api:release` events.

Now your app will send webhook messages to the heroku badger app whenever a new release or build is triggered.

### 2. Add the badge in your repo's readme

Copy/Paste the following and replace `<app_name>` with your apps unique name:

```
[![Build Status](https://heroku-badger.herokuapp.com/build?app_name=<app_name>)](https://heroku-badger.herokuapp.com/build?app_name=<app_name>)
```

## REST API

### GET /build?app_name=<name>

Get a build badge for a specific app by name.

### GET /status?app_name=<name>

Get full status info for a specific app by name.

### POST /build-update

Listener for webhook event updates around build changes in heroku.
