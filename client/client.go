package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	"syscall"
	"time"
	"unicode"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// Variables
var (
	pipedData   []byte
	requestBody string
	responseURL string
	selectRune  string
	selectFile  string
)

// Seems better to store them as a constant
const (
	HttpProtocol  = "http://"
	HttpsProtocol = "https://"
)

/*
Read pipeline input/standard input and store it a s a byte slice
If there is an error with accessing the standard input, it returns an empty byte slice and the error
*/
func readFromStdin() ([]byte, error) {
	var pipedData []byte
	// Get information about the data stream, which is an ugly workaround of the script checking if there even was something piped
	stat, err := os.Stdin.Stat()
	if err != nil {
		// Return empty variable and error if there was a problem getting the info ie: we piped nothing
		return pipedData, err
	}
	// Check if the stdin is a named pipe; FIFO stream
	if stat.Mode()&os.ModeNamedPipe == 0 {
		// Return empty pipedData and a foo as an error if the stdio is not a named pipe ie; there was no |
		return pipedData, errors.New("foo")
	}
	// Read the piped stuff
	pipedData, err = ioutil.ReadAll(os.Stdin)
	if err != nil {
		return pipedData, err
	}
	// Return the data for further processing as 'pipedData'
	return pipedData, nil
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

// Function to post to the fileEndpoint using multipart/form-data encoding
func sendFile(ctx context.Context, endpoint, filePath string) (string, error) {
	// Open the file to be sent
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// Create a new multipart writer to encode the file data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// Create a new form file part for the file being sent, then copy the contents of the file into said form
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	// Close the multipart writer to finalize the encoding
	err = writer.Close()
	if err != nil {
		return "", err
	}

	// Create a new HTTP request with the encoded body and the multipart content type header
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// Send the request and read out returned response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(responseBody), nil
}

// Function to post to the messageEndpoint with a json request body
func sendMessage(ctx context.Context, endpoint, requestBody string) (string, error) {
	// Create a post template with a context timeout (from the timeoutFrame variable in the config or default), as a POST, to the messageEndpoint, and the container for requestBody
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(requestBody))
	if err != nil {
		return "", err
	}
	// Set the content type
	req.Header.Set("Content-Type", "application/json")
	// Send the request and read out the returned response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Read out
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(responseBody), nil
}

// Specific function to generate a qr code out of standard input, be it piped or streamed
func offlineQRgeneration() {
	pipedData, err := readFromStdin()
	if err != nil {
		fmt.Println("No piped data found. Please enter a message:")
		obfuscatedMessage, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Error reading obfuscatedMessage: %v", err)
		}
		pipedData = []byte(obfuscatedMessage)
	}
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
		log.Fatalf("Error getting home directory: %v", err)
	}
	homeDir := usr.HomeDir

	//Set default values for configuration variables
	viper.SetDefault("geberateQRL", false)
	viper.SetDefault("geberateQRC", false)
	viper.SetDefault("timeoutFrame", 30)
	viper.SetDefault("messageEndpoint", "http://localhost:5150/message")
	viper.SetDefault("fileEndpoint", "http://localhost:5150/upload")

	// Read configuration from file and replace any the value of any variable found
	viper.SetConfigFile(homeDir + "/.config/Goshh/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
	}

	// Get settings from the configuration file with viper and set them as the variables' values
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

	// Enable variable/s by invoking with flags
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

	// Declare the timeoutFrame as a context for later use as a timeout
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
			obfuscatedMessage, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatalf("Error reading obfuscatedMessage: %v", err)
			}
			requestBody := map[string]string{
				"message": string(obfuscatedMessage),
				"rune":    normalizeInput(selectRune),
			}
			requestBodyBytes, _ := json.Marshal(requestBody)
			responseURL, err = sendMessage(ctx, messageEndpoint, string(requestBodyBytes))
			if err != nil {
				log.Fatalf("Error sending message request: %v", err)
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
