package myaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func ScanSymbolConfigure(region, profile, tableName string) (*dynamodb.ScanOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials("", profile),
	})
	if err != nil {
		return nil, err
	}
	svc := dynamodb.New(sess)

	input := &dynamodb.ScanInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				S: aws.String("Symbol"),
			},
		},
		FilterExpression:     aws.String("configType = :a"),
		ProjectionExpression: aws.String("id, counterParty, baseCurrency, quoteCurrency, symbol, priceDecimal, quantityDecimal, quantityMin, quantityMax, priceStep, quantityStep"),
		TableName:            aws.String(tableName),
	}
	result, err := svc.Scan(input)
	return result, err
}

func ScanUser(region, profile, tableName string) (*dynamodb.ScanOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials("", profile),
	})
	if err != nil {
		return nil, err
	}
	svc := dynamodb.New(sess)

	filt := expression.Name("__typename").Equal(expression.Value("User"))
	// Or we could get by ratings and pull out those with the right year later
	//    filt := expression.Name("info.rating").GreaterThan(expression.Value(min_rating))

	// Get back the title, year, and rating
	proj := expression.NamesList(expression.Name("id"), expression.Name("name"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err == nil {
		// Build the query input parameters
		params := &dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection(),
			TableName:                 aws.String(tableName),
		}

		// Make the DynamoDB Query API call
		result, err := svc.Scan(params)
		return result, err
	}
	return nil, err
}
