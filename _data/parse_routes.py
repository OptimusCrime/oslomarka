import json
import xml.etree.ElementTree as ET
from dataclasses import dataclass, field
from typing import Optional

from pyproj import Transformer

GML_FILE = "Friluftsliv_03_Oslo_25833_TurOgFriluftsruter_GML.gml"
OUTPUT_FILE = "routes.geojson"

# Coordinates in the GML are EPSG:25833 (UTM Zone 33N / ETRS89).
# GeoJSON requires WGS84 (EPSG:4326), so we reproject on the way out.
TRANSFORMER = Transformer.from_crs("EPSG:25833", "EPSG:4326", always_xy=True)

NS = {
    "gml": "http://www.opengis.net/gml/3.2",
    "app": "http://skjema.geonorge.no/SOSI/produktspesifikasjon/TurOgFriluftsruter/20171210",
}

ROUTE_TYPES = ("Fotrute",)

INFO_ELEMENT = {
    "Fotrute": "FotruteInfo",
}

INFO_ATTR = {
    "Fotrute": "fotruteInfo",
}


@dataclass
class Route:
    id: str
    route_type: str
    name: Optional[str] = None
    route_number: Optional[str] = None
    marked: Optional[str] = None
    surface: Optional[str] = None
    difficulty: Optional[str] = None
    maintainer: Optional[str] = None
    coordinates: list[tuple[float, float]] = field(default_factory=list)  # (lng, lat) WGS84


def _text(element: ET.Element, path: str) -> Optional[str]:
    node = element.find(path, NS)
    return node.text.strip() if node is not None and node.text else None


def parse_poslist(text: str) -> list[tuple[float, float]]:
    """Convert a GML posList (flat easting/northing pairs) to (lng, lat) WGS84 tuples."""
    numbers = [float(n) for n in text.split()]
    coords = []
    for i in range(0, len(numbers) - 1, 2):
        lng, lat = TRANSFORMER.transform(numbers[i], numbers[i + 1])
        coords.append((lng, lat))
    return coords


def parse_route_element(element: ET.Element) -> Optional[Route]:
    tag = element.tag.split("}")[-1]
    if tag not in ROUTE_TYPES:
        return None

    pos_list_node = element.find(".//gml:posList", NS)
    if pos_list_node is None or not pos_list_node.text:
        return None

    info_el = INFO_ELEMENT[tag]
    info_prefix = f".//app:{info_el}/app:"

    return Route(
        id=element.get("{http://www.opengis.net/gml/3.2}id", ""),
        route_type=tag,
        name=_text(element, f"{info_prefix}rutenavn"),
        route_number=_text(element, f"{info_prefix}rutenummer"),
        marked=_text(element, "app:merking"),
        surface=_text(element, "app:ruteFølger"),
        difficulty=_text(element, f"{info_prefix}gradering"),
        maintainer=_text(element, f"{info_prefix}vedlikeholdsansvarlig"),
        coordinates=parse_poslist(pos_list_node.text),
    )


def parse_gml(filepath: str) -> list[Route]:
    print(f"Parsing {filepath} ...")
    tree = ET.parse(filepath)
    root = tree.getroot()

    routes = []
    for member in root.findall("gml:featureMember", NS):
        for child in member:
            route = parse_route_element(child)
            if route:
                routes.append(route)

    print(f"  → {len(routes)} routes loaded")
    return routes


def to_geojson(routes: list[Route]) -> dict:
    features = []
    for r in routes:
        features.append({
            "type": "Feature",
            "geometry": {
                "type": "LineString",
                "coordinates": list(r.coordinates),
            },
            "properties": {
                "id": r.id,
                "route_type": r.route_type,
                "name": r.name,
                "route_number": r.route_number,
                "marked": r.marked,
                "surface": r.surface,
                "difficulty": r.difficulty,
                "maintainer": r.maintainer,
            },
        })
    return {"type": "FeatureCollection", "features": features}


if __name__ == "__main__":
    routes = parse_gml(GML_FILE)

    geojson = to_geojson(routes)
    with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
        json.dump(geojson, f, ensure_ascii=False, indent=2)

    print(f"Wrote {OUTPUT_FILE}")
