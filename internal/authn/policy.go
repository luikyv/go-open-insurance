package authn

import (
	"context"
	"html/template"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/consent"
	"github.com/luikyv/go-open-insurance/internal/oidc"
	"github.com/luikyv/go-open-insurance/internal/user"
)

func Policy(
	userService user.Service,
	consentService consent.Service,
	baseURL string,
) goidc.AuthnPolicy {
	authenticator := authenticator{
		userService:    userService,
		consentService: consentService,
		baseURL:        baseURL,
	}
	return goidc.NewPolicy(
		"main",
		func(r *http.Request, c *goidc.Client, as *goidc.AuthnSession) bool {
			as.StoreParameter(paramStepID, stepIDSetUp)
			return true
		},
		authenticator.authenticate,
	)
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
	loginFormParam    = "login"
	consentFormParam  = "consent"

	correctPassword = "pass"
)

type authnPage struct {
	BaseURL     string
	CallbackID  string
	Permissions []api.ConsentPermission
	Error       string
}

type authenticator struct {
	userService    user.Service
	consentService consent.Service
	baseURL        string
}

func (a authenticator) authenticate(
	w http.ResponseWriter,
	r *http.Request,
	session *goidc.AuthnSession,
) goidc.AuthnStatus {
	ctx := context.WithValue(r.Context(), api.CtxKeyClientID, session.ClientID)
	r = r.WithContext(ctx)

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

func (a authenticator) setUp(
	r *http.Request,
	session *goidc.AuthnSession,
) goidc.AuthnStatus {
	consentID, ok := oidc.ConsentID(session.Scopes)
	if !ok {
		session.SetError("missing consent ID")
		return goidc.StatusFailure
	}

	consent, err := a.consentService.Get(
		r.Context(),
		consentID,
	)
	if err != nil {
		session.SetError(err.Error())
		return goidc.StatusFailure
	}

	if consent.Status != api.ConsentStatusAWAITINGAUTHORISATION {
		session.SetError("consent not awaiting authorization")
		return goidc.StatusFailure
	}

	// Convert permissions to []string for joining.
	strPermissions := make([]string, len(consent.Permissions))
	for i, permission := range consent.Permissions {
		strPermissions[i] = string(permission)
	}

	session.StoreParameter(paramConsentID, consent.ID)
	session.StoreParameter(paramPermissions, strings.Join(strPermissions, " "))
	session.StoreParameter(paramConsentCPF, consent.UserCPF)
	return goidc.StatusSuccess
}

func (a authenticator) login(
	w http.ResponseWriter,
	r *http.Request,
	session *goidc.AuthnSession,
) goidc.AuthnStatus {

	r.ParseForm()

	isLogin := r.PostFormValue(loginFormParam)
	if isLogin == "" {
		w.WriteHeader(http.StatusOK)
		// TODO: Improve this.
		tmpl, _ := template.ParseFiles("../../templates/login.html")
		tmpl.Execute(w, authnPage{
			BaseURL:    a.baseURL,
			CallbackID: session.CallbackID,
		})
		return goidc.StatusInProgress
	}

	if isLogin != "true" {
		consentID := session.Parameter(paramConsentID).(string)
		a.consentService.RejectByID(
			r.Context(),
			consentID,
			consent.RejectionInfo{
				RejectedBy: api.ConsentRejectedByUSER,
				Reason:     api.ConsentRejectedReasonCodeCUSTOMERMANUALLYREJECTED,
			},
		)
		session.SetError("consent not granted")
		return goidc.StatusFailure
	}

	username := r.PostFormValue(usernameFormParam)
	user, err := a.userService.User(username)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		tmpl, _ := template.ParseFiles("../../templates/login.html")
		tmpl.Execute(w, authnPage{
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
		tmpl.Execute(w, authnPage{
			BaseURL:    a.baseURL,
			CallbackID: session.CallbackID,
			Error:      "invalid credentials",
		})
		return goidc.StatusInProgress
	}

	session.StoreParameter(paramUserID, username)
	return goidc.StatusSuccess
}

func (a authenticator) grantConsent(
	w http.ResponseWriter,
	r *http.Request,
	session *goidc.AuthnSession,
) goidc.AuthnStatus {

	r.ParseForm()

	var permissions []api.ConsentPermission
	for _, p := range strings.Split(session.Parameter(paramPermissions).(string), " ") {
		permissions = append(permissions, api.ConsentPermission(p))
	}
	isConsented := r.PostFormValue(consentFormParam)
	if isConsented == "" {
		w.WriteHeader(http.StatusOK)
		tmpl, _ := template.ParseFiles("../../templates/consent.html")
		tmpl.Execute(w, authnPage{
			BaseURL:     a.baseURL,
			CallbackID:  session.CallbackID,
			Permissions: permissions,
		})
		return goidc.StatusInProgress
	}

	consentID := session.Parameter(paramConsentID).(string)

	if isConsented != "true" {
		a.consentService.RejectByID(
			r.Context(),
			consentID,
			consent.RejectionInfo{
				RejectedBy: api.ConsentRejectedByUSER,
				Reason:     api.ConsentRejectedReasonCodeCUSTOMERMANUALLYREJECTED,
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

func (a authenticator) finishFlow(
	session *goidc.AuthnSession,
) goidc.AuthnStatus {
	session.SetUserID(session.Parameter(paramUserID).(string))
	// TODO: Grant scopes based on permissions.
	session.GrantScopes(session.Scopes)
	session.SetIDTokenClaimACR(oidc.ACROpenInsuranceLOA2)
	session.SetIDTokenClaimAuthTime(int(time.Now().Unix()))

	if session.Claims != nil {
		if slices.Contains(session.Claims.IDTokenEssentials(), goidc.ClaimACR) {
			session.SetIDTokenClaimACR(oidc.ACROpenInsuranceLOA2)
		}

		if slices.Contains(session.Claims.UserInfoEssentials(), goidc.ClaimACR) {
			session.SetUserInfoClaimACR(oidc.ACROpenInsuranceLOA2)
		}
	}

	return goidc.StatusSuccess
}
