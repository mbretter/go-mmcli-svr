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
          "$GPGSV,3,1,11,02,52,298,36,03,16,233,28,08,59,191,24,10,33,053,30,1*66",
          "$GPGSV,3,2,11,14,16,307,36,21,71,309,37,22,09,327,21,32,44,087,20,1*6D",
          "$GPGSV,3,3,11,23,00,054,,24,,,,27,29,160,,1*6B",
          "$GPRMC,134442.00,A,4707.382315,N,01517.446408,E,0.0,196.6,250524,1.4,E,A,V*40",
          "$GPGSA,A,2,02,03,10,14,21,,,,,,,,1.7,1.4,0.9,1*22",
          "$GPVTG,196.6,T,195.2,M,0.0,N,0.0,K,A*24",
          "$GPGGA,134443.00,4707.382314,N,01517.446361,E,1,05,1.4,425.4,M,43.8,M,,*6A"
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