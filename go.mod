module github.com/stackitcloud/stackit-cli

go 1.24.0

require (
	github.com/fatih/color v1.18.0
	github.com/goccy/go-yaml v1.19.0
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0
	github.com/inhies/go-bytesize v0.0.0-20220417184213-4913239db9cf
	github.com/jedib0t/go-pretty/v6 v6.7.5
	github.com/lmittmann/tint v1.1.2
	github.com/mattn/go-colorable v0.1.14
	github.com/spf13/cobra v1.10.1
	github.com/spf13/pflag v1.0.10
	github.com/spf13/viper v1.21.0
	github.com/stackitcloud/stackit-sdk-go/core v0.20.0
	github.com/stackitcloud/stackit-sdk-go/services/alb v0.7.2
	github.com/stackitcloud/stackit-sdk-go/services/authorization v0.10.0
	github.com/stackitcloud/stackit-sdk-go/services/dns v0.17.2
	github.com/stackitcloud/stackit-sdk-go/services/git v0.9.1
	github.com/stackitcloud/stackit-sdk-go/services/iaas v1.2.2
	github.com/stackitcloud/stackit-sdk-go/services/intake v0.4.0
	github.com/stackitcloud/stackit-sdk-go/services/mongodbflex v1.5.3
	github.com/stackitcloud/stackit-sdk-go/services/opensearch v0.24.2
	github.com/stackitcloud/stackit-sdk-go/services/postgresflex v1.2.1
	github.com/stackitcloud/stackit-sdk-go/services/resourcemanager v0.18.1
	github.com/stackitcloud/stackit-sdk-go/services/runcommand v1.3.2
	github.com/stackitcloud/stackit-sdk-go/services/secretsmanager v0.13.2
	github.com/stackitcloud/stackit-sdk-go/services/serverbackup v1.3.3
	github.com/stackitcloud/stackit-sdk-go/services/serverupdate v1.2.2
	github.com/stackitcloud/stackit-sdk-go/services/serviceaccount v0.11.2
	github.com/stackitcloud/stackit-sdk-go/services/serviceenablement v1.2.3
	github.com/stackitcloud/stackit-sdk-go/services/ske v1.5.0
	github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex v1.3.3
	github.com/zalando/go-keyring v0.2.6
	golang.org/x/mod v0.30.0
	golang.org/x/oauth2 v0.33.0
	golang.org/x/term v0.37.0
	golang.org/x/text v0.31.0
	k8s.io/apimachinery v0.34.2
	k8s.io/client-go v0.34.2
)

require (
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

require (
	al.essio.dev/pkg/shellescape v1.5.1 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/telemetry v0.0.0-20251111182119-bc8e575c7b54 // indirect
	golang.org/x/tools v0.39.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.0 // indirect
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/danieljoos/wincred v1.2.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/stackitcloud/stackit-sdk-go/services/kms v1.1.1
	github.com/stackitcloud/stackit-sdk-go/services/loadbalancer v1.6.1
	github.com/stackitcloud/stackit-sdk-go/services/logme v0.25.2
	github.com/stackitcloud/stackit-sdk-go/services/mariadb v0.25.2
	github.com/stackitcloud/stackit-sdk-go/services/objectstorage v1.4.1
	github.com/stackitcloud/stackit-sdk-go/services/observability v0.15.1
	github.com/stackitcloud/stackit-sdk-go/services/rabbitmq v0.25.2
	github.com/stackitcloud/stackit-sdk-go/services/redis v0.25.2
	github.com/subosito/gotenv v1.6.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.34.2 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/utils v0.0.0-20250604170112-4c0f3b243397 // indirect
	sigs.k8s.io/json v0.0.0-20241014173422-cfa47c3a1cc8 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

tool golang.org/x/tools/cmd/goimports
