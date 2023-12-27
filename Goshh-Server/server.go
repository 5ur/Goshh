package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type pastebinMessage struct {
	CreatedAt      time.Time
	Contents       string `json:"contents"`
	Rune           string `json:"rune"`
	Lifetime       int    `json:"lifetime"`
	ViewLimit      int    `json:"view_limit"`
	RemainingViews int    `json:"remaining_views"`
}

type Message struct {
	CreatedAt time.Time
	Contents  []byte
	Rune      string `json:"rune"`
	Message   string `json:"message"`
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
	id               string
	mu               sync.Mutex
	messages         = make(map[string]Message)
	pastebinMessages = make(map[string]pastebinMessage)
	uploadedFiles    = make(map[string]UploadedFile)
	defaultValues    = map[string]interface{}{
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
		"pastebinViewLimit":        1,
		"pastebinLifetime":         1,
		"allowAllCors":             false,
		"allowCorsCustom":          false,
		"corsOrigins":              "http://localhost",
		"corsMethods":              []string{"GET", "POST"},
		"corsHeaders":              "Content-Type",
	}
)

/*
This function generates an ID string from the time-stamp at the time it's called.
Reason, being that runes and random strings added symbols, which broke url patterns.
*/
func generateRandomID() string {
	return time.Now().Format("20060102150405")
}

/*
This function is used to normalize a supplied rune, mainly so it doesn't break the returned URL.
*/
func normalizeRune(rune string) string {
	// Remove non-alphanumeric characters using regex pattern
	reg := regexp.MustCompile(`[^\w\d]+`)
	normalizedRune := reg.ReplaceAllString(rune, "")
	// Remove leading and trailing whitespace
	normalizedRune = strings.TrimSpace(normalizedRune)
	return normalizedRune
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
		log.Println("Error reading config file:", err)
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
	pastebinViewLimit := viper.GetInt("pastebinViewLimit")
	pastebinLifetime := viper.GetInt("pastebinLifetime")
	allowAllCors := viper.GetBool("allowAllCors")
	allowCorsCustom := viper.GetBool("allowCorsCustom")
	corsOrigins := viper.GetStringSlice("corsOrigins")
	corsMethods := viper.GetStringSlice("corsMethods")
	corsHeaders := viper.GetStringSlice("corsHeaders")
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
		"pastebinViewLimit":        pastebinViewLimit,
		"pastebinLifetime":         pastebinLifetime,
		"allowAllCors":             allowAllCors,
		"allowCorsCustom":          allowCorsCustom,
		"corsOrigins":              corsOrigins,
		"corsMethods":              corsMethods,
		"corsHeaders":              corsHeaders,
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
		// https://github.com/gin-gonic/gin/blob/master/logger.go#L56
		loggerConfig := gin.LoggerConfig{
			Formatter: func(params gin.LogFormatterParams) string {
				return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
					params.ClientIP,
					params.TimeStamp.Format(time.RFC1123),
					params.Method,
					params.Path,
					params.Request.Proto,
					params.StatusCode,
					params.Latency,
					params.Request.UserAgent(),
					params.ErrorMessage,
				)
			},
			Output: os.Stdout,
			// Skip logging for the root endpoint
			SkipPaths: []string{"/"},
		}
		// Set the maximum size of form values to 32MB
		router.MaxMultipartMemory = 32 << 20
		router.Use(gin.LoggerWithConfig(loggerConfig))
	default:
		// https://github.com/gin-gonic/gin/blob/master/gin.go#LL183C6-L183C9
		router = gin.New()
	}

	// Enable CORS configuration
	switch allowAllCors {
	// Enables basically everything:
	// https://github.com/gin-contrib/cors#default-allows-all-origins
	case true:
		router.Use(cors.Default())

		switch allowCorsCustom {
		case true:
			// Enable CORS with specific configuration settings:
			// https://github.com/gin-contrib/cors#canonical-example
			router.Use(cors.New(cors.Config{
				AllowOrigins: corsOrigins,
				AllowMethods: corsMethods,
				AllowHeaders: corsHeaders,
			}))
		}
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
		content, err := os.ReadFile("root.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading file: %v", err)
			return
		}
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Writer.Write(content)
	})

	// The endpoint for posting longer-form messages which is subject to the user's views on lifetime and view limit frames.
	router.POST("/pastebin", func(c *gin.Context) {
		// Bind the request JSON payload to a pastebinMessage struct
		var pastebinPost pastebinMessage
		if err := c.BindJSON(&pastebinPost); err != nil {
			// Return an error response if the request payload is invalid
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Check if the 'Contents' field is empty
		if pastebinPost.Contents == "" {
			// Return an error response if the 'Contents' field is required but empty
			c.JSON(http.StatusBadRequest, gin.H{"error": "Contents field is required"})
			return
		}

		// Generate an ID for the pastebin message based on the 'Rune' field or a random ID
		if pastebinPost.Rune == "" {
			id = generateRandomID()
		} else {
			id = normalizeRune(pastebinPost.Rune)
		}

		// Set the view limit and lifetime based on the configured default values
		viewLimit := pastebinViewLimit
		lifetime := pastebinLifetime

		// Override the view limit if the request specifies a different value
		if pastebinPost.ViewLimit != pastebinViewLimit {
			viewLimit = pastebinPost.ViewLimit
		}

		// Override the lifetime if the request specifies a different value
		if pastebinPost.Lifetime != pastebinLifetime {
			lifetime = pastebinPost.Lifetime
		}

		// Lock the shared data to safely add the pastebin message to the map
		mu.Lock()
		pastebinMessages[id] = pastebinMessage{
			CreatedAt:      time.Now(),
			Contents:       pastebinPost.Contents,
			Rune:           pastebinPost.Rune,
			Lifetime:       lifetime,
			ViewLimit:      viewLimit,
			RemainingViews: viewLimit,
		}
		// Unlock the shared data after adding the pastebin message
		mu.Unlock()

		// Generate the URL for the pastebin message
		url := fmt.Sprintf("%s://%s/pastebin/%s", getScheme(c.Request), c.Request.Host, id)
		// Return a success response with the URL and other information
		c.JSON(http.StatusOK, gin.H{
			"url":             url,
			"view_limit":      viewLimit,
			"lifetime":        lifetime,
			"remaining_views": viewLimit,
		})
	})

	router.GET("/pastebin/:id", func(c *gin.Context) {
		// Get the ID parameter from the URL
		id := c.Param("id")
		// Lock the shared data to safely access the pastebin message
		mu.Lock()
		pastebinMsg, ok := pastebinMessages[id]
		if !ok {
			// Return a not found response if the pastebin message doesn't exist
			mu.Unlock()
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Check if the pastebin message has expired or reached the view limit
		if pastebinMsg.RemainingViews <= 1 || pastebinMsg.Lifetime <= 1 {
			// Delete the pastebin message if it has expired or reached the view limit
			delete(pastebinMessages, id)
			mu.Unlock()
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Decrement the remaining views of the pastebin message
		pastebinMsg.RemainingViews--

		// Unlock the shared data after accessing the pastebin message
		mu.Unlock()

		// Update the pastebin message in the map with the decremented remaining views
		pastebinMessages[id] = pastebinMsg

		// Return the pastebin message details in the response
		c.JSON(http.StatusOK, gin.H{
			"contents":        pastebinMsg.Contents,
			"view_limit":      pastebinMsg.ViewLimit,
			"lifetime":        pastebinMsg.Lifetime,
			"remaining_views": pastebinMsg.RemainingViews,
		})
	})

	// The endpoint for posting a message
	router.POST("/message", func(c *gin.Context) {
		// Create a struct/variable 'request' with the use of the Message type (at the top) by parsing the JSON request body to the, now 'request' struct.
		var request Message
		if err := c.BindJSON(&request); err != nil {
			// Return a Expectation Failed if the JSON body couldn't be parsed
			c.AbortWithStatus(http.StatusExpectationFailed)
			return
		}
		// If the Message value/contents in the struct is empty, return a teapot
		if request.Message == "" {
			c.AbortWithStatus(http.StatusTeapot)
			return
		}
		// Generate an ID for the message, either as a random string or based on the rune sent by the client (and stored in the 'request' struct's Rune string).
		// If it's empty, just generate one with the generateRandomID function.
		if request.Rune == "" {
			id = generateRandomID()
		} else {
			// If it isn't empty then use the 'Rune' string, but run it by the regex in normalizeRune to mitigate any bad user input.
			id = normalizeRune(request.Rune)
		}
		// Lock then unlock the messages map on individual creation/posts
		mu.Lock()
		messages[id] = Message{
			CreatedAt: time.Now(),
			Contents:  []byte(request.Message),
			// Normalize the "Rune" input
			Rune: string(normalizeRune(request.Rune)),
			// Store the specified lifetime in the message struct
		}
		mu.Unlock()
		// Return the URL as the response body
		url := fmt.Sprintf("%s://%s/message/%s", getScheme(c.Request), c.Request.Host, id)
		c.String(http.StatusOK, url)
	})

	// The the endpoint for retrieving a message
	router.GET("/message/:id", func(c *gin.Context) {
		// Get the ID of the message from the URL parameter
		id := c.Param("id")
		// Place a lock on the shared messages map to prevent concurrent access
		mu.Lock()
		message, ok := messages[id]
		// If the message isn't found, release the lock and return a 404
		if !ok {
			mu.Unlock()
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		// Delete the message from memory once it has been retrieved
		delete(messages, id)
		// Release the lock on the shared messages map
		mu.Unlock()
		// Set the response content type and return the message contents as plain text
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(message.Contents))
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

	// Endpoint for file upload.
	/* Accepts a file in the multipart/form-data format with the "File" key, basically: curl -X POST -F "file=@README.pdf" http://localhost:5150/upload
	Checks if the file type is allowed by checking the extension against allowedFileTypes array of extensions.
	Saves the file to disk at the configured fileSavePath.
	Records information about the uploaded file in the uploadedFiles map.
	Returns a JSON response with a basic success/failure status the upload attempt. */
	router.POST("/upload", func(c *gin.Context) {
		// Extract the uploaded file from the request.
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the file type is allowed
		// TODO: File type check is executed after the upload is complete. -- Wastes time with large uploads.
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

		// Save the uploaded file to disk
		savePath := filepath.Join(fileSavePath, file.Filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Record information about the uploaded file
		uploadedFiles[file.Filename] = UploadedFile{
			Name:        file.Filename,
			Size:        file.Size,
			ContentType: file.Header.Get("Content-Type"),
			UploadTime:  time.Now(),
			DeleteAfter: time.Now().Add(staleFileTTL),
		}

		// Return the download endpoint for the uploaded file as a string
		url := fmt.Sprintf("%s://%s/download/%s", getScheme(c.Request), c.Request.Host, file.Filename)
		c.String(http.StatusOK, url)
	})

	// Endpoint for file download
	router.GET("/download/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join(fileSavePath, filename)
		// Check if the file exists before attempting to download it
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "File not found."})
			return
		}
		// Open the file and check for any errors
		file, err := os.Open(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		defer file.Close()
		// Get information about the file and check for any errors
		fileInfo, err := file.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		// Set headers for the file download.
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		// Copy the file to the http response writer
		if _, err := io.Copy(c.Writer, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		/* If the file was uploaded via the upload endpoint (not manually placed), increment its download count,
		then delete the file if the download count has exceeded the allowed download limit specified in the allowedFileDownloads variable */
		if uploadedFile, exists := uploadedFiles[filename]; exists {
			uploadedFile.Downloaded++
			if uploadedFile.Downloaded >= allowedFileDownloadCount {
				// Delete is called with the full path, defined in uploadedFile
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

	// The goroutine which deletes content which has reached an expiry condition:
	go func() {
		// Create a routine that triggers every second
		secondTick := time.Tick(time.Second)

		for {
			select {
			// Wait for the tick from the secondTick variable
			case <-secondTick:
				// Get the current time
				now := time.Now()
				// Acquire the lock to safely access shared data
				mu.Lock()

				// Cleanup pastebinMessages section:
				for id, pastebinMsg := range pastebinMessages {
					// Check if the pastebin message has expired or reached the view_limit
					if pastebinMsg.Lifetime <= 1 || pastebinMsg.RemainingViews <= 1 {
						// And remove the pastebin message if it has
						delete(pastebinMessages, id)
					} else {
						// Decrease the lifetime and remaining lifetime of the pastebin message with one second
						pastebinMessages[id] = pastebinMessage{
							CreatedAt:      pastebinMsg.CreatedAt,
							Contents:       pastebinMsg.Contents,
							Rune:           pastebinMsg.Rune,
							Lifetime:       pastebinMsg.Lifetime - 1,
							ViewLimit:      pastebinMsg.ViewLimit,
							RemainingViews: pastebinMsg.RemainingViews,
						}
					}
				}

				// Cleanup messages section:
				for id, message := range messages {
					// Check if the message has exceeded the cleanup interval
					if now.Sub(message.CreatedAt) > cleanupInterval {
						// Remove the message
						delete(messages, id)
					}
				}
				// Release the lock to allow other resources to access shared data
				mu.Unlock()
			}
		}
	}()

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
