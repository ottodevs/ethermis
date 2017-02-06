# This is how we want to name the binary output
OUTDIR=$(PWD)/build/bin
BIN=$(OUTDIR)/ethermis
APIDIR=api
# These are the values we want to pass for Version and BuildTime

GIT_COMMIT=`git rev-parse HEAD`
BUILD_TIME=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X constant.GitCommit=${GIT_COMMIT} -X constant.BuildTime=${BUILD_TIME}"

PROTOS := \
	api/ethereum/ethereum.proto

PROTO_SRC := \
	$(subst .proto,.pb.go,$(PROTOS)) \
	$(subst .proto,.pb.gw.go,$(PROTOS))

all: $(OUTDIR) $(PROTO_SRC)
	@go build ${LDFLAGS} -o ${BIN} main.go

$(OUTDIR):
	@mkdir -p $@

%.pb.go: %.proto
	@protoc -I. \
 		-I$(PWD)/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
 		--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
 		$<

%.pb.gw.go: %.proto
	@protoc -I. \
		-I$(PWD)/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
 		--grpc-gateway_out=logtostderr=true:. \
 		$<

api: $(PROTO_SRC)
	@echo "Done"

codegen: FORCE
	@swagger-codegen generate --api-package api -i swagger.yaml -l go-server -o .
	#@swagger generate server --exclude-main --flag-strategy=pflag -t $(APIDIR) -f swagger.yaml

clean:
	rm $(BIN)

FORCE:
