package github

import (
	"sync"
	"time"
)

// TokenCacheManager 管理GitHub安装令牌的缓存
type TokenCacheManager struct {
	cache map[int64]*TokenCache // 安装ID -> 令牌缓存的映射
	mu    sync.RWMutex          // 保护缓存的互斥锁
}

// NewTokenCacheManager 创建一个新的令牌缓存管理器
func NewTokenCacheManager() *TokenCacheManager {
	return &TokenCacheManager{
		cache: make(map[int64]*TokenCache),
	}
}

// Get 获取指定安装ID的令牌，如果令牌存在且未过期则返回
func (m *TokenCacheManager) Get(installationID int64) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if cache, ok := m.cache[installationID]; ok {
		if time.Now().Before(cache.ExpiresAt) {
			return cache.Token, true
		}
	}
	return "", false
}

// Set 设置指定安装ID的令牌和过期时间
func (m *TokenCacheManager) Set(installationID int64, token string, expiresAt time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache[installationID] = &TokenCache{
		Token:     token,
		ExpiresAt: expiresAt,
	}
}

// Clean 清理过期的令牌
func (m *TokenCacheManager) Clean() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, cache := range m.cache {
		if now.After(cache.ExpiresAt) {
			delete(m.cache, id)
		}
	}
}
