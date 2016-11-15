package auth

import (
	"encoding/json"
	"errors"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/happierall/l"
	soso "github.com/happierall/soso-server"
)

func ReadToken(m *soso.Msg, sign string) (string, int64, error) {
	type Other struct {
		Token string `json:"token"`
	}
	other := Other{}
	err := m.ReadOther(&other)
	if err != nil {
		return other.Token, -1, errors.New("Token is does not exists")
	}

	payload, _, err := jose.Decode(other.Token, []byte(sign))
	if err != nil {
		return other.Token, -1, err
	}

	type tokenData struct {
		UID int64
	}
	var td tokenData

	if err = json.Unmarshal([]byte(payload), &td); err != nil {
		return other.Token, -1, err
	}

	return other.Token, td.UID, nil
}

func CreateToken(uid int64, sign string) string {
	payload, err := json.Marshal(map[string]interface{}{
		"UID": uid,
	})
	if err != nil {
		l.Error(err)
		return ""
	}

	token, err := jose.Sign(string(payload), jose.HS256, []byte(sign))
	if err != nil {
		l.Error(err)
		return ""
	}
	return token
}
