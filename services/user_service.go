package services

import (
	"backend/constants"
	"backend/dtos"
	"backend/models"
	"backend/repositories"
	"backend/utils"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthService interface {
	GetUserInfo(accessToken string) (*dtos.GoogleUserInfo, error)
	ProcessGoogleUser(userInfo *dtos.GoogleUserInfo, token *oauth2.Token) (*models.User, error)
	GenerateJWT(user *models.User) (string, time.Time, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) GetUserInfo(accessToken string) (*dtos.GoogleUserInfo, error) {
	userInfoEndpoint := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequest("GET", userInfoEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	var userInfo dtos.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (s *authService) ProcessGoogleUser(userInfo *dtos.GoogleUserInfo, token *oauth2.Token) (*models.User, error) {
	encryptKey64 := os.Getenv("ENCRYPTION_SECRET_KEY")
	if encryptKey64 == "" {
		return nil, errors.New("ENCRYPTION_SECRET_KEY is not set")
	}

	key, err := base64.StdEncoding.DecodeString(encryptKey64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	accessTokenEncrypt, err := utils.EncryptAccessToken(token.AccessToken, key)
	if err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.GetByGoogleID(userInfo.ID)
	if err == nil {
		if existingUser.FolderID == "" {
			folder_id, err := s.CreateFolderId(userInfo.ID)
			if err != nil {
				return nil, err
			}

			existingUser.FolderID = folder_id
		}
		existingUser.Name = userInfo.Name
		existingUser.Email = userInfo.Email
		existingUser.RefreshToken = token.RefreshToken
		existingUser.AccessToken = accessTokenEncrypt
		existingUser.Expiry = token.Expiry
		if err := s.userRepo.Update(existingUser); err != nil {
			return nil, err
		}
		return existingUser, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	folder_id, err := s.CreateFolderId(userInfo.ID)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		GoogleID:     userInfo.ID,
		Email:        userInfo.Email,
		Name:         userInfo.Name,
		RefreshToken: token.RefreshToken,
		AccessToken:  accessTokenEncrypt,
		Expiry:       token.Expiry,
		FolderID:     folder_id,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) GenerateJWT(user *models.User) (string, time.Time, error) {
	expTime := time.Now().Add(time.Hour * 24)

	claims := dtos.Claims{
		UserID:   strconv.FormatUint(uint64(user.ID), 10),
		GoogleID: user.GoogleID,
		FolderID: user.FolderID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	JWTSecretKey := []byte(os.Getenv(constants.JWT_SECRET_KEY))
	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expTime, nil
}

func (s *authService) CreateFolderId(googleId string) (string, error) {
	folderInfoUrl := "http://localhost:8081/storage/folder"
	payload := map[string]interface{}{"title": googleId}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", folderInfoUrl, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create root folder: status %d", resp.StatusCode)
	}

	var userInfo struct {
		Data interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", err
	}

	dataMap, ok := userInfo.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Wrong format")
	}

	folderID, ok := dataMap["id"].(string)
	if !ok {
		return "", fmt.Errorf("Folder id is empty")
	}

	return folderID, nil
}
