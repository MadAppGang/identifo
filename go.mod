module github.com/madappgang/identifo

go 1.16

// https://github.com/etcd-io/etcd/issues/11749#issuecomment-679189808
replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

require (
	github.com/aws/aws-sdk-go v1.38.45
	github.com/boltdb/bolt v1.3.1
	github.com/casbin/casbin v1.9.1
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-redis/redis v0.0.0-20190503082931-75795aa4236d
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mailgun/mailgun-go v1.1.1
	github.com/njern/gonexmo v2.0.0+incompatible
	github.com/pallinder/go-randomdata v1.2.0
	github.com/qiangmzsx/string-adapter v0.0.0-20180323073508-38f25303bb0c
	github.com/rs/cors v1.6.0
	github.com/rs/xid v1.2.1
	github.com/sfreiberg/gotwilio v0.0.0-20201211181435-c426a3710ab5
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/urfave/negroni v1.0.0
	github.com/xlzd/gotp v0.0.0-20181030022105-c8557ba2c119
	go.mongodb.org/mongo-driver v1.3.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9 // indirect
	golang.org/x/sys v0.0.0-20201214210602-f9fddec55a1e // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/yaml.v2 v2.3.0
)
