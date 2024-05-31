package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"nacos/configuration"
	"nacos/util"
	"nacos/util/collection"
	"sync"
	"time"
)

var (
	Manager    manager = nil
	secretKey          = "NamingAndConfigurationService@golang"
	expireTime         = int64(7200)
)

func init() {
	auth := configuration.Configuration.Nacos.Auth
	secretKey = util.ConditionalExpression(auth.SecretKey != "", auth.SecretKey, secretKey)
	expireTime = util.ConditionalExpression(auth.ExpireTime > 0, auth.ExpireTime, expireTime)
	manager := &jwtManager{}
	if auth.Cache {
		Manager = &jwtCacheManager{
			jwtManager: manager,
			mutex:      sync.RWMutex{},
			userCache:  collection.NewExpiredHashMap[string, *jwtCache](),
			tokenCache: collection.NewExpiredHashMap[string, *jwtCache](),
		}
	} else {
		Manager = manager
	}
}

type manager interface {
	CreateToken(username string) (string, *jwt.StandardClaims)
	ParseToken(tokenString string) (*jwt.StandardClaims, error)
}

type jwtManager struct{}

func (manager *jwtManager) CreateToken(username string) (string, *jwt.StandardClaims) {
	now := time.Now()
	claims := &jwt.StandardClaims{
		Subject:   username,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Duration(expireTime) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString, claims
}

func (manager *jwtManager) ParseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, _ := token.Claims.(*jwt.StandardClaims)
	if err = token.Claims.Valid(); err != nil {
		return nil, err
	}
	return claims, nil
}

type jwtCache struct {
	token  string
	claims *jwt.StandardClaims
}

type jwtCacheManager struct {
	*jwtManager
	mutex      sync.RWMutex
	userCache  *collection.ExpiredHashMap[string, *jwtCache]
	tokenCache *collection.ExpiredHashMap[string, *jwtCache]
}

func (manager *jwtCacheManager) CreateToken(username string) (string, *jwt.StandardClaims) {
	value, ok := func() (*jwtCache, bool) {
		manager.mutex.RLock()
		defer manager.mutex.RUnlock()
		return manager.userCache.Get(username)
	}()
	tokenExpired := false
	if !ok || value == nil {
		tokenExpired = true
	} else {
		expiresAt := value.claims.ExpiresAt
		tokenExpireTime := expiresAt - value.claims.IssuedAt
		ttl := expiresAt - time.Now().Unix()
		if float64(ttl)/float64(tokenExpireTime) <= 0.1 {
			tokenExpired = true
		}
	}
	if tokenExpired {
		token, claims := manager.jwtManager.CreateToken(username)
		manager.cache(username, token, claims, expireTime)
		return token, claims
	}
	return value.token, value.claims
}

func (manager *jwtCacheManager) ParseToken(tokenString string) (*jwt.StandardClaims, error) {
	value, ok := func() (*jwtCache, bool) {
		manager.mutex.RLock()
		defer manager.mutex.RUnlock()
		return manager.tokenCache.Get(tokenString)
	}()
	if ok {
		return value.claims, nil
	} else {
		claims, err := manager.jwtManager.ParseToken(tokenString)
		if err != nil {
			return nil, err
		} else {
			if claims.ExpiresAt-claims.IssuedAt > expireTime {
				return nil, errors.New("token invalid")
			}
			manager.cache(claims.Subject, tokenString, claims, claims.ExpiresAt-time.Now().Unix())
			return claims, nil
		}
	}
}

func (manager *jwtCacheManager) cache(username, token string, claims *jwt.StandardClaims, expireTime int64) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	cache := &jwtCache{token: token, claims: claims}
	cacheTime := time.Duration(expireTime) * time.Second
	manager.userCache.Put(username, cache, cacheTime)
	manager.tokenCache.Put(token, cache, cacheTime)
}
