package server

import (
	"fmt"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func generateToken(userID uuid.UUID, t int) (string, error) {
	u := database.User{Id: userID}
	lineId := u.GetLineID()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * time.Duration(t)).Unix(), // token 有效期為 t 小時
		"lineID": lineId,
	})

	tokenString, err := token.SignedString([]byte(setting.JWTsecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getJwtClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("sign error: %v", token.Header["alg"])
		}
		return []byte(setting.JWTsecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func verifyToken(tokenString string) (string, string, error) {
	claims, err := getJwtClaims(tokenString)
	if err != nil {
		return "", "", err
	}

	userID := claims["userID"].(string)
	lineID := claims["lineID"].(string)
	return userID, lineID, nil
}
