package modules

type Tbl_mcht_inf struct {
	MCHT_CD               string
	SN                    string
	AIP_BRAN_CD           string
	GROUP_CD              string
	ORI_CHNL              string
	ORI_CHNL_DESC         string
	BANK_BELONG_CD        string
	DVP_BY                string
	MCC_CD_18             string
	APPL_DATE             string
	UP_BC_CD              string
	UP_AC_CD              string
	UP_MCC_CD             string
	NAME                  string
	NAME_BUSI             string
	BUSI_LICE_NO          string
	BUSI_RANG             string
	BUSI_MAIN             string
	CERTIF                string
	CERTIF_TYPE           string
	CERTIF_NO             string
	NATION_CD             string
	PROV_CD               string
	CITY_CD               string
	AREA_CD               string
	REG_ADDR              string
	CONTACT_NAME          string
	CONTACT_PHONENO       string
	ISGROUP               string
	MONEYTOGROUP          string
	STLM_WAY              string
	STLM_WAY_DESC         string
	STLM_INS_CIRCLE       string
	APPR_DATE             string
	STATUS                string
	DELETE_DATE           string
	UC_BC_CD_32           string
	K2WORKFLOWID          string
	SYSTEMFLAG            string
	APPROVALUSERNAME      string
	FINALARRPOVALUSERNAME string
	IS_UP_STANDARD        string
	BILLINGTYPE           string
	BILLINGLEVEL          string
	SLOGAN                string
	EXT1                  string
	EXT2                  string
	EXT3                  string
	EXT4                  string
	AREA_STANDARD         string
	MCHTCD_AREA_CD        string
	UC_BC_CD_AREA         string
	REC_OPR_ID            string
	REC_UPD_OPR           string
	REC_CRT_TS            string
	REC_UPD_TS            string
	OPER_IN               string
	REC_APLLY_TS          string
	OEM_ORG_CODE          string
}

func (t Tbl_mcht_inf) TableName() string {
	return "TBL_MCHT_INF"
}
