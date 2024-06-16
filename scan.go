package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"syscall"
	"unsafe"

	"github.com/shirou/gopsutil/process"
)

const processName = "starbase.exe"

var handle syscall.Handle

// Custom byte slice comparison pattern
var pattern = []byte{
	0xFC, '?', '?', '?', '?', '?', 0x00, 0x00,
	0x0E, '?', '?', '?', '?', '?', 0x00, 0x00,
	0x24, '?', '?', '?', '?', '?', 0x00, 0x00,
	0x3A, '?', '?', '?', '?', '?', 0x00, 0x00,
	0x53, '?', '?', '?', '?', '?', 0x00, 0x00,
	0x6E, '?', '?', '?', '?', '?', 0x00, 0x00,
	0x89, '?', '?', '?', '?', '?', 0x00, 0x00,
	'?', '?', '?', '?',
}

const patternSize = 60 // Adjusted size of the pattern we are looking for
const windowSize = 60  // Window size for moving window approach
const stepSize = 4     // Step size for moving window

const (
	PROCESS_VM_READ           = 0x0010
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_OPERATION      = 0x0008
	bufferSize                = 4096
	MEM_COMMIT                = 0x1000
	PAGE_GUARD                = 0x100
	PAGE_NOACCESS             = 0x01
)

var (
	kernel32              = syscall.MustLoadDLL("kernel32.dll")
	procReadProcessMemory = kernel32.MustFindProc("ReadProcessMemory")
	procVirtualQueryEx    = kernel32.MustFindProc("VirtualQueryEx")
)

type MEMORY_BASIC_INFORMATION struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

type FoundValue struct {
	Address    uintptr
	Value      uint32
	FloatValue float32
}

func ReadProcessMemory(hProcess syscall.Handle, lpBaseAddress uintptr, lpBuffer *byte, nSize uintptr, lpNumberOfBytesRead *uintptr) (err error) {
	r1, _, e1 := syscall.SyscallN(procReadProcessMemory.Addr(), uintptr(hProcess), lpBaseAddress, uintptr(unsafe.Pointer(lpBuffer)), nSize, uintptr(unsafe.Pointer(lpNumberOfBytesRead)))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func VirtualQueryEx(hProcess syscall.Handle, lpAddress uintptr, lpBuffer *MEMORY_BASIC_INFORMATION, dwLength uintptr) (uintptr, error) {
	r1, _, e1 := syscall.SyscallN(procVirtualQueryEx.Addr(), uintptr(hProcess), lpAddress, uintptr(unsafe.Pointer(lpBuffer)), dwLength)
	if r1 == 0 {
		if e1 != 0 {
			return 0, error(e1)
		} else {
			return 0, syscall.EINVAL
		}
	}
	return r1, nil
}

func matchPattern(data, pattern []byte) bool {
	if len(data) < len(pattern) {
		return false
	}
	for i := 0; i < len(pattern); i++ {
		if pattern[i] != '?' && data[i] != pattern[i] {
			return false
		}
	}
	return true
}

func scanMemory(handle syscall.Handle, memoryRegions []MEMORY_BASIC_INFORMATION, results chan<- FoundValue, wg *sync.WaitGroup) {
	defer wg.Done()

	buffer := make([]byte, bufferSize+windowSize)
	var bytesRead uintptr

	for _, region := range memoryRegions {
		for addr := region.BaseAddress; addr < region.BaseAddress+region.RegionSize; addr += bufferSize {
			err := ReadProcessMemory(handle, addr, &buffer[0], bufferSize+windowSize, &bytesRead)
			if err != nil {
				continue
			}

			for i := 0; i < int(bytesRead)-windowSize; i += stepSize {
				window := buffer[i : i+windowSize]

				if matchPattern(window, pattern) {
					extractedBytes := window[len(window)-4:] // Grab the last 4 bytes of the window
					var value uint32
					binary.Read(bytes.NewReader(extractedBytes), binary.LittleEndian, &value)
					floatValue := math.Float32frombits(value)
					results <- FoundValue{
						Address:    addr + uintptr(i) + uintptr(len(window)-4),
						Value:      value,
						FloatValue: floatValue,
					}
				}
			}
		}
	}
}

func getScanResults() ([]FoundValue, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("error getting processes: %v", err)
	}

	var targetProcess *process.Process
	for _, proc := range procs {
		name, err := proc.Name()
		if err == nil && name == processName {
			targetProcess = proc
			break
		}
	}

	if targetProcess == nil {
		return nil, fmt.Errorf("process %s not found", processName)
	}

	pid := targetProcess.Pid
	fmt.Printf("Found process %s with PID %d\n", processName, pid)

	handle, err = syscall.OpenProcess(PROCESS_VM_READ|PROCESS_QUERY_INFORMATION|PROCESS_VM_OPERATION, false, uint32(pid))
	if err != nil {
		return nil, fmt.Errorf("error opening process: %v", err)
	}
	//defer syscall.CloseHandle(handle)

	var memoryInfo MEMORY_BASIC_INFORMATION
	baseAddress := uintptr(0x010000000000)
	endAddress := uintptr(0x02FFFFFFFFFF) // User-mode address space
	var memoryRegions []MEMORY_BASIC_INFORMATION
	for baseAddress < endAddress {
		_, err := VirtualQueryEx(handle, baseAddress, &memoryInfo, unsafe.Sizeof(memoryInfo))
		if err != nil {
			baseAddress += 0x1000 // Skip to next page on error
			continue
		}
		if memoryInfo.State == MEM_COMMIT && (memoryInfo.Protect&PAGE_GUARD == 0) && (memoryInfo.Protect&PAGE_NOACCESS == 0) {
			memoryRegions = append(memoryRegions, memoryInfo)
		}
		baseAddress += memoryInfo.RegionSize
	}

	var wg sync.WaitGroup
	results := make(chan FoundValue, len(memoryRegions)*bufferSize/patternSize)
	numGoroutines := 16
	chunkSize := (len(memoryRegions) + numGoroutines - 1) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(memoryRegions) {
			end = len(memoryRegions)
		}

		wg.Add(1)
		go scanMemory(handle, memoryRegions[start:end], results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var foundValues []FoundValue
	for result := range results {
		foundValues = append(foundValues, result)
	}

	return foundValues, nil
}

func GetFloat32ValueAtAddress(address uintptr) (float32, error) {
	var value float32
	buffer := make([]byte, 4)
	var bytesRead uintptr

	err := ReadProcessMemory(handle, address, &buffer[0], uintptr(len(buffer)), &bytesRead)
	if err != nil {
		fmt.Println("Error reading memory:", err)
		return 0, fmt.Errorf("failed to read memory at address %x: %v", address, err)
	}

	// print raw bytes
	// fmt.Printf("Raw bytes: %x\n", buffer)

	bits := binary.LittleEndian.Uint32(buffer)
	value = math.Float32frombits(bits)

	return value, nil
}
