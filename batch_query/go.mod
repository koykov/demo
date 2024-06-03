module github.com/koykov/demo/batch_query

go 1.20

require (
	github.com/aerospike/aerospike-client-go v4.5.2+incompatible
	github.com/go-sql-driver/mysql v1.8.1
	github.com/koykov/batch_query v0.0.0-20240603180339-62c0ff204f00
	github.com/koykov/batch_query/mods/aerospike v0.0.0-20240603180339-62c0ff204f00
	github.com/koykov/batch_query/mods/sql v0.0.0-20240603180339-62c0ff204f00
	github.com/koykov/metrics_writers/batch_query v0.0.0-20230904210402-14dadf68561a
	github.com/lib/pq v1.10.9
	github.com/prometheus/client_golang v1.19.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/yuin/gopher-lua v1.1.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
