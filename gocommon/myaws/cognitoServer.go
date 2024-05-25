package myaws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	set "github.com/deckarep/golang-set"
)

var sRegion = ""     // <- your region
var sUserPoolId = "" // <- add your pool id here
var sClientId = ""
var sProfile = ""
var sess *session.Session
var identitySvc *cognitoidentityprovider.CognitoIdentityProvider

func SetupCognitoServer(config map[string]string) error {
	var err error
	var ok bool
	sRegion, ok = config["region"]
	if !ok {
		err = errors.New("Missing region config.")
	}
	sUserPoolId, ok = config["userPoolId"]
	if !ok {
		err = errors.New("Missing poolId config")
	}
	sClientId, ok = config["clientId"]
	if !ok {
		err = errors.New("Missing clientId config")
	}
	sProfile, ok = config["profile"]
	if !ok {
		err = errors.New("Missing profile config")
	}
	if err == nil {
		fmt.Printf("region:%v, sUserPoolId:%v clientId:%v", sRegion, sUserPoolId, sClientId)
		sess = session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(sRegion),
			Credentials: credentials.NewSharedCredentials("", sProfile),
		}))
		identitySvc = cognitoidentityprovider.New(sess, &aws.Config{
			Region: aws.String(sRegion),
		})
	}
	return err
}

func SignUp(input cognitoidentityprovider.SignUpInput) (*cognitoidentityprovider.SignUpOutput, error) {
	output, err := identitySvc.SignUp(&input)
	return output, err
}

func UpdateCustomAttribute(userName, name, value string) error {
	var err error
	attribute := &cognitoidentityprovider.AttributeType{
		Name:  aws.String(name),
		Value: aws.String(value),
	}
	input := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			attribute,
		},
		UserPoolId: aws.String(sUserPoolId),
		Username:   aws.String(userName),
	}
	_, err1 := identitySvc.AdminUpdateUserAttributes(input)
	err = err1
	return err
}

func SyncUserGroup(username string, groups []string) error {
	var err error
	input := cognitoidentityprovider.AdminListGroupsForUserInput{
		Limit:      aws.Int64(50),
		UserPoolId: aws.String(sUserPoolId),
		Username:   aws.String(username),
	}
	output, err1 := identitySvc.AdminListGroupsForUser(&input)
	err = err1
	if err == nil {
		newGroups := set.NewSet()
		for _, g := range groups {
			newGroups.Add(g)
		}
		oldGroups := set.NewSet()
		cogGroups := output.Groups
		for _, g := range cogGroups {
			oldGroups.Add(*g.GroupName)
		}
		newlyAdded := newGroups.Difference(oldGroups)
		for _, g := range newlyAdded.ToSlice() {
			group := g.(string)
			AddUserToGroup(group, username)
		}
		oldRemove := oldGroups.Difference(newGroups)
		for _, g := range oldRemove.ToSlice() {
			group := g.(string)
			RemoveUserFromGroup(group, username)
		}
	}
	return err
}

func AddUserToGroup(group, username string) {
	input := cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  aws.String(group),
		UserPoolId: aws.String(sUserPoolId),
		Username:   aws.String(username),
	}
	identitySvc.AdminAddUserToGroup(&input)
}

func RemoveUserFromGroup(group, username string) {
	removeInput := cognitoidentityprovider.AdminRemoveUserFromGroupInput{
		GroupName:  aws.String(group),
		UserPoolId: aws.String(sUserPoolId),
		Username:   aws.String(username),
	}
	identitySvc.AdminRemoveUserFromGroup(&removeInput)
}
