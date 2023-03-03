spec:
	cp placetypes.json placetypes.json.last
	go run cmd/sfom-compile-spec/main.go > placetypes.json.tmp
	cp placetypes.json.tmp placetypes.json
	go run cmd/sfom-render-spec/main.go -path docs/images/placetypes.png
