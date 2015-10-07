# MySQL test with Go

## Todo

- Cleanup docker script(s) : add start-db.sh
- docker-compose (no volumes)
- Rewrite using `sqlx`
- Sink for InfluxDB
- Batch writes to ted/watt2 - Reader or LOAD INFILE

## Vendoring
See this [Go/Wiki for reference](https://github.com/golang/go/wiki/PackageManagementTools)

We want to use `GO15VENDOREXPERIMENT=1` and place our external dependencies in a `vendor folder`.
We are using [`govend`](https://github.com/gophersaurus/govend) as listed.

To install `govend`, we did a standard `go get -u github.com/gophersaurus/govend`, and makde sure our `GOPATH` was set and `$GOPATH/bin` is on our `$PATH`. also `GO15VENDOREXPERIMENT=1` needs to be set.

## Docker

These will be data volumes..

    docker create -v /data --name teddbdata mysql /bin/true
    docker create -v /data --name tedfluxdata tutum/influxdb /bin/true

## Timing of MySQL reads
For timing of MySQL selects with maxCount results

From goedel to cantor

	3600: 989s
	3600*24: 357s
	3600*24*10: 324s

From Godel to local docker:

    3600*24: 294s,290s  (Read-only)
	10000: --s  (Batch Writes)



## MySQL inserts are ridiculously slow
This is what we di to spead things up:
(all time measured on dirac/docker)

- Initial naive approach: 550 ins/s
- Prepared statements: 650 ins/s
- Transactions: 2500-3000 ins/s (depending on size of batch...)

