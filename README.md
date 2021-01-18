# Project
This is the successor of my previous project [Dart-O-Mat 3000](https://github.com/patrickhener/dart-o-mat-3000). I chost to rewrite it using go as a backend language and svelte as a frontend framework. Also I spiced it up a little and did a full redesign. It is now called `DaSCR - Board` and will be one of three projects within the `DaSCR range`.

It will handle darts games and keep the score. That's its basic function. Detailed documentation will follow.

# Installation

Right now only linux installation is described here. You might adapt it to Windows or Mac. Just be sure to set the env variables before running.

To install and use you will need a running go environment already setup. Also to build the frontend you will need `node` and `yarn` running.

```bash
git clone https://github.com/dascr/dascr-board
cd dascr-board/frontend
yarn install
cd ..
go mod download
```

Next you will need two .env files to use the provided `Makefile`.

**.env in root folder of dascr-board**:
```bash
API_IP=0.0.0.0
API_PORT=8000
DEBUG=TRUE // or false if you do not want to have a lot of output
```

**.env in frontend folder**
```bash
API_BASE=http://localhost:8000/
API_URL=http://localhost:8000/api
WS_URL=ws://localhost:8000/ws
```

This needs to point to the API accordingly and it should use the ip addresses of the deployment scenarion already (also see Deployment below).

You then need to run `make generate` once to embed the static files using parcello.

Then you can either run two terminals with
```bash
make run-dev-backend
```
and
```bash
make run-dev-frontend
```
to run the development version right off the bat or you can build the project with `make build-all` to build the backend executable and to bundle the release of the frontend.

# Usage
When running you need to navigate your browser to `http://ip:port` of the frontend. Basically everything there is explained in detail. But in short you need to create player and then you need two browser windows. One is pointing at `http://ip:port/<gameid>/start`. There the scoreboard will be shown after starting a game. To start a game and input data you point your browser to `http://ip:port/<gameid>/game`.

# Deployment
As of covid-19 I use my scoreboard to host a remote game once a week and therefore deployed it to my hosting server. This is the way I did it.

I used [caddy server v2](https://caddyserver.com/) to host the frontend. Also I used a systemd service file to run the backend service in the background.

**/etc/caddy/Caddyfile**
```bash
example.com {
	root * /var/www/dascr
	encode gzip

	handle /api/* {
		reverse_proxy localhost:8000
	}

	handle /images/* {
		reverse_proxy localhost:8000
	}

	handle /uploads/* {
		reverse_proxy localhost:8000
	}

	handle /ws/* {
		reverse_proxy localhost:8000
	}

	handle {
		try_files {path} {file} /index.html
		file_server
	}

	header {
	# enable HSTS
	Strict-Transport-Security max-age=31536000;

	# disable clients from sniffing the media type
	X-Content-Type-Options nosniff

	# clickjacking protection
	X-Frame-Options DENY

	# keep referrer data off of HTTP connections
	Referrer-Policy no-referrer-when-downgrade
	}

	log {
		output file /var/log/caddy/example.com.access.log {
			roll_size 1gb
			roll_keep 5
			roll_keep_for 720h
		}
	}
}
```

`var/www/dascr` in this case points to `./frontend/public` of the project folder after building it.


**/etc/systemd/system/dascr-board**
```bash
[Unit]
Description=DaSCR Board - Backend API
After=network.target

[Service]
Type=simple
User=dascr
Group=dascr
Restart=always
RestartSec=5s

Environment=API_IP=0.0.0.0
Environment=API_PORT=8000

WorkingDirectory=/var/lib/dascr
ExecStart=/var/lib/dascr/dascr-board
SyslogIdentifier=dascr-board

[Install]
WantedBy=multi-user.target
```

`/var/lib/dascr-board` is the executable resulting from `make build-all` found in `./dist` folder.

When building my frontend I made sure to have my .env to point to my domain instead of a local ip address. This way the clients browser later knows where to fetch the data from API:

**.env in frontend folder**
```bash
API_BASE=https://example.com/
API_URL=https://example.co/api
WS_URL=wss://example.com/ws
```

Also make sure to choose the right protocol here. Caddy server automatically uses https and therefore also *wss* is used instead of *ws*.

# API

The API has a few endpoints.

```bash
[*] Starting Backend Development
DEBUG  [2021-01-18 11:16:19] All routes are
DEBUG  [2021-01-18 11:16:19] GET /api/
DEBUG  [2021-01-18 11:16:19] GET /api/debug/{id}/redirect
DEBUG  [2021-01-18 11:16:19] GET /api/debug/{id}/update
DEBUG  [2021-01-18 11:16:19] GET /api/game/
DEBUG  [2021-01-18 11:16:19] GET /api/game/{id}
DEBUG  [2021-01-18 11:16:19] POST /api/game/{id}
DEBUG  [2021-01-18 11:16:19] DELETE /api/game/{id}
DEBUG  [2021-01-18 11:16:19] GET /api/game/{id}/display
DEBUG  [2021-01-18 11:16:19] POST /api/game/{id}/nextPlayer
DEBUG  [2021-01-18 11:16:19] POST /api/game/{id}/rematch
DEBUG  [2021-01-18 11:16:19] POST /api/game/{id}/throw/{number}/{modifier}
DEBUG  [2021-01-18 11:16:19] POST /api/game/{id}/undo
DEBUG  [2021-01-18 11:16:19] POST /api/player/
DEBUG  [2021-01-18 11:16:19] GET /api/player/
DEBUG  [2021-01-18 11:16:19] GET /api/player/{id}
DEBUG  [2021-01-18 11:16:19] PATCH /api/player/{id}
DEBUG  [2021-01-18 11:16:19] DELETE /api/player/{id}
DEBUG  [2021-01-18 11:16:19] POST /api/player/{id}/image
DEBUG  [2021-01-18 11:16:19] HEAD /images/*
DEBUG  [2021-01-18 11:16:19] PUT /images/*
DEBUG  [2021-01-18 11:16:19] POST /images/*
DEBUG  [2021-01-18 11:16:19] CONNECT /images/*
DEBUG  [2021-01-18 11:16:19] TRACE /images/*
DEBUG  [2021-01-18 11:16:19] PATCH /images/*
DEBUG  [2021-01-18 11:16:19] GET /images/*
DEBUG  [2021-01-18 11:16:19] DELETE /images/*
DEBUG  [2021-01-18 11:16:19] OPTIONS /images/*
DEBUG  [2021-01-18 11:16:19] PATCH /uploads/*
DEBUG  [2021-01-18 11:16:19] PUT /uploads/*
DEBUG  [2021-01-18 11:16:19] CONNECT /uploads/*
DEBUG  [2021-01-18 11:16:19] HEAD /uploads/*
DEBUG  [2021-01-18 11:16:19] GET /uploads/*
DEBUG  [2021-01-18 11:16:19] TRACE /uploads/*
DEBUG  [2021-01-18 11:16:19] OPTIONS /uploads/*
DEBUG  [2021-01-18 11:16:19] DELETE /uploads/*
DEBUG  [2021-01-18 11:16:19] POST /uploads/*
DEBUG  [2021-01-18 11:16:19] GET /ws/{id}
INFO   [2021-01-18 11:16:19] Starting API at: 0.0.0.0:8000
```

Those are basically all endpoints (you can read them when starting with DEBUG=TRUE in the console). The most important ones are `POST /api/game/{id}/nextPlayer` and `POST /api/game/{id}/throw/{number}/{modifier}`. Those are the endpoints a recognition software should send to (and will after finishing dascr-machine and dascr-cam).

# Screenshots

Here are a few screenshots of the games and the UI.

## Setup

Player Setup

![Player Setup](https://raw.githubusercontent.com/patrickhener/image-cdn/main/dascr-board-player.png)

Start Page

![Start](https://raw.githubusercontent.com/patrickhener/image-cdn/main/dascr-board-start.png)

Game Setup

![Game Setup](https://raw.githubusercontent.com/patrickhener/image-cdn/main/dascr-board-game.png)

## X01

Scoreboard

![](https://raw.githubusercontent.com/patrickhener/image-cdn/main/dascr-board-x01-board.png)

Controller

![](https://raw.githubusercontent.com/patrickhener/image-cdn/main/dascr-board-x01-controller.png)

## Cricket

Scoreboard

![](https://raw.githubusercontent.com/patrickhener/image-cdn/main/dascr-board-cricket-board.png)

Controller

![](https://raw.githubusercontent.com/patrickhener/image-cdn/main/dascr-board-cricket-controller.png)

# Roadmap

Right now I am missing a few things I planned on.

* More games (Around the clock, Split-Score, Highscore, Elimination)
* Sound