package storage

type DiscountStoreage interface {
	GetDiscount(category string) (int, error)
}
