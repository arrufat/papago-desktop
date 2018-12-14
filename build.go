// +build ignore

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	win64CC   string = "x86_64-w64-mingw32-gcc"
	win64CXX  string = "x86_64-w64-mingw32-g++"
	darwinCC  string = "clang"
	darwinCXX string = "clang-++"
)

func buildWindows() {
	os.Setenv("GOOS", "windows")
	os.Setenv("CC", win64CC)
	os.Setenv("CXX", win64CXX)
	os.Setenv("CGO_ENABLED", "1")
	cmd := exec.Command("go", "build", "-ldflags", "-H=windowsgui", "-o", "papago.exe", "main.go")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		log.Fatal(err)
	}
	fmt.Println(stdout.String())
}

// unable to build from Linux
func buildDarwin() {
	os.Setenv("GOOS", "darwin")
	os.Setenv("CC", darwinCC)
	os.Setenv("CXX", darwinCXX)
	os.Setenv("CGO_ENABLED", "1")
	os.Setenv("CGO_LDFLAGS", "-shared")
	cmd := exec.Command("go", "build", "-o", "papago.app", "main.go")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		log.Fatal(err)
	}
	fmt.Println(stdout.String())

}

func main() {
	buildWindows()
	buildDarwin()
}
