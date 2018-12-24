package films

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	oktaUtils "github.com/okta/samples-golang/okta-hosted-login/utils"
)

func FilmAnonymousRegister(router *gin.RouterGroup) {
	router.GET("/", FilmList)
}


func FilmList(c *gin.Context) {
	filmModel, err := FindOneFilm()
	if err != nil {
		c.JSON(http.StatusNotFound, errors.New("invalid param"))
		return
	}
	serializer := FilmSerializer{c, filmModel}
	c.JSON(http.StatusOK, gin.H{"films": serializer.Response()})
}

func LoginHandler(c *gin.Context) {
	nonce, _ := oktaUtils.GenerateNonce()

	q := c.Request.URL.Query()
	q.Add("client_id", os.Getenv("CLIENT_ID"))
	q.Add("response_type", "token")
	q.Add("scope", "openid")
	q.Add("redirect_uri", "http://localhost:8080/implicit/callback")
	q.Add("state", "ApplicationState")
	q.Add("nonce", nonce)

	redirectPath := os.Getenv("ISSUER") + "/v1/authorize?" + q.Encode()
	c.Redirect(http.StatusMovedPermanently, redirectPath)
}

func AuthCodeCallbackHandler(c *gin.Context) {
	return
}