package graph

import (
	"context"
	"fmt"
	"sort"

	"github.com/samstringzz/alutamarket-backend/graph/model"
)

// GetStoreTransactions is the resolver for the getStoreTransactions field.
func (r *queryResolver) GetStoreTransactions(ctx context.Context, storeID int) (*model.StoreTransactions, error) {
	// Get DVA deposits
	var dvaTransactions []*model.DepositTransaction
	if err := r.DB.Raw(`
		SELECT
			t.id::text as id,
			'dva' as type,
			t.amount::float as amount,
			t.reference,
			t.status,
			t.created_at,
			'Direct deposit to DVA account' as description
		FROM transactions t
		JOIN dva_accounts dva ON dva.customer_id = t.customer->>'customer_code'
		WHERE dva.store_id = ?
		ORDER BY t.created_at DESC
	`, storeID).Scan(&dvaTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch DVA transactions: %v", err)
	}

	// Get order payments
	var orderTransactions []*model.DepositTransaction
	if err := r.DB.Raw(`
		SELECT
			o.id::text as id,
			'order' as type,
			o.total_amount::float as amount,
			o.reference,
			o.status,
			o.created_at,
			'Payment for order #' || o.id::text as description
		FROM orders o
		WHERE o.store_id = ? AND o.status != 'cancelled'
		ORDER BY o.created_at DESC
	`, storeID).Scan(&orderTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch order transactions: %v", err)
	}

	// Get withdrawals
	var withdrawals []*model.WithdrawalTransaction
	if err := r.DB.Raw(`
		SELECT
			id::text as id,
			amount,
			status,
			bank_name,
			account_number,
			account_name,
			created_at,
			completed_at
		FROM withdrawals
		WHERE store_id = ?
		ORDER BY created_at DESC
	`, storeID).Scan(&withdrawals).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch withdrawals: %v", err)
	}

	// Combine DVA and order transactions
	deposits := append(dvaTransactions, orderTransactions...)
	// Sort deposits by created_at
	sort.Slice(deposits, func(i, j int) bool {
		return deposits[i].CreatedAt.After(deposits[j].CreatedAt)
	})

	return &model.StoreTransactions{
		Deposits:    deposits,
		Withdrawals: withdrawals,
	}, nil
}
