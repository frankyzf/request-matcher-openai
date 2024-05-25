package myaws

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var cognitoSetting map[string]string
var jwkMap map[string]JWKKey
var region string
var userPoolID string

// JWK is json data struct for JSON Web Key
type JWK struct {
	Keys []JWKKey
}

// JWKKey is json data struct for cognito jwk key
type JWKKey struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
}

func SetupCognito(m map[string]string) {
	cognitoSetting = m
	region = "ap-southeast-1"
	if value, ok := cognitoSetting["region"]; ok {
		region = value
	}
	userPoolID = "ap-southeast-1_wLd4XhaVh"
	if value, ok := cognitoSetting["userPoolId"]; ok {
		userPoolID = value
	}
	jwkMap = make(map[string]JWKKey, 0)
	GetKeys()
}

func GetKeys() {
	jwkURL := fmt.Sprintf("https://cognito-idp.%v.amazonaws.com/%v/.well-known/jwks.json", region, userPoolID)
	fmt.Println(jwkURL)
	jwkMap = getJWK(jwkURL)
	fmt.Println(jwkMap)
}

func Verify(tokenStr string) (*jwt.Token, error) {
	token, err := validateToken(tokenStr)
	if err != nil || !token.Valid {
		if err == nil {
			err = errors.New("Unsuccessful verification.")
		}
		return nil, err
	}
	return token, nil
}

func validateToken(tokenStr string) (*jwt.Token, error) {
	// 2. Decode the token string into JWT format.
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// cognito user pool : RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// 5. Get the kid from the JWT token header and retrieve the corresponding JSON Web Key that was stored
		if kid, ok := token.Header["kid"]; ok {
			if kidStr, ok := kid.(string); ok {
				key := jwkMap[kidStr]
				// 6. Verify the signature of the decoded JWT token.
				rsaPublicKey := convertKey(key.E, key.N)
				return rsaPublicKey, nil
			}
		}
		return "", nil
	})
	if err != nil {
		return token, err
	}
	claims := token.Claims.(jwt.MapClaims)
	iss, ok := claims["iss"]
	if !ok {
		return token, fmt.Errorf("token does not contain issuer")
	}
	issStr := iss.(string)
	if strings.Contains(issStr, "cognito-idp") {
		err = validateAWSJwtClaims(claims)
		if err != nil {
			return token, err
		}
	}
	if token.Valid {
		return token, nil
	}
	return token, err
}

func GetOwner(token *jwt.Token) string {
	claims := token.Claims.(jwt.MapClaims)
	owner := claims["cognito:username"].(string)
	return owner
}

func GetEmail(token *jwt.Token) string {
	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	return email
}

func GetFamilyId(token *jwt.Token) string {
	claims := token.Claims.(jwt.MapClaims)
	owner := claims["custom:familyId"].(string)
	return owner
}

func GetGroups(token *jwt.Token) []string {
	res := []string{}
	claims := token.Claims.(jwt.MapClaims)
	groups := claims["cognito:groups"].([]interface{})
	for _, g := range groups {
		gs := g.(string)
		res = append(res, gs)
	}
	return res
}

// validateAWSJwtClaims validates AWS Cognito User Pool JWT
func validateAWSJwtClaims(claims jwt.MapClaims) error {
	var err error
	// 3. Check the iss claim. It should match your user pool.
	issShoudBe := fmt.Sprintf("https://cognito-idp.%v.amazonaws.com/%v", region, userPoolID)
	err = validateClaimItem("iss", []string{issShoudBe}, claims)
	if err != nil {
		return err
	}
	// 4. Check the token_use claim.
	validateTokenUse := func() error {
		if tokenUse, ok := claims["token_use"]; ok {
			if tokenUseStr, ok := tokenUse.(string); ok {
				if tokenUseStr == "id" || tokenUseStr == "access" {
					return nil
				}
			}
		}
		return errors.New("Token_use should be id or access.")
	}
	err = validateTokenUse()
	if err != nil {
		return err
	}
	// 7. Check the exp claim and make sure the token is not expired.
	err = validateExpired(claims)
	if err != nil {
		return err
	}
	return nil
}

func validateClaimItem(key string, keyShouldBe []string, claims jwt.MapClaims) error {
	if val, ok := claims[key]; ok {
		if valStr, ok := val.(string); ok {
			for _, shouldbe := range keyShouldBe {
				if valStr == shouldbe {
					return nil
				}
			}
		}
	}
	return fmt.Errorf("%v does not match any of valid values: %v", key, keyShouldBe)
}

func validateExpired(claims jwt.MapClaims) error {
	if tokenExp, ok := claims["exp"]; ok {
		if exp, ok := tokenExp.(float64); ok {
			now := time.Now().Unix()
			fmt.Printf("current unixtime : %v\n", now)
			fmt.Printf("expire unixtime  : %v\n", int64(exp))
			if int64(exp) > now {
				return nil
			}
		}
		return errors.New("Can't parse token exp.")
	}
	return errors.New("This request has expired.")
}

// https://gist.github.com/MathieuMailhos/361f24316d2de29e8d41e808e0071b13
func convertKey(rawE, rawN string) *rsa.PublicKey {
	decodedE, err := base64.RawURLEncoding.DecodeString(rawE)
	if err != nil {
		panic(err)
	}
	if len(decodedE) < 4 {
		ndata := make([]byte, 4)
		copy(ndata[4-len(decodedE):], decodedE)
		decodedE = ndata
	}
	pubKey := &rsa.PublicKey{
		N: &big.Int{},
		E: int(binary.BigEndian.Uint32(decodedE[:])),
	}
	decodedN, err := base64.RawURLEncoding.DecodeString(rawN)
	if err != nil {
		panic(err)
	}
	pubKey.N.SetBytes(decodedN)
	// fmt.Println(decodedN)
	// fmt.Println(decodedE)
	// fmt.Printf("%#v\n", *pubKey)
	return pubKey
}

func getJWK(jwkURL string) map[string]JWKKey {
	jwk := &JWK{}
	getJSON(jwkURL, jwk)
	jm := make(map[string]JWKKey, 0)
	for _, jwk := range jwk.Keys {
		jm[jwk.Kid] = jwk
	}
	return jm
}

func getJSON(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
