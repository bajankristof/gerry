build:
	@go build -o bin/gerry cmd/main.go

link:
	@ln -sf $(CURDIR)/bin/gerry /usr/local/bin/gerry
