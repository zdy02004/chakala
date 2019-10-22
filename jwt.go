package main

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"log"
)

// request handler
func jwt_request( user_name string, SigningKey string, expire_times int64 ) (token_string string, err error) {

	//mySigningKey := []byte("AllYourBase")
	mySigningKey := []byte(SigningKey)

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: expire_times,
		Issuer:    user_name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", nil
	}

	log.Printf("%v %v", ss, err)
	return ss, nil
}

func at(t time.Time, f func()( bool, error) ) (valid bool, err error) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	token, err :=f()
	jwt.TimeFunc = time.Now
	return token, err
}

func jwt_parse(tokenString string, SigningKey string) (valid bool, err error) {
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}

	// sample token is expired.  override time so it parses as valid
	return at(time.Unix(0, 0), func()( bool, error)   {
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SigningKey), nil
		})

		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			log.Printf("%v %v", claims.Foo, claims.StandardClaims.ExpiresAt)
			return true, nil
		} else {
			log.Println("jwt_parse failed because ", err)
			return false, err
		}
	})

}
