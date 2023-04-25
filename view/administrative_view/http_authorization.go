package administrative_view

import (
	"davinci/domain/administrative_usecases/authorization"
	"davinci/domain/entities"
	"davinci/settings"
	"davinci/view"
	"davinci/view/http_error"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

type newHTTPAuthorizationModule struct {
	useCases authorization.UseCases
	settings settings.Settings
}

func NewHTTPAuthorization(settings settings.Settings, useCases authorization.UseCases) view.HttpModule {
	return &newHTTPAuthorizationModule{
		useCases: useCases,
		settings: settings,
	}
}

func (n newHTTPAuthorizationModule) Setup(router *mux.Router) {
	router.HandleFunc("/login", n.login).Methods(http.MethodPost)
}

func (n newHTTPAuthorizationModule) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[login] Error ReadAll", err)
		http_error.HandleError(w, err)
		return
	}

	var credentials entities.Credential
	if err = json.Unmarshal(b, &credentials); err != nil {
		log.Println("[login] Error Unmarshal", err)
		http_error.HandleError(w, err)
		return
	}

	user, err := n.useCases.Login(ctx, credentials)
	if err != nil {
		log.Println("[login] Error Login", err)
		http_error.HandleError(w, err)
		return
	}

	userClaim := &entities.UserCredential{
		Id:     user.Id,
		Email:  user.Credential.Email,
		RoleId: user.Credential.RoleId,
	}

	userClaimByte, err := json.Marshal(*userClaim)
	if err != nil {
		log.Println("[login] Error Marshal", err)
		http_error.HandleError(w, err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": string(userClaimByte),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("DAVINCI_SECRET_KEY")))
	if err != nil {
		log.Println("[login] Error SignedString", err)
		http_error.HandleError(w, err)
		return
	}

	_, err = w.Write([]byte(tokenString))
	if err != nil {
		log.Println("[login] Error Write", err)
		http_error.HandleError(w, err)
		return
	}
}
