package cmd

import (
	"Mail-Achive/pkg/api"
	"Mail-Achive/pkg/config"
	"Mail-Achive/pkg/log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the http server",
	Run: func(cmd *cobra.Command, args []string) {
		// start the http server
		serve()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

// updateOptions updates the log options
func updateOptions(scope string, options *log.Options, settings *config.Config) error {
	if settings.Log.OutputPath != "" {
		options.OutputPaths = []string{settings.Log.OutputPath}
	}
	if settings.Log.RotationPath != "" {
		options.RotateOutputPath = settings.Log.RotationPath
	}

	options.RotationMaxBackups = settings.Log.RotationMaxBackups
	options.RotationMaxSize = settings.Log.RotationMaxSize
	options.RotationMaxAge = settings.Log.RotationMaxAge
	options.JSONEncoding = settings.Log.JSONEncoding
	level, err := options.ConvertLevel(settings.Log.OutputLevel)
	if err != nil {
		return err
	}
	options.SetOutputLevel(scope, level)
	options.SetLogCallers(scope, true)

	return nil
}

func serve() {
	var settings config.Config
	// parse the config file
	if err := config.ParseYamlFile(configFile, &settings); err != nil {
		panic(err)
	}

	// init and update the log options
	logOptions := log.DefaultOptions()
	if err := updateOptions("default", logOptions, &settings); err != nil {
		panic(err)
	}
	// configure the log options
	if err := log.Configure(logOptions); err != nil {
		panic(err)
	}

	// init the http server
	server := api.NewServer(&settings)

	// start the http server
	server.Start()

	log.Info("server is started")

	sig := make(chan os.Signal, 1024)
	// subscribe signals: SIGINT & SINGTERM
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case s := <-sig:
			log.Infof("receive signal: %v", s)

			// flush the log
			log.Sync()

			start := time.Now()

			// stop the grpc server gracefully
			server.Close()

			log.Info("server is stopped")

			log.Infof("shut down takes time: %v", time.Now().Sub(start))
			return
		}
	}
}
