import json


try:
    f = open('data/nodes.json')
    nodes = json.load(f)
    f.close()
except Exception as e:
    print(e)
    exit(1)

geojson = {
    'type': 'FeatureCollection',
    'features': []
}

for node in nodes:
    geojson['features'].append({
        'type': 'Feature',
        'properties': {
            # 'marker-color': '#FF0000',
            # 'marker-color': '#0000FF',
            'marker-color': '#00FF00',
            'marker-size': 'small',
            'marker-symbol': 'circle',
            'Node ID': str(node['id'])
        },
        'geometry': {
            'type': 'Point',
            'coordinates': [node['position']['longitude'], node['position']['latitude']]
        }
    })

with open('data/nodes.geojson', 'w') as f:
    json.dump(geojson, f, indent=2)