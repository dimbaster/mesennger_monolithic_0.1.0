package tokens

import (
	"Server/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretkey = os.Getenv("TOKEN_KEY")

func CreateToken(id int, login string, password string) (string, error) {
	sub, err := json.Marshal(models.User{
		ID:       id,
		Login:    login,
		Password: password,
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
		"sub": string(sub),
	})

	tokenString, err := token.SignedString([]byte(secretkey))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretkey), nil
	})

	if !token.Valid {
		return nil, err
	}
	return token, nil
}

func VerifyAndParse(r *http.Request) (models.User, error) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := ValidateToken(tokenString)
	if err != nil {
		return models.User{}, err
	}

	var res models.User
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return models.User{}, err
	}

	err = json.Unmarshal([]byte(sub), &res)
	if err != nil {
		return models.User{}, err
	}

	return res, nil
}
