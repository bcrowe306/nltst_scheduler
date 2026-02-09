# NLTST Scheduler
A web app to help with scheduling tasks in the church. May eventually grow into complete CHMS.

## Dependencies and tools
**Go**: Go language v1.25+
> brew install go

**Templ**: Templating system for go
> go get -tool github.com/a-h/templ/cmd/templ@latest

**Tailwindcss**: CSS styling
> npm install tailwindcss @tailwindcss/cli

**Go Air**: Hot reloading for go web projects
> brew install go-air

## Setting up project
1. Install project dependencies
> go mod tidy
> npn install

2. Create .env file

Sample
```bash
MONGODB_URI=
MONGODB_DATABASE=nltst_scheduler
ADMIN_EMAIL=admin@nltst.com
ADMIN_PASSWORD=

# Twilio
TWILIO_ACCOUNT_SID=
TWILIO_AUTH_TOKEN=
TWILIO_FROM_NUMBER=

# ClickSend Credentials
CLICK_SEND_API_KEY=
CLICK_SEND_USERNAME=
CLICK_SEND_FROM_NUMBER=

# Sendgrid
SENDGRID_API_KEY=

# App Port
PORT=8080
```

## Running project
To run the project, perfrom the steps listed above. Then simply run air command.
> air

This will launch go-air and watch for file changes preconfigured in .air.toml.
Air-Go is configured as a reverse proxy to enable triggering of hot-reload. While the APP runs of port 8080 by default, the proxy runs on 8081. To load the website navigate to the proxy port:
>http://localhost:8081

When a change is made to files in the working directory, the project with trigger a hot-reload.