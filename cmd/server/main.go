package main

import (
	"context"
	"flag"
	"log"
	"strconv"

	product_application "mytheresa/internal/application/product"
	"mytheresa/internal/handlers"
	"mytheresa/internal/server"
	"mytheresa/internal/storage/mongodb"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	port := flag.String("port", "8000", "Port on which the server will be listening for incoming requests.")
	portMongodb := flag.String("port_mongodb", "27017", "Port on which the MongoDB server will be serving connections.")
	ipMongodb := flag.String("ip_mongodb", "mongodb", "IP address on which the MongoDB server will be listening.")
	databaseName := flag.String("database_name", "mytheresadb", "Name of the database connection.")
	pageSizeStr := flag.String("page_size_product", "5", "Number the products per page response.")
	flag.Parse()

	pageSize, err := strconv.Atoi(*pageSizeStr)
	if err != nil || pageSize <= 0 {
		log.Fatalf("Invalid page_size_product: %v", err)
	}

	mongoCtx, mongoCancel := context.WithCancel(context.Background())
	defer mongoCancel()

	client, database, err := mongodb.InitMongo(*ipMongodb, *portMongodb, *databaseName)
	if err != nil {
		log.Fatal("Error mongodb:", err)
	}

	defer func() {
		if err := mongodb.CloseMongo(client); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	mongoDiscount := &mongodb.Discount{
		CollectionMDB: database.Collection("DiscountCategory"),
		Ctx:           mongoCtx,
	}

	mongoProduct := &mongodb.Product{
		CollectionMDB: database.Collection("Product"),
		Ctx:           mongoCtx,
		PageSize:      pageSize,
	}

	appProduct := &product_application.ProductApplication{
		StgProduct:  mongoProduct,
		StgDiscount: mongoDiscount,
	}

	hProduct := &handlers.ProductHandler{
		AppProduct: appProduct,
	}
	router := server.NewRouter(hProduct)
	srv := &http.Server{
		Addr:    ":" + *port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Print("Server Started")
	log.Printf("Listening on 0.0.0.0:%s", *port)
	<-done

	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
