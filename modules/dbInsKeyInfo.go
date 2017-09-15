package modules

import "time"

type DbInsKeyInfo struct {
	InsIdCd    string
	KeyId      string
	RsaPrivKey string
	RsaPubKey  string
	Trk        string
	MakLen     int

	RecUpdTs time.Time
	RecCrtTs time.Time
}

func (t DbInsKeyInfo) TableName() string {
	return "TBL_INS_KEY_INF"
}
