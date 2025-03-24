package graph

import (
	"github.com/Chrisentech/aluta-market-api/internals/messages"
	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserHandler    user.Handler
	ProductHandler *product.Handler
	ProductService product.Service
	MessageHandler *messages.Handler
}

func NewResolver(userHandler user.Handler, productService product.Service, productHandler *product.Handler, messageHandler *messages.Handler) *Resolver {
	return &Resolver{
		UserHandler:    userHandler,
		ProductService: productService,
		ProductHandler: productHandler,
		MessageHandler: messageHandler,
	}
}
