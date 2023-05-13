package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
)

var (
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

func sendRequest(ctx context.Context, endpoint string, requestBody string) (string, error) {
	clientRequest := strings.NewReader(requestBody)
	preRequest, err := http.NewRequestWithContext(ctx, "POST", endpoint, clientRequest)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	preRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	postResponse, err := client.Do(preRequest)
	if err != nil {
		return "", fmt.Errorf("error sending POST request: %v", err)
	}
	defer postResponse.Body.Close()

	responseBody, err := ioutil.ReadAll(postResponse.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(responseBody), nil
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

	//Tuka
	// if selectRune != "" {
	// 	match, _ := regexp.MatchString("^[[:alnum:]]+$", strings.TrimSpace(strings.ReplaceAll(selectRune, " ", "")))
	// 	if !match {
	// 		log.Fatalf("Invalid rune: '%s'", selectRune)
	// 		return
	// 	}
	// }

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
	pipedData, err := readFromStdin()
	if err != nil {
		log.Fatalf("Error reading from standard input: %v", err)
	}

	// Create the timeoutFrame as a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutFrame)*time.Second)
	defer cancel()

	// // Build out the request as: Context, Type, Server's message endpoint, and the STDIN as a/the client message.
	// clientRequest := strings.NewReader(fmt.Sprintf("message=%s&rune=%s", pipedData, normalizeInput(selectRune)))
	// preRequest, err := http.NewRequestWithContext(ctx, "POST", messageEndpoint, clientRequest)
	// if err != nil {
	// 	fmt.Println("Error creating request:", err)
	// 	return
	// }
	// // Set request headers
	// preRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//
	// // Send the POST request to messageEndpoint
	// client := &http.Client{}
	// postResponse, err := client.Do(preRequest)
	// if err != nil {
	// 	log.Fatalf("Error sending POST request: %v", err)
	// 	return
	// }
	// defer postResponse.Body.Close()
	//
	// // Read the POST response body
	// responseBody, err := ioutil.ReadAll(postResponse.Body)
	// if err != nil {
	// 	log.Fatalf("Error reading response body: %v", err)
	// 	return
	// }
	//
	// /* In case it's wrapped in quptes
	// responseURL := string(bytes.Trim(responseBody, "\"")) */
	// responseURL := string(responseBody)

	switch {
	case selectRune != "":
		requestBody := fmt.Sprintf("message=%s&rune=%s", pipedData, normalizeInput(selectRune))
		responseURL, err = sendRequest(ctx, messageEndpoint, requestBody)
		if err != nil {
			log.Fatalf("Error sending message request: %v", err)
			return
		}
	case selectFile != "":
		requestBody := fmt.Sprintf("file=@%s", selectFile)
		responseURL, err = sendRequest(ctx, fileEndpoint, requestBody)
		if err != nil {
			log.Fatalf("Error sending file request: %v", err)
			return
		}
	default:
		requestBody := fmt.Sprintf("message=%s&rune=%s", pipedData, normalizeInput(selectRune))
		responseURL, err = sendRequest(ctx, messageEndpoint, requestBody)
		if err != nil {
			log.Fatalf("Error sending message request: %v", err)
			return
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
