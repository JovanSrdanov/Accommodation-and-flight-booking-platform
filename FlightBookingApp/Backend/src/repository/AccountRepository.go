package repository

import (
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"context"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type accountRepository struct {
	base     Repository
	Accounts []model.Account
}

type AccountRepository interface {
	Create(newAccount model.Account) (model.Account, error)
	GetAll() []model.Account
	GetById(id uuid.UUID) (model.Account, error)
	Delete(id primitive.ObjectID) error
}

func NewAccountRepository(ctx context.Context, logger *log.Logger) (*accountRepository, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	return &accountRepository{
		base: Repository{
			client: client,
			logger: logger,
		},
	}, nil
}

func (repo *accountRepository) Create(account model.Account) model.Account {
	repo.base.Connect(context.Background())
	defer repo.base.Disconnect(context.Background())

	repo.Accounts = append(repo.Accounts, account)
	return account
}

func (repo *accountRepository) GetAll() []model.Account {
	repo.base.Connect(context.Background())
	defer repo.base.Disconnect(context.Background())

	return repo.Accounts
}

func (repo *accountRepository) GetById(id uuid.UUID) (model.Account, error) {
	repo.base.Connect(context.Background())
	defer repo.base.Disconnect(context.Background())

	for _, account := range repo.Accounts {
		if account.ID == id {
			return account, nil
		}
	}

	return model.Account{}, &errors.NotFoundError{}
}

func (repo *accountRepository) Delete(id primitive.ObjectID) error {
	return nil
}