package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/madappgang/tayblebackend/controllers"
	"bitbucket.org/madappgang/tayblebackend/controllers/helpers"
	"bitbucket.org/madappgang/tayblebackend/models"
	"bitbucket.org/madappgang/tayblebackend/models/validator"
	"bitbucket.org/madappgang/tayblebackend/services"
	"github.com/astaxie/beego"
)

var (
	// stageServer check we are working stage server
	allowFakeNumbers = os.Getenv("ALLOW_FAKE_PHONE_NUMBERS") == "TRUE"
	// testPhoneRegex is regexpression that help identify test mobile number
	testPhoneRegex = regexp.MustCompile(`^[+][0-9]{9,15}[+][0-9]{1,6}[+]$`)
)

const (
	phoneVerificationCodeLength         = 6
	invalidVerificationCodeErrorMessage = "Sorry, the code you entered is invalid or has expired. Please get a new one."
	supportAnonymousLogin               = false
)

// User to authenticate require to have valid phone number
// to verify the service is valid, we are sending message to the user
func (ar *Router) RequestVerificationCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authData helpers.PhoneLogin
		if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, err.Error(), "RequestVerificationCode.Unmarshal")
			return
		}

		if !allowFakeNumbers || !isTestPhone(authData.PhoneNumber) {
			if err := authData.IsValidPhone(); err != nil {
				ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, err.Error(), "RequestVerificationCode.IsValidPhone")
				return
			}
		}

		code := randStringBytes(phoneVerificationCodeLength)
		err := models.MainVerificationCodeDatastore.CreateVerificationCode(authData.PhoneNumber, code)
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "RequestVerificationCode.CreateVerificationCode")
			beego.Error("Unable to create verification code. Error: " + err.Error())
			ac.CustomAbort(http.StatusInternalServerError, "Unable to create verification code. Please try again. ")
		}
		sms := ""
		if ac.IsAndroidClient() {
			sms = fmt.Sprintf(services.SMSAndroidVerificationCode, code)
		} else {
			sms = fmt.Sprintf(services.SMSVerificationCode, code)
		}

		if allowFakeNumbers && isTestPhone(authData.PhoneNumber) {
			authData.PhoneNumber = parseRealPhoneFromTestPhone(authData.PhoneNumber)
		}

		err = services.MainSMSService.SendSMS(authData.PhoneNumber, sms)
		if err != nil {
			beego.Error("Unable to send SMS. Error: " + err.Error())
			ac.CustomAbort(http.StatusInternalServerError, "Unable to send SMS. Please try again.")
		}

		ac.ServeJSON()
	}
}

//PhoneLogin authenticates user with phone number and verification code
//If user exists - create new session and return token
//If user does not exists - register and then login (create session and return token)
//If code is invalid - return error
func (ar *Router) PhoneLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authData helpers.PhoneLogin
		err := json.Unmarshal(ac.Ctx.Input.RequestBody, &authData)
		if err != nil {
			ac.CustomAbort(http.StatusBadRequest, "Unable to parse data. Error: "+err.Error())
		}
		if allowFakeNumbers && isTestPhone(authData.PhoneNumber) {
			// do nothing - do not validate phone number
		} else if err := authData.IsValidCodeAndPhone(); err != nil {
			ac.CustomAbort(http.StatusBadRequest, err.Error())
		}

		// check verification code
		if !models.MainVerificationCodeDatastore.IsVerificationCodeExists(authData.PhoneNumber, authData.Code) {
			ac.CustomAbort(http.StatusBadRequest, invalidVerificationCodeErrorMessage)
		}

		// if requester sent deviceId - we need to verify what to do with it.
		if authData.DeviceId != "" {
			// TODO: hope soon we will not use device_id
			ac.mergeUserWithPhoneAndDeviceID(authData.PhoneNumber, authData.DeviceId)
		}

		user := ac.findOrCreateUserWithPhone(authData.PhoneNumber)

		token, err := ac.loginUser(*user)
		if err != nil {
			ac.CustomAbort(http.StatusInternalServerError, "Unable to login user. Please try again.")
		}

		ac.Data["json"] = struct {
			Token string      `json:"token"`
			User  models.User `json:"user"`
		}{
			Token: token,
			User:  *user,
		}
		ac.ServeJSON()
	}
}

// login existent user: generate new token and create session
func (ac *AuthUserController) loginUser(user models.User) (string, error) {
	token, err := helpers.RandomString(controllers.TokenStringLenBytes)
	if err != nil {
		return "", err
	}
	sessionExp := time.Now().Add(controllers.SessionExpireTime)
	session := models.NewSession(user.ID, models.UserModel, token).SetExpirationTime(sessionExp)
	session, err = models.MainSessionDatastore.CreateSession(*session)
	if err != nil {
		return "", err
	}

	return token, nil

}

// Generate user code
func randStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		code := "1234567890"
		b[i] = code[rand.Intn(len(code))]
	}
	return string(b)
}

// isTestPhone used to check if client send test phone number in special format
func isTestPhone(phone string) bool {
	return testPhoneRegex.MatchString(phone)
}

// parseRealPhoneFromTestPhone looking for a real phone number in written in special format
func parseRealPhoneFromTestPhone(phone string) string {
	phone = phone[0 : len(phone)-1]
	secondPlusIndex := strings.LastIndex(phone, "+")
	phone = phone[:secondPlusIndex]
	return phone
}

//PhoneLogin is used to parse input data from the client during phone login.
type PhoneLogin struct {
	PhoneNumber         string `json:"phone_number"`
	FacebookAccessToken string `json:"facebook_access_token"`
	Code                string `json:"code"`
}

func (l *PhoneLogin) IsValidCodeAndPhone() error {
	if len(l.Code) == 0 {
		return errors.New("Verificataion code is too short or missing. ")
	}
	return l.IsValidPhone()
}

func (l *PhoneLogin) IsValidPhone() error {
	return validator.ValidatePhone(l.PhoneNumber)
}
