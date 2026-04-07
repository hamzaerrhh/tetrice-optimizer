package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	dir := "../example"

	runAll := false

	// Parse args
	for _, arg := range os.Args[1:] {
		if arg == "-a" || arg == "--all" {
			runAll = true
		}
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	fmt.Println("start testing 1")

	allPassed := true

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		filePath := filepath.Join(dir, name)

		fmt.Printf("Testing: %s\n", name)

		start := time.Now()

		cmd := exec.Command("go", "run", "../main.go", filePath)
		outputBytes, err := cmd.CombinedOutput()

		duration := time.Since(start)
		output := strings.TrimSpace(string(outputBytes))

		if err != nil && output == "" {
			fmt.Println("❌ Program crashed")
			allPassed = false
			if !runAll {
				os.Exit(1)
			}
			continue
		}

		// =========================
		// BAD FILES → expect ERROR
		// =========================
		if isBadFile(name) {
			if output != "ERROR" {
				fmt.Println("❌ Expected ERROR")
				fmt.Println("Got:", output)
				allPassed = false

				if !runAll {
					os.Exit(1)
				}
			} else {
				fmt.Println("✅ Passed")
			}
			continue
		}

		// =========================
		// GOOD FILES → check dots
		// =========================
		dots := strings.Count(output, ".")
		expectedDots := expectedDotCount(name)

		if dots != expectedDots {
			fmt.Printf("❌ Wrong number of dots. Expected %d, got %d\n", expectedDots, dots)
			allPassed = false

			if !runAll {
				os.Exit(1)
			}
			continue
		}

		// =========================
		// TIME CHECK
		// =========================
		if isTimedCase(name) {
			if duration > time.Second*30 {
				fmt.Printf("❌ Too slow: %v\n", duration)
				allPassed = false

				if !runAll {
					os.Exit(1)
				}
				continue
			}
		}

		fmt.Printf("✅ Passed (dots=%d, time=%v)\n", dots, duration)
	}

	if allPassed {
		fmt.Println("test pass successfully")
	} else {
		fmt.Println("some tests failed")
	}
}

func isBadFile(name string) bool {
	return strings.HasPrefix(name, "bad")
}

func expectedDotCount(name string) int {
	switch name {
	case "goodexample00.txt":
		return 0
	case "goodexample01.txt":
		return 9
	case "goodexample02.txt":
		return 4
	case "goodexample03.txt":
		return 5
	case "hardexam.txt":
		return 1
	default:
		return -1
	}
}

func isTimedCase(name string) bool {
	return name == "goodexample02.txt" ||
		name == "goodexample03.txt" ||
		name == "hardexam.txt"
}
