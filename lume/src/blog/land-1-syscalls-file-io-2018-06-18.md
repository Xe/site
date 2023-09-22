---
title: "Land 1: Syscalls & File I/O"
date: 2018-06-18
series: olin
tags:
 - wasm
---

[Webassembly][wasm] is a new technology aimed at being a vendor-independent virtual machine format. It has implementations by all major browser vendors. It looks like there's staying power in webassembly that could outlast staying power in other new technologies.

So, time's perfect to snipe it with something useful that you can target compilers to today. Hence: [Land][land].

Computer programs are effectively a bunch of business logic around function calls that can affect the "real world" outside of the program. These "magic" functions are also known as [system calls][syscall]. Here's an example of a few in C style syntax:

```
int close(int file);
int open(const char *name, int flags);
int read(int file, char *ptr, int len);
int write(int file, char *ptr, int len);
```

These are all fairly low-level file I/O operations (we're not dealing with structures for now, those are for another day) that all also are (simplified forms of) system calls [like the ones the kernel uses][linux-syscalls].

Effectively, the system calls of a program form the "API" with it and the rest of the computer. Commonly this is called the [ABI (Applcation Binary Interface)][abi] and is usually platform-specific. With Land, we are effectively creating a platform-independent ABI that just so happens to target Webassembly.

In Land, we can wrap an [Afero][afero] filesystem, a set of files (file descriptors are addresses in this set), a webassembly virtual machine, its related webassembly module and its filename into a Process. This process will also have some functions on it to access the resources in it, aimed at being used by the webassembly guest code. In Land, we define this like such:

```
// Process is a larger level wrapper around a webassembly VM that gives it
// system call access.
type Process struct {
	id    int32
	vm    *exec.VM
	mod   *wasm.Module
	fs    afero.Fs
	files []afero.File
	name  string
}
```

Creating a new process is done in the NewProcess function:

```
// NewProcess constructs a new webassembly process based on the input webassembly module as a reader.
func NewProcess(fin io.Reader, name string) (*Process, error) {
	p := &Process{}

	mod, err := wasm.ReadModule(fin, p.importer)
	if err != nil {
		return nil, err
	}

	if mod.Memory == nil {
		return nil, errors.New("must declare a memory, sorry :(")
	}

	vm, err := exec.NewVM(mod)
	if err != nil {
		return nil, err
	}

	p.mod = mod
	p.vm = vm
	p.fs = afero.NewMemMapFs()

	return p, nil
}
```

The webassembly importer makes a [little shim module for importing host functions][importer] (not inlined due to size).

Memory operations are implemented on top of each WebAssembly process. The two most basic ones are `writeMem` and `readMem`:

```
// writeMem writes the given data to the webassembly memory at the given pointer offset.
func (p *Process) writeMem(ptr int32, data []byte) (int, error) {
	mem := p.vm.Memory()
	if mem == nil {
		return 0, errors.New("no memory, invalid state")
	}

	for i, d := range data {
		mem[ptr+int32(i)] = d
	}

	return len(data), nil
}

// readMem reads memory at the given pointer until a null byte of ram is read.
// This is intended for reading Cstring-like structures that are terminated
// with null bytes.
func (p *Process) readMem(ptr int32) []byte {
	var result []byte

	mem := p.vm.Memory()[ptr:]
	for _, bt := range mem {
		if bt == 0 {
			return result
		}

		result = append(result, bt)
	}

	return result
}
```

Every system call that deals with C-style strings uses these functions to get arguments out of the WebAssembly virtual machine's memory and to put the results back into the WebAssembly virtual machine.

Below is the [`open(2)`][open2] implementation for Land. It implements the following C-style function type:

```
int open(const char *name, int flags);
```

WebAssembly natively deals with integer and floating point types, so the first argument is the pointer to the memory in WebAssembly linear memory. The second is an integer as normal. The code handles this as such:

```
func (p *Process) open(fnamesP int32, flags int32) int32 {
	str := string(p.readMem(fnamesP))

	fi, err := p.fs.OpenFile(string(str), int(flags), 0666)
	if err != nil {
		if strings.Contains(err.Error(), afero.ErrFileNotFound.Error()) {
			fi, err = p.fs.Create(str)
		}
	}

	if err != nil {
		panic(err)
	}

	fd := len(p.files)
	p.files = append(p.files, fi)

	return int32(fd)
}
```

As you can see, the integer arguments can sufficiently represent the datatype of C: machine words. String pointers are machine words. Integers are machine words. Everything is machine words.

Write is very simple to implement. Its type gives us a bunch of advantages out of the gate:

```
int write(int file, char *ptr, int len);
```

This gives us the address of where to start in memory, and adding the length to the address gives us the end in memory:

```
func (p *Process) write(fd int32, ptr int32, len int32) int32 {
	data := p.vm.Memory()[ptr : ptr+len]
	n, err := p.files[fd].Write(data)
	if err != nil {
		panic(err)
	}

	return int32(n)
}
```

Read is also simple. The type of it gives us a hint on how to implement it:

```
int read(int file, char *ptr, int len);
```

We are going to need a buffer at least as large as `len` to copy data from the file to the WebAssembly process. Implementation is then simply:

```
func (p *Process) read(fd int32, ptr int32, len int32) int32 {
	data := make([]byte, len)
	na, err := p.files[fd].Read(data)
	if err != nil {
		panic(err)
	}

	nb, err := p.writeMem(ptr, data)
	if err != nil {
		panic(err)
	}

	if na != nb {
		panic("did not copy the same number of bytes???")
	}

	return int32(na)
}
```

Close lets us let go of files we don't need anymore. This will also have to have a special case to clear out the last file properly when there's only one file open:

```
func (p *Process) close(fd int32) int32 {
	f := p.files[fd]
	err := f.Close()
	if err != nil {
		panic(err)
	}

	if len(p.files) == 1 {
		p.files = []afero.File{}
	} else {
		p.files = append(p.files[:fd], p.files[fd+1])
	}

	return 0
}
```

These calls are enough to make surprisingly nontrivial programs, considering standard input and standard output exist, but here's an example of a trivial program made with some of these calls (equivalent C-like shown too):

```
(module
 ;; import functions from env
 (func $close (import "env" "close") (param i32)         (result i32))
 (func $open  (import "env" "open")  (param i32 i32)     (result i32))
 (func $read  (import "env" "read")  (param i32 i32 i32) (result i32))
 (func $write (import "env" "write") (param i32 i32 i32) (result i32))

 ;; memory
 (memory $mem 1)

 ;; constants
 (data (i32.const 200) "data")
 (data (i32.const 230) "Hello, world!\n")

 ;; land looks for a function named main that returns a 32 bit integer.
 ;; int $main() {
 (func $main (result i32)
       ;; $fd is the file descriptor of the file we're gonna open
       (local $fd i32)

       ;; $fd = $open("data", O_CREAT|O_RDWR);
       (set_local $fd
                  (call $open
                        ;; pointer to the file name
                        (i32.const 200)
                        ;; flags, 42 for O_CREAT,O_RDWR
                        (i32.const 42)))

       ;; $write($fd, "Hello, World!\n", 14);
       (call $write
             (get_local $fd)
             (i32.const 230)
             (i32.const 14))
       (drop)

       ;; $close($fd);
       (call $close
             (get_local $fd))
       (drop)

       (i32.const 0))
 ;; }
 (export "main" (func $main)))
```

This can be verified outside of the WebAssembly environment, I tested mine with the [pretty package][pretty].

Right now this is very lean and mean, as such all errors instantly result in a panic which will kill the WebAssembly VM. I would like to fix this but I will need to make sure that programs don't use certain bits of memory where Land will communicate with the WebAssembly module. Other good steps are going to be setting up reserved areas of memory for things like error messages, [posix errno][errno] and other superglobals.

A huge other feature is going to be the ability to read C structures out of the WebAssembly memory, this will let Land support calls like `stat()`.

[wasm]: https://webassembly.org
[land]: https://tulpa.dev/cadey/land
[syscall]: https://en.wikipedia.org/wiki/System_call
[linux-syscalls]: https://www.nullmethod.com/syscall-table/
[abi]: https://en.m.wikipedia.org/wiki/Application_binary_interface
[afero]: https://github.com/spf13/afero
[importer]: https://gist.github.com/Xe/a29c86755a04a8096082ec8a32e0c13f
[open2]: https://linux.die.net/man/2/open
[pretty]: https://github.com/kr/pretty
[errno]: https://man7.org/linux/man-pages/man3/errno.3.html
