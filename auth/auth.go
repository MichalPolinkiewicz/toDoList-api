package auth

import (
	"encoding/json"
	"github.com/MichalPolinkiewicz/to-do-api/models"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var JwtKey = []byte("my_secret_key")

var ValidUsers = []models.User{
	{"User", "Pass", "", false},
}

var LoggedUsers []models.User

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func LogIn(res http.ResponseWriter, req *http.Request) {
	var creds Credentials

	// Get the JSON body and decode into credentials
	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// check that user account exists in 'db' and password is valid - user can log in
	if !IsValidUser(&creds) {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(JwtKey)

	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(res, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func LogOut(res http.ResponseWriter, req *http.Request){
	//pobrac z req ciasteczko, z niego wartosc tokena, dodac go do blacklisty
}

func CheckJwtToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// We can obtain the session token from the requests cookies, which come with every request
		c, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get the JWT string from the cookie
		tknStr := c.Value

		// Initialize a new instance of `Claims`
		claims := &Claims{}

		// Parse the JWT string and store the result in `claims`. Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in), or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if !tkn.Valid {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		//refreshing token
		expirationTime := time.Now().Add(5 * time.Minute)
		claims.ExpiresAt = expirationTime.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(JwtKey)

		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(res, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		h.ServeHTTP(res, req)
	})
}

func IsValidUser(c *Credentials) bool {
	if len(c.Username) > 0 && len(c.Password) > 0 {
		for _, u := range ValidUsers {
			if c.Username == u.Login && c.Password == u.Password {
				return true
			}
		}
	}
	return false
}
