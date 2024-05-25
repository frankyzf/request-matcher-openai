package myaws

import (
	"github.com/spf13/viper"
)

var CognitoRegion, CognitoUserPoolId, CognitoClientId, CognitoProfile string
var CognitoVerify bool

func LoadAndSetupCognito(vp *viper.Viper) {
	CognitoRegion = "ap-southeast-1"
	if vp.IsSet("cognito.region") {
		CognitoRegion = vp.GetString("cognito.region")
	}

	CognitoUserPoolId = "ap-southeast-1_0FwlJiyZf"
	if vp.IsSet("cognito.userPoolId") {
		CognitoUserPoolId = vp.GetString("cognito.userPoolId")
	}

	CognitoClientId = "5gdkjogchqdqoppeck0l3d0bjp"
	if vp.IsSet("cognito.clientId") {
		CognitoClientId = vp.GetString("cognito.clientId")
	}

	CognitoProfile = "default"
	if vp.IsSet("cognito.profile") {
		CognitoProfile = vp.GetString("cognito.profile")
	}

	CognitoVerify = false
	if vp.IsSet("cognito.verify") {
		CognitoVerify = vp.GetBool("cognito.verify")
	}

	Cognito := make(map[string]string, 0)
	Cognito["region"] = CognitoRegion
	Cognito["userPoolId"] = CognitoUserPoolId
	Cognito["clientId"] = CognitoClientId
	Cognito["profile"] = CognitoProfile
	SetupCognito(Cognito)
}
