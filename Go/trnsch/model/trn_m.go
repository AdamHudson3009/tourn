package model

//Payloads

type Response struct {
	Msg string `json:"msg"`
}

type League struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Abc         string `json:"abc"`
}

type Trn struct {
	Id int `json:"id"`
}

type TrnBuild struct {
	Id    int `json:"id"`
	Build int `json:"build"`
}

type TrnSch struct {
	Id  int `json:"id"`
	Rnd int `json: "rnd"`
}

type SchSwap struct {
	TrnId int  `json:"tourn_id"`
	Rnd   int  `json:"rnd"`
	Gid   int  `json:"gid"`
	Up    bool `json:"up"`
}

type SchEdit struct {
	Gid    int    `json:"gid"`
	Plyrwn int    `json:"plyrwn"`
	Pf     int    `json:"pf"`
	Pa     int    `json:"pa"`
	Notes  string `json:"notes"`
}

/////////////////////////////////////////////////
//// Sch Res Below
/////////////////////////////////////////////////////

type SchResPlyrtm struct {
	Rnk     int      `json:"rnk"`
	Plyrtm  string   `json:"plyrtm"`
	Columns []string `json:"columns"`
}

type SchResDivi struct {
	Divi    string         `json:"divi"`
	Plyrtms []SchResPlyrtm `json:"plyrtms"`
	Headers []string       `json:"headers"`
}

type SchResConf struct {
	Conf  string       `json:"conf"`
	Divis []SchResDivi `json:"divis"`
}

type PlyrtmOOO struct {
	Plyrtm string   `json:"plyrtm"`
	Wn     []string `json:"wn"`
	Ls     []string `json:"ls"`
	Tie    []string `json:"tie"`
}

type DiviOOO struct {
	Divi      string      `json:"divi"`
	PlyrtmOoo []PlyrtmOOO `json:"plyrtm_ooo"`
}

type ConfOOO struct {
	Conf    string    `json:"conf"`
	DiviOoo []DiviOOO `json:"divi_ooo"`
}

type RnkStnd struct {
	Rnk     int      `json:"rnk"`
	Name    string   `json:"name"`
	Rec     string   `json:"rec"`
	GFields []string `json:"g_fields"`
	DvRsn   []string `json:"dv_rsn"`
	CnRsn   []string `json:"cn_rsn"`
}

type ConfStnd struct {
	Conf      string    `json:"conf"`
	CNRnkStnd []RnkStnd `json:"cn"`
	WCRnkStnd []RnkStnd `json:"wc"`
	DVRnkStnd []RnkStnd `json:"dv"`
}

type PlyrtmGMS struct {
	Plyrtm string    `json:"plyrtm"`
	Wn     float64   `json:"wn"`
	Wns    []float64 `json:"wns"`
	Remain []string  `json:"remain"`
	W      []string  `json:"w"`
	L      []string  `json:"l"`
	T      []string  `json:"t"`
}

type ConfGMS struct {
	Conf      string      `json:"conf"`
	PlyrtmGms []PlyrtmGMS `json:"plyrtm_gms"`
}

type SchRes struct {
	Confs    []SchResConf `json:"confs"`
	ConfOoo  []ConfOOO    `json:"conf_ooo"`
	ConfStnd []ConfStnd   `json:"conf_stnd"`
	ConfGms  []ConfGMS    `json:"conf_gms"`
}

//************************************************************************8

type PlyrtmSch struct {
	PlyrtmId int   `json:"plyrtm_id"`
	Plyrtms  []int `json:"plyrtms"`
}

type Divi struct {
	DiviId     int               `json:"divi_id"`
	Divi       string            `json:"divi"`
	PlyrtmsSch map[int]PlyrtmSch `json:"plyrtms_sch"`
}

type Conf struct {
	ConfId int          `json:"conf_id"`
	Conf   string       `json:"conf"`
	Divis  map[int]Divi `json:"divis"`
}

type TrnTree struct {
	Confs   map[int]Conf       `json:"confs"`
	Plyrtms map[int]PlyrtmInfo `json:"plyrtm_info"`
}

type PlyrtmInfo struct {
	PlyrtmId int    `json:"plyrtm_id"`
	Plyrtm   string `json:"plyrtm"`
	Cid      int    `json:"cid"`
	Conf     string `json:"conf"`
	Did      int    `json:"did"`
	Divi     string `json:"divi"`
	Sds      int    `json:"sds"`
}

type Sch struct {
	Gid      int      `json:"gid"`
	Plyr1    int      `json:"plyr1"`
	Plyr2    int      `json:"plyr2"`
	Ordr     int      `json:"ordr"`
	Plyrwn   int      `json:"plyrwn"`
	Pf       int      `json:"pf"`
	Pa       int      `json:"pa"`
	Plyr1nm  string   `json:"plyr1nm"`
	Plyr2nm  string   `json:"plyr2nm"`
	Plyrwnnm string   `json:"plyrwnnm"`
	Rec1     int      `json:"rec1"`
	Rec2     int      `json:"rec2"`
	Tie1     int      `json:"tie1"`
	Tie2     int      `json:"tie2"`
	Notes    string   `json:"notes"`
	Msc1     []string `json:"msc1"`
	Msc2     []string `json:"msc2"`
	Letter1  string   `json:"letter1"`
	Letter2  string   `json:"letter2"`
	ArchNm   string   `json:"arch_nm"`
	ArchPts  string   `json:"arch_pts"`
}

type AddDeleteMsc struct {
	TrnId     int    `json:"tourn_id"`
	Rnd       int    `json:"rnd"`
	Plyrtm    int    `json:"plyrtm_id"`
	Letters   string `json:"letters"`
	AddDelete int    `json:"add_delete"`
}

type Import struct {
	TrnId    int    `json:"tourn_id"`
	Filename string `json:"filename"`
}
