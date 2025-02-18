from dataclasses import dataclass
import logging
from urllib.parse import urlparse
import httpx


@dataclass
class IMDBEntity:
    imdb_id: str
    is_movie: bool


class IMDB:
    __headers = {
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:135.0) Gecko/20100101 Firefox/135.0",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
        "Accept-Language": "en-US,en;q=0.5",
        "Alt-Used": "www.imdb.com",
        "Upgrade-Insecure-Requests": "1",
        "Sec-Fetch-Dest": "document",
        "Sec-Fetch-Mode": "navigate",
        "Sec-Fetch-Site": "same-origin",
        "Sec-GPC": "1",
        "Priority": "u=0, i",
        "Pragma": "no-cache",
        "Cache-Control": "no-cache",
        "referrer": "https://www.imdb.com/",
    }

    @staticmethod
    def __get_id_from_url(url: str):
        for u in url.split("/"):
            if u.startswith("tt"):
                return u

    @staticmethod
    def get(url: str):
        parsed_url = urlparse(url)
        if parsed_url.hostname not in {"www.imdb.com", "imdb.com"}:
            return

        imdb_id = IMDB.__get_id_from_url(url)
        if not imdb_id:
            return

        logging.info("processing %s", imdb_id)

        res = httpx.get(url, headers=IMDB.__headers, follow_redirects=True)
        if res.status_code >= 400:
            return
        is_tv = "episode-guide-text" in res.text
        is_movie = not is_tv

        return IMDBEntity(imdb_id=imdb_id, is_movie=is_movie)
