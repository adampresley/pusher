package local

var UploadDockerImageCommand = `
   scp {{.ServiceName}}-latest.tar {{.Host}}:~/applications/{{.ServiceName}} && rm {{.ServiceName}}-latest.tar
`
