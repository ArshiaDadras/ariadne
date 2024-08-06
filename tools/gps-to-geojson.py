import pandas as pd
import json, math


def distance(lat1, lon1, lat2, lon2):
    lat1 = lat1 * math.pi / 180
    long1 = lon1 * math.pi / 180
    lat2 = lat2 * math.pi / 180
    long2 = lon2 * math.pi / 180

    dlat = lat2 - lat1
    dlong = long2 - long1

    a = math.pow(math.sin(dlat / 2), 2) + math.cos(lat1) * math.cos(lat2) * math.pow(math.sin(dlong / 2), 2)
    c = 2 * math.atan2(math.sqrt(a), math.sqrt(1 - a))
    return c * 6378137


try:
    df_gps = pd.read_csv('data/gps_data.csv', sep='\t')
except Exception as e:
    print(e)
    exit(1)

geojson = {
    'type': 'FeatureCollection',
    'features': []
}

rows = []
for _, row in df_gps.iterrows():
    if len(rows) == 0:
        rows.append((0, row))
    else:
        last_index, last_row = rows[-1]
        if distance(last_row['Latitude'], last_row['Longitude'], row['Latitude'], row['Longitude']) >= 10:
            rows.append((last_index + 1, row))

for index, row in rows:
    geojson['features'].append({
        'type': 'Feature',
        'properties': {
            'marker-color': '#FF0000',
            'Index': str(index),
            'Latitude': row['Latitude'],
            'Longitude': row['Longitude'],
            'Time': row['Date (UTC)'] + ' ' + row['Time (UTC)']
        },
        'geometry': {
            'type': 'Point',
            'coordinates': [row['Longitude'], row['Latitude']]
        }
    })

with open('data/gps.geojson', 'w') as f:
    json.dump(geojson, f, indent=2)