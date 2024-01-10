# PostgreSql

PostgreSql backend for CoreDNS

## Name
PostgreSql - PostgreSql backend for CoreDNS

## Description

This plugin uses PostgreSql as a backend to store DNS records. These will then can served by CoreDNS. The backend uses a simple, single table data structure that can be shared by other systems to add and remove records from the DNS server. As there is no state stored in the plugin, the service can be scaled out by spinning multiple instances of CoreDNS backed by the same database.

## Syntax
```
postgresql {
    datasource DATA_SOURCE
    [table_prefix TABLE_PREFIX]
    [max_lifetime MAX_LIFETIME]
    [max_open_connections MAX_OPEN_CONNECTIONS]
    [max_idle_connections MAX_IDLE_CONNECTIONS]
    [ttl DEFAULT_TTL]
    [zone_update_interval ZONE_UPDATE_INTERVAL]
}
```

- `datasource` Datasource for PostgreSql, for example `host=127.0.0.1 port=5432 password=coredns sslmode=disable`
- `table_prefix` Prefix for the PostgreSql tables. Defaults to `coredns_`.
- `max_lifetime` Duration (in Golang format) for a SQL connection. Default is 1 minute.
- `max_open_connections` Maximum number of open connections to the database server. Default is 10.
- `max_idle_connections` Maximum number of idle connections in the database connection pool. Default is 10.
- `ttl` Default TTL for records without a specified TTL in seconds. Default is 360 (seconds)
- `zone_update_interval` Maximum time interval between loading all the zones from the database. Default is 10 minutes.

## Supported Record Types

A, AAAA, CNAME, SOA, TXT, NS, MX, CAA and SRV. This backend doesn't support AXFR requests. It also doesn't support wildcard records yet.

## Setup (as an external plugin)

Add this as an external plugin in `plugin.cfg` file:

```
postgresql:github.com/he-deng/coredns_postgresql
```

then run

```shell script
$ go generate
$ go build
```

Add any required modules to CoreDNS code as prompted.

## Build Docker image

Add this `Dockerfile` file:

```shell script
ARG DEBIAN_IMAGE=debian:stable-slim
ARG BASE=debian:stable-slim
FROM ${DEBIAN_IMAGE} AS build
SHELL [ "/bin/sh", "-ec" ]

RUN export DEBCONF_NONINTERACTIVE_SEEN=true \
           DEBIAN_FRONTEND=noninteractive \
           DEBIAN_PRIORITY=critical \
           TERM=linux ; \
    apt-get -qq update ; \
    apt-get -yyqq upgrade ; \
    apt-get -yyqq install ca-certificates libcap2-bin; \
    apt-get clean
COPY coredns /coredns
RUN setcap cap_net_bind_service=+ep /coredns

FROM ${BASE}
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /coredns /coredns
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
```

then run

```shell script
docker build -t coredns:1.11.1-postgresql .
```


## Database Setup
This plugin doesn't create or migrate database schema for its use yet. To create the database and tables, use the following table structure (note the table name prefix):

```sql
CREATE SEQUENCE coredns_records_id_seq
    INCREMENT 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    START 1
    CACHE 1;
CREATE TABLE coredns_records (
    id bigint DEFAULT nextval('coredns_records_id_seq'::regclass) NOT NULL,
    zone VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    ttl INT DEFAULT NULL,
    content TEXT,
    record_type VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
) ;
```

## Record setup
Each record served by this plugin, should belong to the zone it is allowed to server by CoreDNS. Here are some examples:

```sql
-- Insert batch #1
INSERT INTO coredns_records (zone, name, ttl, content, record_type) VALUES
('example.org.', '', 30, '{"ip": "1.1.1.1"}', 'A'),
('example.org.', '', '60', '{"ip": "1.1.1.0"}', 'A'),
('example.org.', 'test', 30, '{"text": "hello"}', 'TXT'),
('example.org.', 'mail', 30, '{"host" : "mail.example.org.","priority" : 10}', 'MX');
```

These can be queries using `dig` like this:

```shell script
$ dig A MX mail.example.org 
```

