package app

import (
	"github.com/Chrisentech/aluta-market-api/internals/cart"
	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/skynet"
	"github.com/Chrisentech/aluta-market-api/internals/store"
	"github.com/Chrisentech/aluta-market-api/internals/user"
)

type PackageType int

const (
	UserPackage PackageType = iota
	ProductPackage
	CartPackage
	StorePackage
	SkynetPackage
)

func CreateRepository(pkgType PackageType) interface{} {

	// var
	switch pkgType {
	case UserPackage:
		return user.NewRepository()
	case ProductPackage:
		return product.NewRepository()
	case CartPackage:
		return cart.NewRepository()
	case StorePackage:
		return store.NewRepository()
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
