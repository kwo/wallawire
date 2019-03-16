package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"wallawire/logging"
	"wallawire/model"
)

// TODO: missing role management: List,Set,Delete Roles
// TODO: missing user List function

type dbUser struct {
	UserID       sql.NullString `db:"id"`
	Disabled     bool           `db:"disabled"`
	Username     sql.NullString `db:"username"`
	Name         sql.NullString `db:"name"`
	PasswordHash sql.NullString `db:"password_hash"`
	Created      sql.NullInt64  `db:"created"`
	Updated      sql.NullInt64  `db:"updated"`
}

type dbUserRole struct {
	ID        sql.NullString `db:"id"`
	Name      sql.NullString `db:"name"`
	ValidFrom sql.NullInt64  `db:"valid_from"`
	ValidTo   sql.NullInt64  `db:"valid_to"`
}

func (z *Repository) GetUser(ctx context.Context, tx model.ReadOnlyTransaction, userID string) (*model.User, error) {

	logger := logging.New(ctx, componentRepo, "GetUser")
	logger.Debug().Msg("invoked")

	query := `
	SELECT id, disabled, username, name, password_hash, created, updated
    FROM users
    WHERE id = :id
	`
	params := map[string]interface{}{
		"id": userID,
	}

	rs, errQuery := tx.Query(query, params)
	if errQuery != nil {
		return nil, errQuery
	}

	defer func() {
		if err := rs.Close(); err != nil {
			logger.Warn().Err(err).Msg("cannot close resultset")
		}
	}()

	var user *model.User

	if rs.Next() {
		u := dbUser{}
		if err := rs.StructScan(&u); err != nil {
			return nil, err
		}
		user = convertToUser(u)
	}

	return user, nil

}

func (z *Repository) GetActiveUserByUsername(ctx context.Context, tx model.ReadOnlyTransaction, username string) (*model.User, error) {

	logger := logging.New(ctx, componentRepo, "GetActiveUserByUsername")
	logger.Debug().Str("username", username).Msg("invoked")

	query := `
	SELECT id, disabled, username, name, password_hash, created, updated
    FROM users
    WHERE username = :username
	`
	params := map[string]interface{}{
		"username": username,
	}

	rs, errQuery := tx.Query(query, params)
	if errQuery != nil {
		return nil, errQuery
	}

	defer func() {
		if err := rs.Close(); err != nil {
			logger.Warn().Err(err).Msg("cannot close resultset")
		}
	}()

	var user *model.User
	if rs.Next() {
		u := dbUser{}
		if err := rs.StructScan(&u); err != nil {
			return nil, err
		}
		user = convertToUser(u)
	}

	if user == nil || user.Disabled {
		return nil, nil
	}

	return user, nil

}

func (z *Repository) IsUsernameAvailable(ctx context.Context, tx model.ReadOnlyTransaction, username string) (bool, error) {

	logger := logging.New(ctx, componentRepo, "IsUsernameAvailable")
	logger.Debug().Msg("invoked")

	query := "SELECT username from usernames WHERE username = :username"
	params := map[string]interface{}{
		"username": strings.ToLower(username),
	}

	rs, errQuery := tx.Query(query, params)
	if errQuery != nil {
		return false, errQuery
	}

	defer func() {
		if err := rs.Close(); err != nil {
			logger.Warn().Err(err).Msg("cannot close resultset")
		}
	}()

	var uname string
	if rs.Next() {
		if err := rs.Scan(&uname); err != nil {
			return false, err
		}
	}

	if uname == "" {
		return true, nil
	}

	return false, nil

}

// SetUser will add or update a user except created and updated which are set automatically.
func (z *Repository) SetUser(ctx context.Context, tx model.WriteOnlyTransaction, user model.User) error {

	logger := logging.New(ctx, componentRepo, "SetUser")
	logger.Debug().Msg("invoked")

	query := `
	INSERT INTO users (id, disabled, username, name, password_hash, created, updated)
	VALUES (:id, :disabled, :username, :name, :passwordHash, EXTRACT('epoch', now()), EXTRACT('epoch', now()))
	ON CONFLICT (id) DO UPDATE SET
	disabled = :disabled,
	username = :username,
	name = :name,
	password_hash = :passwordHash,
	updated = EXTRACT('epoch', now())
	`
	params := userToParams(user)
	if _, err := tx.Exec(query, params); err != nil {
		return err
	}
	return nil
}

func (z *Repository) DeleteUser(ctx context.Context, tx model.WriteOnlyTransaction, userID string) error {

	logger := logging.New(ctx, componentRepo, "DeleteUser")
	logger.Debug().Msg("invoked")

	errRoles := z.deleteUserRoles(tx, userID)
	if errRoles != nil {
		return errRoles
	}

	errUser := z.deleteUser(tx, userID)
	if errUser != nil {
		return errUser
	}

	return nil

}

// GetUserRoles returns all the roles for a user.
// Only roles active at given time will be returned if parameter is non-nil.
func (z *Repository) GetUserRoles(ctx context.Context, tx model.ReadOnlyTransaction, userID string, t *time.Time) ([]model.UserRole, error) {

	logger := logging.New(ctx, componentRepo, "GetUserRoles")
	logger.Debug().Msg("invoked")

	query := `
	SELECT r.id, r.name, ur.valid_from, ur.valid_to
	FROM roles r
	JOIN user_role ur ON (ur.role_id = r.id AND ur.user_id = :userID)
	`

	params := map[string]interface{}{
		"userID": userID,
	}

	if t != nil {
		query += `WHERE (ur.valid_from IS NULL OR ur.valid_from <= :t)`
		query += `AND (ur.valid_to IS NULL OR ur.valid_to > :t)`
		params["t"] = toNullTimeInteger(t)
	}

	query += "ORDER BY r.name"

	rs, errQuery := tx.Query(query, params)
	if errQuery != nil {
		return nil, errQuery
	}

	defer func() {
		if err := rs.Close(); err != nil {
			logger.Warn().Err(err).Msg("cannot close resultset")
		}
	}()

	roles := make([]model.UserRole, 0)
	for rs.Next() {
		var role dbUserRole
		if err := rs.StructScan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, convertToRole(role))
	}

	return roles, nil

}

func (z *Repository) SetUserRoles(ctx context.Context, tx model.Transaction, userID string, roles []model.UserRole) error {

	logger := logging.New(ctx, componentRepo, "SetUserRoles")
	logger.Debug().Msg("invoked")

	currentRoles, errRoles := z.GetUserRoles(ctx, tx, userID, nil)
	if errRoles != nil {
		return errRoles
	}

	// delete old roles
	rolesToDelete := z.findComplement(currentRoles, roles)
	for _, ur := range rolesToDelete {
		if err := z.deleteUserRole(tx, userID, ur.ID); err != nil {
			return err
		}
	}

	// upsert given roles
	for _, role := range roles {
		if err := z.setUserRole(tx, userID, role); err != nil {
			return err
		}
	}

	return nil

}

func (z *Repository) deleteUser(tx model.WriteOnlyTransaction, userID string) error {

	query := "DELETE FROM users WHERE id = :userID"
	params := map[string]interface{}{
		"userID": userID,
	}
	rs, errExec := tx.Exec(query, params)
	if errExec != nil {
		return errExec
	}
	count, errCount := rs.RowsAffected()
	if errCount != nil {
		return errCount
	}
	if count == 0 {
		return errors.New("no records deleted")
	} else if count != 1 {
		return errors.New("multiple records deleted")
	}
	return nil

}

func (z *Repository) deleteUserRole(tx model.WriteOnlyTransaction, userID, roleID string) error {

	query := "DELETE FROM user_role WHERE user_id = :userID AND role_id = :roleID"
	params := map[string]interface{}{
		"userID": userID,
		"roleID": roleID,
	}
	_, errExec := tx.Exec(query, params)
	return errExec

}

func (z *Repository) deleteUserRoles(tx model.WriteOnlyTransaction, userID string) error {

	query := "DELETE FROM user_role WHERE user_id = :userID"
	params := map[string]interface{}{
		"userID": userID,
	}
	_, errExec := tx.Exec(query, params)
	return errExec

}

func (z *Repository) setUserRole(tx model.WriteOnlyTransaction, userID string, role model.UserRole) error {
	query := `
	INSERT INTO user_role (user_id, role_id, valid_from, valid_to)
	VALUES (:userID, :roleID, :validFrom, :validTo)
	ON CONFLICT (user_id, role_id) DO UPDATE SET
	valid_from = :validFrom,
	valid_to = :validTo
	`
	params := userroleToParams(userID, role)
	if _, err := tx.Exec(query, params); err != nil {
		return err
	}
	return nil
}

// findComplement returns the roles in a but not in b
func (z *Repository) findComplement(a, b []model.UserRole) []model.UserRole {

	includes := func(z []model.UserRole, roleID string) bool {
		for _, x := range z {
			if x.ID == roleID {
				return true
			}
		}
		return false
	}

	var result []model.UserRole
	for _, x := range a {
		if !includes(b, x.ID) {
			result = append(result, x)
		}
	}
	return result

}

func convertToUser(u dbUser) *model.User {
	return &model.User{
		ID:           u.UserID.String,
		Disabled:     u.Disabled,
		Username:     u.Username.String,
		Name:         u.Name.String,
		PasswordHash: u.PasswordHash.String,
		Created:      toTime(u.Created),
		Updated:      toTime(u.Updated),
	}
}

func convertToRole(r dbUserRole) model.UserRole {
	return model.UserRole{
		ID:        r.ID.String,
		Name:      r.Name.String,
		ValidFrom: toTimePointer(toTime(r.ValidFrom)),
		ValidTo:   toTimePointer(toTime(r.ValidTo)),
	}
}

func userToParams(user model.User) map[string]interface{} {
	return map[string]interface{}{
		"id":           toNullString(user.ID),
		"disabled":     user.Disabled,
		"username":     toNullString(user.Username),
		"name":         toNullString(user.Name),
		"passwordHash": toNullString(user.PasswordHash),
		// Note: no updated, created as those are handled automatically
	}
}

func userroleToParams(userID string, role model.UserRole) map[string]interface{} {
	return map[string]interface{}{
		"userID":    toNullString(userID),
		"roleID":    toNullString(role.ID),
		"validFrom": toNullTimeInteger(role.ValidFrom),
		"validTo":   toNullTimeInteger(role.ValidTo),
	}
}
