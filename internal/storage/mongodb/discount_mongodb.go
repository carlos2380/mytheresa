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
	// Crear un contexto con un tiempo límite (por ejemplo, 5 segundos)
	ctx, cancel := context.WithCancel(discountMdb.Ctx)
	defer cancel() // Asegura que se liberen los recursos del contexto

	// Filtro para buscar solo por categoría
	filter := bson.D{
		{Key: "category", Value: category},
	}

	// Opciones para ordenar por price en orden descendente (-1)
	opts := options.FindOne().SetSort(bson.D{
		{Key: "price", Value: -1},
	})

	// Ejecutar la consulta y obtener un único documento
	var discount models.DiscountCategory
	err := discountMdb.CollectionMDB.FindOne(ctx, filter, opts).Decode(&discount)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No se encontró ningún documento, retornar 0 como porcentaje
			return 0, nil
		}
		// Manejar otros errores
		return 0, fmt.Errorf("error al buscar descuento: %w", err)
	}

	// Retornar el porcentaje de descuento encontrado
	return discount.Percent, nil
}
