module github.com/stackitcloud/stackit-cli

go 1.25.1

require (
	github.com/fatih/color v1.18.0
	github.com/goccy/go-yaml v1.18.0
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0
	github.com/inhies/go-bytesize v0.0.0-20220417184213-4913239db9cf
	github.com/jedib0t/go-pretty/v6 v6.6.8
	github.com/lmittmann/tint v1.1.2
	github.com/mattn/go-colorable v0.1.14
	github.com/spf13/cobra v1.10.1
	github.com/spf13/pflag v1.0.10
	github.com/spf13/viper v1.21.0
	github.com/stackitcloud/stackit-sdk-go/core v0.17.3
	github.com/stackitcloud/stackit-sdk-go/services/alb v0.6.1
	github.com/stackitcloud/stackit-sdk-go/services/authorization v0.8.1
	github.com/stackitcloud/stackit-sdk-go/services/dns v0.17.1
	github.com/stackitcloud/stackit-sdk-go/services/git v0.7.1
	github.com/stackitcloud/stackit-sdk-go/services/iaas v0.30.0
	github.com/stackitcloud/stackit-sdk-go/services/mongodbflex v1.5.2
	github.com/stackitcloud/stackit-sdk-go/services/opensearch v0.24.1
	github.com/stackitcloud/stackit-sdk-go/services/postgresflex v1.2.1
	github.com/stackitcloud/stackit-sdk-go/services/resourcemanager v0.17.1
	github.com/stackitcloud/stackit-sdk-go/services/runcommand v1.3.1
	github.com/stackitcloud/stackit-sdk-go/services/secretsmanager v0.13.1
	github.com/stackitcloud/stackit-sdk-go/services/serverbackup v1.3.2
	github.com/stackitcloud/stackit-sdk-go/services/serverupdate v1.2.1
	github.com/stackitcloud/stackit-sdk-go/services/serviceaccount v0.11.1
	github.com/stackitcloud/stackit-sdk-go/services/serviceenablement v1.2.2
	github.com/stackitcloud/stackit-sdk-go/services/ske v1.4.0
	github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex v1.3.1
	github.com/zalando/go-keyring v0.2.6
	golang.org/x/mod v0.28.0
	golang.org/x/oauth2 v0.31.0
	golang.org/x/term v0.35.0
	golang.org/x/text v0.29.0
	k8s.io/apimachinery v0.34.1
	k8s.io/client-go v0.34.1
)

require (
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/time v0.13.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

require (
	al.essio.dev/pkg/shellescape v1.6.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.0 // indirect
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/danieljoos/wincred v1.2.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
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
	github.com/stackitcloud/stackit-sdk-go/services/loadbalancer v1.5.1
	github.com/stackitcloud/stackit-sdk-go/services/logme v0.25.1
	github.com/stackitcloud/stackit-sdk-go/services/mariadb v0.25.1
	github.com/stackitcloud/stackit-sdk-go/services/objectstorage v1.4.0
	github.com/stackitcloud/stackit-sdk-go/services/observability v0.14.0
	github.com/stackitcloud/stackit-sdk-go/services/rabbitmq v0.25.1
	github.com/stackitcloud/stackit-sdk-go/services/redis v0.25.1
	github.com/subosito/gotenv v1.6.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.34.1 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/utils v0.0.0-20250820121507-0af2bda4dd1d // indirect
	sigs.k8s.io/json v0.0.0-20250730193827-2d320260d730 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)
