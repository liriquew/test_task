APP_NAME = app
SRC = $(shell find ./internal ./cmd ./pkg -type f -name '*.go')

GOOSE_DBSTRING = "postgres://postgres:passw0rd@localhost:5432/users_db?sslmode=disable"

.PHONY: test clean

build: $(APP_NAME)

$(APP_NAME): $(SRC)
	go build -o $@ ./cmd/main.go

test:
	go test -count 1 -v  ./...

clean:
	rm -f $(APP_NAME)

gen_spec: ./spec/main.tsp
	tsp compile ./spec/main.tsp --output-dir ./spec/tsp-output --emit=@typespec/openapi3 && \
	cp ./spec/tsp-output/schema/openapi.yaml .
	go run github.com/ogen-go/ogen/cmd/ogen@latest --target ./internal/domain --clean openapi.yaml

migrate:
	goose postgres $(GOOSE_DBSTRING) -dir migrations up

check_health:
	curl -X GET -I -u admin:admin  localhost:8080/health
