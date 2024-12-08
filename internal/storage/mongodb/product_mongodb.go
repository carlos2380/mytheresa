package mongodb

import (
	"context"
	"fmt"
	"mytheresa/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	CollectionMDB *mongo.Collection
	Ctx           context.Context
	PageSize      int
}

func (productMdb *Product) GetProducts(category string, priceLessThan *int, cursor string) ([]models.Product, string, error) {

	ctx, cancel := context.WithCancel(productMdb.Ctx)
	defer cancel()

	//Add filters category, price and cursor if is necesarry
	//This could move it on a 3 differents functions to re-use
	filter := bson.D{}
	if category != "" {
		filter = append(filter, bson.E{Key: "category", Value: category})
	}
	if priceLessThan != nil {
		filter = append(filter, bson.E{Key: "price", Value: bson.D{{Key: "$lt", Value: *priceLessThan}}})
	}
	if cursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor: %w", err)
		}
		filter = append(filter, bson.E{Key: "_id", Value: bson.D{{Key: "$gt", Value: cursorID}}})
	}

	opts := options.Find()
	opts.SetLimit(int64(productMdb.PageSize) + 1) //6 to know if it's the finish

	cursorResult, err := productMdb.CollectionMDB.Find(ctx, filter, opts)
	if err != nil {
		return nil, "", fmt.Errorf("error on find product: %w", err)
	}
	defer cursorResult.Close(ctx)

	var products []models.Product
	for cursorResult.Next(ctx) {
		var product models.Product
		if err := cursorResult.Decode(&product); err != nil {
			return nil, "", fmt.Errorf("error on decode product: %w", err)
		}
		products = append(products, product)
	}

	if err := cursorResult.Err(); err != nil {
		return nil, "", fmt.Errorf("error on iterate cursor: %w", err)
	}

	var nextCursor string
	if len(products) == (productMdb.PageSize + 1) {
		nextCursor = products[(productMdb.PageSize - 1)].Id.Hex()
		products = products[:productMdb.PageSize]
	}

	return products, nextCursor, nil
}
