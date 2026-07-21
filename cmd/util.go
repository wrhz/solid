package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

func isDir(path string) (bool, error) {
    info, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, nil
        }
        return false, err
    }
    return info.IsDir(), nil
}

func clearDir(dir string) error {
    entries, err := os.ReadDir(dir)
    if err != nil {
        return fmt.Errorf("读取目录失败: %w", err)
    }

    for _, entry := range entries {
        fullPath := filepath.Join(dir, entry.Name())

        if err := os.RemoveAll(fullPath); err != nil {
            return fmt.Errorf("删除 %s 失败: %w", fullPath, err)
        }
    }
    return nil
}

func exportWasm() error {
	wasmDirPath := filepath.Join("resource", "wasm")
	wasmOutputPath := filepath.Join("output", "wasm")

	clearDir(wasmOutputPath)

	is, err := isDir("./" + wasmDirPath)

	if err != nil {
		return err
	}

	if is {
		data, err := readWasmExecJS()

		if err != nil {
			return err
		}

		data = append(data, "\nexport default Go;"...)

		if err = os.WriteFile(filepath.Join("resource", "lib", "wasm_exec.js"), data, 0644); err != nil {
			return err
		}

		dirs, err := getDirs()

		if err != nil {
			return err
		}

		for _, dir := range dirs {
			var buildCmd *exec.Cmd

			inputPath := "./" + filepath.Join(wasmDirPath, dir.inputDir)
			outputPath := filepath.Join(wasmOutputPath, dir.outputDir)

			buildCmd = exec.Command("go", "build", "-o", outputPath, inputPath)

			buildCmd.Stdout = os.Stdout
			buildCmd.Stderr = os.Stderr
			buildCmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")

			if err := buildCmd.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}

func getSubDirNames(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs, nil
}