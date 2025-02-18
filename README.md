# arr-telegram-bot

A Telegram bot that listens to incoming messages and finds movies and shows
from IMDB links to add in Radarr and Sonarr.

The bot itself is silent, so I suggest you use the same token in Radarr and
Sonarr as well so that you get notifications when shows and movies are
added.

## Running the bot

I build a Docker image you can use.

```sh
docker run --rm tommasoamici/arr-telegram-bot:latest
```

### Configure environment

You'll need to provide the following variables to the container.

- `TELEGRAM_TOKEN`: ask BotFather for one
- `RADARR_HOST`
- `RADARR_TOKEN`
- `RADARR_ROOT`: the directory where you store your media
- `RADARR_QUALITY_PROFILE`: defaults to "HD - 720p/1080p"
- `SONARR_HOST`
- `SONARR_TOKEN`
- `SONARR_ROOT`: the directory where you store your media
- `SONARR_QUALITY_PROFILE`: defaults to "HD - 720p/1080p"
