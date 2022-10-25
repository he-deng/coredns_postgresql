package coredns_postgresql

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

const (
	defaultTtl                = 360
	defaultMaxLifeTime        = 1 * time.Minute
	defaultMaxOpenConnections = 10
	defaultMaxIdleConnections = 10
	defaultZoneUpdateTime     = 10 * time.Minute
)

func init() {
	caddy.RegisterPlugin("postgresql", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	r, err := postgresqlParse(c)
	if err != nil {
		return plugin.Error("postgresql", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		r.Next = next
		return r
	})

	return nil
}

func postgresqlParse(c *caddy.Controller) (*CoreDNSPostgreSql, error) {
	postgresql := CoreDNSPostgreSql{
		TablePrefix: "coredns_",
		Ttl:         300,
	}
	var err error

	c.Next()
	if c.NextBlock() {
		for {
			switch c.Val() {
			case "datasource":
				if !c.NextArg() {
					return &CoreDNSPostgreSql{}, c.ArgErr()
				}
				postgresql.Datasource = c.Val()
			case "table_prefix":
				if !c.NextArg() {
					return &CoreDNSPostgreSql{}, c.ArgErr()
				}
				postgresql.TablePrefix = c.Val()
			case "max_lifetime":
				if !c.NextArg() {
					return &CoreDNSPostgreSql{}, c.ArgErr()
				}
				var val time.Duration
				val, err = time.ParseDuration(c.Val())
				if err != nil {
					val = defaultMaxLifeTime
				}
				postgresql.MaxLifetime = val
			case "max_open_connections":
				if !c.NextArg() {
					return &CoreDNSPostgreSql{}, c.ArgErr()
				}
				var val int
				val, err = strconv.Atoi(c.Val())
				if err != nil {
					val = defaultMaxOpenConnections
				}
				postgresql.MaxOpenConnections = val
			case "max_idle_connections":
				if !c.NextArg() {
					return &CoreDNSPostgreSql{}, c.ArgErr()
				}
				var val int
				val, err = strconv.Atoi(c.Val())
				if err != nil {
					val = defaultMaxIdleConnections
				}
				postgresql.MaxIdleConnections = val
			case "zone_update_interval":
				if !c.NextArg() {
					return &CoreDNSPostgreSql{}, c.ArgErr()
				}
				var val time.Duration
				val, err = time.ParseDuration(c.Val())
				if err != nil {
					val = defaultZoneUpdateTime
				}
				postgresql.zoneUpdateTime = val
			case "ttl":
				if !c.NextArg() {
					return &CoreDNSPostgreSql{}, c.ArgErr()
				}
				var val int
				val, err = strconv.Atoi(c.Val())
				if err != nil {
					val = defaultTtl
				}
				postgresql.Ttl = uint32(val)
			default:
				if c.Val() != "}" {
					return &CoreDNSPostgreSql{}, c.Errf("unknown property '%s'", c.Val())
				}
			}

			if !c.Next() {
				break
			}
		}

	}

	db, err := postgresql.db()
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	postgresql.tableName = postgresql.TablePrefix + "records"

	return &postgresql, nil
}

func (handler *CoreDNSPostgreSql) db() (*sql.DB, error) {
	db, err := sql.Open("postgres", handler.Datasource)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(handler.MaxLifetime)
	db.SetMaxOpenConns(handler.MaxOpenConnections)
	db.SetMaxIdleConns(handler.MaxIdleConnections)

	return db, nil
}
