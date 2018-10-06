# gominhook

`gominhook` is a Golang wrapper of [minhook](https://github.com/TsudaKageyu/minhook).

This stuff heavily relies on `cgo`.

- - -

### Installation

1. Install `gominhook`.
    
    `go get -v github.com/nanitefactory/gominhook`

2. Have `MinHook.x64.dll` with your project.
    
    You can get it from either of these sites:

    - [gominhook/MinHook_133_bin/bin/MinHook.x64.dll](./MinHook_133_bin/bin/MinHook.x64.dll)

    - https://github.com/TsudaKageyu/minhook/releases

3. That's it!

    `import "github.com/nanitefactory/gominhook"`

- - -

### Exports

Here is a list of all supported functions & data by `gominhook`,

```Go
func gominhook.Initialize() error
func gominhook.Uninitialize() error
```

```Go
func gominhook.CreateHook(pTarget, pDetour, ppOriginal uintptr) error
func gominhook.CreateHookAPI(strModule, strProcName string, pDetour, ppOriginal uintptr) error
func gominhook.CreateHookAPIEx(strModule, strProcName string, pDetour, ppOriginal, ppTarget uintptr) error
```

```Go
func gominhook.RemoveHook(pTarget uintptr) error
func gominhook.EnableHook(pTarget uintptr) error
func gominhook.DisableHook(pTarget uintptr) error
```

```Go
func gominhook.QueueEnableHook(pTarget uintptr) error
func gominhook.QueueDisableHook(pTarget uintptr) error
func gominhook.ApplyQueued() error
```

```Go
const gominhook.AllHooks = NULL
const gominhook.NULL = 0
```

which is straightforward & effective enough. xD

- - -

### Sample

This example below tries to hook `user32.MessageBoxW`.

```Go
package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"github.com/nanitefactory/gominhook"
)

/*
#include <Windows.h>

// Put C prototypes here

// Delegate type for calling original MessageBoxW.
typedef int (WINAPI *MESSAGEBOXW)(HWND, LPCWSTR, LPCWSTR, UINT);

// (!) This way you can connect/convert a go function to a c function.
int WINAPI MessageBoxWOverrideHellYeah(HWND hWnd, LPCWSTR lpText, LPCWSTR lpCaption, UINT uType);
*/
import "C"

// Pointer for calling original MessageBoxW.
var fpMessageBoxW C.MESSAGEBOXW

// (!) This way you can connect/convert a go function to a c function.
//export MessageBoxWOverrideHellYeah
func MessageBoxWOverrideHellYeah(hWnd C.HWND, lpText C.LPCWSTR, lpCaption C.LPCWSTR, uType C.UINT) C.int {
	fmt.Println(" - MessageBoxW Override")
	foo()
	ret, _, _ := syscall.Syscall6(
		uintptr(unsafe.Pointer(fpMessageBoxW)),
		4,
		uintptr(unsafe.Pointer(hWnd)),
		uintptr(unsafe.Pointer(lpText)),
		uintptr(unsafe.Pointer(lpCaption)),
		uintptr(uint(uType)),
		0, 0,
	)
	return C.int(ret)
}

func foo() {
	fmt.Println(" - I'm so hooked now.")
}

func main() {
	// Initialize minhook
	err := gominhook.Initialize()
	if err != nil {
		log.Fatalln(err)
	}
	defer gominhook.Uninitialize()

	// Get procedure user32.MessageBoxW
	procedure := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW")
	fmt.Println("-- not hooked yet")
	procedure.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Hello1"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("World1"))),
		1,
	)
	fmt.Println(fmt.Sprintf("0x%X", procedure.Addr()), fmt.Sprintf("0x%X", &fpMessageBoxW), fmt.Sprintf("0x%X", fpMessageBoxW))
	fmt.Println()

	// Create a hook for MessageBoxW.
	err = gominhook.CreateHook(procedure.Addr(), uintptr(C.MessageBoxWOverrideHellYeah), uintptr(unsafe.Pointer(&fpMessageBoxW)))
	if err != nil {
		log.Fatalln(err)
	}

	// Enable the hook for MessageBoxW.
	err = gominhook.EnableHook(gominhook.AllHooks)
	if err != nil {
		log.Fatalln(err)
	}

	// Calling our hooked procedure user32.MessageBoxW.
	fmt.Println("-- after hook")
	procedure.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Hello2"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("World2"))),
		1,
	)
	fmt.Println(fmt.Sprintf("0x%X", procedure.Addr()), fmt.Sprintf("0x%X", &fpMessageBoxW), fmt.Sprintf("0x%X", fpMessageBoxW))
	fmt.Println()

	// Disable the hook for MessageBoxW.
	err = gominhook.DisableHook(gominhook.AllHooks)
	if err != nil {
		log.Fatalln(err)
	}

	// Calling our unhooked procedure user32.MessageBoxW.
	fmt.Println("-- after unhook")
	procedure.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Hello3"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("World3"))),
		1,
	)
	fmt.Println(fmt.Sprintf("0x%X", procedure.Addr()), fmt.Sprintf("0x%X", &fpMessageBoxW), fmt.Sprintf("0x%X", fpMessageBoxW))
	fmt.Println()
}

/* This outputs...

-- not hooked yet
0x7FFE6CA4EE10 0x578180 0x0

-- after hook
 - MessageBoxW Override
 - I'm so hooked now.
0x7FFE6CA4EE10 0x578180 0x&

-- after unhook
0x7FFE6CA4EE10 0x578180 0x&

*/
```

- - -

### More information

See `minhook` ref for C/C++ users.
- https://github.com/TsudaKageyu/minhook
- https://github.com/TsudaKageyu/minhook/blob/master/include/MinHook.h
- https://github.com/TsudaKageyu/minhook/wiki
- https://www.codeproject.com/Articles/44326/MinHook-The-Minimalistic-x-x-API-Hooking-Libra

- - -
