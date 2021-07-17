module github.com/madappgang/identifo/test/runner

go 1.16

replace github.com/madappgang/identifo => ../..

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/madappgang/identifo v1.2.9
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/onsi/gomega v1.13.0
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	gopkg.in/h2non/baloo.v3 v3.0.2
	gopkg.in/h2non/gentleman.v2 v2.0.5 // indirect
)
