import pandas as pd
import json

try:
    df_route = pd.read_csv('data/ground_truth_route.csv', sep='\t')
    df_roads = pd.read_csv('data/road_network.csv', sep='\t')
except Exception as e:
    print(e)
    exit(1)

geojson = {
    'type': 'FeatureCollection',
    'features': []
}
for index, row in df_route.iterrows():
    line = []
    for location in df_roads.loc[row['Edge ID'] == df_roads['Edge ID']]['LINESTRING()'].values[0].replace('LINESTRING(', '').replace(')', '').split(', '):
        lon, lat = map(float, location.split())
        line.append([lon, lat])
    
    geojson['features'].append({
        'type': 'Feature',
        'properties': {
            'stroke': '#FF0000' if index % 2 else '#0000FF',
            'Edge ID': str(row['Edge ID']),
            'Traverse Direction': 'Forward' if row['Traversed From to To'] == 1 else 'Backward'
        },
        'geometry': {
            'type': 'LineString',
            'coordinates': line if row['Traversed From to To'] == 1 else line[::-1]
        }
    })

with open('data/route.geojson', 'w') as f:
    json.dump(geojson, f, indent=2)