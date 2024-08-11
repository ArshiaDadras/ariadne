import json


try:
    f = open('data/edges.json')
    data = json.load(f)
    f.close()
except Exception as e:
    print(e)
    exit(1)

geojson = {
    "type": "FeatureCollection",
    "features": []
}

for index, edge in enumerate(data):
    print(index)
    geojson['features'].append({
        "type": "Feature",
        "properties": {
            'ID': edge['id'],
            # 'stroke': '#FF0000' if index % 2 else '#0000FF',
        },
        "geometry": {
            "type": "LineString",
            "coordinates": [[poly['longitude'], poly['latitude']] for poly in edge['polygon']]
        }
    })

with open("data/result.geojson", "w") as f:
    json.dump(geojson, f, indent=2)