package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Message struct {
	CreatedAt time.Time
	Encrypted []byte
}

type UploadedFile struct {
	Name        string
	Size        int64
	ContentType string
	UploadTime  time.Time
	Downloaded  int
	DeleteAfter time.Time
}

// Maps for messages and file/s as well as the default variables (if non are supplied)
var (
	encryptionKey []byte
	mu            sync.Mutex
	messages      = make(map[string]Message)
	uploadedFiles = make(map[string]UploadedFile)
	defaultValues = map[string]interface{}{
		"debugMode":                false,
		"serverPort":               5150,
		"useDefault":               false,
		"trustedProxies":           []string{"127.0.0.1"},
		"cleanupInterval":          30 * time.Second,
		"allowLocalNetworkAccess":  false,
		"allowedFileTypes":         []string{"txt", "md", "jpg"},
		"fileSavePath":             "/path/to/save/files",
		"staleFileTTL":             30 * time.Second,
		"allowedFileDownloadCount": 1,
	}
)

func init() {
	// Generate a 32 byte slice for the messages' encryptionKey. Pointless, but it was interesting to do
	encryptionKey = make([]byte, 32)
	if _, err := rand.Read(encryptionKey); err != nil {
		panic(err)
	}
}

// More pointless, yet interesting encryption functions
func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Generate random initialization vector
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	// Encrypt the data using counter mode
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(data, data)

	// Concat the init vector and encrypted data, then return as a b64 encoded result
	encrypted := append(iv, data...)
	return []byte(base64.StdEncoding.EncodeToString(encrypted)), nil
}

func decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// Decode b64
	encrypted, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	// Extract the init vector and encrypted data
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]
	// Decrypt the data using CTR mode
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted, nil
}

/*
This function generates an ID string from the time-stamp at the time it's called.
Reason, being that runes and random strings added symbols, which broke url patterns.
*/
func generateRandomID() string {
	return time.Now().Format("20060102150405")
}

/*
Used to get the porper scheme (http or https) if the "X-Forwarded-Proto" header and the value of the TLS supplied in the request.
If the XFP header is present, then its value is returned as the scheme
If XFP header is not present, thne the function checks if the request was made over https by checking if the TLS field is not nil;
If the "TLS" field is not nil, https is returned as the scheme. Otherwise, the function elses out to http.
*/
func getScheme(req *http.Request) string {
	if proto := req.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}
	if req.TLS != nil {
		return "https"
	}
	return "http"
}

func main() {
	// Load configuration from file, located in the same directory.
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error reading config file: %s", err)
	}

	// Set default values for configuration variables (not sure if there is a better way to not waste time setting defaults)
	for key, value := range defaultValues {
		viper.SetDefault(key, value)
	}

	// Get configuration values and log any missing variables
	missingVars := []string{}
	debug := viper.GetBool("debugMode")
	port := viper.GetInt("serverPort")
	useDefault := viper.GetBool("useDefault")
	trustedProxies := viper.GetStringSlice("trustedProxies")
	cleanupInterval := viper.GetDuration("cleanupInterval")
	allowLocalNetworkAccess := viper.GetBool("allowLocalNetworkAccess")
	allowedFileTypes := viper.GetStringSlice("allowedFileTypes")
	fileSavePath := viper.GetString("fileSavePath")
	staleFileTTL := viper.GetDuration("staleFileTTL")
	allowedFileDownloadCount := viper.GetInt("allowedFileDownloadCount")
	for key := range defaultValues {
		if !viper.IsSet(key) {
			missingVars = append(missingVars, key)
		}
	}

	// Set default values for any variables which turned out to be missing
	if len(missingVars) > 0 {
		for _, key := range missingVars {
			viper.Set(key, defaultValues[key])
			log.Printf("Using default value for missing variable %s", key)
		}
	}

	// Create a message which will be printed out during startup (what ended up being used as configuration variables)
	// https://pkg.go.dev/fmt#hdr-Printing
	logMsg := "Loading configuration values:\n"
	vars := map[string]interface{}{
		"debugMode":                debug,
		"serverPort":               port,
		"useDefault":               useDefault,
		"trustedProxies":           trustedProxies,
		"cleanupInterval":          cleanupInterval,
		"allowLocalNetworkAccess":  allowLocalNetworkAccess,
		"allowedFileTypes":         allowedFileTypes,
		"fileSavePath":             fileSavePath,
		"staleFileTTL":             staleFileTTL,
		"allowedFileDownloadCount": allowedFileDownloadCount,
	}

	for v, val := range vars {
		switch val := val.(type) {
		case string:
			logMsg += fmt.Sprintf(" %v=%q\n", v, val)
		case []string:
			logMsg += fmt.Sprintf(" %v=%q\n", v, val)
		case bool:
			logMsg += fmt.Sprintf(" %v=%v\n", v, val)
		case int:
			logMsg += fmt.Sprintf(" %v=%d\n", v, val)
		case time.Duration:
			logMsg += fmt.Sprintf(" %v=%v\n", v, val)
		}
	}
	log.Print(logMsg)

	// Set Gin operational mode
	// https://github.com/gin-gonic/gin/blob/master/mode.go#L15
	switch debug {
	case true:
		gin.SetMode(gin.DebugMode)
	case false:
		gin.SetMode(gin.ReleaseMode)
	}

	/* Create a router engine
	https://github.com/gin-gonic/gin/blob/master/mode.go#L16 */
	var router *gin.Engine
	switch {
	case useDefault:
		// https://github.com/gin-gonic/gin/blob/master/gin.go#L215
		router = gin.Default()
	default:
		// https://github.com/gin-gonic/gin/blob/master/gin.go#LL183C6-L183C9
		router = gin.New()
	}

	/* Define the list of trusted proxies
	https://github.com/gin-gonic/gin/blob/master/gin.go#L205 */
	router.SetTrustedProxies(trustedProxies)

	/* Allowing all private IP ranges to access the server
	Redundant, but it saves time. */
	if allowLocalNetworkAccess {
		trustedProxies = append(trustedProxies, "192.168.0.0/16", "10.0.0.0/8", "172.16.0.0/12")
	}

	/* Root of the http engine
	Currently in use as a documentation serving page/endoiunt */
	router.GET("/", func(c *gin.Context) {
		// Some html here.
		content, err := ioutil.ReadFile("root.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading file: %v", err)
			return
		}
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Writer.Write(content)
	})

	// The post endpoint for sending the server a "message"
	router.POST("/message", func(c *gin.Context) {
		message, err := c.GetRawData()
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		// Encrypt the message using the encryption key
		encrypted, err := encrypt(message, encryptionKey)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// Generate an ID for the message (the timestamp format)
		id := generateRandomID()
		// Store the encrypted message and it's creation time in memory (time stamp is so we can clear out stale messages with the cleanupInterval)
		mu.Lock()
		messages[id] = Message{
			CreatedAt: time.Now(),
			Encrypted: encrypted,
		}
		mu.Unlock()
		// Build the URL for the message
		url := fmt.Sprintf("%s://%s/message/%s", getScheme(c.Request), c.Request.Host, id)
		// Return the value to the client
		c.JSON(http.StatusOK, gin.H{
			"url": url,
		})
	})

	// The the endpoint for retrieving a message
	router.GET("/message/:id", func(c *gin.Context) {
		// Get the ID of the message from the URL
		id := c.Param("id")
		// Find the encrypted message in memory
		mu.Lock()
		message, ok := messages[id]
		if !ok {
			mu.Unlock()
			// If the message isn't found, return a 404 error
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		// Decrypt the message using the encryption key
		decrypted, err := decrypt(message.Encrypted, encryptionKey)
		if err != nil {
			mu.Unlock()
			// If decryption fails, return a 418
			c.AbortWithStatus(http.StatusTeapot)
			return
		}
		// Delete the message from memory after retrieval
		delete(messages, id)
		mu.Unlock()
		// Generate the response with the message as a raw string
		c.Data(http.StatusOK, "text/plain; charset=utf-8", decrypted)
	})

	// Infinite loop go routine to reset the messages map based on the cleanupInterval variable
	go func() {
		for {
			// Wait for the cleanup interval to elapse
			<-time.After(time.Duration(cleanupInterval))
			// Lock the shared memory to prevent concurrent access
			mu.Lock()
			// Iterate over the messages in the map
			now := time.Now()
			for id, message := range messages {
				// If a message is older than the cleanup interval, delete it from the map
				if now.Sub(message.CreatedAt) > time.Duration(cleanupInterval) {
					delete(messages, id)
				}
			}
			// Remove the lock on shared memory
			mu.Unlock()
		}
		// Repeat
	}()

	// Endpoint for uploading files
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		extension := filepath.Ext(file.Filename)[1:]
		allowed := false
		for _, fileType := range allowedFileTypes {
			if fileType == extension {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed."})
			return
		}

		savePath := filepath.Join(fileSavePath, file.Filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		uploadedFiles[file.Filename] = UploadedFile{
			Name:        file.Filename,
			Size:        file.Size,
			ContentType: file.Header.Get("Content-Type"),
			UploadTime:  time.Now(),
			DeleteAfter: time.Now().Add(staleFileTTL),
		}

		c.JSON(http.StatusOK, gin.H{"status": "File uploaded successfully!"})
	})

	// Endpoint for downloading files
	router.GET("/download/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join(fileSavePath, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File not found."})
			return
		}
		file, err := os.Open(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()
		fileInfo, err := file.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		if _, err := io.Copy(c.Writer, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if uploadedFile, exists := uploadedFiles[filename]; exists {
			uploadedFile.Downloaded++
			if uploadedFile.Downloaded >= viper.GetInt("allowedFileDownloads") {
				delete(uploadedFiles, filename)
				file.Close()
				if err := os.Remove(filePath); err != nil {
					fmt.Println("Failed to remove file:", err)
				}
			} else {
				uploadedFiles[filename] = uploadedFile
			}
		}
	})
	// Perpetual go-routine to delete stale/forgotten files (same as for messages/id)
	go func() {
		for {
			// Loop through uploaded files
			for filename, file := range uploadedFiles {
				// Check if current time stamp has surpassed the staleFileTTL
				if time.Now().After(file.DeleteAfter) {
					// Delete all the files from the uploadedFiles map
					delete(uploadedFiles, filename)
					// Create a full path for the individual file/s
					filePath := filepath.Join(fileSavePath, filename)
					// Mace a remove call to the path made above
					if err := os.Remove(filePath); err != nil {
						// Log a specific error for this loop
						fmt.Println("[Looping routine] Failed to remove file:", err)
					}
				}
			}
			// Sleep before restarting the loop all over again. Maybe not the best idea, but it's ok
			time.Sleep(staleFileTTL)
		}
	}()

	// Start the Gin router
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Couldn't start Gin router: %s", err)
	}
}
