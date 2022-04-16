# Heroku Badger

App that allows you to get a badge image to render in markdown for heroku builds.

## REST API

### GET /status?app_name=<name>

Get status info for a specific app by name. The app must be configured with the heroku-badger add on in Heroku
and have had at least 1 build since the add on was added.

### POST /build-update

Listener for webhook event updates around build changes in heroku.
