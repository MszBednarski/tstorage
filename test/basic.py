import requests
import json
import sys

# march 21
_from = 1617193736
# march 22
_to = 1648729736
# 10 sec step
step = 10

metric = "polygon.test"
labels = [{"name": "type", "value": "raw"}]

# data to send
# rows = []
# for i in range(_from, _to, step):
#     rows = [{"metric": metric, "labels": labels, "point": {
#         "value": i - _from,
#         "timestamp": i
#     }}]
#     res = requests.post("http://127.0.0.1:3000/put", json={"rows": rows})
    # print(res.text)

# print("number of datapoints", len(rows))

res = requests.post("http://127.0.0.1:3000/get", json={
    "metric": metric,
    "labels": labels,
    "start": _from,
    "end": _to
})

values = res.json()
print("response size", len(values["points"]))
