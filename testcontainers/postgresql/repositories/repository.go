// Package repositories contém as implementações dos repositórios para acesso a dados.
package repositories

import "database/sql"

// Repository é uma interface que define o contrato para os métodos de acesso a dados.
type Repository interface {
	Ping() error
}

// repository é uma estrutura que contém a conexão com o banco de dados e implementa a
// interface Repository para acessar e manipular os dados armazenados.
type repository struct {
	db *sql.DB
}

// Ping é um método que verifica se a conexão com o banco de dados está ativa e funcionando.
// Retorna um erro, se houver algum problema de conexão.
func (r repository) Ping() error {
	return r.db.Ping()
}

// NewRepository cria uma nova instância do repositório utilizando a conexão de banco de dados fornecida.
// Retorna uma instância que implementa a interface Repository.
func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}
