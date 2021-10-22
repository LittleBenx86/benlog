package service

import (
	"errors"
	"fmt"
	"github.com/LittleBenx86/Benlog/internal/global/consts"
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/LittleBenx86/Benlog/internal/utils/encryptor"
	"github.com/LittleBenx86/Benlog/internal/utils/security"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type SecurityService struct {
	*dependencies.Dependencies
	TokenBuilder
}

func NewSecurityService(isAuth bool, d *dependencies.Dependencies) *SecurityService {
	return &SecurityService{
		Dependencies: d,
		TokenBuilder: &defaultTokenBuilder{
			SigningKey: "benlog",
			PostParseContentValidateFn: func() PostParseContentValidateFn {
				if isAuth {
					return anonymousPostParseValidate
				}
				return authenticatedPostParseValidate
			}(),
			SignFn: func(s string) ([]byte, error) {
				hash, e := encryptor.NewSHA256().SetPlainBytes([]byte(s)).Hash()
				if e != nil {
					return []byte{}, e
				}
				return []byte(hash), nil
			},
		},
	}
}

func (s *SecurityService) AddPolicies(policies []map[string]interface{}) {
	s.DBClient.Table(security.GetCasbinTable()).CreateInBatches(policies, 50)
}

type TokenBuilder interface {
	Generate(c CustomClaims) (string, error)
	Parse(token string) (*CustomClaims, error)
	Refresh(token string, increaseSeconds int64) (string, error)
}

type defaultTokenBuilder struct {
	SigningKey string
	PostParseContentValidateFn
	SignFn
}

func (a *defaultTokenBuilder) Generate(c CustomClaims) (string, error) {
	signedBytes, err := a.SignFn((&c).GeneratePreHashContent())
	if err != nil {
		return "", err
	}
	c.Hash = string(signedBytes)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c) // Generate jwt header and claims part
	return token.SignedString([]byte(a.SigningKey))       // Append signing key at tail
}

func (a *defaultTokenBuilder) Parse(token string) (*CustomClaims, error) {
	t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.SigningKey), nil
	})

	if t == nil {
		return nil, errors.New("invalid token to parse")
	}

	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			if e.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("malformed token")
			} else if e.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("inactive token")
			} else if e.Errors&jwt.ValidationErrorExpired != 0 {
				t.Valid = true
				goto realParseHandle
			} else {
				return nil, errors.New("invalid token to parse")
			}
		}
	}
realParseHandle:
	c, ok := t.Claims.(*CustomClaims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token to parse")
	}

	err = a.PostParseContentValidateFn(c, a.SignFn)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (a *defaultTokenBuilder) Refresh(token string, increaseSeconds int64) (string, error) {
	c, err := a.Parse(token)
	if err != nil {
		return "", err
	}
	c.ExpiresAt = time.Now().Unix() + increaseSeconds
	return a.Generate(*c)
}

type CustomClaims struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
	Hash string `json:"hash"`
	jwt.StandardClaims
}

func (c *CustomClaims) GeneratePreHashContent() string {
	return fmt.Sprintf("%s:[%s:%s]", "benlog", c.UID, c.Name)
}

type PostParseContentValidateFn func(*CustomClaims, SignFn) error

type SignFn func(string) ([]byte, error)

func anonymousPostParseValidate(claims *CustomClaims, fn SignFn) error {
	if claims.Name != consts.USER_ANONYMOUS_NAME {
		return errors.New("unknown user access")
	}
	if bytes, err := fn(claims.GeneratePreHashContent()); err != nil {
		return errors.New("post parse jwt validate hash failed")
	} else {
		if string(bytes) != claims.Hash {
			return errors.New("post parse jwt validate hash incorrect")
		}
	}
	return nil
}

func authenticatedPostParseValidate(claims *CustomClaims, fn SignFn) error {
	return nil
}
