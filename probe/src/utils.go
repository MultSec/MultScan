package main

import (
	"os"
	"fmt"
	"encoding/json"
    "io"
	"net/http"

	"github.com/mgutz/ansi"
)

type Log int64

const (
    logError Log = iota
    logInfo
    logStatus
    logInput
	logSuccess
	logSection
	logSubSection
)

// Function to print logs
func printLog(log Log, text string) {
	switch log {
	case logError:
		fmt.Printf("[%s] %s %s\n", ansi.ColorFunc("red")("!"), ansi.ColorFunc("red")("ERROR:"), ansi.ColorFunc("cyan")(text))
	case logInfo:
		fmt.Printf("[%s] %s\n", ansi.ColorFunc("blue")("i"), text)
	case logStatus:
		fmt.Printf("[*] %s\n", text)
	case logInput:
		fmt.Printf("[%s] %s", ansi.ColorFunc("yellow")("?"), text)
	case logSuccess:
		fmt.Printf("[%s] %s\n", ansi.ColorFunc("green")("+"), text)
	case logSection:
		fmt.Printf("\t[%s] %s\n", ansi.ColorFunc("yellow")("-"), text)
	case logSubSection:
		fmt.Printf("\t\t[%s] %s\n", ansi.ColorFunc("magenta")(">"), text)
	}
}

func GetConf(configPath string) (map[string]string, error) {
	var conf map[string]string

	data, err := os.ReadFile(configPath)
	if err != nil {
		return conf, fmt.Errorf("failed to read file: %v", err)
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return conf, fmt.Errorf("failed to parse file: %v", err)
	}

	return conf, nil
}

// Function to delete the sample file
func DelSample(id string) error {
    filename := fmt.Sprintf("sample_%s.bin", id)
    err := os.Remove(filename)
    if err != nil {
        if os.IsNotExist(err) {
            printLog(logError, fmt.Sprintf("[%s] Sample file not found: %s", ansi.ColorFunc("default+hb")(id), filename))
            return nil
        }
        return fmt.Errorf("failed to delete sample file: %v", err)
    }
    printLog(logSuccess, fmt.Sprintf("[%s] Sample file deleted successfully", ansi.ColorFunc("default+hb")(id)))
    return nil
}

// Function to get sample content from the server
func GetSample(ip string, id string) error {
    // Create a new HTTP request
    req, err := http.NewRequest("GET", "http://"+ip+":8000/api/v1/sample/download/"+id, nil)
    if err != nil {
        printLog(logError, fmt.Sprintf("[%s] Error creating request: %v", ansi.ColorFunc("default+hb")(id), err))
        return err
    }

    // Set the headers for the request
    req.Header.Set("Accept", "application/octet-stream")
    req.Header.Set("Content-Type", "application/octet-stream")

    // Create a new HTTP client
    client := &http.Client{}

    // Send the request
    resp, err := client.Do(req)
    if err != nil {
        printLog(logError, fmt.Sprintf("[%s] Error downloading sample from server: %v", ansi.ColorFunc("default+hb")(id), err))
        return err
    }
    defer resp.Body.Close()

    // Check if the response status is OK
    if resp.StatusCode != http.StatusOK {
        printLog(logError, fmt.Sprintf("[%s] Unexpected status code: %d", ansi.ColorFunc("default+hb")(id), resp.StatusCode))
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    // Save Sample
    err = saveSample(id, resp.Body)
    if err != nil {
        printLog(logError, fmt.Sprintf("[%s] Error saving sample: %v", ansi.ColorFunc("default+hb")(id), err))
        return err
    }

    printLog(logSuccess, fmt.Sprintf("[%s] Sample downloaded and saved successfully", ansi.ColorFunc("default+hb")(id)))
    return nil
}

// Function to save the sample to a file
func saveSample(id string, data io.Reader) error {
    filename := fmt.Sprintf("sample_%s.bin", id)
    out, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file: %v", err)
    }
    defer out.Close()

    _, err = io.Copy(out, data)
    if err != nil {
        return fmt.Errorf("failed to save sample: %v", err)
    }
    return nil
}

func copyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    destFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destFile.Close()

    _, err = io.Copy(destFile, sourceFile)
    if err != nil {
        return err
    }

    return nil
}