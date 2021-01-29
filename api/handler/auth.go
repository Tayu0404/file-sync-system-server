package handler

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	cache "github.com/patrickmn/go-cache"

	"github.com/Tayu0404/file-sync-system-server/api/model"
)

type JWTClaims struct {
	Username      string `json:"username"`
	Locale        string `json:"locale"`
	Avater        string `json:"avater"`
	ID            uint64 `json:"id,string"`
	jwt.StandardClaims
}

// Resposes
type LoginSuccess struct {
	Username      string `json:"Username"`
	Locale        string `json:"locale"`
	Avater        string `json:"avater"`
	ID            uint64 `json:"id,string"`
	JWT           string `json:"jwt"`
}

type TOTPTuthRequired struct {
	Token         string `json:"Token"`
}

func (h *handler) Signup(c echo.Context) error {
	user := model.User{}
	user.UserID         = uint64(h.Node.Generate().Int64())
	user.Email          = c.FormValue("email")
	user.Locale         = c.FormValue("locale")
	user.TwoFAType      = model.TwoFATypeNotUsed
	user.TotpSecret     = ""
	user.Role, _        = strconv.Atoi(c.FormValue("role"))
	user.Profile.Name   = c.FormValue("name")
	//user.Profile.Avatar = c.FormValue("Avatar")

	passwordHash, err := userPassHash(c.FormValue("password"))
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to create user")
	}
	user.Password = passwordHash
	if err = h.Model.SignupUser(&user); err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to create user")
	}
	return c.String(http.StatusOK, "Signed up!!")
}

func (h *handler) Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := h.Model.GetDetailForUserLogin(email)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to Login")
	}

	if !userPassMach(user.Password, password) {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to Login")
	}

	switch user.TwoFAType {
	case model.TwoFATypeNotUsed:
		if err := h.Model.LoginUser(user); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to Login")
		}
	
		token, err := generateJWT(user)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to Login")
		}

		res := LoginSuccess{
			user.Profile.Name,
			user.Locale,
			user.Profile.Avatar,
			user.UserID,
			token,
		}
		return c.JSON(http.StatusOK, res)

	case model.TwoFATypeTOTP:
		token := GenerateTwoFAToken()
		h.Cache.Set(token, user, cache.DefaultExpiration)

		res := TOTPTuthRequired{
			token,
		}
		return c.JSON(http.StatusOK, res) 
	}
	return c.String(http.StatusInternalServerError, "Failed to Login")
}

func (h *handler) TOTPAuth(c echo.Context) error {
	token := c.FormValue("token")
	totpcode := c.FormValue("code")
	user, found := h.Cache.Get(token)
	if found {
		return c.String(http.StatusInternalServerError, "Failed to Login")
	}

	if match := validateTOTP(totpcode, user.(model.User).TotpSecret); !match {
		return c.String(http.StatusInternalServerError, "Failed to Login")
	}

	token, err := generateJWT(user.(*model.User))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to Login")
	}

	res := LoginSuccess{
		user.(model.User).Profile.Name,
		user.(model.User).Locale,
		user.(model.User).Profile.Avatar,
		user.(model.User).UserID,
		token,
	}
	return c.JSON(http.StatusOK, res)
}

func userPassHash(pass string)(string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func userPassMach(hash,pw string)bool{
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

func generateJWT(user *model.User) (token string, err error) {
	claims := &JWTClaims{
		user.Profile.Name,
		user.Locale,
		user.Profile.Avatar,
		user.UserID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func generateTOTPKey(user *model.User) (key *otp.Key, err error) {
	opts := totp.GenerateOpts{}
	opts.Issuer = "fss"
	opts.AccountName = strconv.FormatUint(user.UserID, 10)
	opts.Algorithm = otp.AlgorithmSHA256

	key, err = totp.Generate(opts)
	return
}

func generateTOTPQR(key *otp.Key) (qr string, err error ) {
	width := 256
	height := 256

	qrimg, err := key.Image(width, height)
	if err != nil {
		return
	}

	imgbuf := new(bytes.Buffer)
	if err = png.Encode(imgbuf, qrimg); err != nil {
		return
	}
	
	qr = base64.StdEncoding.EncodeToString(imgbuf.Bytes())
	return
}

func validateTOTP(code string, secret string) (match bool) {	
	match = totp.Validate(code, secret)
	return
}

func GenerateTwoFAToken() string {
	const n = 256
	var randSrc = rand.NewSource(time.Now().UnixNano())

	const (
		rs6Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		rs6LetterIdxBits = 6
		rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
		rs6LetterIdxMax = 63 / rs6LetterIdxBits
	)

    token := make([]byte, n)
    cache, remain := randSrc.Int63(), rs6LetterIdxMax
    for i := n-1; i >= 0; {
        if remain == 0 {
            cache, remain = randSrc.Int63(), rs6LetterIdxMax
        }
        idx := int(cache & rs6LetterIdxMask)
        if idx < len(rs6Letters) {
            token[i] = rs6Letters[idx]
            i--
        }
        cache >>= rs6LetterIdxBits
        remain--
    }
    return string(token)
}