package notificationaway

import (
	"github.com/amrHassanAbdallah/notificationaway/api"
	"github.com/amrHassanAbdallah/notificationaway/consumer"
	"github.com/amrHassanAbdallah/notificationaway/persistence"
	"github.com/amrHassanAbdallah/notificationaway/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/segmentio/kafka-go"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Logger = *zap.SugaredLogger

var logger Logger

func init() {
	core := zap.NewProductionConfig()
	core.EncoderConfig.TimeKey = "timestamp"
	core.EncoderConfig.MessageKey = "message"
	core.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	customLog, _ := core.Build()
	logger = customLog.WithOptions(zap.AddCallerSkip(1)).Sugar()
}

func main() {

	var (
		DBHost              string
		DBUser              string
		DBPassword          string
		DBName              string
		DBConnectionTimeOut int
		DBReplicasetName    *string

		RequestTimeOut int
		Listen         string

		KafkaBroker        string
		KafkaConsumerGroup string
		KafkaTopic         string
		RunConsumer        bool
	)
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "notificationaway-service"
	app.Usage = "manage messages and their notifications"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "mongo-dsn",
			Usage:       "DSN for MongoDB",
			Value:       "mongodb://localhost:27017",
			Destination: &DBHost,
		},
		&cli.StringFlag{
			Name:        "mongo-user",
			Usage:       "mongo user name used during the auth",
			Value:       "",
			Destination: &DBUser,
		},
		&cli.StringFlag{
			Name:        "mongo-password",
			Usage:       "mongo user password used during the auth",
			Value:       "",
			Destination: &DBPassword,
		},
		&cli.StringFlag{
			Name:        "db-name",
			Usage:       "db name",
			Value:       "",
			Destination: &DBName,
		},
		&cli.IntFlag{
			Name:        "db-connection-timeout",
			Usage:       "timeout the db connection after how many seconds",
			Value:       5,
			Destination: &DBConnectionTimeOut,
		},
		&cli.StringFlag{
			Name:        "db-replicaset",
			Usage:       "used to connect to db cluster",
			Destination: DBReplicasetName,
		},
		&cli.IntFlag{
			Name:        "request-timeout",
			Usage:       "http requests timeout",
			Value:       5,
			Destination: &RequestTimeOut,
		},
		&cli.StringFlag{
			Name:        "listen",
			Usage:       "port for the server to recieve requests on",
			Value:       ":7981",
			Destination: &Listen,
		},

		&cli.StringFlag{
			Name:        "kafka-brokers",
			Usage:       "kafka brokers DSN comma separated",
			Value:       "localhost:9092",
			Destination: &KafkaBroker,
		},
		&cli.StringFlag{
			Name:        "kafka-consumer-group",
			Usage:       "kafka consumer group",
			Value:       "notificationaway-consumer",
			Destination: &KafkaConsumerGroup,
		},
		&cli.StringFlag{
			Name:        "kafka-topic",
			Usage:       "kafka topic",
			Value:       "notifications",
			Destination: &KafkaTopic,
		},
		&cli.BoolFlag{
			Name:        "run-consumer",
			Usage:       "to run the kafka consumer",
			Value:       true,
			Destination: &RunConsumer,
		},
	}

	app.Action = func(context *cli.Context) error {
		mctx := context.Context
		mongoHosts := make([]string, 0)
		if strings.Contains(DBHost, ",") {
			mongoHosts = strings.Split(DBHost, ",")
		} else {
			mongoHosts = append(mongoHosts, DBHost)
		}
		extraArgs := map[string]interface{}{
			"timeout":    time.Duration(DBConnectionTimeOut) * time.Second,
			"replicaset": DBReplicasetName,
		}

		persistenceLayer, err := persistence.NewPersistenceLayer(mongoHosts, DBUser, DBPassword, DBName, extraArgs)
		if err != nil {
			logger.Fatalw("failed to initiate persistence layer", "err", err)
		}
		err = persistenceLayer.Connect(mctx)
		if err != nil {
			logger.Fatalw("failed to connect to persistence layer", "err", err)
		}
		logger.Info("connection with mongodb has been established")

		notificationawayService := service.NewService(persistenceLayer)

		server := api.NewServer(notificationawayService)

		r := chi.NewRouter()
		r.Route("/api/v1", func(r chi.Router) {
			r.Use(middleware.Timeout(time.Duration(RequestTimeOut) * time.Second))
			api.HandlerFromMux(server, r)
		})
		srv := &http.Server{
			Handler: r,
			Addr:    Listen,
		}
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// ES is a http.Handler, so you can pass it directly to your mux

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatalw("server failed to start", "error", err)
			}
		}()
		logger.Infof("server started on port %v", Listen)

		if RunConsumer {
			var wg sync.WaitGroup
			kbrokers := make([]string, 0)
			if strings.Contains(KafkaBroker, ",") {
				kbrokers = strings.Split(KafkaBroker, ",")
			} else {
				kbrokers = append(kbrokers, DBHost)
			}
			kafkaReader := &consumer.KafkaHandler{Reader: kafka.NewReader(kafka.ReaderConfig{
				Brokers:     kbrokers,
				GroupID:     KafkaConsumerGroup,
				Topic:       KafkaTopic,
				MinBytes:    1,
				MaxBytes:    100e6, // 100MB
				ErrorLogger: &consumer.ErrLoggerWrapper{Logger: logger},
				Logger:      &consumer.DebLoggerWrapper{Logger: logger}},
			)}
			wconfig := consumer.WorkerConfig{
				MessageBrokerReader: kafkaReader,
				EventMessageHandler: consumer.NewNotificationEventHandler(persistence.ReadLayerInterface(persistenceLayer)),
				RequestTimeOut:      30 * time.Second,
				MaxRetriesOnFailure: 3,
			}
			wg.Add(1)
			go consumer.HandleKafkaMessages(mctx, &wg, wconfig)
		}

		<-done
		logger.Info("server terminating...")
		if err := srv.Shutdown(mctx); err != nil {
			logger.Errorw("server terminating failed", "err", err)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
