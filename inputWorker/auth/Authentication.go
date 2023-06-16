package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var SecretKey = []byte("e-notificationQWE")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	authHeader := r.Header.Get("Authorization")
	// 	if authHeader == "" {
	// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 		return
	// 	}

	// 	tokenString := strings.Split(authHeader, " ")[1]
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	// Check the signing method
	// 	if token.Method != jwt.SigningMethodHS256 {
	// 		// return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	// 		return nil, nil
	// 	}
	// 	// Return the secret key for signature verification
	// 	return SecretKey, nil
	// })

	// if err != nil || !token.Valid {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	// // Pass the request to the next handler
	// next.ServeHTTP(w, r)
	// }
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the token from the cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		tokenString := cookie.Value

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if token.Method != jwt.SigningMethodHS256 {
				// return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				return nil, nil
			}
			// Return the secret key for signature verification
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		// JSON data
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		handleLogin(w, r, user)

	case "application/x-www-form-urlencoded":
		// Form data
		username := r.FormValue("username")
		password := r.FormValue("password")

		user := User{
			Username: username,
			Password: password,
		}

		handleLogin(w, r, user)

	default:
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request, user User) {
	// tur zuuriin username && password
	validUser := User{
		Username: "ohmynotif",
		Password: "qwe123!@#",
	}

	if user.Username == validUser.Username && user.Password == validUser.Password {

		// Create JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
			},
		})
		tokenString, err := token.SignedString(SecretKey)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		// create response
		// response := map[string]interface{}{
		// 	"message": "Login successful",
		// 	"status":  true,
		// 	"token":   tokenString,
		// }
		// jsonResponse, err := json.Marshal(response)
		// if err != nil {
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	return
		// }
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusOK)
		// w.Write(jsonResponse)

		// Set the token in a cookie
		cookie := http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  time.Now().Add(time.Hour * 24),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		// Redirect to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		response := map[string]interface{}{"message": "Invalid username or password", "status": false}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonResponse)
	}
}
