package main

import (
	"os"
	"testing"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/curl-li/protoc-gen-validate/module"
)

func TestGenerator(t *testing.T) {
	f, err := os.Open("./t.out")
	if err != nil {
		panic(err)
	}
	optional := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	pgs.
		Init(pgs.DebugEnv("DEBUG_PGV"), pgs.SupportedFeatures(&optional), pgs.ProtocInput(f), pgs.DebugMode()).
		RegisterModule(module.Validator()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
