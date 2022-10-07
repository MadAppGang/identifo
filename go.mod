module github.com/madappgang/identifo/v2

go 1.19

require (
	github.com/aws/aws-sdk-go v1.38.45
	github.com/casbin/casbin v1.9.1
	github.com/go-redis/redis v0.0.0-20190503082931-75795aa4236d
	github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/sessions v1.2.1
	github.com/hashicorp/go-plugin v1.4.3
	github.com/jszwec/s3fs v0.3.1
	github.com/mailgun/mailgun-go v1.1.1
	github.com/markbates/goth v1.68.0
	github.com/njern/gonexmo v2.0.0+incompatible
	github.com/pallinder/go-randomdata v1.2.0
	github.com/qiangmzsx/string-adapter v0.0.0-20180323073508-38f25303bb0c
	github.com/rs/cors v1.6.0
	github.com/rs/xid v1.2.1
	github.com/sfreiberg/gotwilio v0.0.0-20201211181435-c426a3710ab5
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/spf13/afero v1.6.0
	github.com/stretchr/testify v1.7.0
	github.com/urfave/negroni v1.0.0
	github.com/xlzd/gotp v0.0.0-20181030022105-c8557ba2c119
	go.etcd.io/bbolt v1.3.6
	go.mongodb.org/mongo-driver v1.8.2
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/exp v0.0.0-20221006183845-316c7553db56
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/yaml.v2 v2.4.0
)

require (
	cloud.google.com/go v0.67.0 // indirect
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/yamux v0.0.0-20180604194846-3520598351bb // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lestrrat-go/jwx v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mitchellh/go-testing-interface v0.0.0-20171004221916-a61a99592b77 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.16.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f // indirect
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20200929141702-51c3e5b607fe // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/sfreiberg/gotwilio => github.com/MadAppGang/gotwilio v0.0.0-20210820024906-f91dd2ebe762
