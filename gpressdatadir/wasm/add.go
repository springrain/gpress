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

//tinygo编译wasm命令: tinygo build -o add.wasm -scheduler=none --no-debug -target=wasi .
// 使用 tinygo  go:export add

// go 1.24+ 编译wasm命令: GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o add.wasm
// go 1.24+ windwos编译wasm命令: set GOOS=wasip1&&set GOARCH=wasm&&go build -buildmode=c-shared -o add.wasm

//go:wasmexport add
func add(x, y uint32) uint32 {
	return x + y
}

// main is required for the `wasi` target, even if it isn't used.
// See https://wazero.io/languages/tinygo/#why-do-i-have-to-define-main
func main() {}
