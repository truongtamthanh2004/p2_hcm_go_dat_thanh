package repository

import (
	"payment-service/internal/model"

	"gorm.io/gorm"
)

type TransactionRepository interface {
  Create(tx *model.PaymentTransaction) error
  FindByTxnRef(txnRef string) (*model.PaymentTransaction, error)
  Update(tx *model.PaymentTransaction) error
}

type transactionRepositoryImpl struct {
  db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
  return &transactionRepositoryImpl{db}
}

func (r *transactionRepositoryImpl) Create(tx *model.PaymentTransaction) error {
  return r.db.Create(tx).Error
}

func (r *transactionRepositoryImpl) FindByTxnRef(txnRef string) (*model.PaymentTransaction, error) {
  var tx model.PaymentTransaction
  err := r.db.Where("txn_ref = ?", txnRef).First(&tx).Error
  if err != nil {
    return nil, err
  }
  return &tx, nil
}

func (r *transactionRepositoryImpl) Update(tx *model.PaymentTransaction) error {
  return r.db.Save(tx).Error
}
