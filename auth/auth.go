package auth

import (
	"encoding/json"
	"github.com/MichalPolinkiewicz/to-do-api/db"
	"github.com/MichalPolinkiewicz/to-do-api/models"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var JwtKey = []byte("my_secret_key")
var LoggedOutUsers []string

// struct that will be encoded to a JWT. We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateAccount(res http.ResponseWriter, req *http.Request) {
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)

	//if some params missing return http 400
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if db.CheckIfUserExistsInDb(&user.Username) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	db.SaveUser(&user)
}

func Login(res http.ResponseWriter, req *http.Request) {
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)

	//if some params missing return http 400
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// check that user account exists in 'db' and password is valid - user can log in
	userFromDb := db.GetUserFromDb(&user.Username, &user.Password)

	if userFromDb.Id == 0 {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	//generate and set token cookie
	tokenString, err := createToken(userFromDb.Id, user.Username, "token")

	// if there is an error in creating the JWT return an internal server error
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	//no error - set cookie
	http.SetCookie(res, &http.Cookie{
		Name:  "token",
		Value: tokenString,
	})

	//generate and set token refresh cookie
	tokenString, err = createToken(userFromDb.Id, user.Username, "refresh")

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:  "refresh",
		Value: tokenString,
	})
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

		//if token is valid - return response
		if isValidToken(&t) {
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
			//if is valid - create new token
			if isValidToken(&t) {
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
				t, err := createToken(claims.UserId, claims.Username, "token")

				if err != nil {
					res.WriteHeader(http.StatusUnauthorized)
					return
				}

				c, e := getCookie(req, "token")

				if e != 0 {
					res.WriteHeader(http.StatusUnauthorized)
					return
				}

				c.Value = t
				c.Name = "token"
				h.ServeHTTP(res, req)
			} else {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	})
}

func Logout(res http.ResponseWriter, req *http.Request) {
	c, err := getCookie(req, "refresh")

	if err != 0 {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	if isValidToken(&c.Value) {
		for _, l := range LoggedOutUsers {
			if l == c.Value {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		LoggedOutUsers = append(LoggedOutUsers, c.Value)

		c, err := getCookie(req, "token")
		if err != 0 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		LoggedOutUsers = append(LoggedOutUsers, c.Value)
		return
	} else {
		res.WriteHeader(http.StatusUnauthorized)
	}
}

func GetUserIdFromRequest(req *http.Request) int {
	c, _ := getCookie(req, "refresh")

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(c.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		return 0
	}

	return claims.UserId
}

func isValidToken(t *string) bool {
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`. Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in), or if the signature does not match
	tkn, err := jwt.ParseWithClaims(*t, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		return false
	}

	return tkn.Valid
}

func createToken(id int, u string, t string) (string, error) {
	// Declare the expiration time of the token, depending on type
	var expirationTime time.Time

	if t == "refresh" {
		expirationTime = time.Now().Add(12 * time.Hour)
	} else {
		expirationTime = time.Now().Add(30 * time.Second)
	}

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		UserId:   id,
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
