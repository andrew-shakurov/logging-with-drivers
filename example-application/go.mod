module example.com/example-application

go 1.23.2

replace example.com/log => ../log

replace example.com/logdriverjson => ../log-driver-json

require (
	example.com/log v0.0.0-00010101000000-000000000000
	example.com/logdriverjson v0.0.0-00010101000000-000000000000
)

require github.com/google/uuid v1.6.0 // indirect
