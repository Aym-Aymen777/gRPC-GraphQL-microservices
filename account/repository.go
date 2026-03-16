package account

import (
	"context"
	"database/sql"
	"log"
	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/account/types"
	_ "github.com/go-sql-driver/mysql"
)

type Repository interface {
	Close()
	PutAccount(ctx context.Context, account *types.Account) error
	GetAccount(ctx context.Context, id string) (*types.Account, error)
	ListAccounts(ctx context.Context, skip, limit uint64) ([]*types.Account, error)
}

type mySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(url string) (Repository, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &mySQLRepository{db: db}, nil
}

func (r *mySQLRepository) Close() {
	err := r.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (r *mySQLRepository) Ping() error {
	return r.db.Ping()
}

func (r *mySQLRepository) PutAccount(ctx context.Context, account *types.Account) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts (id, username, email) VALUES (?, ?, ?)", account.ID, account.Username, account.Email)
	return err
}

func (r *mySQLRepository) GetAccount(ctx context.Context, id string) (*types.Account, error) {
	row, err := r.db.QueryContext(ctx, "SELECT id,username,email FROM accounts WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		return nil, sql.ErrNoRows
	}

	account := &types.Account{}
	err = row.Scan(&account.ID, &account.Username, &account.Email)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *mySQLRepository) ListAccounts(ctx context.Context, skip, limit uint64) ([]*types.Account, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id,username,email FROM accounts LIMIT ? OFFSET ?", limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*types.Account
	for rows.Next() {
		account := &types.Account{}
		err = rows.Scan(&account.ID, &account.Username, &account.Email)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
