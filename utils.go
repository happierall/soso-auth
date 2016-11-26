package auth

import (
	"encoding/json"
	"errors"

	jose "github.com/dvsekhvalnov/jose2go"
	soso "github.com/happierall/soso-server"
)

type TokenData struct {
	UID         int64 `json:"uid"`
	IsAnonymous bool  `json:"isAnonymous"`
}

func ReadToken(m *soso.Msg, sign string) (string, TokenData, error) {
	type Other struct {
		Token string `json:"authToken"`
	}
	other := Other{}
	err := m.ReadOther(&other)
	if err != nil {
		return other.Token, TokenData{}, errors.New("Token is does not exists")
	}

	if other.Token == "" {
		return "", TokenData{}, errors.New("Token empty")
	}

	payload, _, err := jose.Decode(other.Token, []byte(sign))
	if err != nil {
		return other.Token, TokenData{}, err
	}

	var td TokenData

	if err = json.Unmarshal([]byte(payload), &td); err != nil {
		return other.Token, TokenData{}, err
	}

	return other.Token, td, nil
}

func CreateToken(data map[string]interface{}, sign string) string {
	payload, err := json.Marshal(data)
	if err != nil {
		Log.Error(err)
		return ""
	}

	token, err := jose.Sign(string(payload), jose.HS256, []byte(sign))
	if err != nil {
		Log.Error(err)
		return ""
	}
	return token
}
