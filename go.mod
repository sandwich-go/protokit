module github.com/sandwich-go/protokit

go 1.22

toolchain go1.22.1

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/golang/protobuf v1.5.4
	github.com/hoisie/mustache v0.0.0-20160804235033-6375acf62c69
	github.com/jhump/protoreflect v1.12.0
	github.com/mattn/go-colorable v0.1.12
	github.com/rs/zerolog v1.26.1
	github.com/sandwich-go/boost v1.3.1
	github.com/sandwich-go/protokit/option v0.0.1
	github.com/smartystreets/goconvey v1.7.2
	google.golang.org/genproto v0.0.0-20220405205423-9d709892a2bf
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
)

replace github.com/sandwich-go/protokit/option => ./option
