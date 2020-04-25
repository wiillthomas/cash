<style>

</style>
# Cash

## What is Cash?
![Cash-Logo](https://cdn.clipart.email/6d189a101415d9170937b3253c31fc24_piggy-bank-transparent-background-money-clipart_900-600.jpeg)

## Cash is an in-memory cache for Docker containers, written in Go & designed to be stupidly easy to set up and use.

### Are you looking for a key-value in memory caching solution for your networked containers that:
- Is really fast at low/medium throughput
- Works out of the box with 0 config
- Allows for custom expiration/cleanup times
- Allows purging
- Works over HTTP
- Automatic sharding

### Do you NOT need all those useful features that Redis provides like:
- Storing anything that isn’t a string
- Backups/Restores
- Message brokering
- LRU replacement

<br />
<br />

---

## Installation

#### **Cash is available as a docker image here.**

Connection examples for `docker-compose` can be found in the examples folder.

The following enviorment variables can be set for config:

Env Variable | Accepts | Default | Description
--- | --- | --- | ---
verbose | Bool | false | Outputs any shards that have values every 2 seconds to the terminal.
expiry | Int | 20 | Default time each item in cache takes to expire - in seconds.
cleanup | Int | 10 | Interval at which the cleanup runs through the cache to find expired items - in seconds.
port | Int | 9192 | Port through which the service is accessable to other containers.
<br />
<br />

---------


## Usage

Url | Query Params | Method | Success Response | Failure Response | Sample Call
--- | --- | --- | --- | --- | ----
/create | **Required:** `value`<br /> **Optional:** `expiry` | GET | (200): `key` | (400): Must supply a value in query string | `http://cache:9192/create?value=CacheThisValue&expiry=600`
/read?key=`key` | **Required:** `key` | GET | (200): `value` | (204): Not Found <br /> <br />(400): Must supply a key in query string  | `http://cache:9192/read?key=e332a76c29654fcb7f6e6b31ced090c7`
/destroy | **Required:** `key` |  GET | (200): Destroyed | None - if not found then it is considered destroyed | `http://cache:9192/destroy?key=332a76c29654fcb7f6e6b31ced090c7`
/purge | | GET | (200): Purged | | `http://cache:9192/purge`
/healthcheck | | GET | (200): OK | | `http://cache:9192/healthcheck `

<br />
<br />

---

## Examples

An example of usage can be found in the 'examples' folder, run `docker-compose up`.

<br />
<br />