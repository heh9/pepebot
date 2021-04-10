<pre align="center">
                           ____        __ 
    ____  ___  ____  ___  / __ )____  / /_
   / __ \/ _ \/ __ \/ _ \/ __  / __ \/ __/
  / /_/ /  __/ /_/ /  __/ /_/ / /_/ / /_  
 / .___/\___/ .___/\___/_____/\____/\__/  
/_/        /_/                            
</pre>

### This project is under construction.

## Pull docker imgae
```bash
$ docker pull mrjoshlab/pepebot:latest
```

## Configurations
* We're using hcl (Hashicorp Config Language)
```hcl
discord {
  token = "<discord_bot_token_here>"
}

sounds {
  type  = "filesystem"
  path  = "./sounds"
  win   = ["gta", "win"]
  loss  = ["loss"]
  runes = ["runes"]
}

db {
  type = "sqlite3"
  path = "./db/data/pepebot.db"
}

steam {
  web_api_token = "<steam_webapi_token_here>"
  # For getting match summeries from valve official apis
}
```

## Using docker image
```bash
docker run -dp 9001:9001 --volume $PWD/db/data:/data --volume $PWD/config.hcl:/config/config.hcl mrjoshlab/pepebot:latest --config-file=/config/config.hcl
```

## docker compose example
```yaml
version: '3'

services:

  bot:
    container_name: pepebot
    image: mrjoshlab/pepebot:latest
    restart: always
    command: ["--config-file", "/config/config.hcl"]
    ports:
      - 9001:9001
    volumes:
      - $PWD/db/data:/data
      - $PWD/config.hcl:/config/config.hcl
```

## Pull git repo
```bash
$ git clone https://github.com/mrjosh/pepe.bot.git
```

## Build the docker image
```bash
$ docker-compose build
```

## Run the docker image
```bash
$ docker-compose up -d
```

### Ez :)
