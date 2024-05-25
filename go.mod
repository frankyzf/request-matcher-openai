module request-matcher-openai

go 1.18

require (
	github.com/aws/aws-sdk-go v1.44.214
	github.com/bsm/go-guid v1.0.0
	github.com/deckarep/golang-set v1.8.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dpapathanasiou/go-recaptcha v0.0.0-20190121160230-be5090b17804
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-contrib/static v0.0.1
	github.com/gin-gonic/gin v1.7.7
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.7.0
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/consul/api v1.12.0
	github.com/itsjamie/gin-cors v0.0.0-20160420130702-97b4a9da7933
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lib/pq v1.10.7
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cast v1.4.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	golang.org/x/crypto v0.11.0
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b
	gorm.io/driver/mysql v1.2.2
	gorm.io/gorm v1.22.4
	moul.io/http2curl v1.0.0
)

require (
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.0.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/serf v0.9.6 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lestrrat-go/strftime v1.0.5 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.27.10 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.4

replace gorm.io/gorm => gorm.io/gorm v1.22.5

replace gorm.io/datatypes => gorm.io/datatypes v1.0.5

replace gorm.io/driver/mysql => gorm.io/driver/mysql v1.2.3

replace gorm.io/driver/postgres => gorm.io/driver/postgres v1.2.3

replace gorm.io/driver/sqlite => gorm.io/driver/sqlite v1.2.6

replace gorm.io/driver/sqlserver => gorm.io/driver/sqlserver v1.2.1
