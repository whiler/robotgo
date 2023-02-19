tidy: go.mod
	go mod tidy

cover:
	go test --cover -coverprofile=coverage.out -v github.com/whiler/robotgo
	go tool cover -html=coverage.out -o coverage.html
