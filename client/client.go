package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
)

var (
	pipedData   []byte
	requestBody string
	responseURL string
	selectRune  string
	selectFile  string
)

const (
	HttpProtocol  = "http://"
	HttpsProtocol = "https://"
)

// Read pipeline input/standard input.
func readFromStdin() ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
}

/*
If a rune is supplied, then do validation, so it doesn't just return a broken link
https://pkg.go.dev/regexp/syntax
Pretty brute force, but it's ok.
*/
func normalizeInput(input string) string {
	// Remove special characters and spaces
	pattern := regexp.MustCompile("[^a-zA-Z0-9]+")
	input = pattern.ReplaceAllString(input, "")

	// Remove unicode in cases where the terminal does some odd windows stuff
	input = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return r
	}, input)

	return input
}

func sendFile(ctx context.Context, endpoint, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func sendMessage(ctx context.Context, endpoint, requestBody string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func offlineQRgeneration() {
	// Just take the standard input and encode it into a QR code without sending it anywhere.
	pipedData, err := readFromStdin()
	qr, err := qrcode.New(string(pipedData), qrcode.Low)
	if err != nil {
		log.Fatalf("Error generating QR code: %v", err)
	}
	fmt.Println(qr.ToSmallString(false))
}

func main() {
	// Get user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting user's home directory: %v", err)
	}
	homeDir := usr.HomeDir

	//Set default values for configuration variables
	viper.SetDefault("geberateQRL", false)
	viper.SetDefault("geberateQRC", false)
	viper.SetDefault("timeoutFrame", 10)
	viper.SetDefault("messageEndpoint", "http://localhost:5150/message")
	viper.SetDefault("fileEndpoint", "http://localhost:5150/upload")

	// Read configuration from file and replace any the value of any variable found
	viper.SetConfigFile(homeDir + "/.config/Goshh/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
	}

	// Get setting from configuration
	messageEndpoint := viper.GetString("messageEndpoint")
	fileEndpoint := viper.GetString("fileEndpoint")
	generateQRL := viper.GetBool("generateQRL")
	generateQRC := viper.GetBool("generateQRC")
	timeoutFrame := viper.GetInt("timeoutFrame")

	// Check if QR code generation is enabled via flag, rather than the config file
	enableQRLFlag := flag.Bool("qrl", false, "Generate QR code for response URL.")
	enableQRCFlag := flag.Bool("qrc", false, "Generate QR code for response Contents.")
	enableQRFlag := flag.Bool("qr", false, "Generate an offline QR code for piped input.")
	flag.StringVar(&selectRune, "rune", "", "Specify a rune to use for the message ID.")
	flag.StringVar(&selectFile, "file", "", "Specify a file for upload.")
	flag.Parse()

	switch {
	case *enableQRLFlag:
		generateQRL = true
	case *enableQRCFlag:
		generateQRC = true
	case *enableQRFlag:
		offlineQRgeneration()
		return
	}

	// Call readFromStdin to store the piped data in the 'pipedData' variable.
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == 0 {
		pipedData, _ = readFromStdin()
	}

	// Create the timeoutFrame as a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutFrame)*time.Second)
	defer cancel()

	switch {
	case selectFile != "":
		responseURL, err = sendFile(ctx, fileEndpoint, selectFile)
		if err != nil {
			log.Fatalf("Error sending file request: %v", err)
		}
	case len(pipedData) == 0:
		if selectFile != "" {
			responseURL, err = sendFile(ctx, fileEndpoint, selectFile)
			if err != nil {
				log.Fatalf("Error sending file request: %v", err)
			}
		} else {
			fmt.Println("Please enter a message:")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				requestBody := map[string]string{
					"message": scanner.Text(),
					"rune":    normalizeInput(selectRune),
				}
				requestBodyBytes, _ := json.Marshal(requestBody)
				responseURL, err = sendMessage(ctx, messageEndpoint, string(requestBodyBytes))
				if err != nil {
					log.Fatalf("Error sending message request: %v", err)
				}
			} else {
				fmt.Println("Error: no input provided.")
			}
		}
	default:
		requestBody := map[string]string{
			"message": string(pipedData),
			"rune":    normalizeInput(selectRune),
		}
		requestBodyBytes, _ := json.Marshal(requestBody)
		responseURL, err = sendMessage(ctx, messageEndpoint, string(requestBodyBytes))
		if err != nil {
			log.Fatalf("Error sending message request: %v", err)
		}
	}

	// If the URL doesn't start with "http://" or "https://", prepend "http://"
	if !strings.HasPrefix(messageEndpoint, HttpProtocol) && !strings.HasPrefix(messageEndpoint, HttpsProtocol) {
		messageEndpoint = HttpProtocol + messageEndpoint
	}

	switch {
	// If QR code generation is enabled, by flag or config
	case generateQRL:
		// generateQRL: Generate QR code for response URL
		qr, err := qrcode.New(responseURL, qrcode.Low)
		if err != nil {
			log.Fatalf("Error generating QR code: %v", err)
		}
		// Print a QR code containing the message GET endpoint
		fmt.Println(qr.ToSmallString(false))

		// Print the message GET endpoint
		fmt.Print(responseURL)
	case generateQRC:
		// generateQRC: Generate QR code for response contents
		// Use context object to manage the lifecycle of the request
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutFrame)*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", responseURL, nil)
		if err != nil {
			log.Fatalf("Error creating GET request: %v", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending GET request: %v", err)
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		// Generate QR code for the response body
		qr, err := qrcode.New(string(body), qrcode.Low)
		if err != nil {
			log.Fatalf("Error generating QR code: %v", err)
		}

		fmt.Println(qr.ToSmallString(false))
	default:
		// Print HTTP URL to standard output regardless of the configuration and flags
		fmt.Printf(responseURL)
	}

}
