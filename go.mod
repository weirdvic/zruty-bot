module github.com/weirdvic/zruty-bot

go 1.24

replace github.com/yanzay/tbot/v2 => github.com/weirdvic/tbot/v2 v2.2.0-patched

require (
	github.com/golang-migrate/migrate/v4 v4.18.2
	github.com/mattn/go-sqlite3 v1.14.27
	github.com/yanzay/tbot/v2 v2.2.0-patched
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
)
