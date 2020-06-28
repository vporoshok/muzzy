all         : test cover lint

test        :
	@echo "Testing..."
	@go test -test.race -test.count=1 -test.bench=. -test.benchmem -test.cover -test.coverprofile=.coverprofile ./...
	@echo ""

cover       :
	@echo "Check coverage..."
	@go tool cover -func=.coverprofile | tail -n 1 | awk '{print "Total coverage:", $$3;}'
	@test `go tool cover -func=.coverprofile | tail -n 1 | awk '{print $$3;}' | sed 's/\..*//'` -ge 90
	@echo ""

lint        :
	@echo "Linting..."
	@golangci-lint run
	@echo "PASS"
	@echo ""
