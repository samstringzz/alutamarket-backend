package app

import (
	"github.com/samstringzz/alutamarket-backend/internals/cart"
	"github.com/samstringzz/alutamarket-backend/internals/messages"
	"github.com/samstringzz/alutamarket-backend/internals/product"
	"github.com/samstringzz/alutamarket-backend/internals/skynet"
	"github.com/samstringzz/alutamarket-backend/internals/store"
	"github.com/samstringzz/alutamarket-backend/internals/user"
)

type PackageType int

const (
	UserPackage PackageType = iota
	ProductPackage
	CartPackage
	StorePackage
	SkynetPackage
	MessagePackage
)

func CreateRepository(pkgType PackageType) interface{} {
	switch pkgType {
	case UserPackage:
		return user.NewRepository()
	case ProductPackage:
		return product.NewRepository()
	case CartPackage:
		return cart.NewRepository()
	case StorePackage:
		return store.NewRepository()
	case MessagePackage:
		return messages.NewRepository()
	case SkynetPackage:
		return skynet.NewRepository()
	// Add more cases for other packages as needed
	default:
		return nil // Handle unknown package types gracefully
	}
}
func InitializePackage(pkgType PackageType) interface{} {
	return CreateRepository(pkgType)
}
