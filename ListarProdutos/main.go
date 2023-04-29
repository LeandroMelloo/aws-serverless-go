package main

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Produto struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Preco int    `json:"preco"`
}

func ListarProdutos(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sessao := session.Must(session.NewSession()) // abro uma nova sessão
	servico := dynamodb.New(sessao)

	entrada := &dynamodb.ScanInput{
		TableName: aws.String("Produtos"),
	}

	resultado, erro := servico.Scan(entrada)
	if erro != nil {
		return events.APIGatewayProxyResponse{Body: erro.Error(), StatusCode: 500}, nil
	}

	var produtos []Produto
	for _, item := range resultado.Items {
		preco, erro := strconv.Atoi(*item["preco"].N) // N -> transforma o preco em um número
		if erro != nil {
			return events.APIGatewayProxyResponse{Body: erro.Error(), StatusCode: 500}, nil
		}

		produtos = append(produtos, Produto{
			ID:    *item["id"].S,
			Nome:  *item["nome"].S,
			Preco: preco,
		})
	}

	// converter em Json para retornar
	body, erro := json.Marshal(produtos)
	if erro != nil {
		return events.APIGatewayProxyResponse{Body: erro.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(ListarProdutos)
}
