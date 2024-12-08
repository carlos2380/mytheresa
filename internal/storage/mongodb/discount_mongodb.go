package mongodb

import (
	"context"
	"fmt"
	"mytheresa/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Discount struct {
	CollectionMDB *mongo.Collection
	Ctx           context.Context
}

func (discountMdb *Discount) GetDiscount(category string) (int, error) {
	ctx, cancel := context.WithCancel(discountMdb.Ctx)
	defer cancel()

	//Set filter category and order price Desc
	filter := bson.D{
		{Key: "category", Value: category},
	}
	opts := options.FindOne().SetSort(bson.D{
		{Key: "price", Value: -1},
	})

	var discount models.DiscountCategory
	err := discountMdb.CollectionMDB.FindOne(ctx, filter, opts).Decode(&discount)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, fmt.Errorf("error al buscar descuento: %w", err)
	}

	return discount.Percent, nil
}
