package authn

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AkashGit21/task-ms/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const TOKEN_EXPIRY_DURATION time.Duration = 24 * time.Hour

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

/** Implement basic username / password authentication **/
func (anh *authnHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := anh.AuthnOps.FetchActiveRecord(req.Username)
	if err != nil || user == nil {
		utils.ErrorLog("Encountered error while fetching user record", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(req.Password))
	if err != nil {
		utils.ErrorLog("Hash comparison failed", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	currentTime := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.UserClaims{
		Username: user.Username,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			Subject:   user.Username,
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(TOKEN_EXPIRY_DURATION)),
		},
	})

	tokenString, err := token.SignedString(anh.JWTSecret)
	if err != nil {
		utils.ErrorLog("Failed to generate token", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := loginResponse{
		Token:     tokenString,
		ExpiresAt: currentTime.Add(TOKEN_EXPIRY_DURATION),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
