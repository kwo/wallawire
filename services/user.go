package services

import (
	"context"
	"net/http"
	"strings"
	"time"

	"wallawire/ctxutil"
	"wallawire/model"
)

const (
	componentUserService = "UserService"
)

type UserRepository interface {
	GetUser(context.Context, model.ReadOnlyTransaction, string) (*model.User, error)
	GetActiveUserByUsername(context.Context, model.ReadOnlyTransaction, string) (*model.User, error)
	GetUserRoles(context.Context, model.ReadOnlyTransaction, string, *time.Time) ([]model.UserRole, error)
	IsUsernameAvailable(context.Context, model.ReadOnlyTransaction, string) (bool, error)
	SetUser(context.Context, model.WriteOnlyTransaction, model.User) error
}

type IdGenerator interface {
	NewID() string
}

// TODO: need policy
func isValidPassword(password string) bool {
	if l := len(strings.TrimSpace(password)); l >= 8 && l <= 72 {
		return true
	}
	return false
}

// TODO: need policy
func isValidUsername(username string) bool {
	if l := len(strings.TrimSpace(username)); l >= 3 && l <= 64 {
		return true
	}
	return false
}

type UserService struct {
	db       model.Database
	userRepo UserRepository
	idgen    IdGenerator
}

func NewUserService(db model.Database, userRepo UserRepository, idgen IdGenerator) *UserService {
	return &UserService{
		db:       db,
		userRepo: userRepo,
		idgen:    idgen,
	}
}

func (z *UserService) ChangeUsername(ctx context.Context, req model.ChangeUsernameRequest) model.ChangeUsernameResponse {

	logger := ctxutil.NewLogger(componentUserService, "ChangeUsername", ctx)

	var user *model.User
	var roles []model.UserRole

	err := z.db.Run(func(tx model.Transaction) error {

		u, errGet := z.userRepo.GetUser(ctx, tx, req.UserID)
		if errGet != nil {
			logger.Error().Err(errGet).Msg("repo GetUser")
			return errGet // 500
		}
		if u == nil {
			return model.NewNotFoundError("user not found") // 404
		}
		if !u.MatchPassword(req.Password) {
			return model.NewValidationError("password incorrect") // 400
		}

		// validation of cur
		if !isValidUsername(req.NewUsername) {
			return model.NewValidationError("username not valid") // 400
		}

		ok, errCheckUsername := z.userRepo.IsUsernameAvailable(ctx, tx, req.NewUsername)
		if errCheckUsername != nil {
			logger.Error().Err(errCheckUsername).Msg("repo IsUsernameAvailable")
			return errCheckUsername // 500
		}
		if !ok {
			return model.NewValidationError("username not available") // 400
		}

		u.Username = req.NewUsername
		if err := z.userRepo.SetUser(ctx, tx, *u); err != nil {
			logger.Error().Err(err).Msg("repo SetUser")
			return err // 500
		}
		now := time.Now()
		rs, errRoles := z.userRepo.GetUserRoles(ctx, tx, u.ID, &now)
		if errRoles != nil {
			return errRoles
		}
		user = u
		roles = rs
		return nil
	})

	rsp := model.ChangeUsernameResponse{}

	if err != nil {
		logger.Debug().Err(err).Msg("cannot change username")
		rsp.Message = err.Error()
		if model.IsValidationError(err) {
			rsp.Code = http.StatusBadRequest
		} else if model.IsNotFoundError(err) {
			rsp.Code = http.StatusNotFound
		} else {
			rsp.Code = http.StatusInternalServerError
		}
	} else {
		logger.Debug().Msg("username updated")
		suc := ctxutil.TokenFromContext(ctx)
		su := model.ToSessionToken(suc.SessionID, user, roles, suc.Issued, suc.Expires)
		rsp.Code = http.StatusOK
		rsp.SessionToken = su
	}

	return rsp

}

func (z *UserService) ChangePassword(ctx context.Context, req model.ChangePasswordRequest) model.ChangePasswordResponse {

	logger := ctxutil.NewLogger(componentUserService, "ChangePassword", ctx)

	var user *model.User
	var roles []model.UserRole

	err := z.db.Run(func(tx model.Transaction) error {
		u, errGet := z.userRepo.GetUser(ctx, tx, req.UserID)
		if errGet != nil {
			logger.Error().Err(errGet).Msg("repo GetUser")
			return errGet // 500
		}
		if u == nil {
			return model.NewNotFoundError("user not found") // 404
		}
		if !u.MatchPassword(req.Password) {
			return model.NewValidationError("password incorrect") // 400
		}
		if !isValidPassword(req.NewPassword) {
			return model.NewValidationError("invalid new password") // 400
		}
		if err := u.SetPassword(req.NewPassword); err != nil {
			logger.Error().Err(err).Msg("user SetPassword")
			return err // 500
		}
		if err := z.userRepo.SetUser(ctx, tx, *u); err != nil {
			logger.Error().Err(err).Msg("repo SetUser")
			return err // 500
		}
		now := time.Now()
		rs, errRoles := z.userRepo.GetUserRoles(ctx, tx, u.ID, &now)
		if errRoles != nil {
			return errRoles
		}
		user = u
		roles = rs
		return nil
	});
	rsp := model.ChangePasswordResponse{}

	if err != nil {
		logger.Debug().Err(err).Msg("password NOT updated")
		rsp.Message = err.Error()
		if model.IsValidationError(err) {
			rsp.Code = http.StatusBadRequest
		} else if model.IsNotFoundError(err) {
			rsp.Code = http.StatusNotFound
		} else {
			rsp.Code = http.StatusInternalServerError
		}
	} else {
		logger.Debug().Msg("password updated")
		suc := ctxutil.TokenFromContext(ctx)
		su := model.ToSessionToken(suc.SessionID, user, roles, suc.Issued, suc.Expires)
		rsp.Code = http.StatusOK
		rsp.SessionToken = su
	}

	return rsp

}

func (z *UserService) ChangeProfile(ctx context.Context, req model.ChangeProfileRequest) model.ChangeProfileResponse {

	logger := ctxutil.NewLogger(componentUserService, "ChangeProfile", ctx)

	var user *model.User
	var roles []model.UserRole

	err := z.db.Run(func(tx model.Transaction) error {

		u, errGet := z.userRepo.GetUser(ctx, tx, req.UserID)
		if errGet != nil {
			logger.Error().Err(errGet).Msg("repo GetUser")
			return errGet // 500
		}
		if u == nil {
			return model.NewNotFoundError("user not found") // 404
		}

		// validation of cur
		if !isValidUsername(req.Displayname) {
			return model.NewValidationError("displayname not valid") // 400
		}

		u.Name = req.Displayname
		if err := z.userRepo.SetUser(ctx, tx, *u); err != nil {
			logger.Error().Err(err).Msg("repo SetUser")
			return err // 500
		}

		now := time.Now()
		rs, errRoles := z.userRepo.GetUserRoles(ctx, tx, u.ID, &now)
		if errRoles != nil {
			return errRoles
		}
		user = u
		roles = rs

		return nil
	})

	rsp := model.ChangeProfileResponse{}

	if err != nil {
		logger.Debug().Err(err).Msg("cannot change profile")
		rsp.Message = err.Error()
		if model.IsValidationError(err) {
			rsp.Code = http.StatusBadRequest
		} else if model.IsNotFoundError(err) {
			rsp.Code = http.StatusNotFound
		} else {
			rsp.Code = http.StatusInternalServerError
		}
	} else {
		logger.Debug().Msg("profile updated")
		rsp.Code = http.StatusOK
		suc := ctxutil.TokenFromContext(ctx)
		su := model.ToSessionToken(suc.SessionID, user, roles, suc.Issued, suc.Expires)
		rsp.SessionToken = su
	}

	return rsp

}

func (z *UserService) Login(ctx context.Context, req model.LoginRequest) model.LoginResponse {

	logger := ctxutil.NewLogger(componentUserService, "Login", ctx)

	// db username field is 64, go bcrypt max is 72
	if x, y := len(req.Username), len(req.Password); x == 0 || y == 0 || x > 64 || y > 72 {
		msg := "invalid username/password"
		logger.Debug().Msg(msg)
		return model.LoginResponse{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	var user *model.User
	var roles []model.UserRole

	err := z.db.Run(func(tx model.Transaction) error {
		usr, errGet := z.userRepo.GetActiveUserByUsername(ctx, tx, req.Username)
		if errGet != nil {
			logger.Error().Err(errGet).Msg("repo GetActiveUserByUsername")
			return errGet // 500
		}
		if usr == nil {
			// 400 - do not give out info that user does not exist at login
			return model.NewValidationError("invalid username/password")
		}
		if !usr.MatchPassword(req.Password) {
			return model.NewValidationError("invalid username/password") // 400
		}
		now := time.Now()
		rs, errRoles := z.userRepo.GetUserRoles(ctx, tx, usr.ID, &now)
		if errRoles != nil {
			return errRoles
		}
		user = usr
		roles = rs
		return nil
	});
	rsp := model.LoginResponse{}

	if err != nil {
		logger.Debug().Err(err).Msg("cannot login")
		rsp.Message = err.Error()
		if model.IsValidationError(err) {
			rsp.Code = http.StatusBadRequest
		} else {
			rsp.Code = http.StatusInternalServerError
		}
	} else {
		sessionID := z.idgen.NewID()
		rsp.Code = http.StatusOK
		issued := time.Now().Truncate(time.Minute)
		expires := issued.Add(model.LoginTimeout)
		rsp.SessionToken = model.ToSessionToken(sessionID, user, roles, issued, expires)
		logger.Info().Str("username", user.Username).Str("UserID", user.ID).Str("SessionID", sessionID).Msg("login")
	}

	return rsp
}
