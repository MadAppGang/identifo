module github.com/madappgang/identifo

go 1.16

// https://github.com/etcd-io/etcd/issues/11749#issuecomment-679189808
replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

require (
	github.com/aws/aws-sdk-go v1.27.0
	github.com/boltdb/bolt v1.3.1
	github.com/casbin/casbin v1.9.1
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-redis/redis v0.0.0-20190503082931-75795aa4236d
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/schema v1.1.0 // indirect
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/mailgun/mailgun-go v1.1.1
	github.com/njern/gonexmo v2.0.0+incompatible
	github.com/pallinder/go-randomdata v1.2.0
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/qiangmzsx/string-adapter v0.0.0-20180323073508-38f25303bb0c
	github.com/rs/cors v1.6.0
	github.com/rs/xid v1.2.1
	github.com/sfreiberg/gotwilio v0.0.0-20190708190155-499f54b30211
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/urfave/negroni v1.0.0
	github.com/xlzd/gotp v0.0.0-20181030022105-c8557ba2c119
	go.etcd.io/etcd v3.3.25+incompatible
	go.mongodb.org/mongo-driver v1.3.0
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/yaml.v2 v2.3.0
)
