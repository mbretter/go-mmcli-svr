[![](https://github.com/mbretter/go-mmcli-svr/actions/workflows/go.yml/badge.svg)](https://github.com/mbretter/go-mmcli-svr/actions/workflows/go.yml)
[![](https://goreportcard.com/badge/mbretter/go-mmcli-svr)](https://goreportcard.com/report/mbretter/go-mmcli-svr "Go Report Card")
[![codecov](https://codecov.io/gh/mbretter/go-mmcli-svr/graph/badge.svg?token=YMBMKY7W9X)](https://codecov.io/gh/mbretter/go-mmcli-svr)

mmcli-svr provides a http/api for accessing ModemManager.

It is possible to send SMS and/or get location information (GNSS). 
It can be easily integrated into home automation systems like HomeAssistant.

## commandline

For a list of available commandline options invoke `mmcli-srv -h`
```
Usage of ./mmcli-srv:
  -gps-refresh int
        gps refresh rate in seconds
  -listen string
        listen: <ip:port|:port> (default "127.0.0.1:8743")
  -location-enable string
        enable location gathering: <all|gps-nmea|gps-raw|3gpp|agps‐msa|agps‐msb>
```
The listen option changes the listening ip/port, by default the service runs on localhost only, due to security reasons.
There is no authentication implemented, anybody with network access could use the service. 
If you want authentication, put the service behind a reverse proxy like nginx.

For a detailed description of the ModemManager related options, see `man mmcli`

`./mmcli-srv -location-enable=gps-raw,gps-nmea -gps-refresh=5`

## api docs

You can access the api docs using the included openapi documentation: http://127.0.0.1:8743/swagger/index.html

## examples

### get location

```
curl -X 'GET' 'http://127.0.0.1:8743/location' -H 'accept: application/json'
```

Response body:
```
{
  "modem": {
    "location": {
      "3gpp": {
        "cid": "000FA908",
        "lac": "FFFE",
        "mcc": "232",
        "mnc": "01",
        "tac": "00003D"
      },
      "cdma-bs": {
        "latitude": "--",
        "longitude": "--"
      },
      "gps": {
        "altitude": "425,400000",
        "latitude": "37,123039",
        "longitude": "5,290773",
        "nmea": [         
        ],
        "utc": "134443.00"
      }
    }
  }
}
```

### send SMS

```
curl -X 'POST' \
  'http://127.0.0.1:8743/sms' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "number": "+436641234567",
  "text": "Ping"
}'
```

Response body:
```
{
  "message": "successfully sent the SMS"
}
```
