package main

import (
	"context"
	_ "embed"
	"io"
	"log"
	"os"
	"path/filepath"
	"write/config"
	"write/server"
	"write/service"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	initLogging()
	err := config.Init()
	if err != nil {
		logrus.Fatal(err)
	}

	tp, err := initTracer()
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			logrus.Error(err)
		}
	}()

	service.Init()
	server.Start()
}

func initLogging() {
	logFolderPath := config.Configuration.LogPath
	os.MkdirAll(logFolderPath, 0755)
	logFilePath := filepath.Join(filepath.Clean(logFolderPath), "write-service.log")
	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	logDest := io.MultiWriter(logFile, os.Stdout)
	if err != nil {
		log.Println("could not open log folder path", err)
		return
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(logDest)
}

func initTracer() (*sdktrace.TracerProvider, error) {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}

	resources, err := resource.New(context.Background(), resource.WithFromEnv())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resources),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)
	return tp, nil
}
