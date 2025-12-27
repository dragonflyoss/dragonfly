package auth

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// Claims is a minimal JWT claims set used for inter-component gRPC auth.
type Claims struct {
	Issuer   string    `json:"iss"`
	Audience string    `json:"aud"`
	IssuedAt time.Time `json:"iat"`
	Expires  time.Time `json:"exp"`
}

// internal registry for per-component server keys (temporary until config wiring is complete).
var serverKeys = map[string]string{}

// SetServerKey sets the shared signing key for a component's server (e.g., "manager", "scheduler").
func SetServerKey(component, key string) {
	serverKeys[component] = key
}

// GetServerKey retrieves the key for a component server.
func GetServerKey(component string) string {
	return serverKeys[component]
}

// SignHS256 signs the provided claims with the given shared secret key using HS256.
func SignHS256(key string, c Claims) (string, error) {
	if key == "" {
		return "", errors.New("jwt: empty signing key")
	}
	claims := jwtlib.MapClaims{
		"iss": c.Issuer,
		"aud": c.Audience,
		"iat": c.IssuedAt.Unix(),
		"exp": c.Expires.Unix(),
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

// ValidateHS256 validates token signature and basic claims. Returns parsed claims.
func ValidateHS256(key string, tokenStr string, expectedAudience string) (Claims, error) {
	var out Claims
	if key == "" {
		return out, errors.New("jwt: empty validation key")
	}
	parser := jwtlib.NewParser(jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}))
	token, err := parser.Parse(tokenStr, func(t *jwtlib.Token) (any, error) {
		return []byte(key), nil
	})
	if err != nil || !token.Valid {
		return out, errors.New("jwt: invalid token")
	}
	claims, ok := token.Claims.(jwtlib.MapClaims)
	if !ok {
		return out, errors.New("jwt: invalid claims type")
	}
	// Audience check
	if audAny, ok := claims["aud"]; ok {
		if audStr, ok := audAny.(string); ok {
			if expectedAudience != "" && audStr != expectedAudience {
				return out, errors.New("jwt: audience mismatch")
			}
			out.Audience = audStr
		}
	}
	// Issuer
	if issAny, ok := claims["iss"]; ok {
		if issStr, ok := issAny.(string); ok {
			out.Issuer = issStr
		}
	}
	// Time checks
	now := time.Now()
	if expAny, ok := claims["exp"]; ok {
		switch v := expAny.(type) {
		case float64:
			out.Expires = time.Unix(int64(v), 0)
		case int64:
			out.Expires = time.Unix(v, 0)
		case uint64:
			out.Expires = time.Unix(int64(v), 0)
		}
		if now.After(out.Expires) {
			return out, errors.New("jwt: token expired")
		}
	}
	if iatAny, ok := claims["iat"]; ok {
		switch v := iatAny.(type) {
		case float64:
			out.IssuedAt = time.Unix(int64(v), 0)
		case int64:
			out.IssuedAt = time.Unix(v, 0)
		}
	}
	return out, nil
}

// DurationClaims constructs Claims with given ttl.
func DurationClaims(issuer, audience string, ttl time.Duration) Claims {
	now := time.Now()
	return Claims{
		Issuer:   issuer,
		Audience: audience,
		IssuedAt: now,
		Expires:  now.Add(ttl),
	}
}
