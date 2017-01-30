# This is how we want to name the binary output
OUTDIR=$(PWD)/build/bin
BIN=$(OUTDIR)/ethermis
APIDIR=api
# These are the values we want to pass for Version and BuildTime

GIT_COMMIT=`git rev-parse HEAD`
BUILD_TIME=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X constant.GitCommit=${GIT_COMMIT} -X constant.BuildTime=${BUILD_TIME}"

all: $(OUTDIR)
	@go build ${LDFLAGS} -o ${BIN} main.go

$(OUTDIR):
	@mkdir -p $@

api: FORCE
	@swagger generate server --exclude-main --flag-strategy=pflag -t $(APIDIR) -f swagger.yaml

codegen: FORCE
	@swagger-codegen generate --api-package api -i swagger.yaml -l go-server -o .

clean:
	rm $(BIN)

FORCE:
