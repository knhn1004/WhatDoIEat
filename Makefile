build:
	@go build -o bin/whatDoIEat
run: build
	@cp .env.local bin/.env.local && ./bin/whatDoIEat
# Not yet implemented
seed: build
	@./bin/whatDoIEat --seed
test:
	@go test -v ./...
