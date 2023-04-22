package infrastructure

import (
	"base/view/administrative_view"
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"time"
)

func Setup(router *mux.Router) error {
	db, err := setupDataBase()
	if err != nil {
		log.Println("[Setup] Error setupDataBase", err)
		return err
	}

	err = setupModules(router, db)
	if err != nil {
		log.Println("[Setup] Error setupModules", err)
		return err
	}

	return nil
}

// SetupDataBase set the connection to the database and set connection settings.
func setupDataBase() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(devserver:3306)/davinci?parseTime=true")
	if err != nil {
		log.Println("[Setup] Error connecting to database", err)
		return nil, err
	}

	// Limit the amount of time the connections are kept in the pool
	db.SetConnMaxLifetime(time.Minute * 10)

	// Limit the number of connections stored in the pool
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)

	return db, nil
}

// setupModules set the MVC structure for the application.
func setupModules(router *mux.Router, db *sql.DB) error {
	router.Use(rootMiddleware)

	err := SetupAdministrativeModules(router, db)
	if err != nil {
		log.Println("[SetupMVC] Error SetupAdministrativeModules", err)
		return err
	}
	err = SetupDeviceModules(router, db)
	if err != nil {
		log.Println("[SetupMVC] Error SetupDeviceModules", err)
		return err
	}

	return nil
}

// rootMiddleware set the response content type for the api as json.
func rootMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Set the response content type for the api as json
		w.Header().Set("Content-Type", "application/json")

		//Set the origin to allow all.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		//Set the valid methods to all.
		w.Header().Set("Access-Control-Allow-Methods", "*")

		//Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// getTokenFromCookie get the token from the cookie.
func getTokenFromCookie(cookie *http.Cookie) (*jwt.Token, error) {
	secureCookie := securecookie.New([]byte(administrative_view.SecretJWTKey), nil)
	var tokenString string
	err := secureCookie.Decode("token", cookie.Value, &tokenString)
	if err != nil {
		log.Println("[login] Error Decode", err)
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Println("[login] token.Method.(*jwt.SigningMethodHMAC) !ok", err)
			return nil, errors.New("error parsing token")
		}
		return []byte(administrative_view.SecretJWTKey), nil
	})
	if err != nil {
		log.Println("[isCookieValid] Error parsing token", err)
		return nil, err
	}

	return token, nil
}
