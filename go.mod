module github.com/open-wanderer/wanderer/plugins/garmin

go 1.25.0

require (
	github.com/extism/go-pdk v1.1.3
	github.com/open-wanderer/wanderer/plugins/sdk v0.0.0
)

replace github.com/open-wanderer/wanderer/plugins/sdk => ./sdk
