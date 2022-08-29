#!/usr/bin/env sh
protoc \
  --proto_path=. \
  --go_out=. \
  --go_opt=module=github.com/go-leo/example \
  --go-grpc_out=. \
  --go-grpc_opt=module=github.com/go-leo/example \
  --go-leo_out=. \
  --go-leo_opt=module=github.com/go-leo/example \
  api/*/*.proto
