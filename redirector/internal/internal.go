package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"gr-blockchain-side/api"
	"gr-blockchain-side/internal/blockchain/signer"
	"gr-blockchain-side/internal/config"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func LoadEnv() (err error) {
	if err = godotenv.Load(os.Args[1]); err != nil {
		return
	}
	config.DefaultConfig.DefaultPort = os.Getenv("DEFAULT_PORT")
	config.DefaultConfig.DeployMode = os.Getenv("MODE")
	config.DefaultConfig.Node = os.Getenv("BLOCKCHAIN_NODE")
	config.DefaultConfig.APIKey = os.Getenv("API_KEY")
	fmt.Println("Default port:", config.DefaultConfig.DefaultPort)
	return nil
}

func InitEngine(mode string) (engine *gin.Engine) {
	gin.SetMode(mode)
	engine = gin.New()

	engine.Use(cors.Default())
	engine.Use(gin.Logger())
	return
}

func SetUpAPIs(router *gin.Engine) {
	api.SetUpSavingsAPI(router)
	api.SetUpPortalAPI(router)
}

func LoadConfiguration() (err error) {
	// create signer
	raw, err := ioutil.ReadFile(os.Getenv("BLOCKCHAIN_CONFIG"))
	if err != nil {
		return
	}
	fmt.Println(string(raw))
	var configuration map[string]string = make(map[string]string)
	err = json.Unmarshal(raw, &configuration)
	if err != nil {
		return
	}
	signer, err := signer.NewSigner(configuration["private_key_path"])
	if err != nil {
		return
	}
	config.DefaultConfig.Signer = signer
	config.DefaultConfig.Contract = configuration["contract"]
	config.DefaultConfig.ABIPath = configuration["abi_path"]
	return
}

func SetupGracefulShutdown(ctx context.Context, port string, engine *gin.Engine) {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}

	defer func() {
		if err := server.Shutdown(ctx); err != nil {
			log.Info("Server Shutdown: ", err)
		}
	}()

	signalForExit := make(chan os.Signal, 1)
	signal.Notify(signalForExit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Application failed", err)
		}
	}()
	log.WithFields(log.Fields{"bind": port}).Info("Running application")

	stop := <-signalForExit
	log.Info("Stop signal Received", stop)
	log.Info("Waiting for all jobs to stop")
}
