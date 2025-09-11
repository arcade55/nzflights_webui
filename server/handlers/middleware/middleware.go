package middleware

import (
	"context"
	"net/http"
)

// A custom type is used for the context key to avoid naming collisions
// with other packages.
type contextKey string

// UserIDKey is the key used to store the User ID in the request context.
// Exporting it allows other packages (like your handlers) to access it.
const UserIDKey = contextKey("userID")

// SessionStore is a placeholder for your actual session management logic.
// In a real application, this would be an interface that communicates
// with a database like Redis or a SQL table to validate session tokens.
type SessionStore interface {
	// GetUserIDFromSession takes a session token (UUID) and returns the
	// associated User ID. If the session is not valid, it returns an error.
	GetUserIDFromSession(sessionToken string) (string, error)
}

// MockSessionStore is a simple, hardcoded session store for demonstration.
// Replace this with your actual session store implementation.
type MockSessionStore struct {
	sessions map[string]string // Maps session UUID -> User ID
}

// GetUserIDFromSession implements the SessionStore interface for our mock store.
func (s *MockSessionStore) GetUserIDFromSession(sessionToken string) (string, error) {
	// In a real implementation, you would query your database here.
	userID, ok := s.sessions[sessionToken]
	if !ok {
		return "", http.ErrNoCookie // Using a standard error is fine here.
	}
	return userID, nil
}

// NewMockSessionStore creates a new mock store with some dummy data.
func NewMockSessionStore() *MockSessionStore {
	return &MockSessionStore{
		sessions: map[string]string{
			"a-valid-session-uuid-from-a-cookie": "user123",
			"another-valid-session-uuid":         "user456",
		},
	}
}

// Auth is the middleware handler. It takes a SessionStore to validate sessions.
func Auth(store SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get the session cookie from the request.
			cookie, err := r.Cookie("session_token")
			if err != nil {
				// If no cookie is present, the user is not authenticated.
				http.Error(w, "Unauthorized: Missing session token", http.StatusUnauthorized)
				return
			}

			// 2. Get the session UUID from the cookie's value.
			sessionUUID := cookie.Value

			// 3. Look up the User ID in the session store.
			userID, err := store.GetUserIDFromSession(sessionUUID)
			if err != nil {
				// If the session is not valid or expired, the user is not authorized.
				http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
				return
			}

			// 4. Add the validated User ID to the request's context. This is the
			//    key step that makes the User ID available to downstream handlers.
			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			// 5. Call the next handler in the chain, passing the new context.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
