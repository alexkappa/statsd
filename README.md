# statsd

Port of Etsy's [statsd](https://github.com/etsy/statsd), written in Go.

Supports

* Timing (with optional percentiles)
* Counters (with optional sampling)
* Gauges

## Installing

```bash
go get github.com/alexkappa/statsd
```

## Usage

```bash
statsd config.json
```