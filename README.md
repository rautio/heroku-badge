# Heroku Badger

App that allows you to get a badge image to render in markdown for heroku builds.

Note: The app must be configured with the heroku-badger add on in Heroku
and have had at least 1 build since the add on was added.

## REST API

### GET /build?app_name=<name>

Get a build badge for a specific app by name.

### GET /status?app_name=<name>

Get full status info for a specific app by name.

### POST /build-update

Listener for webhook event updates around build changes in heroku.
