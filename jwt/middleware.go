package jwt

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-study/lib/httprouter"
	"log"
	"net/http"
	"time"
)

const (
	SecretKey = "felix jwt demo"
)
type Token struct {
	Token string `json:"token"`
}
func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//生成token
func GenerateToken(username string,pwd string) (Token,error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["username"] = username
	claims["pwd"] = pwd
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		fmt.Println("Error while signing the token")
		fatal(err)
	}

	response := Token{tokenString}

	return response,err
}

/**
 * 解析 token
 */
func ParseToken(tokenSrt string) (token *jwt.Token, err error) {
	//var token *jwt.Token
	token, err = jwt.Parse(tokenSrt, func(*jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	//claims = token.Claims
	//fmt.Println(token.Valid)
	//fmt.Println(claims)
	//fmt.Println(err)
	return
}


//JWT 鉴权中间件
func ValidateTokenMiddleware()func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				params := httprouter.Params{}
				login(w, r, params)
			}else {
				token, err := ParseToken(r.Header.Get("Authorization"))
				if err == nil {
					if token.Valid{
						next.ServeHTTP(w, r)
					} else {
						w.WriteHeader(http.StatusUnauthorized)
						fmt.Fprint(w, "Token is not valid")
					}
				} else {
					//w.WriteHeader(http.StatusUnauthorized)
					//fmt.Fprint(w, "Unauthorized access to this resource")

					//假设永远验证成功
					next.ServeHTTP(w, r)

				}
			}
		}

		return http.HandlerFunc(fn)
	}
}


func ValidateTokenMiddleware2(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := ParseToken(r.Header.Get("Authorization"))

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}
}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//*
// 为了演示JWT 中间件 与请求Token
// 使用GET方法 任何参数都会返回相同token
// *//
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handle  httprouter.Handle
}

type Routes []Route

func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	response := map[string]string{}
	response["username"] = ps.ByName("username")
	Token,_ := GenerateToken(ps.ByName("username"),ps.ByName("pwd"))
	response["token"] = Token.Token

	JsonResponse(response,w)
}

func NewRouter() *httprouter.Router {

	router := httprouter.New()
	route :=  Route{
		"login",
		"GET",
		"/login",
		login,
	}
	router.Handle(route.Method, route.Pattern, route.Handle)
	return router
}