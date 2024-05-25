package myaws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/spf13/viper"
)

var SnsSvc *sns.SNS

func SetupSNSService(config map[string]string) error {
	var err error
	var ok bool
	myRegion := "ap-southeast-1" // <- your region
	myProfile := ""
	if myRegion, ok = config["region"]; !ok {
		err = errors.New("Missing region config.")
		return err
	}
	if myProfile, ok = config["profile"]; !ok {
		err = errors.New("Missing profile config")
		return err
	}
	if err == nil {
		fmt.Printf("SNS service region:%v, profile:%v", myRegion, myProfile)
		sess = session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(myRegion),
			Credentials: credentials.NewSharedCredentials("", myProfile),
		}))
		SnsSvc = sns.New(sess, &aws.Config{
			Region: aws.String(myRegion),
		})
	}
	return err
}

func SendSNS(topic, message string) (string, error) {
	result, err := SnsSvc.Publish(&sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: &topic,
	})
	if err != nil {
		return "", err
	}

	return *result.MessageId, err
}

func LoadAndSetupSNS(vp *viper.Viper) {
	myRegion := "ap-southeast-1"
	if vp.IsSet("sns.region") {
		myRegion = vp.GetString("sns.region")
	}

	myProfile := "default"
	if vp.IsSet("sns.profile") {
		myProfile = vp.GetString("sns.profile")
	}

	mm := make(map[string]string)
	mm["region"] = myRegion
	mm["profile"] = myProfile
	SetupSNSService(mm)
}
