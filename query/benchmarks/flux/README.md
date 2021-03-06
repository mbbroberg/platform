## Benchmarks

The output generated by `curl` uses `curl-format.txt`, present along side this readme.

Telegraf is configured with the supplied `telegraf.conf` to capture metrics from fluxd and influxdb using the 
Prometheus `/metrics` HTTP endpoint and machine metrics including CPU usage and disk I/O. Note that `influxd` is running
on port `8186`, allowing a separate `influxd` on the default port to receive metrics from Telegraf.

## Dataset #1

|       |       |
| ----- | ----- |
| series        | 100,000     |
| pps           | 3,000       |
| shards        | 12          |
| pps / shard   | 250         |
| total points  | 300,000,000 |

**pps**: points per series


### Hardware

|       |       |
| ----- | ----- |
| AWS instance type     | c3.4xlarge     |


### Generate dataset

1. Use [ingen][ingen] to populate a database with data.

   ```sh
   $ ingen -p=250 -t=1000,100 -shards=12 -start-time="2017-11-01T00:00:00Z" -data-path=~/.influxdb/data -meta-path=~/.influxdb/meta
   ```

   The previous command will
   
    * populate a database named `db` (default), 
    * create 100,000 series (1000×100),
    * made up of 2 tag keys (`tag0` and `tag1`) each with 1000 and 100 tag values respectively.
    * 250 points per series, per shard, for a total of 3,000 points per series.
    * Points will start from `2017-11-01 00:00:00 UTC` and
    * span 12 shards.


### Flux queries

Query #1

```sh
HOST=localhost:8093; curl -w "@curl-format.txt" -H 'Accept: text/plain' -o /dev/null -s http://${HOST}/query \
    --data-urlencode 'q=from(db:"db").range(start:2017-11-01T00:00:00Z, stop:2017-11-02T00:00:00Z).filter(exp:{"_measurement" == "m0" and "_field" == "v0" and $ > 0}).sum()'

    
time_starttransfer: 0.138
size_download:      5800000
time_total:         7.578

```

Query #2

```sh
HOST=localhost:8093; curl -w "@curl-format.txt" -H 'Accept: text/plain' -o /dev/null -s http://${HOST}/query \
    --data-urlencode 'q=from(db:"db").range(start:2017-11-01T00:00:00Z, stop:2017-11-05T00:00:00Z).filter(exp:{"_measurement" == "m0" and "_field" == "v0" and $ > 0}).sum()'

    
time_starttransfer: 0.305
size_download:      5900000
time_total:         17.909
```

Query #3

```sh
HOST=localhost:8093; curl -w "@curl-format.txt" -H 'Accept: text/plain' -o /dev/null -s http://${HOST}/query \
    --data-urlencode 'q=from(db:"db").range(start:2017-11-01T00:00:00Z, stop:2017-11-05T00:00:00Z).filter(exp:{"_measurement" == "m0" and "_field" == "v0" and $ > 0}).group(by:["tag0"]).sum()'

    
time_starttransfer: 22.727
size_download:      60000
time_total:         22.730
```

Query #4

```sh
HOST=localhost:8093; curl -w "@curl-format.txt" -H 'Accept: text/plain' -o /dev/null -s http://${HOST}/query \
    --data-urlencode 'q=from(db:"db").range(start:2017-11-01T00:00:00Z, stop:2017-11-13T00:00:00Z).filter(exp:{"_measurement" == "m0" and "_field" == "v0" and $ > 0}).sum()'

    
time_starttransfer: 0.713
size_download:      5900000
time_total:         44.159
```

Query #5

```sh
HOST=localhost:8093; curl -w "@curl-format.txt" -H 'Accept: text/plain' -o /dev/null -s http://${HOST}/query \
    --data-urlencode 'q=from(db:"db").range(start:2017-11-01T00:00:00Z, stop:2017-11-13T00:00:00Z).filter(exp:{"_measurement" == "m0" and "_field" == "v0" and $ > 0}).group(by:["tag0"]).sum()'

    
time_starttransfer: 56.257
size_download:      60000
time_total:         56.261
```

## Dataset #2

|       |       |
| ----- | ----- |
| series        | 10,000,000     |
| pps           | 1,000          |
| shards        | 4              |
| pps / shard   | 250            |
| total points  | 10,000,000,000 |

**pps**: points per series


### Hardware

|       |       |
| ----- | ----- |
| AWS instance type     | c5.4xlarge     |


### Generate dataset

1. Use [ingen][ingen] to populate a database with data.

   ```sh
   $ ingen -p=250 -t=10000,100,10 -shards=4 -start-time="2017-11-01T00:00:00Z" -data-path=~/.influxdb/data -meta-path=~/.influxdb/meta
   ```

   The previous command will
   
    * populate a database named `db` (default), 
    * create 10,000,000 series (10000×100×10),
    * made up of 3 tag keys (`tag0`, `tag1`, `tag2`) each with 10000, 100 and 10 tag values respectively.
    * 250 points per series, per shard, for a total of 1,000 points per series.
    * Points will start from `2017-11-01 00:00:00 UTC` and
    * span 4 shards.


### Flux queries

Query #1

```sh
HOST=localhost:8093; curl -w "@curl-format.txt" -H 'Accept: text/plain' -o /dev/null -s http://${HOST}/query \
    --data-urlencode 'q=from(db:"db").range(start:2017-11-01T00:00:00Z, stop:2017-11-05T00:00:00Z).filter(exp:{"_measurement" == "m0" and "_field" == "v0" and "tag1" == "value00"}).group(by:["tag0"]).sum()'

    
time_starttransfer: 0.325
size_download:      7200000
time_total:         11.437
```

Query #2

```sh
HOST=localhost:8093; curl -w "@curl-format.txt" -H 'Accept: text/plain' -o /dev/null -s http://${HOST}/query \
    --data-urlencode 'q=from(db:"db").range(start:2017-11-01T00:00:00Z, stop:2017-11-05T00:00:00Z).filter(exp:{"_measurement" == "m0" and "_field" == "v0" and "tag1" == "value00"}).group(by:["tag0"]).sum()'

    
time_starttransfer: 13.174
size_download:      600000
time_total:         13.215
```

Query #3

```sh
HOST=localhost:8093; curl -w "@curl-format.txt" -H 'Accept: text/plain' -o /dev/null -s http://${HOST}/query \
    --data-urlencode 'q=from(db:"db").range(start:2017-11-01T00:00:00Z, stop:2017-11-05T00:00:00Z).filter(exp:{"_measurement" == "m0" and "_field" == "v0"}).group(by:["tag0"]).sum()'

    
time_starttransfer: 1190.204
size_download:      620000
time_total:         1190.244
```

Query #4

```sh
HOST=localhost:8093; curl -w "@curl-format.txt" -H 'Accept: text/plain' -o /dev/null -s http://${HOST}/query \
    --data-urlencode 'q=from(db:"db").range(start:2017-11-01T00:00:00Z, stop:2017-11-05T00:00:00Z).filter(exp:{"_measurement" == "m0" and "_field" == "v0"}).sum()'

    
time_starttransfer: 23.975
size_download:      720000000
time_total:         803.254
```



[ingen]: https://github.com/influxdata/ingen
