package models

import (
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type UserSession struct {
	ID          uint   `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        User
	UserID      uint   `gorm:"not null"`
	TokenDigest string `gorm:"not null"`
}

var (
	errorLogin = errors.New("incorrect username or password")

	lifetimeTicks = time.Duration(24 * time.Hour)
)

// emailとpasswordが正しければ新しいUserSessionとTokenを返す
func NewSession(email, password string) (*UserSession, string, error) {
	user := &User{Email: email}
	db.Where(user).First(user)

	if user.ID == 0 || !user.IsCorrectPassword(password) {
		return nil, "", errorLogin
	}

	token := []byte(GenerateSecretToken(32))
	digest, _ := bcrypt.GenerateFromPassword(token, GetBcryptCost())

	oldSession := getSession(user.ID)
	if oldSession != nil {
		db.Delete(oldSession)
	}
	session := &UserSession{
		User:        *user,
		TokenDigest: string(digest),
	}
	db.Create(session)

	return session, string(token), nil
}

func CheckLogin(userID uint, token string) *User {
	session := getSession(userID)
	if session == nil {
		return nil
	}
	duration := session.CreatedAt.Sub(time.Now())
	if lifetimeTicks < duration {
		return nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(session.TokenDigest), []byte(token))
	if err != nil {
		return nil
	}

	return &session.User
}

func getSession(userID uint) *UserSession {
	session := &UserSession{UserID: userID}
	db.Where(session).First(session)
	if session.ID == 0 {
		return nil
	}
	session.fetchUser()
	return session
}

// gormの挙動があやしくて、これをしないと動かない(最悪)
func (s *UserSession) fetchUser() {
	s.User.ID = s.UserID
	db.Where(&s.User).First(&s.User)
}
