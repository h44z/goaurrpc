# goaurrpc
[![Release](https://img.shields.io/github/v/release/moson-mo/goaurrpc)](https://github.com/moson-mo/goaurrpc/releases) [![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/moson-mo/goaurrpc/go.yml?branch=main)](https://github.com/moson-mo/goaurrpc/actions) [![Coverage](https://img.shields.io/badge/Coverage-98.2%25-brightgreen)](https://github.com/moson-mo/goaurrpc/blob/main/test_coverage.out) [![Go Report Card](https://goreportcard.com/badge/github.com/moson-mo/goaurrpc)](https://goreportcard.com/report/github.com/moson-mo/goaurrpc)

### An implementation of the [aurweb](https://gitlab.archlinux.org/archlinux/aurweb) (v6) - /rpc - REST API service in go

goaurrpc allows you to run your own self-hosted aurweb /rpc endpoint.  
This project implements the /rpc interface (REST API; version 5) as described [here](https://aur.archlinux.org/rpc/).  

In it's default configuration, package data is being downloaded/refreshed from the AUR every 5 minutes.  
The data is entirely held in-memory as opposed to storing it in a database for example.  
This avoids the need to make heavy database queries for each request.  
For a performance comparison, see [Benchmarks](BENCHMARKS.md)

### How to build

- Download repository `git clone https://github.com/moson-mo/goaurrpc.git`
- `cd goaurrpc`
- Build with: `./build.sh`
- This will create a binary `goaurrpc`

### Configuration

The whole configuration is done via environment variables.  
The following environment variables are supported: 


| Environment Variable            | Description                                                                                                                                                           |
|---------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| PORT                            | The port number our service is listening on                                                                                                                           |
| AUR_FILE_LOCATION               | Either the URL to the full metadata archive `packages-meta-ext-v1.json.gz` or a local copy of the file                                                                |
| MAX_RESULTS                     | The maximum number of package results that are being returned to the client                                                                                           |
| REFRESH_INTERVAL                | The interval (in seconds) in which the metadata file is being reloaded                                                                                                |
| RATE_LIMIT                      | The maximum number of requests that are allowed within the time-window                                                                                                |
| LOAD_FROM_FILE                  | Set to true when using a local file instead of a URL for `AurFileLocation`                                                                                            |
| RATE_LIMIT_CLEANUP_INTERVAL     | The interval (in seconds) in which rate-limits are being cleaned up                                                                                                   |
| RATE_LIMIT_TIME_WINDOW          | Defines the length of the time window for rate-limiting (in seconds)                                                                                                  |
| Trusted TRUSTED_REVERSE_PROXIES | A list of trusted IP-Addresses, in case you use a reverse proxy and need to rely on `X-Real-IP` or `X-Forwarded-For` headers to identify a client (for rate-limiting) |
| ENABLE_SSL                      | Enables internal SSL/TLS. You'll need to provide `CertFile`and `KeyFile` when enabling it. I'd recommend to use nginx as reverse proxy to add encryption instead      |
| CERT_FILE                       | Path to the cert file (if SSL is enabled)                                                                                                                             |
| KEY_FILE                        | Path to the corresponding key file (if SSL is enabled)                                                                                                                |
| ENABLE_SEARCH_CACHE             | Caches data for search queries that have been performed by clients                                                                                                    |
| CACHE_CLEANUP_INTERVAL          | The interval (in seconds) for performing cleanup of search-cache entries                                                                                              |
| CACHE_EXPIRATION_TIME           | The number of seconds an entry should stay in the search-cache                                                                                                        |
| ENABLE_METRICS                  | Enables Prometheus metrics at /metrics                                                                                                                                |
| ENABLE_ADMIN_API                | Enables the administrative endpoint at /admin                                                                                                                         |
| ADMIN_API_KEY                   | The API Key that is to be provided in the header for the /admin endpoint                                                                                              |
| LOG_FILE                        | File path to which logs should be written. Empty means logs are printed to stdout                                                                                     |
| LOG_LEVEL                       | Sets the logging level (valid values: debug, info, warn, error), defaults to info                                                                                     |

### Public endpoint

Feel free to make use of the following public instance of goaurrpc:   

[HTTP](http://server.moson.rocks/rpc) / [HTTPS](https://server.moson.rocks/rpc)

### Using this service with yay

```shell
aur --aurrpcurl http://localhost:10666 <all other actions>
```

### Future plans / ideas

- Extend request types (see [v6-proposal branch](https://github.com/moson-mo/goaurrpc/tree/v6-proposal))
- Admin REST-API to be able to control goaurrpc at runtime, for example:
  - reload data
  - get statistics (memory consumption, rate limits, etc.)
  - manage rate-limits
  - manage search-cache
- CLI/TUI tool for administration (making use of the admin api)