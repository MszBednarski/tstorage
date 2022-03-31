import requests

# march 21
_from = 1617193736
# march 22
_to = 1648729736
# 10 sec step
step = 10

metric = "polygon.test"
labels = [{"name": "type", "value": "raw"}]


def rpc(method, *params):
    res = requests.post("http://127.0.0.1:3000",
                        json={
                            "id": 1,
                            "method": method,
                            "params": params
                        })
    return res.json()


print(rpc("ts.InsertRows", [
    {
        "metric": metric,
        "labels": labels,
        "point": {
            "value": 0,
            "timestamp": _from + step
        }
    }]))

print(rpc("ts.Select",
          {
              "metric": metric,
              "labels": labels,
              "start": _from,
              "end": _to
          }
          ))
