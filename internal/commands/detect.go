package commands

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/elastic/go-sysinfo"
	"github.com/matishsiao/goInfo"
	"github.com/shirou/gopsutil/host"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/pkg/pretty"
)

func is64Bit() bool {
	return ^uint(0)>>32 != 0
}

func (app *App) detect() {
	fmt.Println("----- runtime -----")

	fmt.Printf("runtime.GOOS: %q\n", runtime.GOOS)
	fmt.Printf("runtime.GOARCH: %q\n", runtime.GOARCH)

	fmt.Println("----- ------ -----")

	fmt.Println("----- shirou/gopsutil/host -----")

	info, _ := host.Info()
	fmt.Printf("OS: %q\n", info.OS)
	fmt.Printf("Platform: %q\n", info.Platform)
	fmt.Printf("KernelArch: %q\n", info.KernelArch)
	fmt.Printf("KernelVersion: %q\n", info.KernelVersion)

	pretty.PrintJSON(info)

	fmt.Println("----- ----------------- -----")

	fmt.Println("----- uname -----")

	architecture, _ := getArchitecture()

	fmt.Printf("Architecture: %q\n", architecture)

	fmt.Println("----- ----- -----")

	fmt.Println("----- matishsiao/goInfo -----")

	gi, _ := goInfo.GetInfo()

	fmt.Printf("OS: %q\n", gi.GoOS)
	fmt.Printf("Core: %q\n", gi.Core)
	fmt.Printf("Platform: %q\n", gi.Platform)
	fmt.Printf("Kernel: %q\n", gi.Kernel)

	fmt.Println("----- ----------------- -----")

	fmt.Println("----- elastic/go-sysinfo -----")

	host, _ := sysinfo.Host()

	einfo := host.Info()

	fmt.Printf("OS: %q\n", einfo.OS.Name)
	fmt.Printf("Architecture: %q\n", einfo.Architecture)
	fmt.Printf("NativeArchitecture: %q\n", einfo.NativeArchitecture)
	fmt.Printf("KernelVersion: %q\n", einfo.KernelVersion)

	fmt.Println("----- ----------------- -----")

	fmt.Println("----- is64Bit -----")

	fmt.Printf("is64Bit: %t\n", is64Bit())

	fmt.Println("----- ------ -----")

	fmt.Println("----- strconv.IntSize -----")

	if strconv.IntSize == 32 {
		fmt.Println("32-bit")
	}
	if strconv.IntSize == 64 {
		fmt.Println("64-bit")
	}

	fmt.Println("----- -------------- -----")

	fmt.Println("----- internal/detect/platform -----")

	detect := detect.Platform{}

	if err := detect.Detect(); err != nil {
		fmt.Println(err)
	}

	pretty.PrintJSON(detect)

	fmt.Println("----- ------ -----")

	// fmt.Println("----- golang.org/x/sys/cpu:ARM -----")

	// fmt.Println("cpu.Initialized", cpu.Initialized)
	// pretty.PrintJSON(cpu.ARM)

	// fmt.Println("----- ----------------- -----")

	// fmt.Println("----- golang.org/x/sys/cpu:ARM64 -----")

	// fmt.Println("cpu.Initialized", cpu.Initialized)
	// pretty.PrintJSON(cpu.ARM64)

	// fmt.Println("----- ----------------- -----")
}

func getArchitecture() (string, error) {
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
