package dal

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ferryvg/hiring-test-go-users-api/internal/dal/entity"

	"github.com/ferryvg/hiring-test-go-users-api/internal/config"
	"github.com/ferryvg/hiring-test-go-users-api/internal/db"
	"github.com/ferryvg/hiring-test-go-users-api/internal/domain"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UsersManager interface {
	// Authenticate user by credentials and returns access token.
	Authenticate(ctx context.Context, idUser, secret string, options ...ManagerOpt) (*domain.JwtAccessToken, error)

	// VerifyToken checks provided access token and returns associated user.
	VerifyToken(ctx context.Context, accessToken string, options ...ManagerOpt) (*domain.User, error)

	// Get returns user information by id.
	Get(ctx context.Context, idUser string, options ...ManagerOpt) (*domain.User, error)

	// GetList returns list of users.
	GetList(ctx context.Context, options ...ManagerOpt) ([]domain.User, error)

	// Create creates new user.
	Create(ctx context.Context, user *domain.User, options ...ManagerOpt) error

	// ChangeRoles updates user roles.
	ChangeRoles(ctx context.Context, idUser string, roles map[domain.Role]bool, options ...ManagerOpt) error
}

type users struct {
	jwtConfig *config.JwtConfig
	baseManager
}

func NewUsersManager(jwtConfig *config.JwtConfig, manager db.Manager, logger logrus.FieldLogger) UsersManager {
	return &users{
		jwtConfig: jwtConfig,
		baseManager: baseManager{
			manager: manager,
			logger:  logger,
		},
	}
}

func (m *users) Get(ctx context.Context, idUser string, options ...ManagerOpt) (*domain.User, error) {
	if idUser == "" {
		return nil, nil
	}

	opts, err := m.resolveOptions(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare operation: %w", err)
	}

	query := "SELECT id_user, secret FROM users WHERE id_user = ?"

	var row *sqlx.Row
	if opts.Tx != nil {
		row = opts.Tx.QueryRowxContext(opts.Ctx, query, idUser)
	} else {
		row = opts.Conn.QueryRowxContext(opts.Ctx, query, idUser)
	}

	var userEntity entity.User
	if err := row.StructScan(&userEntity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("failed to load user with '%s' id: %w", idUser, err)
	}

	user := &domain.User{
		ID:     userEntity.ID,
		Secret: userEntity.Secret,
	}

	if err := m.adjustUser(user, opts); err != nil {
		return nil, fmt.Errorf("failed to adjust user details: %w", err)
	}

	return user, nil
}

func (m *users) Authenticate(
	ctx context.Context,
	idUser, secret string,
	options ...ManagerOpt,
) (*domain.JwtAccessToken, error) {
	user, err := m.Get(ctx, idUser, options...)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Secret), []byte(secret)); err != nil {
		m.logger.WithError(err).
			WithField("user", idUser).WithField("secret", secret).
			Info("Failed to check user secret")
		return nil, ErrInvalidCredentials
	}

	accessToken, err := m.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	token := &domain.JwtAccessToken{
		IDUser:      user.ID,
		AccessToken: accessToken,
		CreatedAt:   time.Now(),
		ExpiredAt:   time.Now().Add(m.jwtConfig.TTL),
	}

	if err := m.registerAccessToken(ctx, token, options...); err != nil {
		return nil, fmt.Errorf("failed to register user access token: %w", err)
	}

	return token, nil
}

func (m *users) VerifyToken(ctx context.Context, accessToken string, options ...ManagerOpt) (*domain.User, error) {
	if accessToken == "" {
		return nil, ErrNotFound
	}

	opts, err := m.resolveOptions(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare operation: %w", err)
	}

	query := "SELECT id_user FROM jwt_access_tokens WHERE access_token = ? AND expired_at > ?"

	var row *sqlx.Row
	now := time.Now()

	if opts.Tx != nil {
		row = opts.Tx.QueryRowxContext(opts.Ctx, query, accessToken, now)
	} else {
		row = opts.Conn.QueryRowxContext(opts.Ctx, query, accessToken, now)
	}

	var idUser string
	if err := row.Scan(&idUser); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("failed to load access token '%s' info: %w", accessToken, err)
	}

	return m.Get(ctx, idUser, options...)
}

func (m *users) Create(ctx context.Context, user *domain.User, options ...ManagerOpt) error {
	if user == nil || user.ID == "" || user.Secret == "" {
		return fmt.Errorf("invalid user")
	}

	opts, err := m.resolveOptions(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to prepare operation: %w", err)
	}

	if _, err := m.Get(ctx, user.ID, options...); errors.Is(err, ErrNotFound) && errors.Is(err, sql.ErrNoRows) {
		if err == nil {
			return ErrUserExists
		}

		return err
	}

	hashedSecret, err := m.hashSecret(user.Secret)
	if err != nil {
		return fmt.Errorf("failed to hash user secret: %w", err)
	}

	userEntity := entity.User{
		ID:     user.ID,
		Secret: hashedSecret,
	}

	query := "INSERT INTO users (id_user, secret) VALUES(:id_user, :secret)"

	if opts.Tx != nil {
		_, err = opts.Tx.NamedExecContext(ctx, query, userEntity)
	} else {
		_, err = opts.Conn.NamedExecContext(ctx, query, userEntity)
	}

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return m.ChangeRoles(ctx, user.ID, user.Roles, options...)
}

func (m *users) ChangeRoles(
	ctx context.Context,
	idUser string,
	roles map[domain.Role]bool,
	options ...ManagerOpt,
) error {
	if idUser == "" {
		return nil
	}

	opts, err := m.resolveOptions(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to prepare operation: %w", err)
	}

rolesLoop:
	for role, enabled := range roles {
		switch role {
		case domain.GuestRole, domain.UnknownRole:
			continue rolesLoop
		}

		if enabled {
			err = m.enableRole(idUser, role, opts)
		} else {
			err = m.disableRole(idUser, role, opts)
		}

		if err != nil {
			return fmt.Errorf("failed to change user role '%s' status: %w", role, err)
		}
	}

	return nil
}

func (m *users) GetList(ctx context.Context, options ...ManagerOpt) ([]domain.User, error) {
	opts, err := m.resolveOptions(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare operation: %w", err)
	}

	query := "SELECT id_user FROM users"

	var userEntitiesList []entity.User

	if opts.Tx != nil {
		err = opts.Tx.SelectContext(opts.Ctx, &userEntitiesList, query)
	} else {
		err = opts.Conn.SelectContext(opts.Ctx, &userEntitiesList, query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load users list: %w", err)
	}

	usersList := make([]domain.User, 0, len(userEntitiesList))
	for _, user := range userEntitiesList {
		usersList = append(usersList, domain.User{ID: user.ID, Secret: user.Secret})
	}

	if err := m.adjustUsersList(usersList, opts); err != nil {
		return nil, fmt.Errorf("failed to adjust users details: %w", err)
	}

	return usersList, nil
}

func (m *users) adjustUsersList(usersList []domain.User, opts *Options) error {
	userIds := make([]string, 0, len(usersList))
	for _, user := range usersList {
		userIds = append(userIds, user.ID)
	}

	usersRoles, err := m.loadUsersRoles(userIds, opts)
	if err != nil {
		return fmt.Errorf("failed to load users roles: %w", err)
	}

	for idx, user := range usersList {
		roles := usersRoles[user.ID]

		if len(roles) == 0 {
			roles = m.buildDomainRoles(m.getDefaultRoles())
		}

		user.Roles = roles

		usersList[idx] = user
	}

	return nil
}

func (m *users) adjustUser(user *domain.User, opts *Options) error {
	roles, err := m.loadRoles(user.ID, opts)
	if err != nil {
		return fmt.Errorf("failed to load '%s' user roles: %w", user.ID, err)
	}

	user.Roles = roles

	return nil
}

func (m *users) loadRoles(idUser string, opts *Options) (map[domain.Role]bool, error) {
	query := "SELECT id_role FROM user_roles WHERE id_user = ?"

	var userRoles []int

	var err error
	if opts.Tx != nil {
		err = opts.Tx.SelectContext(opts.Ctx, &userRoles, query, idUser)
	} else {
		err = opts.Conn.SelectContext(opts.Ctx, &userRoles, query, idUser)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.buildDomainRoles(m.getDefaultRoles()), nil
		}

		return nil, fmt.Errorf("failed to load user roles: %w", err)
	}

	return m.buildDomainRoles(userRoles), nil
}

func (m *users) loadUsersRoles(usersIds []string, opts *Options) (map[string]map[domain.Role]bool, error) {
	if len(usersIds) == 0 {
		return nil, nil
	}

	query := "SELECT id_user, id_role FROM user_roles WHERE id_user IN (?)"

	query, args, err := sqlx.In(query, usersIds)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare users roles query: %w", err)
	}

	var usersRoles []entity.UserRole

	if opts.Tx != nil {
		err = opts.Tx.SelectContext(opts.Ctx, &usersRoles, query, args...)
	} else {
		err = opts.Conn.SelectContext(opts.Ctx, &usersRoles, query, args...)
	}

	result := make(map[string]map[domain.Role]bool, len(usersIds))

	for _, userRole := range usersRoles {
		if _, exists := result[userRole.IDUser]; !exists {
			result[userRole.IDUser] = make(map[domain.Role]bool)
		}

		result[userRole.IDUser][domain.Role(userRole.Role)] = true
	}

	return result, nil
}

func (m *users) buildDomainRoles(rolesList []int) map[domain.Role]bool {
	roles := make(map[domain.Role]bool)

	for _, role := range rolesList {
		roles[domain.Role(role)] = true
	}

	return roles
}

func (m *users) getDefaultRoles() []int {
	return []int{
		int(domain.GuestRole),
	}
}

func (m *users) hashSecret(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate secret hash: %w", err)
	}

	return string(hash), nil
}

func (m *users) generateAccessToken(user *domain.User) (string, error) {
	identity := user.ID + ":" + user.Secret

	hash, err := bcrypt.GenerateFromPassword([]byte(identity), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return base64.StdEncoding.EncodeToString(hash), nil
}

func (m *users) registerAccessToken(
	ctx context.Context,
	token *domain.JwtAccessToken,
	options ...ManagerOpt,
) error {
	opts, err := m.resolveOptions(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to prepare operation: %w", err)
	}

	query := `INSERT INTO jwt_access_tokens(id_user, access_token, created_at, expired_at) 
          VALUES(:id_user, :access_token, :created_at, :expired_at)`

	row := &entity.JwtAccessToken{
		IDUser:      token.IDUser,
		AccessToken: token.AccessToken,
		CreatedAt:   token.CreatedAt,
		ExpiredAt:   token.ExpiredAt,
	}

	if opts.Tx != nil {
		_, err = opts.Tx.NamedExecContext(opts.Ctx, query, row)
	} else {
		_, err = opts.Conn.NamedExecContext(opts.Ctx, query, row)
	}

	if err != nil {
		return fmt.Errorf("failed to execute db query: %w", err)
	}

	return nil
}

func (m *users) enableRole(idUser string, role domain.Role, opts *Options) error {
	query := "INSERT IGNORE INTO user_roles(id_user, id_role) VALUES(?, ?)"

	var err error
	if opts.Tx != nil {
		_, err = opts.Tx.ExecContext(opts.Ctx, query, idUser, role)
	} else {
		_, err = opts.Conn.ExecContext(opts.Ctx, query, idUser, role)
	}

	return err
}

func (m *users) disableRole(idUser string, role domain.Role, opts *Options) error {
	query := "DELETE FROM user_roles WHERE id_user = ? AND id_role = ?"

	var err error
	if opts.Tx != nil {
		_, err = opts.Tx.ExecContext(opts.Ctx, query, idUser, role)
	} else {
		_, err = opts.Conn.ExecContext(opts.Ctx, query, idUser, role)
	}

	return err
}
