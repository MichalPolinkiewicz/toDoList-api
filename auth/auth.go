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
	{"User22", "Pass", "", false},
}

var LoggedOutUsers []string

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT. We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(res http.ResponseWriter, req *http.Request) {
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

	//generate and set token cookie
	tokenString, err := createToken(creds.Username, "token")

	// If there is an error in creating the JWT return an internal server error
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:  "token",
		Value: tokenString,
	})

	//generate and set token refresh cookie
	tokenString, err = createToken(creds.Username, "refresh")

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:  "refresh",
		Value: tokenString,
	})
}

func Logout(res http.ResponseWriter, req *http.Request) {

	c, err := getCookie(req, "token")
	if err != 0 || c == nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	isValid, err := isValidToken(&c.Value)
	if err != 0 {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	if isValid {
		for _, l := range LoggedOutUsers {
			if l == c.Value {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		LoggedOutUsers = append(LoggedOutUsers, c.Value)
		return
	} else {
		res.WriteHeader(http.StatusUnauthorized)
	}
}

func CheckJwtToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// get token cookie
		c, err := getCookie(req, "token")

		if err != 0 {
			if err == 401 {
				// If the cookie is not set, return an unauthorized status
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get the JWT string from the cookie
		t := c.Value

		//check if user don't log out
		for _, it := range LoggedOutUsers {
			if it == t {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		isValid, err := isValidToken(&t)

		//if token is valid - return response
		if isValid && err == 0 {
			h.ServeHTTP(res, req)
		} else {
			//token is expired - check for refresh token
			c, err := getCookie(req, "refresh")

			if err != 0 {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}

			t := c.Value

			//check if refresh token is valid
			isValid, err := isValidToken(&t)

			if err != 0 {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}

			//if is valid - create new token and new refresh token
			if isValid {
				//get username
				claims := &Claims{}
				_, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
					return JwtKey, nil
				})

				if err != nil {
					res.WriteHeader(http.StatusUnauthorized)
					return
				}

				//creating new token
				t, err := createToken(claims.Username, "token")

				if err != nil {
					res.WriteHeader(http.StatusUnauthorized)
					return
				}

				//set new token cookie
				http.SetCookie(res, &http.Cookie{
					Name:  "token",
					Value: t,
				})
				h.ServeHTTP(res, req)
			} else {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	})
}

func isValidToken(t *string) (bool, int) {
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`. Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in), or if the signature does not match
	tkn, err := jwt.ParseWithClaims(*t, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	i := 0

	if err != nil {
		i = http.StatusUnauthorized
	}

	return tkn.Valid, i
}

func createToken(u string, t string) (string, error) {
	// Declare the expiration time of the token, depending on type
	var expirationTime time.Time

	if t == "refresh" {
		expirationTime = time.Now().Add(12 * time.Hour)
	} else {
		expirationTime = time.Now().Add(30 * time.Second)
	}

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: u,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(JwtKey)

	return tokenString, err
}

func getCookie(req *http.Request, cn string) (*http.Cookie, int) {
	c, err := req.Cookie(cn)
	ec := 0

	if err != nil {
		ec = http.StatusUnauthorized
		return nil, ec
	}
	return c, ec
}

func IsValidUser(c *Credentials) bool {
	if userExists(&c.Username) && isCorrectPassword(&c.Username, &c.Password) {
		return true
	}
	return false
}

func userExists(l *string) bool {
	for _, e := range ValidUsers {
		if e.Login == *l {
			return true
		}
	}
	return false
}

func isCorrectPassword(l *string, p *string) bool {
	for _, e := range ValidUsers {
		if e.Login == *l && e.Password == *p {
			return true
		}
	}
	return false
}
