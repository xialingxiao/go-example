test: 
	go test $$(go list ./... | grep -v /vendor/) -cover

run: 
	PORT=8080 go run *.go