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
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
)

const (
	HttpProtocol  = "http://"
	HttpsProtocol = "https://"
)

// Read pipeline input/standard input.
func readFromStdin() ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
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

	// Read configuration from file and replace any the value of any variable found
	viper.SetConfigFile(homeDir + "/.config/Goshh/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
	}

	// Get setting from configuration
	messageEndpoint := viper.GetString("messageEndpoint")
	generateQRL := viper.GetBool("generateQRL")
	generateQRC := viper.GetBool("generateQRC")
	timeoutFrame := viper.GetInt("timeoutFrame")

	// Check if QR code generation is enabled via flag, rather than the config file
	enableQRLFlag := flag.Bool("qrl", false, "Generate QR code for response URL.")
	enableQRCFlag := flag.Bool("qrc", false, "Generate QR code for response Contents.")
	enableQRFlag := flag.Bool("qr", false, "Generate an offline QR code for piped input.")
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
	pipedData, err := readFromStdin()
	if err != nil {
		log.Fatalf("Error reading from standard input: %v", err)
	}

	// Create the timeoutFrame as a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutFrame)*time.Second)
	defer cancel()

	// Build out the request as: Context, Type, Server's message endpoint, and the STDIN as a/the client message.
	clientMessage := strings.NewReader(fmt.Sprintf("message=%s", pipedData))
	preRequest, err := http.NewRequestWithContext(ctx, "POST", messageEndpoint, clientMessage)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// Set request headers
	preRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the POST request to messageEndpoint
	client := &http.Client{}
	postResponse, err := client.Do(preRequest)
	if err != nil {
		log.Fatalf("Error sending POST request: %v", err)
		return
	}
	defer postResponse.Body.Close()

	// Read the POST response body
	responseBody, err := ioutil.ReadAll(postResponse.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return
	}

	/* In case it's wrapped in quptes
	messageURL := string(bytes.Trim(responseBody, "\"")) */
	messageURL := string(responseBody)

	// If the URL doesn't start with "http://" or "https://", prepend "http://"
	if !strings.HasPrefix(messageEndpoint, HttpProtocol) && !strings.HasPrefix(messageEndpoint, HttpsProtocol) {
		messageEndpoint = HttpProtocol + messageEndpoint
	}

	switch {
	// If QR code generation is enabled, by flag or config
	case generateQRL:
		// generateQRL: Generate QR code for response URL
		qr, err := qrcode.New(messageURL, qrcode.Low)
		if err != nil {
			log.Fatalf("Error generating QR code: %v", err)
		}
		// Print a QR code containing the message GET endpoint
		fmt.Println(qr.ToSmallString(false))

		// Print the message GET endpoint
		fmt.Print(messageURL)
	case generateQRC:
		// generateQRC: Generate QR code for response contents
		// Use context object to manage the lifecycle of the request
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutFrame)*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", messageURL, nil)
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
		fmt.Printf(messageURL)
	}

}
