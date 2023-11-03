package main

//编译wasm命令: tinygo build -o add.wasm -scheduler=none --no-debug -target=wasi .

//go:export add
func add(x, y uint32) uint32 {
	return x + y
}

// main is required for the `wasi` target, even if it isn't used.
// See https://wazero.io/languages/tinygo/#why-do-i-have-to-define-main
func main() {}
