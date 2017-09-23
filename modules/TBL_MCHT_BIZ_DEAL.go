package modules

import "time"

type TBL_MCHT_BIZ_DEAL struct {
	MCHT_CD     string `sql:"MCHT_CD"`
	PROD_CD     string  `sql:"PROD_CD"`
	BIZ_CD      string `sql:"BIZ_CD"`
	TRANS_CD    string `sql:"TRANS_CD"`
	OPER_IN     string `sql:"OPER_IN"`
	REC_OPR_ID  string `sql:"REC_OPR_ID"`
	REC_UPD_OPR string `sql:"REC_UPD_OPR"`
	REC_CRT_TS  time.Time `sql:"REC_CRT_TS"`
	REC_UPD_TS  time.Time `sql:"REC_UPD_TS"`
}

func (t TBL_MCHT_BIZ_DEAL)TableName()string{
	return "TBL_MCHT_BIZ_DEAL"
}