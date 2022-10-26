# PostgreSql

PostgreSql backend for CoreDNS

## Name
PostgreSql - PostgreSql backend for CoreDNS

## Description

This plugin uses PostgreSql as a backend to store DNS records. These will then can served by CoreDNS. The backend uses a simple, single table data structure that can be shared by other systems to add and remove records from the DNS server. As there is no state stored in the plugin, the service can be scaled out by spinning multiple instances of CoreDNS backed by the same database.

## Syntax
```
postgresql {
    dsn DSN
    [table_prefix TABLE_PREFIX]
    [max_lifetime MAX_LIFETIME]
    [max_open_connections MAX_OPEN_CONNECTIONS]
    [max_idle_connections MAX_IDLE_CONNECTIONS]
    [ttl DEFAULT_TTL]
    [zone_update_interval ZONE_UPDATE_INTERVAL]
}
```

- `dsn` DSN for PostgreSql, for example `host=10.0.0.80 port=5432 user=py password=312 dbname=haspinfodb sslmode=disable`
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

## Database Setup
This plugin doesn't create or migrate database schema for its use yet. To create the database and tables, use the following table structure (note the table name prefix):

```sql
CREATE SEQUENCE coredns_records_myid_seq
    INCREMENT 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    START 1
    CACHE 1;
CREATE TABLE coredns_records (
    id bigint DEFAULT nextval('coredns_records_myid_seq'::regclass) NOT NULL,
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
('example.org.', 'foo', 30, '{"ip": "1.1.1.1"}', 'A'),
('example.org.', 'foo', '60', '{"ip": "1.1.1.0"}', 'A'),
('example.org.', 'foo', 30, '{"text": "hello"}', 'TXT'),
('example.org.', 'foo', 30, '{"host" : "foo.example.org.","priority" : 10}', 'MX');
```

These can be queries using `dig` like this:

```shell script
$ dig A MX foo.example.org 
```

