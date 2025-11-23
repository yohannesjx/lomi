package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// TelegramUser represents the user data inside initData
type TelegramUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsPremium    bool   `json:"is_premium"`
}

// ValidateTelegramInitData validates the initData string from Telegram WebApp
func ValidateTelegramInitData(initData string, botToken string) (*TelegramUser, error) {
	// Parse query string
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}

	// Extract hash
	hash := values.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("hash is missing")
	}
	values.Del("hash")

	// Sort keys alphabetically
	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create data check string
	var dataCheckArr []string
	for _, k := range keys {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", k, values.Get(k)))
	}
	dataCheckString := strings.Join(dataCheckArr, "\n")

	// Compute HMAC-SHA256
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	secret := secretKey.Sum(nil)

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(dataCheckString))
	computedHash := hex.EncodeToString(h.Sum(nil))

	// Compare hashes
	if computedHash != hash {
		return nil, fmt.Errorf("invalid hash")
	}

	// Parse user data
	userJSON := values.Get("user")
	if userJSON == "" {
		return nil, fmt.Errorf("user data missing")
	}

	var user TelegramUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, err
	}

	return &user, nil
}
