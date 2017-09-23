package modules

type TBL_TXN_STLM_CFG struct {
	INS_ID_CD string `sql:"INS_ID_CD" tname:"TBL_TXN_STLM_CFG"`
	CARD_TP   string
	TXN_NUM   string
	BUS_CD    string
	STLM_FLG  string
	STLM_DESC string
}

func (t TBL_TXN_STLM_CFG) TableName() string {
	return "TBL_TXN_STLM_CFG"
}
