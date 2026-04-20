package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

const SessionName = "clipper-session"

var store *sessions.CookieStore

func InitSessionStore(secret string) {
	store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, SessionName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session error"})
			return
		}

		userID, ok := session.Values["user_id"].(string)
		if !ok || userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func SetUserID(c *gin.Context, userID string) error {
	session, err := store.Get(c.Request, SessionName)
	if err != nil {
		return err
	}

	session.Values["user_id"] = userID
	return session.Save(c.Request, c.Writer)
}

func ClearSession(c *gin.Context) error {
	session, err := store.Get(c.Request, SessionName)
	if err != nil {
		return err
	}

	session.Values["user_id"] = ""
	session.Options.MaxAge = -1
	return session.Save(c.Request, c.Writer)
}

func GetUserID(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return userID.(string)
}