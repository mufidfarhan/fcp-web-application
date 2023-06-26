package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("session_token")
		if err != nil {
			if ctx.GetHeader("Content-Type") == "application/json" {
				ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse("Unauthorized"))
			} else {
				ctx.Redirect(http.StatusSeeOther, "/client/login")
			}
			ctx.Abort()
			return
		}

		tokenStr := cookie.Value
		claims := &model.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})
		if err != nil {
			// Jika parsing token gagal
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
			ctx.Set("email", claims.Email)
		} else {
			ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse("Unauthorized"))
			ctx.Abort()
			return
		}

		ctx.Next()
		// TODO: answer here
	})
}
