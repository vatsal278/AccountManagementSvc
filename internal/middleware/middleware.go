package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/PereRohit/util/log"
	"github.com/PereRohit/util/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/vatsal278/AccountManagmentSvc/internal/codes"
	svcCfg "github.com/vatsal278/AccountManagmentSvc/internal/config"
	"github.com/vatsal278/AccountManagmentSvc/internal/model"
	"github.com/vatsal278/AccountManagmentSvc/internal/repo/authentication"
	"github.com/vatsal278/AccountManagmentSvc/pkg/session"
	"github.com/vatsal278/go-redis-cache"
	"github.com/vatsal278/msgbroker/pkg/sdk"
	"io"
	"net/http"
	"strings"
)

type AccMgmtMiddleware struct {
	cfg    *svcCfg.Config
	jwt    authentication.JWTService
	msg    func(io.ReadCloser) (string, error)
	cacher redis.Cacher
}

type respWriterWithStatus struct {
	status   int
	response string
	http.ResponseWriter
}

func (w *respWriterWithStatus) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *respWriterWithStatus) Write(d []byte) (int, error) {
	w.response = string(d)
	return w.ResponseWriter.Write(d)
}

func NewAccMgmtMiddleware(cfg *svcCfg.SvcConfig) *AccMgmtMiddleware {
	msgQueue := sdk.NewMsgBrokerSvc(cfg.Cfg.MessageQueue.SvcUrl)
	msg := msgQueue.ExtractMsg(&cfg.MsgBrokerSvc.PrivateKey)
	return &AccMgmtMiddleware{
		cfg:    cfg.Cfg,
		jwt:    cfg.JwtSvc.JwtSvc,
		msg:    msg,
		cacher: cfg.Cacher.Cacher,
	}
}

func (u AccMgmtMiddleware) ExtractUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			log.Error(err)
			response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrUnauthorized), nil)
			return
		}
		if cookie.Value == "" {
			log.Error(err)
			response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrUnauthorized), nil)
			return
		}
		token, err := u.jwt.ValidateToken(cookie.Value)
		if err != nil {
			log.Error(err)
			if strings.Contains(err.Error(), "Token is expired") {
				response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrTokenExpired), nil)
				return
			}
			response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrMatchingToken), nil)
			return
		}
		if !token.Valid {
			response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrUnauthorized), nil)
			return
		}
		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.ToJson(w, http.StatusInternalServerError, codes.GetErr(codes.ErrAssertClaims), nil)
			return
		}
		userId, ok := mapClaims["user_id"]
		if !ok {
			response.ToJson(w, http.StatusInternalServerError, codes.GetErr(codes.ErrAssertUserid), nil)
			return
		}
		ctx := session.SetSession(r.Context(), userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (u AccMgmtMiddleware) ScreenRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var urlMatch bool
		if r.UserAgent() != u.cfg.MessageQueue.UserAgent {
			log.Error(codes.GetErr(codes.ErrUnauthorizedAgent))
			response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrUnauthorizedAgent), nil)
			return
		}
		if u.cfg.MessageQueue.UrlCheck != false {
			for _, v := range u.cfg.MessageQueue.AllowedUrl {
				if v == r.RemoteAddr {
					urlMatch = true
				}
			}
			if urlMatch != true {
				log.Error(codes.GetErr(codes.ErrUnauthorizedUrl))
				response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrUnauthorizedUrl), nil)
				return
			}
		}
		decryptMsg, err := u.msg(r.Body)
		if err != nil {
			log.Error(err)
			response.ToJson(w, http.StatusInternalServerError, codes.GetErr(codes.ErrExtractMsg), nil)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer([]byte(decryptMsg)))
		next.ServeHTTP(w, r)
	})
}
func (u AccMgmtMiddleware) Cacher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := session.GetSession(r.Context())
		idStr, ok := id.(string)
		if !ok {
			response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrAssertUserid), nil)
			return
		}
		Cacher := u.cacher
		by, err := Cacher.Get("accounts/summary/user_id/" + idStr)
		if err != nil {
			log.Error()
			next.ServeHTTP(w, r)
			return
		}
		var summary model.AccountSummary
		err = json.Unmarshal(by, &summary)
		if err != nil {
			log.Error()
			next.ServeHTTP(w, r)
			return
		}
		response.ToJson(w, http.StatusOK, "SUCCESS", summary)
	})
}
