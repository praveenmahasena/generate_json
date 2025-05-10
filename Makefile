ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
all: test
all: vet
all: package
all: package_race
all: reader
all: reader_race

test: vet
test: base_test
test: staticcheck
test: shadow


base_test:
	go test ./... -v

vet:
	go vet ./...

staticcheck: staticcheck_bin
	bin/staticcheck ./...

staticcheck_bin:
	GOBIN=${ROOT_DIR}/bin go install honnef.co/go/tools/cmd/staticcheck@latest


shadow: shadow_bin
	bin/shadow ./...

shadow_bin:
	GOBIN=${ROOT_DIR}/bin go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest

package: generate_json

package_race: generate_json_race

reader: jsonreader

reader_race: jsonreader_race

generate_json:
	go build -o ./bin/generate_json ./cmd/generate_json/

generate_json_race:
	go build --race -o ./bin/generate_json_race ./cmd/generate_json/

jsonreader:
	go build -o ./bin/jsonreader ./cmd/read_json/

jsonreader_race:
	go build --race -o ./bin/jsonreader_race ./cmd/read_json/
