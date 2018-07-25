package auth

import (
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"golang-chat/src/models/config"
	"log"
	"net/http"
	"strings"
)

const (
	CONFPATH = "../keys/auth.json"
	CALLBACK = "http://localhost:8080/auth/callback/google"
)

type authHandler struct {
	next http.Handler
}

func (h authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// URL Path -> /auth/{action}/{provider}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	pathes := strings.Split(r.URL.Path, "/")

	if len(pathes) != 4 {
		return
	}

	action := pathes[2]
	providerName := pathes[3]

	switch action {
	case "login":
		provider, err := gomniauth.Provider(providerName)
		if err != nil {
			log.Fatalln("failed getting auth providerName", provider, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("err: calling GrtBeginAuthURL", provider, "-", err)
		}
		w.Header().Set("location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":
		provider, err := gomniauth.Provider(providerName)
		if err != nil {
			log.Fatalln("err: calling GrtBeginAuthURL", provider, "-", err)
		}
		cred, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("failed complete auth", provider, "-", err)
		}
		user, err := provider.GetUser(cred)
		if err != nil {
			log.Fatalln("failed getting user infomation", provider, "-", err)
		}

		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header()["location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

func SetupProvider() error {
	conf, err := config.Perse(CONFPATH)
	if err != nil {
		return err
	}

	gomniauth.SetSecurityKey(conf.GOMNIAUTH.SECURITYKEY)
	gomniauth.WithProviders(
		google.New(conf.GOOGLE.CLIENTID, conf.GOOGLE.CLIENTSECRET, CALLBACK),
	)
	return nil
}
