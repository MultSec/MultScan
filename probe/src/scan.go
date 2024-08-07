package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"runtime"
	"encoding/hex"
	"path/filepath"
	"time"
	//"math/rand"
    "encoding/base64"

	"github.com/mgutz/ansi"
)

func scanFile(binaryPath string, conf map[string]string) (bool, error) {
	// Get absolute path
	absPath, err := filepath.Abs(binaryPath)
    if err != nil {
		return false, fmt.Errorf("failed to get absolute path with error: %v", err)
    }

	// Replace placeholder with actual file path
	cmd := strings.Replace(conf["cmd"], "{{file}}", absPath, -1)

	// Execute the scanner command
	output, err := exec.Command("powershell.exe", "-Command", cmd).CombinedOutput()
	if runtime.GOOS != "windows" {
		return false, fmt.Errorf("program only works on windows")
	}

	// Check if the output contains the positive detection
	if strings.Contains(string(output), conf["out"]) {
		return true, nil
	} else {
		return false, nil
	}
}

func scanSlice(fileData []byte, conf map[string]string) (bool, error) {
	// Create a temp file to scan
	tempFile, err := os.CreateTemp("", "slice_scan_")
	if err != nil {
		return false, fmt.Errorf("failed to create temp file: %v", err)
	}

	// Defer cleanup
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.Write(fileData)
	if err != nil {
		return false, fmt.Errorf("failed to write to temp file: %v", err)
	}

	// Scan the file slice
	scanResult, err := scanFile(tempFile.Name(), conf)
	if err != nil {
		return false, fmt.Errorf("failed to scan temp file: %v", err)
	}

	return scanResult, nil
}

func checkStatic(binaryPath string, conf map[string]string) (string, error) {
	// Read the files content
	data, err := os.ReadFile(binaryPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Set the initial values
	lastGood, mid, upperBound := 0, len(data)/2, len(data)
	threatFound := false

	// Binary search for the malicious content
	for upperBound-lastGood > 1 {
		// Check the slice for malware
		scanResult , err := scanSlice(data[0:mid], conf)
		if err != nil {
			return "", err
		}

		if scanResult {
			threatFound = true
			upperBound = mid
		} else {
			lastGood = mid
		}

		mid = lastGood + (upperBound-lastGood)/2
	}

	// Return the result
	if threatFound {

		// Get the start and end of the malicious content
		start := lastGood - 32
		if start < 0 {
			start = 0
		}

		// Get the start and end of the malicious content
		end := mid + 32
		if end > len(data) {
			end = len(data)
		}

		return fmt.Sprintf("%08x\n%s\n", lastGood, hex.Dump(data[start:end])), nil
	}

	return "", nil
}

func CheckMal(binaryPath string, conf map[string]string) (string, string, error) {
	// Check for Detection
	scanResult, err := scanFile(binaryPath, conf)
	if err != nil {
		return "", "", err
	}
	
	if !scanResult {
		return "Undetected", "", nil
	}

	static, err := checkStatic(binaryPath, conf)
	if err != nil {
		return "", "", err
	}
	
	if static != "" {
		return "Detected", base64.StdEncoding.EncodeToString([]byte(static)), nil
	}

	return "Detected", "V2hvbGUgZmlsZSBpcyBtYWxpY2lvdXM=", nil
}


// Function that scans the sample
func scanSample(id string, scanConf map[string]string) (string,string) {
    printLog(logInfo, fmt.Sprintf("[%s] Scanning sample...", ansi.ColorFunc("default+hb")(id)))
    
    // Retrieve filepath
    filename := fmt.Sprintf("sample_%s.bin", id)
    
    // Get the absolute path
    absPath, err := filepath.Abs(filename)
    if err != nil {
        return "", fmt.Sprintf("failed to get absolute path: %v", err)
    }
    
	// Write program into the public dir and wait 3 minutes to see if it gets deleted
	publicDir := os.Getenv("TEMP") // Using Windows temp directory
	publicPath := filepath.Join(publicDir, filepath.Base(filename))
	err = copyFile(absPath, publicPath)
	if err != nil {
		return "", fmt.Sprintf("failed to copy file to public directory: %v", err)
	}

	time.Sleep(3 * time.Minute)

	// Check if file still exists
	_, err = os.Stat(publicPath)
	if os.IsNotExist(err) {
		// Check for Static detections
		result, badBytes, err := CheckMal(absPath, scanConf)
		if err != nil {
			return "", fmt.Sprintf("failed to scan sample: %v", err)
		}

		return result, badBytes
	}

	// Run the program and wait 3 minutes to see if it gets deleted
	cmd := exec.Command("cmd", "/C", publicPath)
	err = cmd.Start()
	if err != nil {
		return "", fmt.Sprintf("failed to start program: %v", err)
	}

	time.Sleep(3 * time.Minute)

	// Check if file still exists
	_, err = os.Stat(publicPath)
	if os.IsNotExist(err) {
		return "Deleted", "RHluYW1pYyBEZXRlY3Rpb24="
	}

	return "Undetected", ""

	/*
    // Simulate scanning
    time.Sleep(15 * time.Second)
    
    // In a real scenario, you would implement actual scanning logic here
    // For now, we'll use a simple random result
    if rand.Float32() < 0.7 { // 70% chance of detection
        result := "Detected"
        badBytes := "MDAwNDhlM2QKMDAwMDAwMDAgIDY1IDc0IDVmIDYxIDY0IDY0IDY5IDc0ICA2OSA2ZiA2ZSA2MSA2YyA1ZiA3NCA2OSAgfGV0X2FkZGl0aW9uYWxfdGl8CjAwMDAwMDEwICA2MyA2YiA2NSA3NCA3MyAwMCA2NyA2NSAgNzQgNWYgNzQgNjkgNjMgNmIgNjUgNzQgIHxja2V0cy5nZXRfdGlja2V0fAowMDAwMDAyMCAgNzMgMDAgNzMgNjUgNzQgNWYgNzQgNjkgIDYzIDZiIDY1IDc0IDczIDAwIDUzIDc5ICB8cy5zZXRfdGlja2V0cy5TeXwKMDAwMDAwMzAgIDczIDc0IDY1IDZkIDJlIDRlIDY1IDc0ICAyZSA1MyA2ZiA2MyA6YiA2NSA3NCA3MyAgfHN0ZW0uTmV0LlNvY2tldHN8"
        return result, badBytes
    } else {
        return "Undetected", ""
    }
	*/
}