package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Contact struct {
	UserID  string
	Name    string
	Phone   string
	Address string
}

func main() {
	s := session.Must(session.NewSession())
	db := dynamodb.New(s)

	contact := Contact{
		UserID:  "0d50ab52",
		Name:    "Homer Simpson",
		Phone:   "555555556",
		Address: "125 Fake Ave",
	}

	av, err := dynamodbattribute.MarshalMap(contact)

	if err != nil {
		// handle error
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("contacts"),
		Item:      av,
	})

	if err != nil {
		// handle error
	}

	o, err := db.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String("0d50ab52"),
			},
			"Name": {
				S: aws.String("Homer Simpson"),
			},
		},
		TableName: aws.String("contacts"),
	})

	if err != nil {
		// handle error
	}

	var c Contact

	err = dynamodbattribute.UnmarshalMap(o.Item, &c)

	if err != nil {
		// handle error
	}

	_, err = db.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#A": aws.String("Address"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				S: aws.String("742 Evergreen Terrace"),
			},
		},
		UpdateExpression: aws.String("SET #A = :a"),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String("0d50ab52"),
			},
			"Name": {
				S: aws.String("Homer Simpson"),
			},
		},
		TableName: aws.String("contacts"),
	})

	if err != nil {
		// handle error
	}

	res, err := db.Query(&dynamodb.QueryInput{
		KeyConditionExpression: aws.String("Address = :a"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				S: aws.String("742 Evergreen Terrace"),
			},
		},
		IndexName: aws.String("contacts-address_index"),
		TableName: aws.String("contacts"),
	})

	if err != nil {
		// handle error
	}

	// Will only contain attributes projected into GSI!
	var cIDs []Contact

	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &cIDs)

	if err != nil {
		// handle error
	}
}
