import pandas as pd
import json

try:
    df = pd.read_csv('data/gps_data.csv', sep='\t')
except Exception as e:
    print(e)
    exit(1)

df = df.sort_values(by=['Date (UTC)', 'Time (UTC)'])
line = [location for location in zip(df['Longitude'], df['Latitude'])]

geojson = {
    'type': 'FeatureCollection',
    'features': [
        {
            'type': 'Feature',
            'properties': {
                'stroke': '#0000FF',
                'Start': f'{df.iloc[0]["Date (UTC)"]} {df.iloc[0]["Time (UTC)"]}',
                'End': f'{df.iloc[-1]["Date (UTC)"]} {df.iloc[-1]["Time (UTC)"]}',
            },
            'geometry': {
                'type': 'LineString',
                'coordinates': line
            }
        },
        {
            'type': 'Feature',
            'properties': {
                'Position': 'Start',
                'Time': f'{df.iloc[0]["Date (UTC)"]} {df.iloc[0]["Time (UTC)"]}'
            },
            'geometry': {
                'type': 'Point',
                'coordinates': [df.iloc[0]['Longitude'], df.iloc[0]['Latitude']]
            }
        },
        {
            'type': 'Feature',
            'properties': {
                'Position': 'End',
                'Time': f'{df.iloc[-1]["Date (UTC)"]} {df.iloc[-1]["Time (UTC)"]}'
            },
            'geometry': {
                'type': 'Point',
                'coordinates': [df.iloc[-1]['Longitude'], df.iloc[-1]['Latitude']]
            }
        }
    ]
}

with open('data/road.geojson', 'w') as f:
    json.dump(geojson, f, indent=2)