module github.com/massn/ManualAccounter

go 1.15

replace github.com/massn/ManualAccounter/pkg/chart => ./pkg/chart

replace github.com/massn/ManualAccounter/pkg/json => ./pkg/json

require (
	github.com/go-echarts/go-echarts/v2 v2.2.3 // indirect
	github.com/go-openapi/errors v0.20.0 // indirect
	github.com/massn/ManualAccounter/pkg/chart v0.0.0-00010101000000-000000000000
)
