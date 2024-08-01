package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/services"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidAuthToken    = errors.New("invalid auth token")
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrIncorrectPassword   = errors.New("incorrect password")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
	ErrUnknownClaimsType   = errors.New("unknown claims type")
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var creationParams = params{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

var authTokenDuration = time.Minute
var refreshTokenDuration = time.Hour * 24 * 30

var recipient = "https://example.com"

type AuthUser struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	IsVerified bool   `json:"is_verified"`
	Avatar     string `json:"avatar"`
}

type claims struct {
	jwt.RegisteredClaims
	User AuthUser `json:"user"`
}

type refreshClaims struct {
	jwt.RegisteredClaims
	UseCount int `json:"use_count"`
}

func DisableRefreshToken(ctx context.Context, refreshToken string) (err error) {
	parsedToken, err := jwt.ParseWithClaims(refreshToken, &refreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return services.PublicKey, nil
	})
	if err != nil {
		return
	}
	if rc, ok := parsedToken.Claims.(*refreshClaims); ok {
		err = deleteRefreshToken(ctx, rc.ID)
	} else {
		err = ErrUnknownClaimsType
	}
	return
}

func CreateNewUser(ctx context.Context, user util.User, password string) (authToken, refreshToken string, err error) {
	hash, err := generateFromPassword(password)
	if err != nil {
		return
	}
	user.PasswordHash = hash

	err = queries.PostUser(ctx, user)
	if err != nil {
		return
	}
	return CreateNewSession(ctx, user.Email, password)
}

func CreateNewSession(ctx context.Context, username, password string) (authToken, refreshToken string, err error) {
	user, err := queries.GetUserByUsernameOrEmail(ctx, username)
	if err != nil {
		return
	}

	match, err := comparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return
	}
	if !match {
		err = ErrIncorrectPassword
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims{
		jwt.RegisteredClaims{
			Issuer:    "livefetcher-auth",
			Subject:   fmt.Sprintf("%d", user.ID),
			Audience:  jwt.ClaimStrings{"https://example.com"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(authTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		AuthUser{
			ID:         user.ID,
			Email:      user.Email,
			Username:   user.Username,
			Nickname:   user.Nickname,
			IsVerified: user.IsVerified,
			Avatar:     user.Avatar,
		},
	})
	authToken, err = t.SignedString(services.PrivateKey)
	if err != nil {
		return
	}

	rtidRaw, err := generateRandomBytes(32)
	if err != nil {
		return
	}
	rtid := base64.StdEncoding.EncodeToString(rtidRaw)
	rt := jwt.NewWithClaims(jwt.SigningMethodEdDSA, refreshClaims{
		jwt.RegisteredClaims{
			Issuer:    "livefetcher-auth",
			Subject:   fmt.Sprintf("%d", user.ID),
			Audience:  jwt.ClaimStrings{recipient},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        rtid,
		},
		0,
	})
	refreshToken, err = rt.SignedString(services.PrivateKey)
	if err != nil {
		return
	}
	err = registerRefreshToken(ctx, rtid)
	return
}

func verifyAuthToken(authToken string) (user AuthUser, err error) {
	parsedToken, err := jwt.ParseWithClaims(authToken, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return services.PublicKey, nil
	})
	if err != nil {
		return
	}

	if rc, ok := parsedToken.Claims.(*claims); ok {
		user = rc.User
	} else {
		err = ErrInvalidAuthToken
	}
	return
}

func RefreshSession(ctx context.Context, oldRefreshToken string) (authToken, refreshToken string, err error) {
	parsedToken, err := jwt.ParseWithClaims(oldRefreshToken, &refreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return services.PublicKey, nil
	})
	if err != nil {
		return
	}
	if rc, ok := parsedToken.Claims.(*refreshClaims); ok {
		if !checkRefreshToken(ctx, rc.ID, rc.UseCount) {
			err = ErrInvalidRefreshToken
			return
		}

		var id int
		id, err = strconv.Atoi(rc.Subject)
		if err != nil {
			return
		}

		var user util.User
		user, err = queries.GetUserByID(ctx, id)
		if err != nil {
			return
		}

		t := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims{
			jwt.RegisteredClaims{
				Issuer:    "livefetcher-auth",
				Subject:   fmt.Sprintf("%d", user.ID),
				Audience:  jwt.ClaimStrings{"https://example.com"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(authTokenDuration)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			AuthUser{
				ID:         user.ID,
				Email:      user.Email,
				Username:   user.Username,
				Nickname:   user.Nickname,
				IsVerified: user.IsVerified,
				Avatar:     user.Avatar,
			},
		})
		authToken, err = t.SignedString(services.PrivateKey)
		if err != nil {
			return
		}

		rt := jwt.NewWithClaims(jwt.SigningMethodEdDSA, refreshClaims{
			jwt.RegisteredClaims{
				Issuer:    "livefetcher-auth",
				Subject:   fmt.Sprintf("%d", user.ID),
				Audience:  jwt.ClaimStrings{"https://example.com"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenDuration)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        rc.ID,
			},
			rc.UseCount + 1,
		})
		refreshToken, err = rt.SignedString(services.PrivateKey)
	} else {
		err = ErrUnknownClaimsType
	}
	return
}

func EndSession(ctx context.Context, oldRefreshToken string) (err error) {
	parsedToken, err := jwt.ParseWithClaims(oldRefreshToken, &refreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return services.PublicKey, nil
	})
	if err != nil {
		return
	}

	if rc, ok := parsedToken.Claims.(*refreshClaims); ok {
		deleteRefreshToken(ctx, rc.ID)
	} else {
		err = ErrInvalidRefreshToken
	}
	return
}

func generateFromPassword(password string) (encodedHash string, err error) {
	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(creationParams.saltLength)
	if err != nil {
		return
	}

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey([]byte(password), salt, creationParams.iterations, creationParams.memory, creationParams.parallelism, creationParams.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, creationParams.memory, creationParams.iterations, creationParams.parallelism, b64Salt, b64Hash)
	return
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
