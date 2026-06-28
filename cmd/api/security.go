package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type contextKey string

const userContextKey contextKey = "userID"

// AuthMiddleware intercepts requests to validate Bearer tokens
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token format")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		//TODO: Implement proper token parsing and validation using github.com/golang-jwt/jwt/v5

		// NOTE: In production, parse and validate your token cryptographically here
		// using github.com/golang-jwt/jwt/v5.
		userID := "user_12345" // Placeholder for an extracted claim

		//---
		// Generate a sample token for demonstration (typically done in an auth service)
		hmacSampleSecret := []byte("my_secret_key")
		// tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// 	"foo": "bar",
		// 	"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		// })
		// // Sign and get the complete encoded token as a string using the secret
		// tokenString, err := tkn.SignedString(hmacSampleSecret)
		println("Received token:", tokenString) // Debug log for received token

		// Parse takes the token string and a function for looking up the key. The latter is especially
		// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
		// head of the token to identify which key to use, but the parsed token (head and claims) is provided
		// to the callback, providing flexibility.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return hmacSampleSecret, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			fmt.Println("Error parsing token:", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			fmt.Println(claims["foo"], claims["nbf"])
		} else {
			fmt.Println(err)
		}
		//---

		// Inject identity metadata securely into request context
		ctx := context.WithValue(c.Request().Context(), userContextKey, userID)
		// next.ServeHTTP(w, r.WithContext(ctx))
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

// SecureHeaderMiddleware adds baseline web defense properties
func SecureHeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("X-Frame-Options", "DENY")
		// c.Response().Header().Set("Content-Type", "application/json")
		return next(c)
	}
}
