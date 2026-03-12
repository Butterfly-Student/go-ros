module github.com/Butterfly-Student/go-ros

go 1.21

require (
	github.com/go-routeros/routeros/v3 v3.0.0
	go.uber.org/zap v1.27.1
)

require go.uber.org/multierr v1.10.0 // indirect

replace (
	github.com/Butterfly-Student/go-ros/pkg/domain => ./pkg/domain
	github.com/Butterfly-Student/go-ros/pkg/mikrotik => ./pkg/mikrotik
	github.com/Butterfly-Student/go-ros/pkg/scripts => ./pkg/scripts
	github.com/Butterfly-Student/go-ros/pkg/services => ./pkg/services
	github.com/go-routeros/routeros/v3 => ./pkg/routeros
)
