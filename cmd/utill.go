package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func npmBuild() error {
	viteCmd := exec.Command("npm", "run", "build")

	viteCmd.Stdout = os.Stdout
	viteCmd.Stderr = os.Stderr

	err := viteCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func getDirs() ([]*WasmDir, error) {
	rootPath := filepath.Join("resource", "wasm")

	paths := []*WasmDir{}

	entries, err := os.ReadDir(rootPath)

	if err != nil {
		return paths, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			return []*WasmDir{}, nil
		}

		name := entry.Name()

		file := &WasmDir{ inputDir: name, outputDir: name + ".wasm" }

		paths = append(paths, file)
	}

	return paths, nil
}

func findWasmExecJS() (string, error) {
	goroot, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get GOROOT: %w", err)
	}
	root := strings.TrimSpace(string(goroot))

	candidates := []string{
		filepath.Join(root, "lib", "wasm", "wasm_exec.js"),
		filepath.Join(root, "misc", "wasm", "wasm_exec.js"),
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("wasm_exec.js not found in GOROOT (%s)", root)
}

func readWasmExecJS() ([]byte, error) {
	path, err := findWasmExecJS()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}