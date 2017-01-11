test:
	go test -v -coverprofile=cover.out -covermode=count
	go tool cover -html=cover.out -o=cover.html
