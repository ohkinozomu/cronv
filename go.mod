module github.com/ohkinozomu/cronv

go 1.15

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/ohkinozomu/cronv/server v0.0.0-00010101000000-000000000000
	github.com/robfig/cron/v3 v3.0.1
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/stretchr/testify v1.3.0
)

replace github.com/ohkinozomu/cronv/server => ./pkg/server
