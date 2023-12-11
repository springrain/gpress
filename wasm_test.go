// Copyright (c) 2023 gpress Authors.
//
// This file is part of gpress.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
