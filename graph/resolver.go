package graph

import (
	"github.com/samstringzz/alutamarket-backend/internals/messages"
	"github.com/samstringzz/alutamarket-backend/internals/product"
	"github.com/samstringzz/alutamarket-backend/internals/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserHandler    *user.Handler
	ProductHandler *product.Handler
	MessageHandler *messages.Handler
}

func NewResolver(userHandler *user.Handler, productHandler *product.Handler, messageHandler *messages.Handler) *Resolver {
	return &Resolver{
		UserHandler:    userHandler,
		ProductHandler: productHandler,
		MessageHandler: messageHandler,
	}
}
