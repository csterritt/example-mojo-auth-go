package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"mojo-auth-test-1/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RunWebService(port int) error {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	router := gin.Default()
	config := cors.DefaultConfig()
	serverHost := os.Getenv("SERVER_HOST")
	secureCookies := true
	if strings.Index(serverHost, "localhost") != -1 {
		config.AllowOrigins = []string{"http://" + serverHost}
		secureCookies = false
	} else {
		config.AllowOrigins = []string{"https://" + serverHost}
	}
	router.Use(cors.New(config))
	cookieStore := cookie.NewStore([]byte(os.Getenv("MAIN_AUTH_SECRET")), []byte(os.Getenv("MAIN_ENC_SECRET")))
	cookieStore.Options(sessions.Options{
		Path:     "",
		Domain:   serverHost,
		MaxAge:   6 * 30 * 24 * 3600,
		Secure:   secureCookies,
		HttpOnly: true,
		SameSite: 1,
	})
	router.Use(sessions.Sessions("mysession", cookieStore))
	router.Use(limits.RequestSizeLimiter(16384))

	// Define a limit rate to about 10 requests per second (asking for 20, which on the mac limits to about 10).
	rate, err := limiter.NewRateFromFormatted("20-S")
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Create a new middleware with the limiter instance.
	store := memory.NewStore()
	middleware := mgin.NewMiddleware(limiter.New(store, rate))
	router.Use(middleware)

	// Initialize the routes
	routes.InitializeMainRoutes(router)
	//routes.InitializeAuthRoutes(router)
	//router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	log.Println("Starting on port", port)

	// Start serving the application
	err = router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Failed to router.Run with error %v\n", err)
	}

	return err
}

func verifyEnvironment() {
	problemFound := false
	for _, str := range []string{"SERVER_HOST", "MAIN_AUTH_SECRET", "MAIN_ENC_SECRET",
		"SMTP_USERNAME", "SMTP_PASS", "SMTP_HOST"} {
		val := os.Getenv(str)
		if len(strings.Trim(val, " \t\n\r")) == 0 {
			problemFound = true
			_, _ = fmt.Fprintf(os.Stderr, "Error: Unable to find environmental variable %s\n", str)
		}
	}

	if problemFound {
		panic("Set variables above and restart.")
	}
}

func main() {
	verifyEnvironment()
	port := os.Getenv("PORT")
	portId := 4000
	if port != "" {
		var err error
		portId, err = strconv.Atoi(port)
		if err != nil || portId < 1 {
			panic("Cannot convert " + port + " to an integer (or it's a bad number).")
		}
	}
	err := RunWebService(portId)
	if err != nil {
		panic(err)
	}
}
