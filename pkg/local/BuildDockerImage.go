package local

var BuildDockerImageCommand = `
	docker build --cache-from={{.ServiceName}}:latest --tag {{.ServiceName}}:latest --platform linux/amd64 . && \
	docker save -o {{.ServiceName}}-latest.tar {{.ServiceName}}
`
