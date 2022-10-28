build-hugo:
	rm -rf ./docs && cp -r charts-contrib/* web/commander-charts/static/charts && cd web/commander-charts &&  hugo -D -d ../../docs
build-bin:
	go build -o cr ./cmd/main.go 
