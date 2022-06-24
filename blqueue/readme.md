# Balanced Leaky Queue demo app

Demo stage for [blqueue](https://github.com/koykov/blqueue) package.

## Installation

Default way is running app in Docker container. Before start, you need to make `.env` file with admin login/password for
builtin Grafana server:
```dotenv
ADMIN_USER=<login>
ADMIN_PASSWORD=<password>
```
Then just run `sudo docker-compose up` in the package directory.

API will available at `http://localhost:8080`,
Grafana at `http://localhost:3000/`.

You may start the app from host machine, but you need to configure Prometheus and Grafana services yourself.

## API

#### Init the queue

```shell
curl -i -X POST -H "Content-Type:application/json" -d '{<configBody>}' '/api/v1/init?key=<queueKey>'
```

`<queueKey>` is a machine-readable key that uses to perform modification/stop requests for the qey. That key also uses
in metrics writers to filter queue metrics.

`<configBody>` describes queue config in JSON format:
```json lines
{
  "size": 1e5, // Queue size in items. Exceeding this param will block the queue or leak extra items dependent of allow_leak param (see below).
  "allow_leak": true, // Enable leaky feature of the queue. On false queue will block on QFR == 1.
  "heartbeat": 100000000, // Delay between heartbeats in nanosecond.
  "wakeup_factor": 0.05, // Reaching this QFR (queue fullness rate) value will trigger wakeup of idle/sleeping workers.
  "sleep_factor": 0.005, // Reaching this QFR value will sleep all available worker.
  "workers_min": 2, // Minimum workers that always will work.
  "workers_max": 16, // Maximum workers available to use to balance queue.
  "worker_delay": 5000000, // Worker's load emulation duration in nanoseconds.
  "workers_schedule": [ // Workers schedule. Allow you to specify workers min/max and factor params for certain time ranges.
    {"rel_range":"30s-1m","workers_min":4,"workers_max":50},
    {"rel_range":"1m-1m30s","workers_min":8,"workers_max":50},
    {"rel_range":"1m30s-2m","workers_min":12,"workers_max":50},
    {"rel_range":"2m-2m30s","workers_min":16,"workers_max":50},
    {"rel_range":"2m30s-3m","workers_min":12,"workers_max":50},
    {"rel_range":"3m-3m30s","workers_min":8,"workers_max":50},
    {"rel_range":"3m30s-4m","workers_min":4,"workers_max":50}
  ],
  "producers_min": 1, // Minimum producers that always will work.
  "producers_max": 20, // Maximum producers available to emulate queue load.
  "producer_delay": 5000000, // Producer's load emulation duration in nanoseconds.
  "producers_schedule": [ // Producers schedule. Allow you to specify how many producers should work on certain time ranges.
    {"rel_range":"30s-1m","producers":4},
    {"rel_range":"1m-1m30s","producers":8},
    {"rel_range":"1m30s-2m","producers":12},
    {"rel_range":"2m-2m30s","producers":16},
    {"rel_range":"2m30s-3m","producers":12},
    {"rel_range":"3m-3m30s","producers":8},
    {"rel_range":"3m30s-4m","producers":4}
  ]
}
```

#### Producers up

```shell
curl -i -X GET '/api/v1/producer-up?key=<queueKey>&delta=<delta>'
```

Activate `delta` producers. If currently active producers + `delta` will exceed `producers_max` param, then request will
fail with error.

#### Producers down

```shell
curl -i -X GET '/api/v1/producer-down?key=<queueKey>&delta=<delta>'
```

Stop `delta` producers. If currently active producers - `delta` will smaller than `producers_min` param, then request
will fail with error.
