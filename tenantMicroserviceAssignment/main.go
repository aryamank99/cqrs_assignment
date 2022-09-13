package main

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	watermillMiddleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/alexflint/go-arg"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"tenantMicroserviceAssignment/api"
	"tenantMicroserviceAssignment/cmd"
	"tenantMicroserviceAssignment/cqrshandler"
	"tenantMicroserviceAssignment/database"
	_ "tenantMicroserviceAssignment/docs"
)

var appName = "tenant-microservice"

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

// @title Tenant Microservice APIs
// @version 1.0
// @description The Tenant Microservice APIs
// @BasePath /
func main() {
	logrus.Infof("Starting %v\n", appName)

	// Initialize config struct and populate it from env vars and flags.
	cfg := cmd.DefaultConfiguration()
	arg.MustParse(cfg)

	port := ":" + cfg.ServerConfig.Port
	router := chi.NewRouter()

	//router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	// initialize database
	database.Connect(cfg.StoreDatabaseConfig, cfg.CommandDatabaseConfig)

	// setup CQRS
	amqpAddress := cfg.AMQPConfig.RabbitMQConnectionString
	logger := watermill.NewStdLogger(false, false)
	commandsAMQPConfig := amqp.NewDurableQueueConfig(amqpAddress)
	commandsPublisher, err := amqp.NewPublisher(commandsAMQPConfig, logger)
	if err != nil {
		panic(err)
	}
	commandsSubscriber, err := amqp.NewSubscriber(commandsAMQPConfig, logger)
	if err != nil {
		panic(err)
	}

	eventsPublisher, err := amqp.NewPublisher(amqp.NewDurablePubSubConfig(amqpAddress, nil), logger)
	if err != nil {
		panic(err)
	}
	cqrsRouter, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	cqrsRouter.AddMiddleware(watermillMiddleware.Recoverer)
	// cqrshandler.Facade is facade for Command and Event buses and processors.
	// You can use facade, or create buses and processors manually (you can inspire with cqrshandler.NewFacade)
	cqrsFacade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return commandName
		},
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			return []cqrs.CommandHandler{
				cqrshandler.CreateTenantHandler{EventBus: eb},
				cqrshandler.UpdateTenantCredentialsHandler{EventBus: eb},
				cqrshandler.UpdateTenantLicenseHandler{EventBus: eb},
				cqrshandler.UpdateTenantStatusHandler{EventBus: eb},
			}
		},
		CommandsPublisher: commandsPublisher,
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			// we can reuse subscriber, because all commands have separated topics
			return commandsSubscriber, nil
		},
		GenerateEventsTopic: func(eventName string) string {
			// because we are using PubSub RabbitMQ config, we can use one topic for all events
			return "events"

			// we can also use topic per event type
			// return eventName
		},

		EventsPublisher: eventsPublisher,
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			config := amqp.NewDurablePubSubConfig(
				amqpAddress,
				amqp.GenerateQueueNameTopicNameWithSuffix(handlerName),
			)

			return amqp.NewSubscriber(config, logger)
		},
		Router:                cqrsRouter,
		CommandEventMarshaler: cqrs.JSONMarshaler{},
		Logger:                logger,
	})
	if err != nil {
		panic(err)
	}
	go startRouter(cqrsRouter)
	// public routes
	router.Group(func(r chi.Router) {
		router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("Up and running " + appName + " :)"))
		})
		router.Get("/tenant", func(writer http.ResponseWriter, request *http.Request) {
			api.TenantApi{}.GetAllTenants(writer, request)
		})

		router.Get("/tenant/{tenant_id}", func(writer http.ResponseWriter, request *http.Request) {
			api.TenantApi{}.GetTenantById(writer, request)
		})
		router.Post("/tenant", func(writer http.ResponseWriter, request *http.Request) {
			api.TenantApi{}.AddNewTenant(writer, request, cqrsFacade)
		})
	})
	// private routes
	router.Group(func(r chi.Router) {
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens.
		r.Use(jwtauth.Authenticator)

		r.Put("/tenant/credentials", func(writer http.ResponseWriter, request *http.Request) {
			api.TenantApi{}.UpdateTenantCredentials(writer, request, cqrsFacade)
		})

		r.Delete("/tenant/license", func(writer http.ResponseWriter, request *http.Request) {
			api.TenantApi{}.UpdateTenantLicense(writer, request, cqrsFacade)
		})

		r.Post("/tenant/status", func(writer http.ResponseWriter, request *http.Request) {
			api.TenantApi{}.UpdateTenantStatus(writer, request, cqrsFacade)
		})
	})

	// start the server
	router.Mount("/swagger", httpSwagger.WrapHandler)

	logrus.Println(appName + " serving on port " + port + " with profile " + cfg.Environment)
	err = http.ListenAndServe(port, router)
	if err != nil {
		logrus.Fatalf("Error starting the server: %s", err)
		return
	}
}

func startRouter(cqrsRouter *message.Router) {
	// processors are based on router, so they will work when router will start
	if err := cqrsRouter.Run(context.Background()); err != nil {
		panic(err)
	}
}
