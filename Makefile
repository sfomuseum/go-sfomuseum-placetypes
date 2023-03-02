spec:
	go run cmd/sfom-compile-spec/main.go > placetypes.json
	go run cmd/sfom-render-spec/main.go -path docs/images/placetypes.png
