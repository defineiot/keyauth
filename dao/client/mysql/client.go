package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/client"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

func (s *store) CreateClient(clientType client.Type, redirectURI string) (*client.Client, error) {
	id := tools.MakeUUID(24)
	secret := tools.MakeUUID(32)
	cli := client.Client{ID: id, Secret: secret, RedirectURI: redirectURI, Type: clientType}

	_, err := s.stmts[CreateClient].Exec(cli.ID, cli.Secret, string(cli.Type), cli.RedirectURI)
	if err != nil {
		return nil, exception.NewInternalServerError("insert client exec sql err, %s", err)
	}

	return &cli, nil
}

func (s *store) ListClients() ([]*client.Client, error) {
	return nil, nil
}

func (s *store) GetClient(id string) (*client.Client, error) {
	cli := new(client.Client)
	err := s.stmts[FindOneByID].QueryRow(id).Scan(&cli.ID, &cli.Secret, &cli.Type, &cli.RedirectURI, &cli.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("client %s not find", id)
		}

		return nil, exception.NewInternalServerError("query single client error, %s", err)
	}

	return cli, nil
}

func (s *store) DeleteClient(id string) error {
	ret, err := s.stmts[DeleteClient].Exec(id)
	if err != nil {
		return exception.NewInternalServerError("delete client exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("client %s not exist", id)
	}

	return nil
}
