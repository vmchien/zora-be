module vn.vato.zora.be.api/apps/zalo

go 1.26

replace (
	vn.vato.zora.be.api/api/zalo => ../../api/zalo
	// ### public folder ###
	vn.vato.zora.be.api/pkg => ../../pkg
)

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/google/wire v0.7.0
	google.golang.org/protobuf v1.36.11
	vn.vato.zora.be.api/api/zalo v0.0.0-00010101000000-000000000000
	vn.vato.zora.be.api/pkg v0.0.0-00010101000000-000000000000
)

require (
	ariga.io/atlas v1.2.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-kratos/kratos/contrib/log/zap/v2 v2.0.0-20260404020628-f149714c1d54 // indirect
	github.com/go-kratos/kratos/contrib/log/zerolog/v2 v2.0.0-20260404020628-f149714c1d54 // indirect
	github.com/go-openapi/inflect v0.21.5 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/hashicorp/hcl/v2 v2.24.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.9.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.21 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/redis/go-redis/v9 v9.18.0 // indirect
	github.com/rs/zerolog v1.35.0 // indirect
	github.com/zclconf/go-cty v1.18.0 // indirect
	github.com/zclconf/go-cty-yaml v1.2.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/mod v0.35.0 // indirect
	golang.org/x/tools v0.44.0 // indirect
)

require (
	cel.dev/expr v0.25.2 // indirect
	dario.cat/mergo v1.0.2 // indirect
	entgo.io/ent v0.14.6
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260523011958-0a33c5d7ca68 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260523011958-0a33c5d7ca68 // indirect
	google.golang.org/grpc v1.81.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
