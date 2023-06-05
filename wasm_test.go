package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func TestWasmAdd(t *testing.T) {
	ctx := context.Background()
	// JIT模式,性能一般,但是兼容性好
	//r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter())
	// AOT编译模式,性能好,存在操作系统兼容性问题
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	bytes, err := os.ReadFile(datadir + "/wasm/add.wasm")
	if err != nil {
		t.Error(err)
	}
	// Instantiate WASI, which implements host functions needed for TinyGo to
	// implement `panic`.
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.Instantiate(ctx, bytes)
	if err != nil {
		t.Error(err)
	}
	res, err := mod.ExportedFunction("add").Call(ctx, 100, 200)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res[0])
}
