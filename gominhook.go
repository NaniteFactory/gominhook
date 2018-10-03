package gominhook

import (
	"errors"
	"syscall"
	"unsafe"
)

/*
#cgo CFLAGS: -I${SRCDIR}/MinHook_133_bin/include
#cgo LDFLAGS: -L${SRCDIR}/MinHook_133_bin/bin -lMinHook.x64
#include "MinHook.h"

#if defined _M_X64
#pragma comment(lib, "libMinHook.x64.lib")
#elif defined _M_IX86
#pragma comment(lib, "libMinHook.x86.lib")
#endif
*/
import "C"

// Status for MH_STATUS
// MinHook Error Codes.
type _Status int

// const for enum STAUS
// MinHook Error Codes.
const (
	_Unknown                  _Status = C.MH_UNKNOWN                    // -1	// Unknown error. Should not be returned.
	_OK                       _Status = C.MH_OK                         // 0	// Successful.
	_ErrorAlreadyInitialized  _Status = C.MH_ERROR_ALREADY_INITIALIZED  // 1	// MinHook is already initialized.
	_ErrorNotInitialized      _Status = C.MH_ERROR_NOT_INITIALIZED      // 2	// MinHook is not initialized yet, or already uninitialized.
	_ErrorAlreadyCreated      _Status = C.MH_ERROR_ALREADY_CREATED      // 3	// The hook for the specified target function is already created.
	_ErrorNotCreated          _Status = C.MH_ERROR_NOT_CREATED          // 4	// The hook for the specified target function is not created yet.
	_ErrorEnabled             _Status = C.MH_ERROR_ENABLED              // 5	// The hook for the specified target function is already enabled.
	_ErrorDisabled            _Status = C.MH_ERROR_DISABLED             // 6	// The hook for the specified target function is not enabled yet, or already disabled.
	_ErrorNotExecutable       _Status = C.MH_ERROR_NOT_EXECUTABLE       // 7	// The specified pointer is invalid. It points the address of non-allocated and/or non-executable region.
	_ErrorUnsupportedFunction _Status = C.MH_ERROR_UNSUPPORTED_FUNCTION // 8	// The specified target function cannot be hooked.
	_ErrorMemoryAlloc         _Status = C.MH_ERROR_MEMORY_ALLOC         // 9	// Failed to allocate memory.
	_ErrorMemoryProtect       _Status = C.MH_ERROR_MEMORY_PROTECT       // 10	// Failed to change the memory protection.
	_ErrorModuleNotFound      _Status = C.MH_ERROR_MODULE_NOT_FOUND     // 11	// The specified module is not loaded.
	_ErrorFunctionNotFound    _Status = C.MH_ERROR_FUNCTION_NOT_FOUND   // 12	// The specified function is not found.
)

// ToError () for const char * WINAPI MH_StatusToString(MH_STATUS status);
// MH_StatusToString() - Translates the MH_STATUS to its name as a string.
// ------------------------------------------------------------------------
// _Status#ToError() - Converts 'MinHook Error Codes' to 'Golang error'.
func (minhookStatus _Status) ToError() error {
	var errMsg string
	switch minhookStatus {
	case _Unknown:
		errMsg = "MH_UNKNOWN"
	case _OK:
		return nil
	case _ErrorAlreadyInitialized:
		errMsg = "MH_ERROR_ALREADY_INITIALIZED"
	case _ErrorNotInitialized:
		errMsg = "MH_ERROR_NOT_INITIALIZED"
	case _ErrorAlreadyCreated:
		errMsg = "MH_ERROR_ALREADY_CREATED"
	case _ErrorNotCreated:
		errMsg = "MH_ERROR_NOT_CREATED"
	case _ErrorEnabled:
		errMsg = "MH_ERROR_ENABLED"
	case _ErrorDisabled:
		errMsg = "MH_ERROR_DISABLED"
	case _ErrorNotExecutable:
		errMsg = "MH_ERROR_NOT_EXECUTABLE"
	case _ErrorUnsupportedFunction:
		errMsg = "MH_ERROR_UNSUPPORTED_FUNCTION"
	case _ErrorMemoryAlloc:
		errMsg = "MH_ERROR_MEMORY_ALLOC"
	case _ErrorMemoryProtect:
		errMsg = "MH_ERROR_MEMORY_PROTECT"
	case _ErrorModuleNotFound:
		errMsg = "MH_ERROR_MODULE_NOT_FOUND"
	case _ErrorFunctionNotFound:
		errMsg = "MH_ERROR_FUNCTION_NOT_FOUND"
	}
	return errors.New(errMsg)
}

// AllHooks for #define MH_ALL_HOOKS NULL
// Can be passed as a parameter to MH_EnableHook, MH_DisableHook, MH_QueueEnableHook or MH_QueueDisableHook.
const AllHooks = NULL

// NULL = 0x00000000
const NULL = 0

// Initialize () for MH_STATUS WINAPI MH_Initialize(VOID)
// Initialize the MinHook library. You must call this function EXACTLY ONCE at the beginning of your program.
func Initialize() (err error) {
	ret, _, _ := syscall.Syscall(uintptr(C.MH_Initialize), 0, 0, 0, 0)
	return _Status(ret).ToError()
}

// Uninitialize () for MH_STATUS WINAPI MH_Uninitialize(VOID)
// Uninitialize the MinHook library. You must call this function EXACTLY ONCE at the end of your program.
func Uninitialize() (err error) {
	ret, _, _ := syscall.Syscall(uintptr(C.MH_Uninitialize), 0, 0, 0, 0)
	return _Status(ret).ToError()
}

// CreateHook () for MH_STATUS WINAPI MH_CreateHook(LPVOID pTarget, LPVOID pDetour, LPVOID *ppOriginal);
// Creates a Hook for the specified target function, in disabled state.
// Parameters:
//   pTarget    [in]  A pointer to the target function, which will be
//                    overridden by the detour function.
//   pDetour    [in]  A pointer to the detour function, which will override
//                    the target function.
//   ppOriginal [out] A pointer to the trampoline function, which will be
//                    used to call the original target function.
//                    This parameter can be NULL.
func CreateHook(pTarget, pDetour, ppOriginal uintptr) (err error) {
	ret, _, _ := syscall.Syscall(uintptr(C.MH_CreateHook), 3, pTarget, pDetour, ppOriginal)
	return _Status(ret).ToError()
}

// CreateHookAPI () for MH_STATUS WINAPI MH_CreateHookApi(LPCWSTR pszModule, LPCSTR pszProcName, LPVOID pDetour, LPVOID *ppOriginal);
// Creates a Hook for the specified API function, in disabled state.
// Parameters:
//   pszModule  [in]  A pointer to the loaded module name which contains the
//                    target function.
//   pszTarget  [in]  A pointer to the target function name, which will be
//                    overridden by the detour function.
//   pDetour    [in]  A pointer to the detour function, which will override
//                    the target function.
//   ppOriginal [out] A pointer to the trampoline function, which will be
//                    used to call the original target function.
//                    This parameter can be NULL.
// ------------------------------------------------------------------------
// strModule: Module name in Go's string. Replace of pszModule.
// strProcName: Procedure (target function) name in Go's string. Replace of pszTarget.
func CreateHookAPI(strModule, strProcName string, pDetour, ppOriginal uintptr) (err error) {
	ret, _, _ := syscall.Syscall6(
		uintptr(C.MH_CreateHookApi),
		4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(strModule))),
		uintptr(unsafe.Pointer(C.CString(strProcName))),
		pDetour,
		ppOriginal,
		0, 0,
	)
	return _Status(ret).ToError()
}

// CreateHookAPIEx () for MH_STATUS WINAPI MH_CreateHookApiEx(LPCWSTR pszModule, LPCSTR pszProcName, LPVOID pDetour, LPVOID *ppOriginal, LPVOID *ppTarget);
// Creates a Hook for the specified API function, in disabled state.
// Parameters:
//   pszModule  [in]  A pointer to the loaded module name which contains the
//                    target function.
//   pszTarget  [in]  A pointer to the target function name, which will be
//                    overridden by the detour function.
//   pDetour    [in]  A pointer to the detour function, which will override
//                    the target function.
//   ppOriginal [out] A pointer to the trampoline function, which will be
//                    used to call the original target function.
//                    This parameter can be NULL.
//   ppTarget   [out] A pointer to the target function, which will be used
//                    with other functions.
//                    This parameter can be NULL.
// ------------------------------------------------------------------------
// strModule: Module name in Go's string. Replace of pszModule.
// strProcName: Procedure (target function) name in Go's string. Replace of pszTarget.
func CreateHookAPIEx(strModule, strProcName string, pDetour, ppOriginal, ppTarget uintptr) (err error) {
	ret, _, _ := syscall.Syscall6(
		uintptr(C.MH_CreateHookApiEx),
		5,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(strModule))),
		uintptr(unsafe.Pointer(C.CString(strProcName))),
		pDetour,
		ppOriginal,
		ppTarget,
		0,
	)
	return _Status(ret).ToError()
}

// RemoveHook () for MH_STATUS WINAPI MH_RemoveHook(LPVOID pTarget);
// Removes an already created hook.
// Parameters:
//   pTarget [in] A pointer to the target function.
func RemoveHook(pTarget uintptr) (err error) {
	ret, _, _ := syscall.Syscall(uintptr(C.MH_RemoveHook), 1, pTarget, 0, 0)
	return _Status(ret).ToError()
}

// EnableHook () for MH_STATUS WINAPI MH_EnableHook(LPVOID pTarget);
// Enables an already created hook.
// Parameters:
//   pTarget [in] A pointer to the target function.
//                If this parameter is MH_ALL_HOOKS, all created hooks are
//                enabled in one go.
// ------------------------------------------------------------------------
// gominhook.AllHooks is equivalent to MH_ALL_HOOKS.
// gominhook.AllHooks can be used as an argument to this function to enable all created hooks.
func EnableHook(pTarget uintptr) (err error) {
	ret, _, _ := syscall.Syscall(uintptr(C.MH_EnableHook), 1, pTarget, 0, 0)
	return _Status(ret).ToError()
}

// DisableHook () for MH_STATUS WINAPI MH_DisableHook(LPVOID pTarget);
// Disables an already created hook.
// Parameters:
//   pTarget [in] A pointer to the target function.
//                If this parameter is MH_ALL_HOOKS, all created hooks are
//                disabled in one go.
// ------------------------------------------------------------------------
// gominhook.AllHooks is equivalent to MH_ALL_HOOKS.
// gominhook.AllHooks can be used as an argument to this function to disable all created hooks.
func DisableHook(pTarget uintptr) (err error) {
	ret, _, _ := syscall.Syscall(uintptr(C.MH_DisableHook), 1, pTarget, 0, 0)
	return _Status(ret).ToError()
}

// QueueEnableHook () for MH_STATUS WINAPI MH_QueueEnableHook(LPVOID pTarget);
// Queues to enable an already created hook.
// Parameters:
//   pTarget [in] A pointer to the target function.
//                If this parameter is MH_ALL_HOOKS, all created hooks are
//                queued to be enabled.
// ------------------------------------------------------------------------
// gominhook.AllHooks is equivalent to MH_ALL_HOOKS.
// gominhook.AllHooks can be used as an argument to this function to queue all created hooks to be enabled.
func QueueEnableHook(pTarget uintptr) (err error) {
	ret, _, _ := syscall.Syscall(uintptr(C.MH_QueueEnableHook), 1, pTarget, 0, 0)
	return _Status(ret).ToError()
}

// QueueDisableHook () for MH_STATUS WINAPI MH_QueueDisableHook(LPVOID pTarget);
// Queues to disable an already created hook.
// Parameters:
//   pTarget [in] A pointer to the target function.
//                If this parameter is MH_ALL_HOOKS, all created hooks are
//                queued to be disabled.
// ------------------------------------------------------------------------
// gominhook.AllHooks is equivalent to MH_ALL_HOOKS.
// gominhook.AllHooks can be used as an argument to this function to queue all created hooks to be disabled.
func QueueDisableHook(pTarget uintptr) (err error) {
	ret, _, _ := syscall.Syscall(uintptr(C.MH_QueueDisableHook), 1, pTarget, 0, 0)
	return _Status(ret).ToError()
}

// ApplyQueued () for MH_STATUS WINAPI MH_ApplyQueued(VOID);
// Applies all queued changes in one go.
func ApplyQueued() (err error) {
	ret, _, _ := syscall.Syscall(uintptr(C.MH_ApplyQueued), 0, 0, 0, 0)
	return _Status(ret).ToError()
}
