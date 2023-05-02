// Package repositories contém as implementações dos repositórios para acesso a dados.
package repositories

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

// criarConfigurarContainer é uma função auxiliar que cria e configura um contêiner Docker
// para testes usando a imagem Docker, as portas expostas e as variáveis de ambiente fornecidas.
// Retorna um contêiner testcontainers.Container e um erro, se houver.
func criarConfigurarContainer(imagemDocker string, portasExpostas []string, variaveisAmbiente map[string]string) (testcontainers.Container, error) {
	// Configuração do contêiner
	containerRequest := testcontainers.ContainerRequest{
		Image:        imagemDocker,
		ExposedPorts: portasExpostas,
		Env:          variaveisAmbiente,
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithStartupTimeout(2 * time.Second),
	}

	// Criação do contêiner
	conteiner, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao criar o contêiner: %w", err)
	}

	return conteiner, nil
}

// TestRepository_Ping é um teste que verifica se o repositório é capaz de se comunicar
// com o banco de dados. Ele utiliza a função criarConfigurarContainer para criar e configurar
// um contêiner Docker com um banco de dados PostgreSQL para testes.
func TestRepository_Ping(t *testing.T) {
	portasExpostas := []string{"5432/tcp"}
	variaveisAmbiente := map[string]string{"POSTGRES_PASSWORD": "postgres", "POSTGRES_DB": "postgres", "POSTGRES_USER": "postgres"}
	container, err := criarConfigurarContainer("docker.io/postgres:15.2-alpine", portasExpostas, variaveisAmbiente)
	if err != nil {
		t.Fatalf("erro ao criar o contêiner: %s", err)
	}
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("erro ao finalizar o contêiner: %s", err)
		}
	}(container, context.Background())

	// Obtenha a porta mapeada no host
	portaMapeada, err := container.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		t.Fatalf("erro ao obter a porta mapeada: %s", err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", portaMapeada.Port(), "postgres", "postgres", "postgres"))
	if err != nil {
		t.Fatalf("erro ao conectar no banco de dados: %s", err)
	}

	repository := NewRepository(db)
	time.Sleep(1 * time.Second)
	err2 := repository.Ping()
	defer db.Close()
	assert.NoError(t, err2)
}
