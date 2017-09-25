package modules

import "time"

type TBL_MCHT_INF struct {
	MCHT_CD               string `json:"mchtCd"`
	SN                    string `json:"sn"`
	AIP_BRAN_CD           string `json:"aipBranCd,omitempty"`
	GROUP_CD              string `json:"groupCd"`
	ORI_CHNL              string `json:"oriChnl,omitempty"`
	ORI_CHNL_DESC         string `json:"oriChnlDesc,omitempty"`
	BANK_BELONG_CD        string `json:"insIdCd"` //insIdCd bankBelongCd
	DVP_BY                string `json:"managerName"`
	MCC_CD_18             string `json:"mccCd18,omitempty"`
	APPL_DATE             time.Time `json:"applDate,omitempty"`
	UP_BC_CD              string `json:"upBcCd,omitempty"`
	UP_AC_CD              string `json:"up_ac_cd,omitempty"`
	UP_MCC_CD             string `json:"up_mcc_cd,omitempty"`
	NAME                  string `json:"name,omitempty"`
	NAME_BUSI             string `json:"nameBusi"`
	BUSI_LICE_NO          string `json:"busiLiceNo"`
	BUSI_RANG             string `json:"busiRang"`
	BUSI_MAIN             string `json:"busiMain,omitempty"`
	CERTIF                string `json:"policyholdersName"`
	CERTIF_TYPE           string `json:"certifType"`
	CERTIF_NO             string `json:"certifNo"`
	NATION_CD             string `json:"nation_cd,omitempty"`
	PROV_CD               string `json:"provCd"`
	CITY_CD               string `json:"cityCd"`
	AREA_CD               string `json:"countyCd"`
	REG_ADDR              string `json:"regAddr"`
	CONTACT_NAME          string `json:"contact"`
	CONTACT_PHONENO       string `json:"mobile"`
	ISGROUP               string `json:"isgroup,omitempty"`
	MONEYTOGROUP          string `json:"moneytogroup,omitempty"`
	STLM_WAY              string `json:"stlm_way,omitempty"`
	STLM_WAY_DESC         string `json:"stlm_way_desc,omitempty"`
	STLM_INS_CIRCLE       string `json:"stlm_ins_circle,omitempty"`
	APPR_DATE             time.Time `json:"appr_date,omitempty"`
	STATUS                string `json:"status"`
	DELETE_DATE           time.Time `json:"delete_date,omitempty"`
	UC_BC_CD_32           string `json:"uc_bc_cd_32,omitempty"`
	K2WORKFLOWID          string `json:"k_2_workflowid,omitempty"`
	SYSTEMFLAG            string `json:"systemflag"`
	APPROVALUSERNAME      string `json:"approvalusername,omitempty"`
	FINALARRPOVALUSERNAME string `json:"finalarrpovalusername,omitempty"`
	IS_UP_STANDARD        string `json:"is_up_standard,omitempty"`
	BILLINGTYPE           string `json:"billingtype,omitempty"`
	BILLINGLEVEL          string `json:"billinglevel,omitempty"`
	SLOGAN                string `json:"slogan,omitempty"`
	EXT1                  string `json:"ext_1,omitempty"`
	EXT2                  string `json:"ext_2,omitempty"`
	EXT3                  string `json:"ext_3,omitempty"`
	EXT4                  string `json:"ext_4,omitempty"`
	AREA_STANDARD         string `json:"area_standard,omitempty"`
	MCHTCD_AREA_CD        string `json:"mchtcdAreaCd"`
	UC_BC_CD_AREA         string `json:"uc_bc_cd_area,omitempty"`
	REC_OPR_ID            string `json:"rec_opr_id"`
	REC_UPD_OPR           string `json:"recUpdOpr"`
	REC_CRT_TS            time.Time `xorm:"created"`
	REC_UPD_TS            time.Time `xorm:"updated"`
	OPER_IN               string    `json:"recOprId,omitempty"`
	REC_APLLY_TS          time.Time `json:"rec_aplly_ts,omitemty"`
	OEM_ORG_CODE          string `json:"oem_org_code,omitemty"`
}

func (d TBL_MCHT_INF) TableName() string {
	return "TBL_MCHT_INF"
}
