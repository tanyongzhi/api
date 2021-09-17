Room
=====

GET /room/occupancy/
----------

Returns remaining spaces for all rooms

Response format:
```
[
  {
    "roomId": "testroomid",
    "remainingSpaces": 4,
    "maxCapacity": 4
  },
  {
    "roomId": "ROOM_SIEBEL_1",
    "remainingSpaces": 2,
    "maxCapacity": 4
  },
  {
    "roomId": "ROOM_SIEBEL_2",
    "remainingSpaces": 4,
    "maxCapacity": 4
  },
  {
    "roomId": "ROOM_ECEB_1",
    "remainingSpaces": 4,
    "maxCapacity": 4
  },
  {
    "roomId": "ROOM_ECEB_2",
    "remainingSpaces": 10,
    "maxCapacity": 20
  }
]
```

GET /room/ROOMID/
----------------------

Returns remaining spaces for room with `ROOMID`.

Response format:
```
{
  "roomId": "ROOM_ECEB_1",
  "remainingSpaces": 4,
  "maxCapacity": 4
}
```

POST /room/update/

Updates the remaining spaces of a room specified in the roomId field of the request. Note that a negative number for numPeople represents people leaving the room (ie. the number of remaining slots increasing). A positive number for means the opposite.

Request format:
```
{
  "roomId": "ROOM_ECEB_2",
  "numPeople": -10
}
```

Response format:
```
{
  "roomId": "ROOM_ECEB_2",
  "remainingSpaces": 20,
  "maxCapacity": 20
}
```

This microservice also supports integration with Prometheus, which allows system and custom metrics collection. As a demo, we have instrumented Prometheus and the room microservice to record the room occupancy over time.

To view the metrics on Prometheus, navigate to the `/config` directory and run 

```
prometheus --config.file=prometheus.yml
```
This is assuming that you have Prometheus installed on your system. If you don't, you can download the respective Prometheus version [here](https://prometheus.io/download/).