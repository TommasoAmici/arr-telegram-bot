import os
from telegram import Update
from telegram.ext import Application, MessageHandler, filters
from telegram.ext import ContextTypes
import arrapi
import logging

from src.lib.imdb import IMDB
from src.lib.tvdb import TVDB

logging.basicConfig(
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    level=logging.WARN,
)

TELEGRAM_TOKEN = os.environ["TELEGRAM_TOKEN"]
RADARR_HOST = os.environ["RADARR_HOST"]
RADARR_TOKEN = os.environ["RADARR_TOKEN"]
RADARR_ROOT = os.environ["RADARR_ROOT"]
RADARR_QUALITY_PROFILE = os.getenv("RADARR_QUALITY_PROFILE", "HD - 720p/1080p")
SONARR_HOST = os.environ["SONARR_HOST"]
SONARR_TOKEN = os.environ["SONARR_TOKEN"]
SONARR_ROOT = os.environ["SONARR_ROOT"]
SONARR_QUALITY_PROFILE = os.getenv("SONARR_QUALITY_PROFILE", "HD - 720p/1080p")

radarr = arrapi.RadarrAPI(RADARR_HOST, RADARR_TOKEN)
sonarr = arrapi.SonarrAPI(SONARR_HOST, SONARR_TOKEN)


async def add_to_library(update: Update, context: ContextTypes.DEFAULT_TYPE):
    if not update.message:
        return

    for entity in update.message.entities:
        url: str | None
        if entity.type == entity.URL and update.message.text:
            url = update.message.text[entity.offset : entity.offset + entity.length]
        else:
            url = entity.url
        if not url:
            continue

        data = IMDB.get(url)
        if not data:
            continue
        if data.is_movie:
            radarr.add_movie(
                RADARR_ROOT,
                quality_profile=RADARR_QUALITY_PROFILE,
                imdb_id=data.imdb_id,
                minimum_availability="released",
            )
        else:
            tvdb_id = TVDB.get_id_by_imdb_id(data.imdb_id)
            if not tvdb_id:
                return
            sonarr.add_series(
                SONARR_ROOT,
                quality_profile=SONARR_QUALITY_PROFILE,
                tvdb_id=tvdb_id,
                language_profile="",
            )


if __name__ == "__main__":
    logging.info("initializing bot")
    application = Application.builder().token(TELEGRAM_TOKEN).build()

    application.add_handler(MessageHandler(filters.ALL, add_to_library))

    logging.info("starting bot")
    application.run_polling(allowed_updates=Update.ALL_TYPES)
