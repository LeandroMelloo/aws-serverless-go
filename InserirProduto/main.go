package main

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type Produto struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Preco int    `json:"preco"`
}

func InserirProduto(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var produto Produto
	err := json.Unmarshal([]byte(request.Body), &produto)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	produto.ID = uuid.New().String()

	sessao := session.Must(session.NewSession()) // abro uma nova sessão
	servicos := dynamodb.New(sessao)

	entrada := &dynamodb.PutItemInput{
		TableName: aws.String("Produtos"),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(produto.ID),
			},
			"nome": {
				S: aws.String(produto.Nome),
			},
			"price": {
				S: aws.String(strconv.Itoa(produto.Preco)),
			},
		},
	}

	_, err = servicos.PutItem(entrada)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	body, err := json.Marshal(produto)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(InserirProduto)
}
