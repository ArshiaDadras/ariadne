import sys
import math
import json


def dist(lon1, lat1, lon2, lat2):
    lat1 = lat1 * math.pi / 180
    lon1 = lon1 * math.pi / 180
    lat2 = lat2 * math.pi / 180
    lon2 = lon2 * math.pi / 180

    dlat = lat2 - lat1
    dlon = lon2 - lon1

    a = math.pow(math.sin(dlat/2), 2) + math.cos(lat1) * math.cos(lat2) * math.pow(math.sin(dlon/2), 2)
    c = 2 * math.atan2(math.sqrt(a), math.sqrt(1-a))
    return c * 6378137

def d_up(lon, lat, d):
    return lon, lat + (180 * d / (math.pi * 6378137))

def d_right(lon, lat, d):
    return lon + (180 * d / (math.pi * 6378137 * math.cos(lat * math.pi / 180))), lat

def move(lon, lat, dx, dy):
    lon, lat = d_up(lon, lat, dy)
    lon, lat = d_right(lon, lat, dx)
    return lon, lat

def circ(lon, lat, d):
    cur_lon = lon
    while dist(lon, lat, cur_lon, lat) <= d:
        cur_lon += radius
    cur_lon -= radius

    upper_circle = []
    while dist(lon, lat, cur_lon, lat) <= d:
        L, R = lat, 90
        while R - L > radius:
            M = (L + R) / 2
            if dist(lon, lat, cur_lon, M) <= d:
                L = M
            else:
                R = M

        upper_circle.append((cur_lon, L))
        cur_lon -= radius
    cur_lon += radius

    lower_circle = []
    while dist(lon, lat, cur_lon, lat) <= d:
        L, R = lat, -90
        while L - R > radius:
            M = (L + R) / 2
            if dist(lon, lat, cur_lon, M) <= d:
                L = M
            else:
                R = M

        lower_circle.append((cur_lon, L))
        cur_lon += radius

    return upper_circle + lower_circle + [upper_circle[0]]

def rect(lon, lat, d):
    lon1, lat1 = move(lon, lat, -d, -d)
    lon2, lat2 = move(lon, lat, d, d)

    cur = lat1
    upper_polygon = []
    lower_polygon = []
    while cur < lat2:
        upper_polygon.append((lon1, cur))
        lower_polygon.append((lon2, lat2 - cur + lat1))
        cur += radius
    cur = lon1
    while cur < lon2:
        upper_polygon.append((cur, lat2))
        lower_polygon.append((lon2 - cur + lon1, lat1))
        cur += radius

    return upper_polygon + lower_polygon


if len(sys.argv) < 4:
    print("Usage: python d-poly.py lon lat d [radius]")
    sys.exit(1)
lon, lat, d = map(float, sys.argv[1:4])
radius = 1e-5 if len(sys.argv) < 5 else float(sys.argv[4])

polygon = circ(lon, lat, d)
print('Total Length:', len(polygon))

geojson = {
    "type": "FeatureCollection",
    "features": [
        {
            "type": "Feature",
            "properties": {},
            "geometry": {
                "type": "LineString",
                "coordinates": polygon
            }
        }
    ]
}

with open("data/polygon.geojson", "w") as f:
    json.dump(geojson, f, indent=2)