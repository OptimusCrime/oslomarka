# Data

GML source downloaded from Geonorge (Turrutebasen):
https://kartkatalog.geonorge.no/metadata/turrutebasen/d1422d17-6d95-4ef1-96ab-8af31744dd63

- Format: GML
- Projection: EUREF89 UTM sone 33 (EPSG:25833)

## Generating routes.geojson

`parse_routes.py` reads the GML file and writes `routes.geojson` (WGS84 / EPSG:4326).

### Setup

```bash
bash setup.sh           # creates .venv and installs dependencies
source .venv/bin/activate
```

### Run

```bash
python parse_routes.py
```
