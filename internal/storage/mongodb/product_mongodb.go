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
}

func (productMdb *Product) GetProducts(category string, priceLessThan *int, cursor string) ([]models.Product, string, error) {
	// Crear un contexto con un tiempo límite (por ejemplo, 5 segundos)
	// Crear un contexto con un tiempo límite (por ejemplo, 5 segundos)
	ctx, cancel := context.WithCancel(productMdb.Ctx)
	defer cancel() // Asegura que se liberen los recursos del contexto
	// Filtro para buscar por categoría y precio
	filter := bson.D{}
	if category != "" {
		filter = append(filter, bson.E{Key: "category", Value: category})
	}
	if priceLessThan != nil {
		filter = append(filter, bson.E{Key: "price", Value: bson.D{{Key: "$lt", Value: *priceLessThan}}})
	}
	if cursor != "" {
		// Convertir el cursor (_id) de string a ObjectId
		cursorID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, "", fmt.Errorf("cursor inválido: %w", err)
		}
		// Agregar filtro para paginar a partir del _id
		filter = append(filter, bson.E{Key: "_id", Value: bson.D{{Key: "$gt", Value: cursorID}}})
	}

	// Opciones para limitar la cantidad de resultados
	opts := options.Find()
	opts.SetLimit(6) // Limitar a 5 resultados por página

	// Ejecutar la consulta
	cursorResult, err := productMdb.CollectionMDB.Find(ctx, filter, opts)
	if err != nil {
		return nil, "", fmt.Errorf("error al buscar productos: %w", err)
	}
	defer cursorResult.Close(ctx)

	// Decodificar los resultados
	var products []models.Product

	for cursorResult.Next(ctx) {
		var product models.Product
		if err := cursorResult.Decode(&product); err != nil {
			return nil, "", fmt.Errorf("error al decodificar producto: %w", err)
		}
		products = append(products, product)
	}

	// Verificar errores del cursor
	if err := cursorResult.Err(); err != nil {
		return nil, "", fmt.Errorf("error durante la iteración del cursor: %w", err)
	}

	// Generar el nuevo cursor basado en el último _id
	var nextCursor string
	if len(products) == 6 { // Si se obtienen 5 resultados, hay más páginas
		nextCursor = products[4].Id.Hex() // El 6º producto determina el próximo cursor
		products = products[:5]
	}

	return products, nextCursor, nil
}
