GH_NAME:=project-item-list
GH_FULL_NAME:=gh-${GH_NAME}

build:
	go build -o gh-${GH_NAME} main.go

install: build
	gh extension remove ${GH_NAME} || echo
	gh extension install .

LOG_DIR:=/tmp/${GH_FULL_NAME}
start: install
	rm -rf ${LOG_DIR}
	gh ${GH_NAME} -project-url=https://github.com/users/notomo/projects/1 -limit=10 -log=${LOG_DIR} -jq='.[] | select(.fieldValues.nodes|any(.field.name == "Status" and .name == "Todo"))'

test:
	go test -v ./...
