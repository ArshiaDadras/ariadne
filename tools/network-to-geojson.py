import pandas as pd
import json

try:
    df_roads = pd.read_csv('data/road_network.csv', sep='\t')
except Exception as e:
    print(e)
    exit(1)

geojson = {
    'type': 'FeatureCollection',
    'features': []
}

for _, row in df_roads.iterrows():
    line = []
    for location in row['LINESTRING()'].replace('LINESTRING(', '').replace(')', '').split(', '):
        lon, lat = map(float, location.split())
        line.append([lon, lat])

    geojson['features'].append({
        'type': 'Feature',
        'properties': {
            'stroke': '#00FF00',
            'Edge ID': str(row['Edge ID'])
        },
        'geometry': {
            'type': 'LineString',
            'coordinates': line
        }
    })

with open('data/network.geojson', 'w') as f:
    json.dump(geojson, f, indent=2)