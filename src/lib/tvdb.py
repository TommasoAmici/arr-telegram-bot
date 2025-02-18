import httpx

import xml.etree.ElementTree as ET


class TVDB:
    @staticmethod
    def get_id_by_imdb_id(imdb_id: str):
        res = httpx.get(
            f"https://thetvdb.com/api/GetSeriesByRemoteID.php?imdbid={imdb_id}"
        )
        if res.status_code >= 400:
            return
        root = ET.fromstring(res.text)
        id_node = root.find(".//id")
        if id_node is None:
            return
        series_id = id_node.text
        return series_id
