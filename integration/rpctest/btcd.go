// Copyright (c) 2017 The btcsuite developers
// Copyright (c) 2018 The Flo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpctest

import (
	"fmt"
	"go/build"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	// compileMtx guards access to the executable path so that the project is
	// only compiled once.
	compileMtx sync.Mutex

	// executablePath is the path to the compiled executable. This is the empty
	// string until flod is compiled. This should not be accessed directly;
	// instead use the function flodExecutablePath().
	executablePath string
)

// flodExecutablePath returns a path to the flod executable to be used by
// rpctests. To ensure the code tests against the most up-to-date version of
// flod, this method compiles flod the first time it is called. After that, the
// generated binary is used for subsequent test harnesses. The executable file
// is not cleaned up, but since it lives at a static path in a temp directory,
// it is not a big deal.
func flodExecutablePath() (string, error) {
	compileMtx.Lock()
	defer compileMtx.Unlock()

	// If flod has already been compiled, just use that.
	if len(executablePath) != 0 {
		return executablePath, nil
	}

	testDir, err := baseDir()
	if err != nil {
		return "", err
	}

	// Determine import path of this package. Not necessarily bitspill/flod if
	// this is a forked repo.
	_, rpctestDir, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("Cannot get path to flod source code")
	}
	flodPkgPath := filepath.Join(rpctestDir, "..", "..", "..")
	flodPkg, err := build.ImportDir(flodPkgPath, build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("Failed to build flod: %v", err)
	}

	// Build flod and output an executable in a static temp path.
	outputPath := filepath.Join(testDir, "flod")
	if runtime.GOOS == "windows" {
		outputPath += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", outputPath, flodPkg.ImportPath)
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Failed to build flod: %v", err)
	}

	// Save executable path so future calls do not recompile.
	executablePath = outputPath
	return executablePath, nil
}
