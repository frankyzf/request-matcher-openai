package myauth

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func (p *MyAuth) ProduceToken(claims map[string]string) (string, error) {
	mm := map[string]interface{}{}
	for key, value := range claims {
		mm[key] = value
	}
	return CommonProduceToken(mm, p.Secret)
}

func ProduceRS256Token(claims map[string]interface{}, privateKey []byte) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to load private key: %w", err)
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claims)).SignedString(key)
	return token, err
}

func CommonProduceToken(claims map[string]interface{}, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

func (p *MyAuth) JustParseToken(tokenString string) (map[string]interface{}, error) {
	mylogger.Infof("just parse token: %v and secret: %v", tokenString, p.Secret)
	return CommonParseUnverified(tokenString)
}

func (p *MyAuth) ParseToken(tokenString string) (map[string]interface{}, error) {
	mylogger.Infof("parse token: %v and secret: %v", tokenString, p.Secret)
	return CommonParseToken(tokenString, p.Secret)
}

func CommonParseToken(tokenString string, secret string) (map[string]interface{}, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		return map[string]interface{}{}, err
	}
	if token == nil {
		return map[string]interface{}{}, errors.New("token parse is nil")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return map[string]interface{}{}, err
	}
}

func CommonParseUnverified(tokenString string) (map[string]interface{}, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return map[string]interface{}{}, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return map[string]interface{}{}, errors.New("failed to parse unverified token")
	}
}
