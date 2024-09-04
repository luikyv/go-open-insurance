package authn

import (
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"strings"

	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-opf/internal/consent"
	"github.com/luikyv/go-opf/internal/oidc"
	"github.com/luikyv/go-opf/internal/sec"
	"github.com/luikyv/go-opf/internal/time"
	"github.com/luikyv/go-opf/internal/user"
)

type Authenticator struct {
	userService    user.Service
	consentService consent.Service
	baseURL        string
}

func New(
	userService user.Service,
	consentService consent.Service,
	baseURL string,
) Authenticator {
	return Authenticator{
		userService:    userService,
		consentService: consentService,
		baseURL:        baseURL,
	}
}

func (a Authenticator) Authenticate(
	w http.ResponseWriter,
	r *http.Request,
	session *goidc.AuthnSession,
) goidc.AuthnStatus {
	if _, ok := session.Store[paramStepID]; !ok {
		session.StoreParameter(paramStepID, stepIDSetUp)
	}

	if session.Parameter(paramStepID) == stepIDSetUp {
		if status := a.setUp(r, session); status != goidc.StatusSuccess {
			return status
		}
		session.StoreParameter(paramStepID, stepIDLogin)
	}

	if session.Parameter(paramStepID) == stepIDLogin {
		if status := a.login(w, r, session); status != goidc.StatusSuccess {
			return status
		}
		session.StoreParameter(paramStepID, stepIDConsent)
	}

	if session.Parameter(paramStepID) == stepIDConsent {
		if status := a.grantConsent(w, r, session); status != goidc.StatusSuccess {
			return status
		}
		session.StoreParameter(paramStepID, stepIDFinishFlow)
	}

	if session.Parameter(paramStepID) == stepIDFinishFlow {
		return a.finishFlow(session)
	}

	return goidc.StatusFailure
}

func (a Authenticator) setUp(
	r *http.Request,
	session *goidc.AuthnSession,
) goidc.AuthnStatus {
	consentID := consentID(session.Scopes)
	if consentID == "" {
		session.SetError("missing consent ID")
		return goidc.StatusFailure
	}

	consent, err := a.consentService.Get(
		r.Context(),
		consentID,
		sec.Meta{
			ClientID: session.ClientID,
		},
	)
	if err != nil {
		session.SetError(err.Error())
		return goidc.StatusFailure
	}

	session.StoreParameter(paramConsentID, consent.ID)
	permissions := ""
	for _, p := range consent.Permissions {
		permissions += fmt.Sprintf("%s ", p)
	}
	permissions = permissions[:len(permissions)-1]
	session.StoreParameter(paramPermissions, permissions)
	session.StoreParameter(paramConsentCPF, consent.UserCPF)
	return goidc.StatusSuccess
}

func (a Authenticator) login(
	w http.ResponseWriter,
	r *http.Request,
	session *goidc.AuthnSession,
) goidc.AuthnStatus {

	r.ParseForm()

	username := r.PostFormValue(usernameFormParam)
	if username == "" {
		w.WriteHeader(http.StatusOK)
		tmpl, _ := template.ParseFiles("../../templates/login.html")
		tmpl.Execute(w, AuthnPage{
			BaseURL:    a.baseURL,
			CallbackID: session.CallbackID,
		})
		return goidc.StatusInProgress
	}

	user, err := a.userService.User(username)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		// TODO: relative path.
		tmpl, _ := template.ParseFiles("../../templates/login.html")
		tmpl.Execute(w, AuthnPage{
			BaseURL:    a.baseURL,
			CallbackID: session.CallbackID,
			Error:      "invalid username",
		})
		return goidc.StatusInProgress
	}

	password := r.PostFormValue(passwordFormParam)
	if user.CPF != session.Parameter(paramConsentCPF) || password != correctPassword {
		w.WriteHeader(http.StatusOK)
		tmpl, _ := template.ParseFiles("../../templates/login.html")
		tmpl.Execute(w, AuthnPage{
			BaseURL:    a.baseURL,
			CallbackID: session.CallbackID,
			Error:      "invalid credentials",
		})
		return goidc.StatusInProgress
	}

	session.StoreParameter(paramUserID, username)
	return goidc.StatusSuccess
}

func (a Authenticator) grantConsent(
	w http.ResponseWriter,
	r *http.Request,
	session *goidc.AuthnSession,
) goidc.AuthnStatus {

	r.ParseForm()

	var permissions []consent.Permission
	for _, p := range strings.Split(session.Parameter(paramPermissions).(string), " ") {
		permissions = append(permissions, consent.Permission(p))
	}
	isConsented := r.PostFormValue(consentFormParam)
	if isConsented == "" {
		w.WriteHeader(http.StatusOK)
		tmpl, _ := template.ParseFiles("../../templates/consent.html")
		tmpl.Execute(w, AuthnPage{
			BaseURL:     a.baseURL,
			CallbackID:  session.CallbackID,
			Permissions: permissions,
		})
		return goidc.StatusInProgress
	}

	consentID := session.Parameter(paramConsentID).(string)

	if isConsented != "true" {
		a.consentService.Reject(
			r.Context(),
			consentID,
			sec.Meta{
				ClientID: session.ClientID,
			},
		)
		session.SetError("consent not granted")
		return goidc.StatusFailure
	}

	if err := a.consentService.Authorize(r.Context(), consentID, permissions...); err != nil {
		session.SetError(err.Error())
		return goidc.StatusFailure
	}
	return goidc.StatusSuccess
}

func (a Authenticator) finishFlow(
	session *goidc.AuthnSession,
) goidc.AuthnStatus {
	session.SetUserID(session.Parameter(paramUserID).(string))
	session.GrantScopes(session.Scopes)
	session.SetIDTokenClaimACR(oidc.ACROpenInsuranceLOA2)
	session.SetIDTokenClaimAuthTime(int(time.Now().Unix()))
	handleClaimsObject(session)

	return goidc.StatusSuccess
}

func handleClaimsObject(session *goidc.AuthnSession) {
	if session.Claims == nil {
		return
	}

	if slices.Contains(session.Claims.IDTokenEssentials(), goidc.ClaimAuthenticationContextReference) {
		session.SetIDTokenClaimACR(oidc.ACROpenInsuranceLOA2)
	}

	if slices.Contains(session.Claims.UserInfoEssentials(), goidc.ClaimAuthenticationContextReference) {
		session.SetUserInfoClaimACR(oidc.ACROpenInsuranceLOA2)
	}
}

func consentID(scopes string) string {
	for _, s := range strings.Split(scopes, " ") {
		if oidc.ScopeConsent.Matches(s) {
			return strings.Replace(s, "consent:", "", 1)
		}
	}
	return ""
}

type AuthnPage struct {
	BaseURL     string
	CallbackID  string
	Permissions []consent.Permission
	Error       string
}

const (
	paramConsentID   = "consent_id"
	paramPermissions = "permissions"
	paramConsentCPF  = "consent_cpf"
	paramUserID      = "user_id"
	paramStepID      = "step_id"

	stepIDSetUp      = "setup"
	stepIDLogin      = "login"
	stepIDConsent    = "consent"
	stepIDFinishFlow = "finish_flow"

	usernameFormParam = "username"
	passwordFormParam = "password"
	consentFormParam  = "consent"

	correctPassword = "pass"
)
