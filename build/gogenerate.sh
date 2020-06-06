
export GO111MODULE=off
export GOFLAGS=

go get -u github.com/fjl/gencodec
go get -u golang.org/x/tools/cmd/stringer
go get -u github.com/go-bindata/go-bindata/...

export GO111MODULE=on
export GOFLAGS=-mod=vendor

go generate range/core/gen3/core/types
go generate range/core/gen3/core/vm
go generate range/core/gen3/core
go generate range/core/gen3/eth/tracers/internal/tracers/
go generate range/core/gen3/eth/
go generate range/core/gen3/internal/jsre/deps/
go generate range/core/gen3/p2p/discv5
go generate range/core/gen3/signer/rules/deps
go generate range/core/gen3/whisper/whisperv6/

