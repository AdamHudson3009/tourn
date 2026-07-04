package trn

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"trnsch/db"
	"trnsch/model"
)

type PlrStr struct {
	trn int
	cn  int
	dv  int
	cm1 int
	cm2 int
	sds int
	sch []string
}

var plrs map[int]PlrStr
var nme_attr_int map[string]int
var nme_attr_flt map[string]float64
var nme_attr_tme map[string]string

// var cn_to_dv map[string][]string
var cndv_to_plrs map[string][]string
var plr_plr_w map[string][]string
var plr_plr_l map[string][]string
var plr_plr_t map[string][]string
var plr_rsn_dv map[string][]string
var plr_rsn_cn map[string][]string
var plyr_rnk_dsc map[string]string
var gme_mp map[string]string
var gme_mp_rev map[string]string
var is_gme_mp bool
var wkplyr_depc map[string]string
var curplyr_depc map[string]string
var plr_plr_gms_w map[string]string
var plr_plr_gms_l map[string]string
var plr_plr_gms_t map[string]string
var g_string map[string]string
var g_string_ord []string
var plyr_stnd map[string]int
var plyr_stnd_rplc map[string]string
var plyr_extd_sch map[string]string
var plyr_sds map[string]int

var clsgm_hsh map[string]string
var strk_hsh map[string]string

var trn string
var LetterArr []string
var LetterArrIndx int
var LetterArrMax int
var ply_wns map[string]int
var potw map[string]int
var HasAnArch bool
var tbrkRnk bool
var arch_wk_yt bool

var cnHash map[int]string
var dvHash map[int]string
var plyrtmHash map[int]string
var plyr_id_to_name map[int]string
var plyr_id_arch map[int]int
var plyr_id_arch_w map[int]int
var plrsNameToId map[string]int

var stnds model.SchRes

func empty() {
	plrs = nil
	nme_attr_int = nil
	nme_attr_flt = nil
	nme_attr_tme = nil

	// var cn_to_dv map[string][]string
	cndv_to_plrs = nil
	plr_plr_w = nil
	plr_plr_l = nil
	plr_plr_t = nil
	plr_rsn_dv = nil
	plr_rsn_cn = nil
	plyr_rnk_dsc = nil
	gme_mp = nil
	gme_mp_rev = nil
	is_gme_mp = false
	wkplyr_depc = nil
	curplyr_depc = nil
	plr_plr_gms_w = nil
	plr_plr_gms_l = nil
	plr_plr_gms_t = nil
	g_string = nil
	g_string_ord = nil
	plyr_stnd = nil
	plyr_stnd_rplc = nil
	plyr_extd_sch = nil
	plyr_sds = nil

	clsgm_hsh = nil
	strk_hsh = nil

	trn = ""
	LetterArr = nil
	LetterArrIndx = 0
	LetterArrMax = 0
	ply_wns = nil
	potw = nil
	HasAnArch = false
	tbrkRnk = false
	arch_wk_yt = false

	cnHash = nil
	dvHash = nil
	plyrtmHash = nil
	plyr_id_to_name = nil
	plyr_id_arch = nil
	plyr_id_arch_w = nil
	plrsNameToId = nil
	stnds = model.SchRes{}
}

func BuildSchres(trnId int, rnd int) ([]byte, error) {
	empty()
	is_gme_mp = false
	plrs = make(map[int]PlrStr)
	nme_attr_int = make(map[string]int)
	nme_attr_flt = make(map[string]float64)
	nme_attr_tme = make(map[string]string)
	//cn_to_dv = make(map[string][]string)
	cndv_to_plrs = make(map[string][]string)
	plr_plr_w = make(map[string][]string)
	plr_plr_l = make(map[string][]string)
	plr_plr_t = make(map[string][]string)
	plr_rsn_dv = make(map[string][]string)
	plr_rsn_cn = make(map[string][]string)
	plyr_rnk_dsc = make(map[string]string)
	wkplyr_depc = make(map[string]string)
	curplyr_depc = make(map[string]string)
	plr_plr_gms_w = make(map[string]string)
	plr_plr_gms_l = make(map[string]string)
	plr_plr_gms_t = make(map[string]string)
	g_string = make(map[string]string)
	plyr_stnd = make(map[string]int)
	plyr_stnd_rplc = make(map[string]string)
	plyr_extd_sch = make(map[string]string)
	plyr_sds = make(map[string]int)
	clsgm_hsh = make(map[string]string)
	strk_hsh = make(map[string]string)
	potw = make(map[string]int)
	plyr_id_to_name = make(map[int]string)
	plyr_id_arch = make(map[int]int)
	plyr_id_arch_w = make(map[int]int)
	arch_wk_yt = false
	HasAnArch = false

	cnHash = make(map[int]string)
	dvHash = make(map[int]string)
	plyrtmHash = make(map[int]string)

	var dv_tbr_str string
	var cn_tbr_str string
	var hdr_spr_ln []string
	var err error

	wk := 0
	maxwk := 0
	//var plyr string
	//plyr1_c := 0
	//plyr2_c := 1
	//arw_c := -1
	//w_c := 2
	gn_c := -1
	gn_is_tbkr := false
	ot := 0
	ot_is_tbkr := false
	sht_is_tbkr := false
	pts_c := -1
	pts_is_tbkr := false
	//nt := 0
	//nt_c := -1
	//nt_is_tbkr := false
	tme_is_tbkr := false
	cndv_wc := ""
	is_sovors := false
	is_clsgm := false
	is_strk := false
	is_rnk := false
	ptsf_c := -1
	ptsa_c := -1
	wk_done := 0
	plrsNameToId = make(map[string]int)

	if LetterArrIndx < 2 {
		ply_wns = make(map[string]int)
	}

	sql := "select dv_cn_other, description " +
		"from tbrk " +
		"where tourn_id = ? " +
		"order by dv_cn_other, rnk"
	rows, err := db.DB.Query(sql, trnId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dv_cn_oth, description string
	for rows.Next() {
		if err = rows.Scan(&dv_cn_oth, &description); err != nil {
			return nil, err
		}
		if dv_cn_oth == "dv" {
			dv_tbr_str += ";" + description
			if description == "arch" {
				HasAnArch = true
			}
		} else if dv_cn_oth == "cn" {
			cn_tbr_str += ";" + description
		} else if dv_cn_oth == "hdr-spr-ln" {
			hdr_spr_ln = append(hdr_spr_ln, description)
		} else if dv_cn_oth == "wc" {
			cndv_wc += " " + description
		}
	}
	if cndv_wc == "" {
		cndv_wc = "cn 4"
	}

	dv_tbr_str = dv_tbr_str[1:]
	cn_tbr_str = cn_tbr_str[1:]
	cndv_wc = cndv_wc[1:]

	for hdr_i := 0; hdr_i < len(hdr_spr_ln); hdr_i++ {
		switch hdr_spr_ln[hdr_i] {
		/*case "T1":
			plyr1_c = hdr_i
		case "T2":
			plyr2_c = hdr_i
		case "W":
			w_c = hdr_i
		case "Nt":
			nt_c = hdr_i
		*/
		case "GN":
			gn_c = hdr_i
			if strings.Contains(cn_tbr_str, ";gn-pt") || strings.Contains(dv_tbr_str, ";gn-pt") {
				gn_is_tbkr = true
			}
		case "Pt":
			pts_c = hdr_i
			if strings.Contains(cn_tbr_str, ";pts-") || strings.Contains(dv_tbr_str, ";pts-") {
				pts_is_tbkr = true
			}
		case "Pf":
			ptsf_c = hdr_i
		case "Pa":
			ptsa_c = hdr_i
			//case "Rec":
			//	arw_c = hdr_i + 1
		}
	}

	sql = "select p.id pid, p2.description " +
		"from plyrtm p inner join divi d on p.divi_id = d.id " +
		"inner join conf c on d.conf_id = c.id " +
		"inner join gms_plyrtm gp1 on gp1.plyrtm_id = p.id " +
		"inner join gms g on gp1.gms_id = g.id and coalesce(g.has_stnd,'Y') = 'Y' " +
		"inner join gms_plyrtm gp2 on gp2.gms_id = g.id and gp1.plyrtm_id <> gp2.plyrtm_id " +
		"inner join plyrtm p2 on p2.id = gp2.plyrtm_id " +
		"where c.tourn_id = ? " +
		"order by c.description, d.description, p.sds, g.rnd, gp2.plyrtm_id"
	rows, err = db.DB.Query(sql, trnId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	plyr1 := make(map[int]string)
	var cid, did, pid, sds, conf_num, i_num, prev_pid int
	var conf, divi, plyrtm string
	for rows.Next() {
		if err = rows.Scan(&pid, &plyrtm); err != nil {
			return nil, err
		}
		plyr1[pid] += ";" + plyrtm
	}

	sql = "select c.id, d.id, p.id, p.sds, coalesce(cmm.conf_num,0) conf_num, coalesce(cmm.i_num,0) i_num, a.maxRnd, " +
		"c.description, d.description, p.description " +
		"from plyrtm p inner join divi d on p.divi_id = d.id " +
		"inner join conf c on d.conf_id = c.id " +
		"inner join (select max(rnd) maxRnd from gms where tourn_id = ? and coalesce(has_stnd,'Y') = 'Y') a " +
		"left join cmm on cmm.tourn_id = c.tourn_id and d.id = cmm.divi_id " +
		"where c.tourn_id = ? " +
		"order by c.id, d.id, p.sds "
	rows, err = db.DB.Query(sql, trnId, trnId)
	if err != nil {
		fmt.Println("Query err", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&cid, &did, &pid, &sds, &conf_num, &i_num, &maxwk,
			&conf, &divi, &plyrtm); err != nil {
			fmt.Println("abc", err)
			return nil, err
		}
		plyr_id_to_name[pid] = plyrtm
		if sds%2 == 0 {
			plyr_id_arch[pid] = prev_pid
			plyr_id_arch[prev_pid] = pid
			prev_pid = plyr_id_arch[pid]
		}
		var plrs_tmp PlrStr
		plrs_tmp.trn = trnId
		plrs_tmp.cn = cid
		plrs_tmp.dv = did
		if len(plyr1[pid]) > 1 {
			plrs_tmp.sch = strings.Split(plyr1[pid][1:], ";")
		}
		plrs_tmp.cm1 = conf_num
		plrs_tmp.cm2 = i_num
		plrs_tmp.sds = sds
		plrs[pid] = plrs_tmp
		plrsNameToId[plyrtm] = pid
		cndv_to_plrs[conf+";"+divi] = append(cndv_to_plrs[conf+";"+divi], plyrtm)
		prev_pid = pid
	}
	plyr1 = nil
	/*for line_i := 0; line_i < len(lines); line_i++ {
		if regexp.MustCompile(`^arch;\d+;\d+`).MatchString(l) {
			fields := strings.Split(l, ";")
			arch_wk1 = StrToInt(fields[1])
			arch_wk2 = StrToInt(fields[2])

		}
	}*/

	if strings.Contains(cn_tbr_str, ";ot") || strings.Contains(dv_tbr_str, ";ot") {
		if strings.Contains(cn_tbr_str, ";ot-pt") || strings.Contains(dv_tbr_str, ";ot-pt") {
			ot_is_tbkr = true
		}
	}
	if strings.Contains(cn_tbr_str, ";sht") || strings.Contains(dv_tbr_str, ";sht") {
		if strings.Contains(cn_tbr_str, ";sht-") || strings.Contains(dv_tbr_str, ";sht-") {
			sht_is_tbkr = true
		}
	}
	if strings.Contains(cn_tbr_str, ";tme-pt") || strings.Contains(dv_tbr_str, ";tme-pt") {
		tme_is_tbkr = true
	}
	/*if strings.Contains(cn_tbr_str, ";nt-pt") || strings.Contains(dv_tbr_str, ";nt-pt") {
		nt_is_tbkr = true
	}*/
	if strings.Contains(cn_tbr_str, ";sov") || strings.Contains(dv_tbr_str, ";sov") || strings.Contains(cn_tbr_str, ";sos") || strings.Contains(dv_tbr_str, ";sos") {
		is_sovors = true
	}
	if strings.Contains(cn_tbr_str, ";rnk") || strings.Contains(dv_tbr_str, ";rnk") {
		is_rnk = true
	}
	if strings.Contains(cn_tbr_str, ";clsgm") || strings.Contains(dv_tbr_str, ";clsgm") {
		is_clsgm = true
	}
	if strings.Contains(cn_tbr_str, ";strk") || strings.Contains(dv_tbr_str, ";strk") {
		is_strk = true
	}

	fmt.Println("HEreAAAsdsddsdAAAA")

	sql = "select rnd, coalesce(letters,''), plyrtm_id " +
		"from msc where tourn_id = ? and rnd <= ? " +
		"order by rnd, letters"
	rows, err = db.DB.Query(sql, trnId, rnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tempRnd int
	var depc, pname string
	for rows.Next() {
		if err = rows.Scan(&tempRnd, &depc, &pid); err != nil {
			return nil, err
		}

		pname = plyr_id_to_name[pid]
		if depc == "potw" {
			potw[pname]++
		} else if (depc != "h") && (depc != "b") && (depc != "d") && (depc != "p") && (depc != "c") && (depc != "r") && (depc != "x") && (depc != "e") && (depc != "ee") && (depc != "a") {
			continue
		} else {
			if (depc == "d") || (depc == "e") || (depc == "b") || (depc == "h") || (depc == "ee") || (depc == "a") { //(depc != "r") ||
				wkplyr_depc[strconv.Itoa(rnd)+";"+pname] = depc
			} else if (wkplyr_depc[strconv.Itoa(tempRnd)+";"+pname] == "") || ((wkplyr_depc[strconv.Itoa(tempRnd)+";"+pname] != "d") && (wkplyr_depc[strconv.Itoa(tempRnd)+";"+pname] != "e") && (wkplyr_depc[strconv.Itoa(tempRnd)+";"+pname] != "ee")) {
				if (depc == "p") || (depc == "c") || (depc == "x") || (depc == "r") {
					wkplyr_depc[strconv.Itoa(rnd)+";"+pname] += depc
				}
			}
		}
	}
	/*if n == "wc" {
		wc_attr_only = true
	}

	if wc_attr_only {
		//attr_wk_assn(dir, wk, wk == maxwk, schxls_fbase, attr_w, attr_blw, true, attr_orig)
		return
	}*/

	sorted_cndv := SortNamesInMap(cndv_to_plrs)

	dv_tbr_arr := strings.Split(dv_tbr_str, ";")
	//dv_tbr_arr = dv_tbr_arr[1:]
	dv_tbr_bsc := tbr_bsc(dv_tbr_arr)
	dv_tbr_expnd := tbr_exp(dv_tbr_arr, dv_tbr_bsc)
	cn_tbr_arr := strings.Split(cn_tbr_str, ";")
	//cn_tbr_arr = cn_tbr_arr[1:]
	cn_tbr_bsc := tbr_bsc(cn_tbr_arr)
	cn_tbr_expnd := tbr_exp(cn_tbr_arr, cn_tbr_bsc)
	var sht_mp_arr []string
	running_sht := ""

	sql = "select rnd, coalesce(plyrtm_wn,-1), min(gp1.plyrtm_id) plyr1, " +
		"max(gp2.plyrtm_id) plyr2, coalesce(pf,0), coalesce(pa,0) " +
		"from gms g inner join gms_plyrtm gp1 on gp1.gms_id = g.id " +
		"inner join gms_plyrtm gp2 on gp2.gms_id = g.id and gp2.plyrtm_id <> gp1.plyrtm_id " +
		"where tourn_id = ? and rnd <= ? and coalesce(g.has_stnd,'Y') = 'Y' " +
		"group by rnd, coalesce(plyrtm_wn,-1), coalesce(pf,0), coalesce(pa,0) " +
		"order by rnd"
	rows, err = db.DB.Query(sql, trnId, rnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plw, pid1, pid2, pf, pa int
	for rows.Next() {
		if err = rows.Scan(&wk, &plw, &pid1, &pid2, &pf, &pa); err != nil {
			return nil, err
		}
		if plw == -1 {
			if wk < rnd {
				resp := model.Response{Msg: "prev week not done"}
				jsonOut, _ := json.MarshalIndent(resp, "", "  ")
				return jsonOut, nil
			}
			continue
		}

		plyr1 := plyr_id_to_name[pid1]
		plyr2 := plyr_id_to_name[pid2]
		plyrw := plyr_id_to_name[plw]

		pts := 0
		tme := "0:00"
		sht := 0
		gn := 0

		var w_w, l_l, plyrl string
		var is_t bool
		ptsf := 0
		ptsa := 0
		if plyr1 == plyrw {
			plyrl = plyr2
			is_t = false
			w_w = plyr2
			l_l = plyr1
		} else if plyr2 == plyrw {
			plyrl = plyr1
			l_l = plyr2
			w_w = plyr1
			is_t = false
		} else if plw == 0 {
			is_t = true
			inc_nme_attr_int(plyr1+";t;o", 1)
			inc_nme_attr_int(plyr2+";t;o", 1)
			w_w = plyr1
			l_l = plyr2
			if ptsf_c > 0 {
				w_w += fmt.Sprintf("0,%d,%d", pf, pf)
				l_l += fmt.Sprintf("0,%d,%d", pf, pf)
				ptsf = pf
				ptsa = ptsf
				inc_nme_attr_int(plyr1+";pf;", ptsf)
				inc_nme_attr_int(plyr1+";pa;", -ptsf)
				inc_nme_attr_int(plyr2+";pf;", ptsf)
				inc_nme_attr_int(plyr2+";pa;", -ptsf)
			}
		}
		if !is_t {
			inc_nme_attr_int(plyrw+";w;o", 1)
			inc_nme_attr_int(plyrl+";l;o", 1)
			if gn_c > 0 {
				tme = "0.0" //for now since there is no tme
				if tme_f, err := strconv.ParseFloat(tme, 64); err == nil {
					tme = FloatToTme(tme_f)
				}
				if tme == "0:00" {
					gn = 0
				} else {
					gn = 1
				}
				tmps := strconv.Itoa(gn)
				w_w += "," + tmps
				l_l += "," + tmps
			}
			if pts_c > 0 {
				pts = pf - pa
				if (pts == 10) && tme != "0:00" {
					sht = 1
				} else {
					sht = 0
				}
				if (pts == 0) && tme == "0:00" {
					ot = -1
				} else {
					ot = 0
				}
			}
			if ptsf_c > 0 {
				w_w += "," + fmt.Sprintf("%d", pf)
				l_l += "," + fmt.Sprintf("%d", pf)
				ptsf = pf
				inc_nme_attr_int(plyrw+";pf;", ptsf)
				inc_nme_attr_int(plyrl+";pa;", -ptsf)
			}
			if ptsa_c > 0 {
				w_w += "," + fmt.Sprintf("%d", pa)
				l_l += "," + fmt.Sprintf("%d", pa)
				ptsa = pa
				inc_nme_attr_int(plyrw+";pa;", -ptsa)
				inc_nme_attr_int(plyrl+";pf;", ptsa)
			}
			if (ptsf > 0) && (ptsa == 0) {
				sht = 1
			} else {
				sht = 0
			}
			wlt_tbrkr_i(gn_c > -1, plyrw, plyrl, "gn", "", gn)
			wlt_tbrkr_t(gn_c > -1, plyrw, plyrl, "tme", "", tme)
			wlt_tbrkr_i(true, plyrw, plyrl, "ot", "", ot)
			wlt_tbrkr_i(true, plyrw, plyrl, "sht", "", sht)
			wlt_tbrkr_i(true, plyrw, plyrl, "pts", "", pts)
			if is_clsgm {
				add_clsgm(plyrw, pts, -1)
				add_clsgm(plyrl, pts, 1)
			}
			if is_strk {
				add_strk(plyrw, "W")
				add_strk(plyrl, "L")
			}
		}
		/*if nt_c > 0 {
			cell, _ = sheet.Cell(row, nt_c)
			w_w += cell.Value
			l_l += cell.Value
			if cell.Value == "" {
				nt = 0
			} else {
				nt = StrToInt(cell.Value)
			}
			wlt_tbrkr_i(true, plyrw, plyrl, "nt", "", nt)
		}*/
		if !is_t {
			plr_plr_w[plyrw] = append(plr_plr_w[plyrw], w_w)
			plr_plr_l[plyrl] = append(plr_plr_l[plyrl], l_l)
		} else {
			plr_plr_t[plyr2] = append(plr_plr_t[plyr2], w_w)
			plr_plr_t[plyr1] = append(plr_plr_t[plyr1], l_l)
		}
		if HasAnArch && plyr_id_arch[pid1] == pid2 {
			arch_wk_yt = true
			inc_nme_attr_int(plyrw+";w;arch", 1)
			inc_nme_attr_int(plyrl+";l;arch", 1)
			wlt_tbrkr_i(gn_is_tbkr, plyrw, plyrl, "gn", "arch", gn)
			wlt_tbrkr_t(tme_is_tbkr, plyrw, plyrl, "tme", "arch", tme)
			wlt_tbrkr_i(ot_is_tbkr, plyrw, plyrl, "ot", "arch", ot)
			wlt_tbrkr_i(sht_is_tbkr, plyrw, plyrl, "sht", "arch", sht)
			wlt_tbrkr_i(pts_is_tbkr, plyrw, plyrl, "pts", "arch", pts)
			var pidw, pidl int
			if plyr_id_to_name[pid1] == plyrw {
				pidw = pid1
				pidl = pid2
			} else {
				pidw = pid2
				pidl = pid1
			}
			plyr_id_arch_w[pidw] += pts
			plyr_id_arch_w[pidl] -= pts
			if wk == maxwk {
				if plyr_id_arch_w[pidw] > 0 {
					nme_attr_int[plyrw+";arw"] = 1
					nme_attr_int[plyrl+";arw"] = -1
				} else if plyr_id_arch_w[pidl] > 0 {
					nme_attr_int[plyrw+";arw"] = -1
					nme_attr_int[plyrl+";arw"] = 1
				} else {
					//separate tiebreaker
				}
			}
		}
		if plrs[pid1].cn == plrs[pid2].cn {
			if !is_t {
				inc_nme_attr_int(plyrw+";w;cn", 1)
				inc_nme_attr_int(plyrl+";l;cn", 1)
				wlt_tbrkr_i(gn_is_tbkr, plyrw, plyrl, "gn", "cn", gn)
				wlt_tbrkr_t(tme_is_tbkr, plyrw, plyrl, "tme", "cn", tme)
				wlt_tbrkr_i(ot_is_tbkr, plyrw, plyrl, "ot", "cn", ot)
				wlt_tbrkr_i(sht_is_tbkr, plyrw, plyrl, "sht", "cn", sht)
				wlt_tbrkr_i(pts_is_tbkr, plyrw, plyrl, "pts", "cn", pts)
				if ptsf_c > 0 {
					inc_nme_attr_int(plyrw+";pf;cn", ptsf)
					inc_nme_attr_int(plyrl+";pa;cn", -ptsf)
				}
				if ptsa_c > 0 {
					inc_nme_attr_int(plyrw+";pa;cn", -ptsa)
					inc_nme_attr_int(plyrl+";pf;cn", ptsa)
				}
			} else {
				inc_nme_attr_int(plyr1+";t;cn", 1)
				inc_nme_attr_int(plyr2+";t;cn", 1)
				if ptsf_c > 0 {
					inc_nme_attr_int(plyr1+";pf;cn", ptsf)
					inc_nme_attr_int(plyr1+";pa;cn", -ptsf)
					inc_nme_attr_int(plyr2+";pf;cn", ptsf)
					inc_nme_attr_int(plyr2+";pa;cn", -ptsf)
				}
			}
			//wlt_tbrkr_i(nt_is_tbkr, plyrw, plyrl, "nt", "cn", nt)
			if plrs[pid1].dv == plrs[pid2].dv {
				if !is_t {
					inc_nme_attr_int(plyrw+";w;dv", 1)
					inc_nme_attr_int(plyrl+";l;dv", 1)
					wlt_tbrkr_i(gn_is_tbkr, plyrw, plyrl, "gn", "dv", gn)
					wlt_tbrkr_t(tme_is_tbkr, plyrw, plyrl, "tme", "dv", tme)
					wlt_tbrkr_i(ot_is_tbkr, plyrw, plyrl, "ot", "dv", ot)
					wlt_tbrkr_i(sht_is_tbkr, plyrw, plyrl, "sht", "dv", sht)
					wlt_tbrkr_i(pts_is_tbkr, plyrw, plyrl, "pts", "dv", pts)
				} else {
					inc_nme_attr_int(plyr1+";t;dv", 1)
					inc_nme_attr_int(plyr2+";t;dv", 1)
				}
				//wlt_tbrkr_i(nt_is_tbkr, plyrw, plyrl, "nt", "dv", nt)
			}
			if (plrs[pid1].dv == plrs[pid2].dv) || ((plrs[pid1].cm1 == plrs[pid2].cm1) && (plrs[pid1].cm1 > 0)) {
				if !is_t {
					inc_nme_attr_int(plyrw+";w;cm", 1)
					inc_nme_attr_int(plyrl+";l;cm", 1)
					wlt_tbrkr_i(gn_is_tbkr, plyrw, plyrl, "gn", "cm", gn)
					wlt_tbrkr_t(tme_is_tbkr, plyrw, plyrl, "tme", "cm", tme)
					wlt_tbrkr_i(ot_is_tbkr, plyrw, plyrl, "ot", "cm", ot)
					wlt_tbrkr_i(sht_is_tbkr, plyrw, plyrl, "sht", "cm", sht)
					wlt_tbrkr_i(pts_is_tbkr, plyrw, plyrl, "pts", "cm", pts)
				} else {
					inc_nme_attr_int(plyr1+";w;cm", 1)
					inc_nme_attr_int(plyr2+";l;cm", 1)
				}
				//wlt_tbrkr_i(nt_is_tbkr, plyrw, plyrl, "nt", "cm", nt)
			}
		} else if (plrs[pid1].cm2 == plrs[pid2].cm2) && (plrs[pid1].cm2 > 0) {
			if !is_t {
				inc_nme_attr_int(plyrw+";w;cm", 1)
				inc_nme_attr_int(plyrl+";l;cm", 1)
				wlt_tbrkr_i(gn_is_tbkr, plyrw, plyrl, "gn", "cm", gn)
				wlt_tbrkr_t(tme_is_tbkr, plyrw, plyrl, "tme", "cm", tme)
				wlt_tbrkr_i(ot_is_tbkr, plyrw, plyrl, "ot", "cm", ot)
				wlt_tbrkr_i(sht_is_tbkr, plyrw, plyrl, "sht", "cm", sht)
				wlt_tbrkr_i(pts_is_tbkr, plyrw, plyrl, "pts", "cm", pts)
			} else {
				inc_nme_attr_int(plyr1+";w;cm", 1)
				inc_nme_attr_int(plyr2+";l;cm", 1)
			}
			//wlt_tbrkr_i(nt_is_tbkr, plyrw, plyrl, "nt", "cm", nt)
		}
	}

	int_wk := wk
	cur_wk := strconv.Itoa(wk)

	lp_cnt := 1
	if int_wk <= wk_done {
		lp_cnt = 0
	}

	for lp_i := 0; lp_i < lp_cnt; lp_i++ {
		wk_done = int_wk
		if is_sovors {
			sovs()
		}
		if is_rnk {
			rnk(sorted_cndv, "")
			rnk(sorted_cndv, "cn")
		}
		prev1_cn := ""

		for cndv_i := 0; cndv_i < len(sorted_cndv); cndv_i++ {
			ocn1 := strings.Index(sorted_cndv[cndv_i], ";")
			if sorted_cndv[cndv_i][:ocn1] != prev1_cn {
				prev1_cn = sorted_cndv[cndv_i][:ocn1]
			}
			cndv_to_plrs[sorted_cndv[cndv_i]] = sort_plrs(cndv_to_plrs[sorted_cndv[cndv_i]], dv_tbr_expnd, 0, true, false, false, false)
		}

		prev1_cn = ""
		conf_indx := -1
		divi_indx := 0
		for cndv_i := 0; cndv_i < len(sorted_cndv); cndv_i++ {
			ocn1 := strings.Index(sorted_cndv[cndv_i], ";")
			if sorted_cndv[cndv_i][:ocn1] != prev1_cn {
				var tempConf model.SchResConf
				tempConf.Conf = sorted_cndv[cndv_i][:ocn1]
				stnds.Confs = append(stnds.Confs, tempConf)
				conf_indx++
				divi_indx = 0
				prev1_cn = sorted_cndv[cndv_i][:ocn1]
			} else {
				divi_indx++
			}
			var tempDivi model.SchResDivi
			tempDivi.Divi = sorted_cndv[cndv_i][ocn1+1:]
			stnds.Confs[conf_indx].Divis = append(stnds.Confs[conf_indx].Divis, tempDivi)

			stnds.Confs[conf_indx].Divis[divi_indx].Headers = append(stnds.Confs[conf_indx].Divis[divi_indx].Headers, "Place", "Player")
			for dv_tbr_i := 0; dv_tbr_i < len(dv_tbr_arr); dv_tbr_i++ {
				tmp_s := strings.ReplaceAll(dv_tbr_arr[dv_tbr_i], "-pt", "")
				//if (tmp_s == "swp") || (tmp_s == "arch") || (tmp_s == "flg") || strings.HasPrefix(tmp_s, "hh") {
				if (tmp_s == "swp") || (tmp_s == "flg") || strings.HasPrefix(tmp_s, "hh") {
					continue
				} else if strings.Contains(tmp_s, "(") {
					tmp_arr := strings.SplitN(tmp_s, "(", 2)
					attr := tmp_arr[0]
					tmp_arr[1] = strings.Trim(tmp_arr[1], ")")
					tmp_arr = strings.Split(tmp_arr[1], "-")
					for tmp_i := 0; tmp_i < len(tmp_arr); tmp_i++ {
						stnds.Confs[conf_indx].Divis[divi_indx].Headers = append(stnds.Confs[conf_indx].Divis[divi_indx].Headers, gme_mp_prnt(attr+"-"+tmp_arr[tmp_i]))
					}
				} else if tmp_s == "arch-d" {
					stnds.Confs[conf_indx].Divis[divi_indx].Headers = append(stnds.Confs[conf_indx].Divis[divi_indx].Headers, "arch")
				} else {
					stnds.Confs[conf_indx].Divis[divi_indx].Headers = append(stnds.Confs[conf_indx].Divis[divi_indx].Headers, gme_mp_prnt(tmp_s))
				}
			}
			if len(plyr_stnd) > 0 {
				stnds.Confs[conf_indx].Divis[divi_indx].Headers = append(stnds.Confs[conf_indx].Divis[divi_indx].Headers, gme_mp_prnt("stnd wk "+cur_wk))
			}
			stnds.Confs[conf_indx].Divis[divi_indx].Headers = append(stnds.Confs[conf_indx].Divis[divi_indx].Headers, gme_mp_prnt("sds"))

			remove_c := false
			for nms_i := 0; nms_i < len(cndv_to_plrs[sorted_cndv[cndv_i]]); nms_i++ {
				var tempPlyrtm model.SchResPlyrtm
				tempPlyrtm.Rnk = nms_i + 1
				p := cndv_to_plrs[sorted_cndv[cndv_i]][nms_i]
				ply_wns[p] = nme_attr_int[p+";w;o"]
				if _, hasKey := wkplyr_depc[cur_wk+";"+p]; hasKey {
					if wkplyr_depc[cur_wk+";"+p] == "h" {
						curplyr_depc[p] = wkplyr_depc[cur_wk+";"+p]
						remove_c = true
					} else if wkplyr_depc[cur_wk+";"+p] == "b" {
						curplyr_depc[p] = wkplyr_depc[cur_wk+";"+p]
						remove_c = true
					} else if wkplyr_depc[cur_wk+";"+p] == "d" {
						curplyr_depc[p] = wkplyr_depc[cur_wk+";"+p]
						remove_c = true
					} else if (wkplyr_depc[cur_wk+";"+p] == "e") || (wkplyr_depc[cur_wk+";"+p] == "ee") || (wkplyr_depc[cur_wk+";"+p] == "x") || (wkplyr_depc[cur_wk+";"+p] == "a") || (wkplyr_depc[cur_wk+";"+p] == "pc") {
						curplyr_depc[p] = wkplyr_depc[cur_wk+";"+p]
					} else if curplyr_depc[p] == "r" {
						wkplyr_depc[cur_wk+";"+p] = wkplyr_depc[cur_wk+";"+p]
					} else if curplyr_depc[p] != wkplyr_depc[cur_wk+";"+p] {
						curplyr_depc[p] += wkplyr_depc[cur_wk+";"+p]
					}
				}
				if remove_c {
					curplyr_depc[p] = strings.ReplaceAll(curplyr_depc[p], "c", "")
				}
				curplyr_depc[p] = strings.ReplaceAll(curplyr_depc[p], "cp", "pc")
				if curplyr_depc[p] == "" {
					tempPlyrtm.Plyrtm = p
				} else {
					tempPlyrtm.Plyrtm = curplyr_depc[p] + "-" + p
				}
				for dv_tbr_i := 0; dv_tbr_i < len(dv_tbr_arr); dv_tbr_i++ {
					a := dv_tbr_arr[dv_tbr_i]
					if (a == "o") || (a == "dv") || (a == "cm") || (a == "cn") {
						t := ""
						if nme_attr_int[p+";t;"+a] == 1 {
							t = ".5"
						}
						v := strconv.Itoa(nme_attr_int[p+";w;"+a]) + t + "--"
						v += strconv.Itoa(nme_attr_int[p+";l;"+a]) + t
						tempPlyrtm.Columns = append(tempPlyrtm.Columns, v)
						//} else if (a == "swp") || (a == "arch") || (a == "flg") || strings.HasPrefix(a, "hh") {
					} else if a == "arch" {
						wl := ""
						vint := nme_attr_int[p+";w;arch"] - nme_attr_int[p+";l;arch"]
						v := ""
						if vint < 0 {
							v = fmt.Sprintf("(%d)", vint*-1)
							if wk == maxwk {
								wl = "L"
							}
						} else {
							v = fmt.Sprintf("%d", vint)
							if (wk == maxwk) && (vint > 0) {
								wl = "W"
							}
						}

						for _, a2 := range dv_tbr_arr {
							if strings.Contains(a2, "-") {
								a2 = a2[0:strings.Index(a2, "-")]
							}
							if strings.HasPrefix(a2, "ot") || strings.HasPrefix(a2, "gn") || strings.HasPrefix(a2, "sht") || strings.HasPrefix(a2, "pts") {
								vint = nme_attr_int[p+";"+a2+";arch"]
								v2 := ""
								if vint < 0 {
									v2 = fmt.Sprintf("(%d)", vint*-1)
									if (wk == maxwk) && (wl == "") {
										wl = "L"
									}
								} else {
									v2 = fmt.Sprintf("%d", vint)
									if (wk == maxwk) && (wl == "") && (vint > 0) {
										wl = "W"
									}
								}
								v += "--" + v2
							}
						}
						if wl != "" {
							v = wl + "--" + v
						}
						tempPlyrtm.Columns = append(tempPlyrtm.Columns, v)
					} else if (a == "swp") || (a == "flg") || strings.HasPrefix(a, "hh") {
						continue
					} else if strings.HasPrefix(a, "clsgm") {
						tempPlyrtm.Columns = append(tempPlyrtm.Columns, clsgm_hsh[p])
					} else if a == "strk" {
						tempPlyrtm.Columns = append(tempPlyrtm.Columns, strk_hsh[p])
					} else if a == "tme" {
						vs := strings.SplitN(a, "-", 2)
						v := nme_attr_tme[p+";"+vs[0]+";"]
						v = "-" + v + "-"
						if (len(vs) > 1) && (vs[1] == "pt") {
							for bsc_i := 0; bsc_i < len(dv_tbr_bsc); bsc_i++ {
								v2 := nme_attr_tme[p+";"+vs[0]+";"+dv_tbr_bsc[bsc_i]]
								v2 = "-" + v2 + "-"
								v += "**" + v2
							}
						}
						tempPlyrtm.Columns = append(tempPlyrtm.Columns, v)
					} else if strings.HasPrefix(a, "rnk") {
						vs := strings.Split(a, "-")
						v := ""
						if len(vs) > 1 {
							v = plyr_rnk_dsc[p+";"+vs[1]]
						} else {
							v = plyr_rnk_dsc[p+";"]
						}
						tempPlyrtm.Columns = append(tempPlyrtm.Columns, v)
					} else if strings.HasPrefix(a, "potw") {
						v := potw[p]
						tempPlyrtm.Columns = append(tempPlyrtm.Columns, fmt.Sprintf("%d", v))
					} else if strings.Contains(a, "(") {
						tmp_arr := strings.SplitN(a, "(", 2)
						attr := tmp_arr[0]
						tmp_arr[1] = strings.Trim(tmp_arr[1], ")")
						tmp_arr = strings.Split(tmp_arr[1], "-")
						v := 0
						for tmp_i := 0; tmp_i < len(tmp_arr); tmp_i++ {
							if tmp_arr[tmp_i] == "o" {
								v = nme_attr_int[p+";"+attr+";"]
							} else {
								v = nme_attr_int[p+";"+attr+";"+tmp_arr[tmp_i]]
							}
							tempPlyrtm.Columns = append(tempPlyrtm.Columns, strconv.Itoa(v))
						}
					} else if strings.Contains(a, "arch-d") {
						v := nme_attr_int[p+";w;arch"] - nme_attr_int[p+";l;arch"]
						arch_wn := ""
						if int_wk >= maxwk {
							if v > 0 {
								arch_wn = "W "
							} else if v < 0 {
								arch_wn = "L "
							} else if nme_attr_int[p+";arw"] > 0 {
								arch_wn = "W "
							} else if nme_attr_int[p+";arw"] < 0 {
								arch_wn = "L "
							}
						}
						v2 := strconv.Itoa(v)
						if strings.HasPrefix(v2, "-") {
							v2 = "(" + v2[1:] + ")"
						}
						for dv_tbr_j := 0; dv_tbr_j < len(dv_tbr_arr); dv_tbr_j++ {
							if strings.HasPrefix(dv_tbr_arr[dv_tbr_j], "gn") || strings.HasPrefix(dv_tbr_arr[dv_tbr_j], "ot") || strings.HasPrefix(dv_tbr_arr[dv_tbr_j], "sht") || strings.HasPrefix(dv_tbr_arr[dv_tbr_j], "pts") {
								onetwo := strings.Split(dv_tbr_arr[dv_tbr_j], "-")
								v = nme_attr_int[p+";"+onetwo[0]+";arch"]
								if (int_wk >= maxwk) && (arch_wn == "") {
									if v > 0 {
										arch_wn = "W "
									} else if v < 0 {
										arch_wn = "L "
									} else if nme_attr_int[p+";arw"] > 0 {
										arch_wn = "W "
									} else if nme_attr_int[p+";arw"] < 0 {
										arch_wn = "L "
									}
								}
								vs := strconv.Itoa(v)
								if strings.HasPrefix(vs, "-") {
									vs = "(" + vs[1:] + ")"
								}
								v2 += "--" + vs
							}
						}
						if (int_wk >= maxwk) && (arch_wn == "") {
							v = nme_attr_int[p+";arw"]
							vs := strconv.Itoa(v)
							if strings.HasPrefix(vs, "-") {
								vs = "(" + vs[1:] + ")"
							}
							v2 += "--" + vs
						}
						tempPlyrtm.Columns = append(tempPlyrtm.Columns, arch_wn+v2)
					} else {
						vs := strings.SplitN(a, "-", 2)
						v := ""
						if (len(vs) == 1) || (vs[1] == "pt") || (vs[1] == "o") {
							v = strconv.Itoa(nme_attr_int[p+";"+vs[0]+";"])
						} else {
							v = strconv.Itoa(nme_attr_int[p+";"+vs[0]+";"+vs[1]])
						}
						if strings.HasPrefix(v, "-") {
							v = "(" + v[1:] + ")"
						} else if nme_attr_int[p+";"+vs[0]+"-t;"] == 1 {
							v += ".5"
						}
						if (len(vs) > 1) && (vs[1] == "pt") {
							for bsc_i := 0; bsc_i < len(dv_tbr_bsc); bsc_i++ {
								if strings.HasPrefix(dv_tbr_bsc[bsc_i], "cm") && (dv_tbr_bsc[bsc_i] != "cm") {
									continue
								}
								v2 := strconv.Itoa(nme_attr_int[p+";"+vs[0]+";"+dv_tbr_bsc[bsc_i]])
								if strings.HasPrefix(v2, "-") {
									v2 = "(" + v2[1:] + ")"
								}
								v += "--" + v2
							}
						}
						tempPlyrtm.Columns = append(tempPlyrtm.Columns, v)
					}
				}
				if len(plyr_stnd) > 0 {
					tempPlyrtm.Columns = append(tempPlyrtm.Columns, strconv.Itoa(plyr_stnd[p]))
				} else if len(plyr_sds) > 0 {
					tempPlyrtm.Columns = append(tempPlyrtm.Columns, strconv.Itoa(plyr_sds[p]))
				}

				tempPlyrtm.Columns = append(tempPlyrtm.Columns, strconv.Itoa(plrs[plrsNameToId[p]].sds))
				stnds.Confs[conf_indx].Divis[divi_indx].Plyrtms = append(stnds.Confs[conf_indx].Divis[divi_indx].Plyrtms, tempPlyrtm)
			}
		}
	}
	/*if HasAnArch && (wk == (maxwk - 1)) {
		sheet = spr_in_file.Sheets[sheetnum[maxwk]]
		cnt := 0
		row := 1
		for cnt < expected[maxwk] {
			cell, _ = sheet.Cell(row, 0)
			plyr1 := cell.Value
			cell, _ = sheet.Cell(row, 1)
			plyr2 := cell.Value
			if (plyr1 != "") && (plyr2 != "") {
				lst_wks_plyrs += ";" + plyr1 + ";" + plyr2
				cnt++
			}
			row++
			if row == 500 {
				break
			}
		}
		lst_wk(maxwk, dir+"TMaster\\trn"+LetterArr[LetterArrIndx]+"\\", lst_wks_plyrs, dv_tbr_arr)
	}*/

	fmt.Println(wk, maxwk)

	row := 0
	if running_sht != "" {
		sht_mp_arr = append(sht_mp_arr, running_sht)
		running_sht = ""
	}

	var tempConfOOO model.ConfOOO
	prevConf := ""
	for cndv_i := 0; cndv_i < len(sorted_cndv); cndv_i++ {
		cndDiv := strings.Split(sorted_cndv[cndv_i], ";")
		if cndDiv[0] != prevConf {
			if prevConf != "" {
				tempConfOOO.Conf = prevConf
				stnds.ConfOoo = append(stnds.ConfOoo, tempConfOOO)
				tempConfOOO = model.ConfOOO{}
			}
			prevConf = cndDiv[0]
		}

		var tempDiviOOO model.DiviOOO
		tempDiviOOO.Divi = cndDiv[1]

		row++
		for nms_i := 0; nms_i < len(cndv_to_plrs[sorted_cndv[cndv_i]]); nms_i++ {
			var tempPlyrtmOOO model.PlyrtmOOO
			p := cndv_to_plrs[sorted_cndv[cndv_i]][nms_i]
			tempPlyrtmOOO.Plyrtm = p
			tmp_slc := strings.Split(strings.Join(plr_plr_w[p], ";"), ";")
			for i := 0; len(tmp_slc) > 0; i++ {
				indx := lowest_nme_attr(tmp_slc)
				tempPlyrtmOOO.Wn = append(tempPlyrtmOOO.Wn, tmp_slc[indx])
				plr_plr_gms_w[p] += ";" + tmp_slc[indx]
				tmp_slc = append(tmp_slc[:indx], tmp_slc[indx+1:]...)
			}
			row++
			tmp_slc = strings.Split(strings.Join(plr_plr_l[p], ";"), ";")
			for i := 0; len(tmp_slc) > 0; i++ {
				indx := lowest_nme_attr(tmp_slc)
				tempPlyrtmOOO.Ls = append(tempPlyrtmOOO.Ls, tmp_slc[indx])
				plr_plr_gms_l[p] += ";" + tmp_slc[indx]
				tmp_slc = append(tmp_slc[:indx], tmp_slc[indx+1:]...)
			}
			row++
			if len(plr_plr_t[p]) > 0 {
				tmp_slc = strings.Split(strings.Join(plr_plr_t[p], ";"), ";")
				for i := 0; len(tmp_slc) > 0; i++ {
					indx := lowest_nme_attr(tmp_slc)
					tempPlyrtmOOO.Tie = append(tempPlyrtmOOO.Tie, tmp_slc[indx])
					plr_plr_gms_t[p] += ";" + tmp_slc[indx]
					tmp_slc = append(tmp_slc[:indx], tmp_slc[indx+1:]...)
				}
				row++
			}
			tempDiviOOO.PlyrtmOoo = append(tempDiviOOO.PlyrtmOoo, tempPlyrtmOOO)
		}

		tempConfOOO.DiviOoo = append(tempConfOOO.DiviOoo, tempDiviOOO)
	}
	if prevConf != "" {
		tempConfOOO.Conf = prevConf
		stnds.ConfOoo = append(stnds.ConfOoo, tempConfOOO)
	}

	/*if running_sht != "" {
		sht_mp_arr = append(sht_mp_arr, running_sht)
		running_sht = ""
	}
	Sprsht_out(schxls_fbase+"Res.xlsx", sht_mp_arr)
	*/

	//dir1 := dir
	//dir2 := dir
	stnd(wk, cndv_wc, sorted_cndv, cn_tbr_expnd)
	if rnd < maxwk {
		err = nxt_wk(trnId, rnd+1)
		if err != nil {
			return nil, err
		}
	}
	if wk < maxwk {
		err = gms(wk, cndv_wc, sorted_cndv, dv_tbr_arr)
		if err != nil {
			return nil, err
		}
	}
	jsonOut, _ := json.MarshalIndent(stnds, "", "  ")
	return jsonOut, nil
}

func sovs() {
	for plrInt := range plrs {
		p_sovors := 0.0
		plr := plyr_id_to_name[plrInt]
		for plr_i := 0; plr_i < len(plr_plr_w[plr]); plr_i++ {
			plr1str := strings.SplitN(plr_plr_w[plr][plr_i], ",", 2)
			p_sovors += float64(nme_attr_int[plr1str[0]+";w;o"])
			p_sovors += float64(nme_attr_int[plr1str[0]+";t;o"]) / 2
		}
		nme_attr_int[plr+";sov;"] = int(math.Floor(p_sovors))
		nme_attr_int[plr+";sov-t;"] = int(math.Round(p_sovors*2)) % 2

		for plr_i := 0; plr_i < len(plr_plr_t[plr]); plr_i++ {
			plr1str := strings.SplitN(plr_plr_t[plr][plr_i], ",", 2)
			p_sovors += float64(nme_attr_int[plr1str[0]+";w;o"])
			p_sovors += float64(nme_attr_int[plr1str[0]+";t;o"]) / 2
		}
		for plr_i := 0; plr_i < len(plr_plr_l[plr]); plr_i++ {
			plr1str := strings.SplitN(plr_plr_l[plr][plr_i], ",", 2)
			p_sovors += float64(nme_attr_int[plr1str[0]+";w;o"])
			p_sovors += float64(nme_attr_int[plr1str[0]+";t;o"]) / 2
		}

		nme_attr_int[plr+";sos;"] = int(math.Floor(p_sovors))
		nme_attr_int[plr+";sos-t;"] = int(math.Round(p_sovors*2)) % 2
	}
}

func add_clsgm(Aplt string, Ascr int, Aadd int) {
	add_cls_prts := strings.Split(clsgm_hsh[Aplt], "|")
	add_cls_done := false
	add_cls_rtn := ""
	for add_cls_i := 0; (clsgm_hsh[Aplt] != "") && (add_cls_i < len(add_cls_prts)); add_cls_i++ {
		onetwo := strings.SplitN(add_cls_prts[add_cls_i], ":", 2)
		if !add_cls_done {
			clst1 := StrToInt(onetwo[0])
			if Ascr < clst1 {
				add_cls_rtn += "|" + strconv.Itoa(Ascr) + ":" + strconv.Itoa(Aadd)
				add_cls_done = true
			} else if Ascr == clst1 {
				clst1 = StrToInt(onetwo[1])
				clst1 += Aadd
				onetwo[1] = strconv.Itoa(clst1)
				add_cls_done = true
			}
		}
		if onetwo[1] != "0" {
			add_cls_rtn += "|" + onetwo[0] + ":" + onetwo[1]
		}
	}
	if !add_cls_done {
		add_cls_rtn += "|" + strconv.Itoa(Ascr) + ":" + strconv.Itoa(Aadd)
	}
	if add_cls_rtn != "" {
		clsgm_hsh[Aplt] = add_cls_rtn[1:]
	} else {
		clsgm_hsh[Aplt] = ""
	}
}

func add_strk(Aplt string, worl string) {
	if strings.HasPrefix(strk_hsh[Aplt], worl) {
		strk_prt := strings.SplitN(strk_hsh[Aplt][1:], "|", 2)
		tempnum := StrToInt(strk_prt[0])
		if worl == "W" {
			strk_hsh[Aplt] = "W" + strconv.Itoa(tempnum+1)
		} else {
			strk_hsh[Aplt] = "L" + strconv.Itoa(tempnum+1)
		}
		if len(strk_prt) > 1 {
			strk_hsh[Aplt] += "|" + strk_prt[1]
		}
	} else if strk_hsh[Aplt] == "" {
		strk_hsh[Aplt] = worl + "1"
	} else {
		strk_hsh[Aplt] = worl + "1|" + strk_hsh[Aplt]
	}
}

func rnk(sorted_cndv []string, has_cn string) {
	plr_ptsf := make(map[int][]string)
	plr_ptsa := make(map[int][]string)
	prev_cn := ""

	for cndv_i := 0; cndv_i < len(sorted_cndv); cndv_i++ {
		ocn := strings.Index(sorted_cndv[cndv_i], ";")
		if prev_cn == "" {
			prev_cn = sorted_cndv[cndv_i][:ocn]
		} else if (has_cn != "") && !strings.HasPrefix(sorted_cndv[cndv_i], prev_cn) {
			prev_cn = sorted_cndv[cndv_i][:ocn]
			assgn_rnk(";rnk-f;"+has_cn, plr_ptsf)
			assgn_rnk(";rnk-a;"+has_cn, plr_ptsa)
			plr_ptsf = nil
			plr_ptsa = nil
			plr_ptsf = make(map[int][]string)
			plr_ptsa = make(map[int][]string)
		}
		for plr_i := 0; plr_i < len(cndv_to_plrs[sorted_cndv[cndv_i]]); plr_i++ {
			p := cndv_to_plrs[sorted_cndv[cndv_i]][plr_i]
			pts := nme_attr_int[p+";pf;"+has_cn]
			plr_ptsf[pts] = append(plr_ptsf[pts], p)
			pts = nme_attr_int[p+";pa;"+has_cn]
			plr_ptsa[pts] = append(plr_ptsa[pts], p)
		}
	}
	if prev_cn != "" {
		assgn_rnk(";rnk-f;"+has_cn, plr_ptsf)
		assgn_rnk(";rnk-a;"+has_cn, plr_ptsa)
		add_rnks(has_cn)
	}
}

func assgn_rnk(rnk_type string, rnk_plr map[int][]string) {
	rnk_int := SortIntsInMap(rnk_plr)
	cur_rnk := 0.5
	for rnk_i := 0; rnk_i < len(rnk_int); rnk_i++ {
		ln := float64(len(rnk_plr[rnk_int[rnk_i]]))
		nw_rnk := cur_rnk + (ln / 2)
		for pi := 0; pi < len(rnk_plr[rnk_int[rnk_i]]); pi++ {
			nme_attr_flt[rnk_plr[rnk_int[rnk_i]][pi]+rnk_type] = nw_rnk
		}

		cur_rnk += ln
	}
}

func add_rnks(has_cn string) {
	for cndv := range cndv_to_plrs {
		for pi := 0; pi < len(cndv_to_plrs[cndv]); pi++ {
			nme_attr_flt[cndv_to_plrs[cndv][pi]+";rnk;"+has_cn] = nme_attr_flt[cndv_to_plrs[cndv][pi]+";rnk-f;"+has_cn]
			nme_attr_flt[cndv_to_plrs[cndv][pi]+";rnk;"+has_cn] += nme_attr_flt[cndv_to_plrs[cndv][pi]+";rnk-a;"+has_cn]
			plyr_rnk_dsc[cndv_to_plrs[cndv][pi]+";"+has_cn] = fmt.Sprintf("%g", nme_attr_flt[cndv_to_plrs[cndv][pi]+";rnk;"+has_cn]) + "#"
			plyr_rnk_dsc[cndv_to_plrs[cndv][pi]+";"+has_cn] += fmt.Sprintf("%g", nme_attr_flt[cndv_to_plrs[cndv][pi]+";rnk-f;"+has_cn]) + "#"
			plyr_rnk_dsc[cndv_to_plrs[cndv][pi]+";"+has_cn] += strconv.Itoa(nme_attr_int[cndv_to_plrs[cndv][pi]+";pf;"+has_cn]) + "#"
			plyr_rnk_dsc[cndv_to_plrs[cndv][pi]+";"+has_cn] += fmt.Sprintf("%g", nme_attr_flt[cndv_to_plrs[cndv][pi]+";rnk-a;"+has_cn]) + "#"
			plyr_rnk_dsc[cndv_to_plrs[cndv][pi]+";"+has_cn] += strconv.Itoa(nme_attr_int[cndv_to_plrs[cndv][pi]+";pa;"+has_cn] * -1)
		}
	}
}

func tbr_bsc(tbr_arr []string) []string {
	start_i := len(tbr_arr)
	end_i := -1
	for i := 0; i < len(tbr_arr); i++ {
		a := tbr_arr[i]
		if (a == "dv") || strings.HasPrefix(a, "cm") || (a == "cn") || strings.HasPrefix(a, "act") {
			if start_i == len(tbr_arr) {
				start_i = i
			}
		} else if start_i < i {
			end_i = i
			break
		}
	}
	if start_i == len(tbr_arr) {
		return nil
	}
	return tbr_arr[start_i:end_i]
}

func tbr_exp(tbr_arr, tbr_bsc []string) []string {
	var tbr_rslt []string
	for tbr_i := 0; tbr_i < len(tbr_arr); tbr_i++ {
		if strings.HasSuffix(tbr_arr[tbr_i], "-pt") {
			tmps := tbr_arr[tbr_i][:len(tbr_arr[tbr_i])-3]
			tbr_rslt = append(tbr_rslt, tmps)
			for bsc_i := 0; bsc_i < len(tbr_bsc); bsc_i++ {
				tbr_rslt = append(tbr_rslt, tmps+"-"+tbr_bsc[bsc_i])
			}
		} else if regexp.MustCompile(`^.+\(.+\)$`).MatchString(tbr_arr[tbr_i]) {
			tmps := strings.SplitN(tbr_arr[tbr_i], "(", 2)
			one := strings.TrimSpace(tmps[0])
			tmps = strings.Split(strings.Trim(tmps[1], ")"), "-")
			for tbr_j := 0; tbr_j <= len(tmps[1]); tbr_j++ {
				if tmps[tbr_j] == "o" {
					tbr_rslt = append(tbr_rslt, one)
				} else {
					tbr_rslt = append(tbr_rslt, one+"-"+tmps[tbr_j])
				}
			}
		} else {
			tbr_rslt = append(tbr_rslt, tbr_arr[tbr_i])
		}
	}
	return tbr_rslt
}

func sort_plrs(orig []string, tbr []string, ttl int, isd bool, add_to_g_ord bool, d_to_c bool, has_elm bool) []string {
	var new_p []string
	var chs []string
	var hh_cur []string

	frst_lp_strt := strings.Split(strings.Join(orig, "#"), "#")
	if ttl == 0 {
		ttl = len(orig)
	}

	if add_to_g_ord || (len(g_string_ord) == 0) {
		for tbr_i := 0; tbr_i < len(tbr); tbr_i++ {
			g_string_ord = append(g_string_ord, tbr[tbr_i])
		}
	}
	for len(new_p) < ttl {
		rsn := ""
		tmp_rsn := ""
		var remain []string
		for flp_i := 0; flp_i < len(frst_lp_strt); flp_i++ {
			onetwo := strings.SplitN(frst_lp_strt[flp_i], ";", 2)
			remain = append(remain, onetwo[0])
		}
		tbr_i := 0
		for tbr_i < len(tbr) {
			if len(remain) < 2 {
				break
			}
			a := tbr[tbr_i]
			if a == "swp" {
				tmp_rsn, chs = swp(remain)
			} else if a == "arch" {
				if !arch_wk_yt && !HasAnArch {
					chs = remain
				} else {
					tmp_rsn, chs = arch(remain, tbr[tbr_i+1:])
				}
			} else if a == "arch-elm" {
				if !has_elm {
					chs = remain
				} else {
					tmp_rsn, chs = arch2(remain, tbr[tbr_i+1:])
				}
			} else if strings.HasPrefix(a, "hhcn") {
				tmp_rsn, chs = hhcn(remain)
			} else if strings.HasPrefix(a, "hh-") {
				t_arr := strings.Split(a, "-")
				tmp_rsn, chs = hh_(remain, t_arr[1])
			} else if strings.HasPrefix(a, "hh") {
				if !strings.HasPrefix(a, "hhx") || HH_Dif(remain, hh_cur) {
					tmp_rsn, chs = hh(remain, strings.Contains(a, "a"), strings.Contains(a, "b"))
					if len(chs) == len(remain) && strings.HasPrefix(a, "hhx") && (tmp_rsn != "np_excl") {
						hh_cur = remain
					}
					if tmp_rsn == "np_excl" {
						tmp_rsn = ""
					}
				}
			} else if (a == "o") || (a == "dv") || (a == "cm") || (a == "cn") {
				tmp_rsn, chs = tbr_wlt(remain, a)
				//} else if strings.HasPrefix(a, "act") {
				//	tmp_rsn, chs = tbr_act(remain, a)
			} else if strings.HasPrefix(a, "gn") || strings.HasPrefix(a, "ot") || strings.HasPrefix(a, "sht") || strings.HasPrefix(a, "pts") || strings.HasPrefix(a, "so") || strings.HasPrefix(a, "nt") {
				tmp_rsn, chs = tbr_num(remain, a, true, false)
			} else if strings.HasPrefix(a, "rnk") {
				tmp_rsn, chs = tbr_num(remain, a, false, false)
			} else if strings.HasPrefix(a, "tme") {
				tmp_rsn, chs = tbr_num(remain, a, false, true)
			} else if strings.HasPrefix(a, "clsgm") {
				tmp_rsn, chs = tbr_clsgme(remain)
			} else if strings.HasPrefix(a, "strk") {
				tmp_rsn, chs = tbr_strk(remain)
			} else if regexp.MustCompile(`^cm\d+$`).MatchString(a) {
				tmp_i, err := strconv.Atoi(a[2:])
				if err != nil {
					fmt.Println("err cm\\d:" + a[2:])
					return nil
				}
				tmp_rsn, chs = tbr_cm_min(remain, tmp_i, false)
			} else if regexp.MustCompile(`^cmnd\d+$`).MatchString(a) {
				tmp_i, err := strconv.Atoi(a[4:])
				if err != nil {
					fmt.Println("err cmnd\\d:" + a[4:])
					return nil
				}
				tmp_rsn, chs = tbr_cm_min(remain, tmp_i, true)
				//} else if a == "arch-d" {
				//	chs = remain
			} else if a != "flg" {
				fmt.Println("tbr_dv not there: " + a)
				return nil
			}
			if len(chs) != len(remain) {
				remain = chs
				if tmp_rsn != "" {
					rsn += tmp_rsn + "\n"
				}
				tbr_i = 0
			} else {
				tbr_i++
			}
		}
		if len(remain) > 1 {
			remain = remain[:1]
			//rsn += tmp_rsn + "\n"
		}
		if isd {
			plr_rsn_dv[remain[0]] = append(plr_rsn_dv[remain[0]], rsn)
		} else {
			plr_rsn_cn[remain[0]] = append(plr_rsn_cn[remain[0]], rsn)
		}
		new_p = append(new_p, remain[0])
		frst_lp_strt = Splice(frst_lp_strt, remain[0])
		/*if d_to_c {
			c_remain := plrs[remain[0]].cn
			for c_ri := len(frst_lp_strt) - 1; (c_ri >= 0) && (c_ri < len(frst_lp_strt)); c_ri-- {
				if plrs[frst_lp_strt[c_ri]].cn == c_remain {
					frst_lp_strt = Splice(frst_lp_strt, frst_lp_strt[c_ri])
					ttl--
				}

			}
		}*/
	}
	return new_p
}

func Splice(arr []string, str string) []string {
	for i := 0; i < len(arr); i++ {
		if arr[i] == str {
			return append(arr[:i], arr[i+1:]...)
		} else if strings.Contains(arr[i], str+";") {
			arr[i] = strings.Replace(arr[i], str+";", "", 1)
			break
		}
	}
	return arr
}

func CN(CN_orig []string, tbr []string, ttl int, isd bool, add_to_g_ord bool, d_to_c bool, has_elm bool) []string {
	var CN_wn []string
	var CN_ls []string
	for i := 0; i < len(CN_orig); i++ {
		for j := i + 1; j < len(CN_orig); j++ {
			if plrs[plrsNameToId[CN_orig[i]]].cn == plrs[plrsNameToId[CN_orig[j]]].cn {
				var CN_tmp []string
				CN_tmp = append(CN_tmp, CN_orig[i])
				CN_tmp = append(CN_tmp, CN_orig[j])
				CN_tmp = sort_plrs(CN_tmp, tbr, 2, false, add_to_g_ord, d_to_c, has_elm)
				CN_wn = append(CN_wn, CN_tmp[0])
				CN_ls = append(CN_ls, CN_tmp[1])
			}
			break
		}
	}
	CN_wn = sort_plrs(CN_wn, tbr, ttl, isd, add_to_g_ord, d_to_c, has_elm)
	CN_ls = sort_plrs(CN_ls, tbr, ttl, isd, add_to_g_ord, d_to_c, has_elm)
	return append(CN_wn, CN_ls...)
}

func swp(swp_orig []string) (string, []string) {
	var swp_wn []string
	var swp_ls []string
	var swp_nn []string
	high_prc := -1.0

	for i := 0; i < len(swp_orig); i++ {
		w := 0.0
		tot := 0.0
		for j := 0; j < len(swp_orig); j++ {
			if i == j {
				continue
			}
			once := false
			for k := 0; k < len(plr_plr_w[swp_orig[i]]); k++ {
				if strings.HasPrefix(plr_plr_w[swp_orig[i]][k], swp_orig[j]+",") {
					if !once {
						once = true
					} else {
						w++
						tot++
						break
					}
				}
			}
			once = false
			for k := 0; k < len(plr_plr_l[swp_orig[i]]); k++ {
				if strings.HasPrefix(plr_plr_l[swp_orig[i]][k], swp_orig[j]+",") {
					if !once {
						once = true
					} else {
						tot++
						break
					}
				}
			}
		}
		if tot == 0 {
			swp_nn = append(swp_nn, swp_orig[i])
		} else {
			prc := w / tot
			if math.Abs(high_prc-prc) < 0.000001 {
				swp_wn = append(swp_wn, swp_orig[i])
			} else if prc < high_prc {
				swp_ls = append(swp_ls, swp_orig[i])
			} else {
				swp_ls = append(swp_ls, swp_wn...)
				swp_wn = nil
				swp_wn = append(swp_wn, swp_orig[i])
				high_prc = prc
			}
		}
	}
	rsn := ""
	if len(swp_ls) > 0 {
		rsn = strings.Join(swp_wn, ", ") + " SWP o " + strings.Join(swp_ls, ", ")
		put_in_g_string("SWP", swp_wn, swp_ls)
	}
	return rsn, append(swp_wn, swp_nn...)
}

func arch(arch_orig []string, rm_tbr []string) (string, []string) {
	var arch_wn []string
	var arch_oth []string
	var arch_nn []string
	cmpr_found := ";"
	rsn := ""
	for i := 0; i < len(arch_orig); i++ {
		if strings.Contains(cmpr_found, ";"+arch_orig[i]+";") {
			continue
		}
		WNR := ""
		LSR := ""
		a2 := ""
		for j := 0; j < len(arch_orig); j++ {
			if i == j {
				continue
			}
			aa := plrs[plrsNameToId[arch_orig[i]]].sds
			if aa%2 == 1 {
				aa++
			}
			aa2 := plrs[plrsNameToId[arch_orig[j]]].sds
			if aa2%2 == 1 {
				aa2++
			}
			if aa != aa2 {
				continue
			}
			cmpr_found += arch_orig[i] + ";" + arch_orig[j] + ";"
			for k := 0; k < len(rm_tbr); k++ {
				a := rm_tbr[k]
				if strings.HasPrefix(a, "hh") {
					j_plr_plr_w := ";" + strings.Join(plr_plr_w[arch_orig[i]], ";") + ";"
					j_plr_plr_l := ";" + strings.Join(plr_plr_l[arch_orig[i]], ";") + ";"
					if strings.Contains(j_plr_plr_w, ";"+arch_orig[j]+",") {
						if strings.Contains(j_plr_plr_l, ";"+arch_orig[j]+",") {
							continue
						}
						WNR = arch_orig[i]
						LSR = arch_orig[j]
						a2 = "HH"
					} else if strings.Contains(j_plr_plr_l, ";"+arch_orig[j]+",") {
						WNR = arch_orig[j]
						LSR = arch_orig[i]
						a2 = "HH"
					}
				} else if (a == "gn") || (a == "ot") || (a == "sht") || (a == "pts") || (a == "nt") {
					attr2 := a
					prc1 := float64(nme_attr_int[arch_orig[i]+";"+attr2+";arch"])
					prc2 := float64(nme_attr_int[arch_orig[j]+";"+attr2+";arch"])
					if prc1 > prc2 {
						WNR = arch_orig[i]
						LSR = arch_orig[j]
						a2 = a
					} else if prc2 > prc1 {
						WNR = arch_orig[j]
						LSR = arch_orig[i]
						a2 = a
					}
				}

				if WNR != "" {
					arch_wn = append(arch_wn, WNR)
					arch_oth = append(arch_oth, WNR)
					arch_oth = append(arch_oth, LSR)
					rsn += "\r\n" + WNR + " ARCH " + a2 + " o " + LSR
					put_in_g_string("ARCH "+a2, arch_oth, nil)
					arch_oth = nil
					break
				}
			}
			if WNR == "" {
				if nme_attr_int[arch_orig[i]+";arw"] > nme_attr_int[arch_orig[j]+";arw"] {
					arch_wn = append(arch_wn, arch_orig[i])
					rsn += "\r\n" + arch_orig[i] + " ARCH TB o " + arch_orig[j]
				} else if nme_attr_int[arch_orig[i]+";arw"] < nme_attr_int[arch_orig[j]+";arw"] {
					arch_wn = append(arch_wn, arch_orig[j])
					rsn += "\r\n" + arch_orig[j] + " ARCH TB o " + arch_orig[i]
				} else {
					arch_nn = append(arch_nn, arch_orig[i])
					arch_nn = append(arch_nn, arch_orig[j])
				}
			}
		}
		if !strings.Contains(cmpr_found, ";"+arch_orig[i]+";") {
			arch_nn = append(arch_nn, arch_orig[i])
		}
	}
	strings.Trim(rsn, "\r\n")
	return rsn, append(arch_wn, arch_nn...)
}

func arch2(arch_orig []string, rm_tbr []string) (string, []string) {
	var arch_wn []string
	var arch_ls []string
	arch_diff := false
	arch_wh := ""

	for k := 0; k < len(rm_tbr); k++ {
		var arch_wnr_arr []int
		for i := 0; i < len(arch_orig); i++ {
			arch_wnr_arr = append(arch_wnr_arr, 0)
			for j := 0; j < len(arch_orig); j++ {
				if i == j {
					continue
				}
				j_plr_plr_w := ";" + strings.Join(plr_plr_w[arch_orig[i]], ";") + ";"
				j_plr_plr_w_arr := strings.Split(j_plr_plr_w, ";"+arch_orig[j]+",")
				j_plr_plr_l := ";" + strings.Join(plr_plr_l[arch_orig[i]], ";") + ";"
				j_plr_plr_l_arr := strings.Split(j_plr_plr_l, ";"+arch_orig[j]+",")
				a := rm_tbr[k]

				if strings.HasPrefix(a, "hh") {
					arch_wnr_arr[i] += len(j_plr_plr_w_arr) - 1
					arch_wnr_arr[i] -= len(j_plr_plr_l_arr) - 1
				} else {
					for i2 := 1; i2 < len(j_plr_plr_w_arr); i2++ {
						tmp_tmp := strings.SplitN(j_plr_plr_w_arr[i2], ";", 2)
						tmp_arr := strings.Split(tmp_tmp[0], ",")
						pf := StrToInt(tmp_arr[1])
						pa := StrToInt(tmp_arr[2])
						if (a == "ot") && (pf == pa) {
							arch_wnr_arr[i]--
						} else if (a == "sht") && (pf > 0) && (pa == 0) {
							arch_wnr_arr[i]++
						} else if a == "pts" {
							arch_wnr_arr[i] += (pf - pa)
						}
					}
					for i2 := 1; i2 < len(j_plr_plr_l_arr); i2++ {
						tmp_arr := strings.Split(j_plr_plr_l_arr[i2], ",")
						pf := StrToInt(tmp_arr[1])
						pa := StrToInt(tmp_arr[2])
						if (a == "ot") && (pf == pa) {
							arch_wnr_arr[i]++
						} else if (a == "sht") && (pf > 0) && (pa == 0) {
							arch_wnr_arr[i]--
						} else if a == "pts" {
							arch_wnr_arr[i] -= (pf - pa)
						}
					}
				}
			}
			if (i > 0) && (arch_wnr_arr[i] != arch_wnr_arr[i-1]) {
				arch_diff = true
			}
		}

		if arch_diff {
			arch_max := -9999
			for i := 0; i < len(arch_orig); i++ {
				if arch_max == arch_wnr_arr[i] {
					arch_wn = append(arch_wn, arch_orig[i])
				} else if arch_wnr_arr[i] < arch_max {
					arch_ls = append(arch_ls, arch_orig[i])
				} else {
					arch_ls = append(arch_ls, arch_wn...)
					arch_wn = nil
					arch_wn = append(arch_wn, arch_orig[i])
					arch_max = arch_wnr_arr[i]
				}
			}
			arch_wh = rm_tbr[k]
			break
		}
	}

	if arch_diff {
		rsn := ""
		rsn = strings.Join(arch_wn, ", ") + " arch " + arch_wh + " over " + strings.Join(arch_ls, ", ")
		put_in_g_string("arch", arch_wn, arch_ls)
		return rsn, append(arch_wn)
	} else {
		return "", arch_orig
	}
}

func hhcn(hhcn_orig []string) (string, []string) {
	var hhcn_wn []string
	var hhcn_ls []string
	var hhcn_nn []string

	for i := 0; i < len(hhcn_orig); i++ {
		for j := 0; j < len(hhcn_orig); j++ {
			if i == j {
				continue
			}
			if plrs[plrsNameToId[hhcn_orig[i]]].cn == plrs[plrsNameToId[hhcn_orig[j]]].cn {
				j_plr_plr_w := ";" + strings.Join(plr_plr_w[hhcn_orig[i]], ";") + ";"
				if strings.Contains(j_plr_plr_w, ";"+hhcn_orig[j]+",") {
					hhcn_wn = append(hhcn_wn, hhcn_orig[i])
					hhcn_ls = append(hhcn_ls, hhcn_orig[j])
					break
				}
			}
		}
	}
	rsn := ""
	for i := 0; i < len(hhcn_wn); i++ {
		rsn += "\r\n" + hhcn_wn[i] + " " + gme_mp_prnt("HHCN") + " over " + hhcn_ls[i]
	}
	if len(hhcn_wn) > 0 {
		put_in_g_string("HHCN", hhcn_wn, hhcn_ls)
		rsn = strings.Trim(rsn, "\r\n")
	}
	for i := 0; i < len(hhcn_orig); i++ {
		if strings.Contains(";"+strings.Join(hhcn_wn, ";")+";", ";"+hhcn_orig[i]+";") {
			continue
		}
		if strings.Contains(";"+strings.Join(hhcn_ls, ";")+";", ";"+hhcn_orig[i]+";") {
			continue
		}
		hhcn_nn = append(hhcn_nn, hhcn_orig[i])
	}

	return rsn, append(hhcn_wn, hhcn_nn...)
}

func hh(hh_orig []string, np_excl bool, all_or_none bool) (string, []string) {
	var hh_wn []string
	var hh_ls []string
	var hh_nn []string
	high_prc := -1.0
	excl := false

	for i := 0; i < len(hh_orig); i++ {
		w := 0.0
		tot := 0.0
		for j := 0; j < len(hh_orig); j++ {
			if i == j {
				continue
			}
			for k := 0; k < len(plr_plr_w[hh_orig[i]]); k++ {
				if strings.HasPrefix(plr_plr_w[hh_orig[i]][k], hh_orig[j]+",") {
					w++
					tot++
				}
			}
			for k := 0; k < len(plr_plr_t[hh_orig[i]]); k++ {
				if strings.HasPrefix(plr_plr_t[hh_orig[i]][k], hh_orig[j]+",") {
					w += 0.5
					tot++
				}
			}
			for k := 0; k < len(plr_plr_l[hh_orig[i]]); k++ {
				if strings.HasPrefix(plr_plr_l[hh_orig[i]][k], hh_orig[j]+",") {
					tot++
				}
			}
		}
		if tot == 0 {
			if np_excl || all_or_none {
				return "np_excl", hh_orig
			} else {
				hh_nn = append(hh_nn, hh_orig[i])
			}
		} else {
			if all_or_none {
				if (w == float64(len(hh_orig)-1)) || ((w == 0) && (tot == float64(len(hh_orig)-1))) {
					excl = true
				}
			}
			prc := w / tot
			if math.Abs(high_prc-prc) < 0.000001 {
				hh_wn = append(hh_wn, hh_orig[i])
			} else if prc < high_prc {
				hh_ls = append(hh_ls, hh_orig[i])
			} else {
				hh_ls = append(hh_ls, hh_wn...)
				hh_wn = nil
				hh_wn = append(hh_wn, hh_orig[i])
				high_prc = prc
			}
		}
	}
	if all_or_none && !excl {
		fmt.Println(hh_orig, "np_excl")
		return "np_excl", hh_orig
	}
	rsn := ""
	if len(hh_ls) > 0 {
		rsn = strings.Join(hh_wn, ", ") + " " + gme_mp_prnt("HH") + " over " + strings.Join(hh_ls, ", ")
		put_in_g_string("HH", hh_wn, hh_ls)
	}
	return rsn, append(hh_wn, hh_nn...)
}

func hh_(hh_orig []string, which string) (string, []string) {
	var hh_wn []string
	var hh_ls []string
	var hh_nn []string
	high_pts := -1

	for i := 0; i < len(hh_orig); i++ {
		pts := 0
		for j := 0; j < len(hh_orig); j++ {
			if i == j {
				continue
			}
			for k := 0; k < len(plr_plr_w[hh_orig[i]]); k++ {
				if strings.HasPrefix(plr_plr_w[hh_orig[i]][k], hh_orig[j]+",") {
					arr := strings.Split(plr_plr_w[hh_orig[i]][k], ",")
					p1, _ := strconv.Atoi(arr[1])
					p2, _ := strconv.Atoi(arr[2])
					if which == "sht" {
						if (p1 > 0) && (p2 == 0) {
							pts += 1
						}
					} else {
						pts += p1 - p2
					}
				}
			}
			for k := 0; k < len(plr_plr_l[hh_orig[i]]); k++ {
				if strings.HasPrefix(plr_plr_l[hh_orig[i]][k], hh_orig[j]+",") {
					arr := strings.Split(plr_plr_l[hh_orig[i]][k], ",")
					p1, _ := strconv.Atoi(arr[1])
					p2, _ := strconv.Atoi(arr[2])
					if which == "sht" {
						if (p1 > 0) && (p2 == 0) {
							pts -= 1
						}
					} else {
						pts += p2 - p1
					}
				}
			}
		}

		if high_pts == pts {
			hh_wn = append(hh_wn, hh_orig[i])
		} else if pts < high_pts {
			hh_ls = append(hh_ls, hh_orig[i])
		} else {
			hh_ls = append(hh_ls, hh_wn...)
			hh_wn = nil
			hh_wn = append(hh_wn, hh_orig[i])
			high_pts = pts
		}
	}

	rsn := ""
	if len(hh_ls) > 0 {
		rsn = strings.Join(hh_wn, ", ") + " " + gme_mp_prnt("HH") + " over " + strings.Join(hh_ls, ", ")
		put_in_g_string("HH "+which, hh_wn, hh_ls)
	}
	return rsn, append(hh_wn, hh_nn...)
}

func HH_Dif(remain []string, hh_cur []string) bool {
	for i := 0; i < len(remain); i++ {
		found := false
		for j := 0; j < len(hh_cur); j++ {
			if remain[i] == hh_cur[j] {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}
	return false
}

func tbr_wlt(wlt_orig []string, attr string) (string, []string) {
	high_prc := 0.0
	wn := 0.0
	tot := 0.0
	prc := 0.0
	var wlt_ls []string
	var wlt_wn []string
	var wlt_nn []string
	for i := 0; i < len(wlt_orig); i++ {
		wn = float64(nme_attr_int[wlt_orig[i]+";w;"+attr])
		wn += (float64(nme_attr_int[wlt_orig[i]+";t;"+attr]) / 2)
		tot = float64(nme_attr_int[wlt_orig[i]+";w;"+attr])
		tot += float64(nme_attr_int[wlt_orig[i]+";t;"+attr])
		tot += float64(nme_attr_int[wlt_orig[i]+";l;"+attr])
		if tot > 0 {
			prc = wn / tot
		} else {
			if attr == "o" {
				prc = 0.5
			} else {
				wlt_nn = append(wlt_nn, wlt_orig[i])
				continue
			}
		}
		if math.Abs(high_prc-prc) < 0.000001 {
			wlt_wn = append(wlt_wn, wlt_orig[i])
		} else if prc < high_prc {
			wlt_ls = append(wlt_ls, wlt_orig[i])
		} else {
			wlt_ls = append(wlt_ls, wlt_wn...)
			wlt_wn = nil
			wlt_wn = append(wlt_wn, wlt_orig[i])
			high_prc = prc
		}
	}
	rsn := ""
	if (attr != "o") && (len(wlt_ls) > 0) {
		rsn = strings.Join(wlt_wn, ", ") + " " + gme_mp_prnt(attr) + " over " + strings.Join(wlt_ls, ", ")
		put_in_g_string(attr, wlt_wn, wlt_ls)
	}
	return rsn, append(wlt_wn, wlt_nn...)
}

/*
func tbr_act(act_orig []string, attr string) (string, []string) {
	high_prc := 0.0
	wn := 0.0
	tot := 0.0
	prc := 0.0
	var act_ls []string
	var act_wn []string
	var act_nn []string
	for i := 0; i < len(act_orig); i++ {
		wn = float64(len(act_str[act_orig[i]+";w;"+attr])) //nme_attr_int[wlt_orig[i]+";w;"+attr])
		tot = wn + float64(len(act_str[act_orig[i]+";l;"+attr]))
		if tot > 0 {
			prc = wn / tot
		} else {
			act_nn = append(act_nn, act_orig[i])
			continue
		}
		if math.Abs(high_prc-prc) < 0.000001 {
			act_wn = append(act_wn, act_orig[i])
		} else if prc < high_prc {
			act_ls = append(act_ls, act_orig[i])
		} else {
			act_ls = append(act_ls, act_wn...)
			act_wn = nil
			act_wn = append(act_wn, act_orig[i])
			high_prc = prc
		}
	}
	rsn := ""
	if (attr != "o") && (len(act_ls) > 0) {
		rsn = strings.Join(act_wn, ", ") + " " + gme_mp_prnt(attr) + " over " + strings.Join(act_ls, ", ")
		put_in_g_string(attr, act_wn, act_ls)
	}
	return rsn, append(act_wn, act_nn...)
}*/

func tbr_num(tbrn_orig []string, attr string, IsInt, IsTme bool) (string, []string) {
	high_prc := -99999.0
	prc := 0.0
	var tbrn_ls []string
	var tbrn_wn []string
	var attr2 string
	var attrt string

	if strings.Contains(attr, "-") {
		attr2 = strings.ReplaceAll(attr, "-", ";")
		attrt = strings.ReplaceAll(attr, "-", "-t;")
	} else {
		attr2 = attr + ";"
	}
	for i := 0; i < len(tbrn_orig); i++ {
		if IsInt {
			prc = float64(nme_attr_int[tbrn_orig[i]+";"+attr2])
			prc += (float64(nme_attr_int[tbrn_orig[i]+";"+attrt]) / 2)
		} else if IsTme {
			prc = float64(TmeInt(nme_attr_tme[tbrn_orig[i]+";"+attr2]))
		} else {
			prc = nme_attr_flt[tbrn_orig[i]+";"+attr2]
			prc += (nme_attr_flt[tbrn_orig[i]+";"+attrt] / 2)
		}
		if math.Abs(high_prc-prc) < 0.000001 {
			tbrn_wn = append(tbrn_wn, tbrn_orig[i])
		} else if prc < high_prc {
			tbrn_ls = append(tbrn_ls, tbrn_orig[i])
		} else {
			tbrn_ls = append(tbrn_ls, tbrn_wn...)
			tbrn_wn = nil
			tbrn_wn = append(tbrn_wn, tbrn_orig[i])
			high_prc = prc
		}
	}
	rsn := ""
	if len(tbrn_ls) > 0 {
		rsn = strings.Join(tbrn_wn, ", ") + " " + gme_mp_prnt(attr) + " over " + strings.Join(tbrn_ls, ", ")
		put_in_g_string(attr, tbrn_wn, tbrn_ls)
	}
	return rsn, tbrn_wn
}

func tbr_clsgme(clsgme_orig []string) (string, []string) {
	high_num := clsgm_hsh[clsgme_orig[0]]

	var clsgme_ls []string
	var clsgme_wn []string
	clsgme_wn = append(clsgme_wn, clsgme_orig[0])
	for i := 1; i < len(clsgme_orig); i++ {
		tmp_cmpr_num := clsgme_cmpr(high_num, clsgm_hsh[clsgme_orig[i]])
		if tmp_cmpr_num == 0 {
			clsgme_wn = append(clsgme_wn, clsgme_orig[i])
		} else if tmp_cmpr_num == 1 {
			clsgme_ls = append(clsgme_ls, clsgme_orig[i])
		} else {
			clsgme_ls = append(clsgme_ls, clsgme_wn...)
			clsgme_wn = nil
			clsgme_wn = append(clsgme_wn, clsgme_orig[i])
			high_num = clsgm_hsh[clsgme_orig[i]]
		}
	}
	rsn := ""
	if len(clsgme_ls) > 0 {
		rsn = strings.Join(clsgme_wn, ", ") + " " + gme_mp_prnt("clsgme") + " over " + strings.Join(clsgme_ls, ", ")
		put_in_g_string("clsgme", clsgme_wn, clsgme_ls)
	}
	return rsn, clsgme_wn
}

func clsgme_cmpr(clsgme_one, clsgme_two string) int {
	last_pipe := true
	posA := 0
	i2 := 0
	for i2 < len(clsgme_one) && i2 < len(clsgme_two) {
		if clsgme_one[i2] != clsgme_two[i2] {
			break
		}
		if string(clsgme_one[i2]) == "|" {
			last_pipe = true
			posA = i2 + 1
		} else if string(clsgme_one[i2]) == ":" {
			last_pipe = false
			posA = i2 + 1
		}
		i2++
	}
	if clsgme_one == clsgme_two {
		return 0
	}
	if !last_pipe {
		tempArr := strings.SplitN(clsgme_one[posA:], "|", 2)
		clsgme_one = tempArr[0]
		tempArr = strings.SplitN(clsgme_two[posA:], "|", 2)
		clsgme_two = tempArr[0]
		if StrToInt(clsgme_one) < StrToInt(clsgme_two) {
			return 2
		} else {
			return 1
		}
	} else {
		onetwo := 0
		if (i2 >= len(clsgme_one)) && (i2 >= len(clsgme_two)) {
			return 0
		} else if i2 >= len(clsgme_one) {
			tempArr := strings.SplitN(clsgme_two[(posA):], ":", 2)
			clsgme_two = tempArr[1]
			onetwo = 2
		} else if i2 >= len(clsgme_two) {
			tempArr := strings.SplitN(clsgme_one[posA:], ":", 2)
			clsgme_one = tempArr[1]
			onetwo = 1
		} else {
			tempArr := strings.SplitN(clsgme_one[(posA):], ":", 2)
			one1 := tempArr[0]
			clsgme_one = tempArr[1]
			tempArr = strings.SplitN(clsgme_two[(posA):], ":", 2)
			two2 := tempArr[0]
			clsgme_two = tempArr[1]
			if StrToInt(one1) < StrToInt(two2) {
				onetwo = 1
			} else {
				onetwo = 2
			}
		}

		if onetwo == 2 {
			if strings.HasPrefix(clsgme_two, "-") {
				return 1
			} else {
				return 2
			}
		} else {
			if strings.HasPrefix(clsgme_one, "-") {
				return 2
			} else {
				return 1
			}
		}
	}
}

func tbr_strk(strk_orig []string) (string, []string) {
	high_num := strk_hsh[strk_orig[0]]

	var strk_ls []string
	var strk_wn []string
	strk_wn = append(strk_wn, strk_orig[0])
	for i := 1; i < len(strk_orig); i++ {
		tmp_cmpr_num := strk_cmpr(high_num, strk_hsh[strk_orig[1]])
		if tmp_cmpr_num == 0 {
			strk_wn = append(strk_wn, strk_orig[i])
		} else if tmp_cmpr_num == 1 {
			strk_ls = append(strk_ls, strk_orig[i])
		} else {
			strk_ls = append(strk_ls, strk_wn...)
			strk_wn = nil
			strk_wn = append(strk_wn, strk_orig[i])
			high_num = strk_hsh[strk_orig[i]]
		}
	}
	rsn := ""
	if len(strk_ls) > 0 {
		rsn = strings.Join(strk_wn, ", ") + " " + gme_mp_prnt("strk") + " over " + strings.Join(strk_ls, ", ")
		put_in_g_string("strk", strk_wn, strk_ls)
	}
	return rsn, strk_wn
}

func strk_cmpr(strk_one, strk_two string) int {
	strk_arr1 := strings.Split(strk_one, "|")
	strk_arr2 := strings.Split(strk_two, "|")
	i2 := 0
	for i2 < len(strk_arr1) && i2 < len(strk_arr2) {
		strk_num1 := StrToInt(strings.ReplaceAll(strings.ReplaceAll(strk_arr1[i2], "L", "-"), "W", ""))
		strk_num2 := StrToInt(strings.ReplaceAll(strings.ReplaceAll(strk_arr2[i2], "L", "-"), "W", ""))
		if strk_num1 > strk_num2 {
			return 1
		} else if strk_num1 < strk_num2 {
			return 2
		}
		i2++
	}
	if len(strk_arr1) > len(strk_arr2) {
		if strings.HasPrefix(strk_arr1[i2], "W") {
			return 1
		} else {
			return 2
		}
	} else if len(strk_arr1) < len(strk_arr2) {
		if strings.HasPrefix(strk_arr2[i2], "W") {
			return 2
		} else {
			return 1
		}
	} else {
		return 0
	}
}

func tbr_cm_wb(tbr_cm_orig []string) (string, []string) {
	return "rsn", tbr_cm_orig
}

func tbr_cm_min(tbr_cm_orig []string, min int, nd bool) (string, []string) {
	high_prc := -99999.0
	prc := 0.0
	var tbr_cm_ls []string
	var tbr_cm_wn []string
	cmm_sch1_arr := strings.Split(strings.Join(plrs[plrsNameToId[tbr_cm_orig[0]]].sch, "|"), "|")
	for i := 0; i < len(cmm_sch1_arr); i++ {
		cmm_sch1_arr[i] = strings.TrimSpace(cmm_sch1_arr[i])
	}
	cmm_sch := "|" + strings.Join(cmm_sch1_arr, "|") + "|"
	cmm_sch1_arr = cmm_sch1_arr[:0]
	cmm_sch = strings.ReplaceAll(cmm_sch, "|open|", "|")
	for i := 1; i < len(tbr_cm_orig); i++ {
		for j := 0; j < len(plrs[plrsNameToId[tbr_cm_orig[i]]].sch); j++ {
			if nd {
				ndtemp := strings.TrimSpace(plrs[plrsNameToId[tbr_cm_orig[i]]].sch[j])
				if (plrs[plrsNameToId[ndtemp]].trn == plrs[plrsNameToId[tbr_cm_orig[i]]].trn) && (plrs[plrsNameToId[ndtemp]].cn == plrs[plrsNameToId[tbr_cm_orig[i]]].cn) && (plrs[plrsNameToId[ndtemp]].dv == plrs[plrsNameToId[tbr_cm_orig[i]]].dv) {
					continue
				}
			}
			if strings.Contains(cmm_sch, "|"+strings.TrimSpace(plrs[plrsNameToId[tbr_cm_orig[i]]].sch[j])+"|") {
				cmm_sch1_arr = append(cmm_sch1_arr, strings.TrimSpace(plrs[plrsNameToId[tbr_cm_orig[i]]].sch[j]))
			}
		}
		if len(cmm_sch1_arr) < min {
			return "", tbr_cm_orig
		}
		cmm_sch = "|" + strings.Join(cmm_sch1_arr, "|") + "|"
		cmm_sch1_arr = cmm_sch1_arr[:0]
	}
	for i := 0; i < len(tbr_cm_orig); i++ {
		tbr_cm_w := 0.0
		tbr_cm_t := 0.0
		for j := 0; j < len(plr_plr_w[tbr_cm_orig[i]]); j++ {
			tmp1_arr := strings.Split(plr_plr_w[tbr_cm_orig[i]][j], ",")
			if strings.Contains(cmm_sch, "|"+tmp1_arr[0]+"|") {
				tbr_cm_w++
				tbr_cm_t++
			}
		}
		for j := 0; j < len(plr_plr_l[tbr_cm_orig[i]]); j++ {
			tmp1_arr := strings.Split(plr_plr_l[tbr_cm_orig[i]][j], ",")
			if strings.Contains(cmm_sch, "|"+tmp1_arr[0]+"|") {
				tbr_cm_t++
			}
		}
		for j := 0; j < len(plr_plr_t[tbr_cm_orig[i]]); j++ {
			tmp1_arr := strings.Split(plr_plr_t[tbr_cm_orig[i]][j], ",")
			if strings.Contains(cmm_sch, "|"+tmp1_arr[0]+"|") {
				tbr_cm_w += 0.5
				tbr_cm_t++
			}
		}
		if tbr_cm_t == 0 {
			return "", tbr_cm_orig
		}

		prc = tbr_cm_w / tbr_cm_t

		if math.Abs(high_prc-prc) < 0.000001 {
			tbr_cm_wn = append(tbr_cm_wn, tbr_cm_orig[i])
		} else if prc < high_prc {
			tbr_cm_ls = append(tbr_cm_ls, tbr_cm_orig[i])
		} else {
			tbr_cm_ls = append(tbr_cm_ls, tbr_cm_wn...)
			tbr_cm_wn = nil
			tbr_cm_wn = append(tbr_cm_wn, tbr_cm_orig[i])
			high_prc = prc
		}
	}

	rsn := ""
	if len(tbr_cm_ls) > 0 {
		rsn = strings.Join(tbr_cm_wn, ", ") + " " + gme_mp_prnt("CMM") + " over " + strings.Join(tbr_cm_ls, ", ")
		put_in_g_string("CMM", tbr_cm_wn, tbr_cm_ls)
	}
	return rsn, tbr_cm_wn
}

func put_in_g_string(att string, arr_w []string, arr_l []string) {
	arr_tmp := append(arr_w, arr_l...)
	for g_i := 0; g_i < len(arr_tmp); g_i++ {
		if !strings.Contains(g_string[arr_tmp[g_i]], ";"+att+";") {
			g_string[arr_tmp[g_i]] += ";" + att + ";"
		}
	}
}

func SortNamesInMap(mymap map[string][]string) []string {
	var arr_s []string
	for one_s := range mymap {
		arr_s = append(arr_s, one_s)
	}
	sort.Strings(arr_s)
	return arr_s
}

func SortIntsInMap(mymap map[int][]string) []int {
	var arr_i []int
	for one_i := range mymap {
		arr_i = append(arr_i, one_i)
	}
	sort.Ints(arr_i)
	return arr_i
}

func lowest_nme_attr(arr []string) int {
	arr_a1 := strings.Split(arr[0], ",")
	c_indx := 0
	for arr_i := 1; arr_i < len(arr); arr_i++ {
		arr_a2 := strings.Split(arr[arr_i], ",")
		for arr_i2 := 1; arr_i2 < len(arr_a2); arr_i2++ {
			one, err1 := strconv.Atoi(arr_a1[arr_i2])
			two, err2 := strconv.Atoi(arr_a2[arr_i2])
			if (err1 == nil) && (err2 == nil) {
				if one < two {
					break
				} else if one > two {
					c_indx = arr_i
					arr_a1 = strings.Split(arr[c_indx], ",")
					break
				}
			} else {
				if arr_a1[arr_i2] < arr_a2[arr_i2] {
					break
				} else if arr_a1[arr_i2] > arr_a2[arr_i2] {
					c_indx = arr_i
					arr_a1 = strings.Split(arr[c_indx], ",")
					break
				}
			}
		}
	}
	return c_indx
}

func inc_nme_attr_int(key string, value int) {
	if _, haskey := nme_attr_int[key]; !haskey {
		nme_attr_int[key] = value
	} else {
		nme_attr_int[key] += value
	}
	if (strings.Contains(key, "t")) && (nme_attr_int[key] == 2) {
		nme_attr_int[key] = 0
		nme_attr_int[strings.ReplaceAll(key, ";t;", ";w;")] += 1
		nme_attr_int[strings.ReplaceAll(key, ";t;", ";l;")] += 1
	}
}

func inc_nme_attr_tme(key string, value string) {
	if _, haskey := nme_attr_tme[key]; !haskey {
		nme_attr_tme[key] = TmeConv("0:00", value)
	} else {
		nme_attr_tme[key] = TmeConv(nme_attr_tme[key], value)
	}
}

func StrToInt(str string) int {
	var itmp int
	var err error
	if itmp, err = strconv.Atoi(str); err != nil {
		fmt.Printf("Invalid int %s\n", str)
		return 0
	}
	return itmp
}

func FloatToTme(flt float64) string {
	flt *= 1440
	flt = math.Round(flt)
	min := int(flt) / 60
	sec := int(flt) % 60
	if sec < 10 {
		return strconv.Itoa(min) + ":0" + strconv.Itoa(sec)
	} else {
		return strconv.Itoa(min) + ":" + strconv.Itoa(sec)
	}
}

func wlt_tbrkr_i(is_tbrk bool, plyrw, plyrl, tbrs1, tbrs2 string, value int) {
	if is_tbrk {
		inc_nme_attr_int(plyrw+";"+tbrs1+";"+tbrs2, value)
		inc_nme_attr_int(plyrl+";"+tbrs1+";"+tbrs2, -value)
	}
}

func wlt_tbrkr_t(is_tbrk bool, plyrw, plyrl, tbrs1, tbrs2, value string) {
	if is_tbrk {
		inc_nme_attr_tme(plyrw+";"+tbrs1+";"+tbrs2, value)
		inc_nme_attr_tme(plyrl+";"+tbrs1+";"+tbrs2, "("+value+")")
	}
}

func TmeConv(tme1, tme2 string) string {
	tme1_i := TmeInt(tme1)
	tme2_i := TmeInt(tme2)
	tmer_i := tme1_i + tme2_i
	tmer := TmeStr(tmer_i)
	return tmer
}

func TmeInt(tme string) int {
	isneg := false
	if tme == "" {
		tme = "0:00"
	}
	if regexp.MustCompile(`^\(\d+:\d+\)$`).MatchString(tme) {
		isneg = true
		tme = strings.Replace(tme, ")", "", 1)
		tme = strings.Replace(tme, "(", "", 1)
	} else if regexp.MustCompile(`^\d+:\d+$`).MatchString(tme) {
		isneg = false
	} else {
		fmt.Println("Invalid time:", tme)
		return 0
	}
	minsec := strings.Split(tme, ":")
	itme := StrToInt(minsec[0])*60 + StrToInt(minsec[1])
	if isneg {
		itme *= -1
	}
	return itme
}

func TmeStr(nmbr int) string {
	isneg := false
	if nmbr < 0 {
		isneg = true
		nmbr *= -1
	}
	min := nmbr / 60
	sec := nmbr % 60
	var num_str string
	if sec < 10 {
		num_str = strconv.Itoa(min) + ":0" + strconv.Itoa(sec)
	} else {
		num_str = strconv.Itoa(min) + ":" + strconv.Itoa(sec)
	}
	if isneg {
		num_str = "(" + num_str + ")"
	}
	return num_str
}

func nxt_wk(trnId int, nxtRnd int) error {
	for nme, id := range plrsNameToId {
		//plyr1 := plyr_all[plyr_i]
		//plyr2 := plyr_all[plyr_i+1]
		ww := float64(nme_attr_int[nme+";w;o"]) + float64(nme_attr_int[nme+";t;o"])/2
		wwInt := math.Floor(ww)
		wwTie := 0
		if (wwInt - ww) > 0.1 {
			wwTie = 1
		}

		sql := "update gms_plyrtm set rec = ?, tie = ? " +
			"where plyrtm_id = ? and " +
			"gms_id in (select id from gms where tourn_id = ? and rnd = ? and coalesce(has_stnd,'Y') = 'Y')"
		params := []interface{}{wwInt, wwTie, id, trnId, nxtRnd}
		_, err := db.ChangDB(sql, params)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
func lst_wk(wk int, dir string, plrs_all_str string, tbrk_arr []string) {
	buffer := ""
	plyr_all := strings.Split(plrs_all_str, ";")
	for plyr_i := 1; plyr_i < (len(plyr_all) - 1); plyr_i += 2 {
		plyr1 := plyr_all[plyr_i]
		plyr2 := plyr_all[plyr_i+1]
		rnng_str := ""
		wnr_p := ""
		if float64(nme_attr_int[plyr1+";w;arch"]) > 0 {
			wnr_p = plyr1
		} else {
			wnr_p = plyr2
		}
		rnng_str += "\t" + wnr_p
		for tbrk_i := 0; tbrk_i < len(tbrk_arr); tbrk_i++ {
			if strings.HasPrefix(tbrk_arr[tbrk_i], "gn") || strings.HasPrefix(tbrk_arr[tbrk_i], "ot") || strings.HasPrefix(tbrk_arr[tbrk_i], "sht") || strings.HasPrefix(tbrk_arr[tbrk_i], "pts") {
				onetwo := strings.Split(tbrk_arr[tbrk_i], "-")
				rnng_str += "\t" + strconv.Itoa(nme_attr_int[wnr_p+";"+onetwo[0]+";arch"])
			}
		}
		buffer += plyr1 + "\t" + plyr2 + "\t" + rnng_str + "\r\n"
	}
	if err := os.WriteFile(dir+strconv.Itoa(wk)+"-arch.txt", []byte(buffer), 0644); err != nil {
		fmt.Println(err)
	}
}*/

func gms(wk int, cndv_wc string, sorted_cndv []string, dv_tbr_arr []string) error {
	cndvwc := strings.SplitN(cndv_wc, " ", 2)

	if len(cndvwc) < 2 {
		cndvwc = append(cndvwc, "2")
	}
	if strings.Contains(cndvwc[1], "t") {
		cndvwc[1] = strings.ReplaceAll(cndvwc[1], "t", "")
	}
	wc := StrToInt(cndvwc[1])
	var srt_indxs []int

	//var gFile *xlsx.File
	//var gSheet *xlsx.Sheet
	//gFile = xlsx.NewFile()
	prev := ""
	curr := ""

	for cndv_i := 0; cndv_i < len(sorted_cndv); cndv_i++ {
		switch cndvwc[0] {
		case "cn":
			cnanddv := strings.Split(sorted_cndv[cndv_i], ";")
			curr = cnanddv[0]
		case "dv":
			curr = sorted_cndv[cndv_i]
		}
		if curr != prev {
			if prev != "" {
				err := gms_sht(wk, prev, wc, srt_indxs, sorted_cndv, dv_tbr_arr)
				if err != nil {
					return err
				}
			}
			prev = curr
			srt_indxs = nil
		}
		srt_indxs = append(srt_indxs, cndv_i)
	}
	if strings.HasPrefix(cndvwc[0], "A") || (cndvwc[0] == "CD") {
		prev = "A"
	}
	if prev != "" {
		err := gms_sht(wk, prev, wc, srt_indxs, sorted_cndv, dv_tbr_arr)
		if err != nil {
			return err
		}
	}
	return nil
}

func gms_sht(wk int, cnordv string, wc int, srt_indxs []int, sorted_cndv []string, dv_tbr_arr []string) error {
	//var err error
	var col_lens []int
	var col_indxs []int
	total := 0

	dv_nums := make(map[string]float64)
	wc_cnt := 0 //wc
	wc_num := -1.0
	elm_str := ""

	var tempConfGms model.ConfGMS
	tempConfGms.Conf = cnordv
	for i := 0; i < len(srt_indxs); i++ {
		col_lens = append(col_lens, len(cndv_to_plrs[sorted_cndv[srt_indxs[i]]]))
		col_indxs = append(col_indxs, 0)
		total += col_lens[i]
	}
	for i := 0; i < total; i++ {
		//wns := -1
		plr := ""
		chs_j := 0

		high_v := -1.0
		for j := 0; j < len(srt_indxs); j++ {
			if col_indxs[j] >= col_lens[j] {
				continue
			}
			plr = cndv_to_plrs[sorted_cndv[srt_indxs[j]]][col_indxs[j]]
			v_num := float64(nme_attr_int[plr+";w;o"])

			if nme_attr_int[plr+";t;o"] == 1 {
				v_num += .5
			}
			if v_num > high_v {
				high_v = v_num
				chs_j = j
			}
		}
		plr = cndv_to_plrs[sorted_cndv[srt_indxs[chs_j]]][col_indxs[chs_j]]
		plr_i := plrsNameToId[plr]
		v := strconv.Itoa(nme_attr_int[plr+";w;o"])
		v_num := float64(nme_attr_int[plr+";w;o"])
		if nme_attr_int[plr+";t;o"] == 1 {
			v += ".5"
			v_num += .5
		}

		elm := false
		plyr_extd_sch_arr := strings.Split(strings.Trim(plyr_extd_sch[plr], ";"), ";")
		var plr_extd_arr []string
		if len(plrs[plr_i].sch) != len(plyr_extd_sch_arr) {
			plr_extd_arr = append(plrs[plr_i].sch, plyr_extd_sch_arr...)
		} else {
			plr_extd_arr = plrs[plr_i].sch
		}
		gms_lft := float64(len(plr_extd_arr) - wk + 1)
		if dv_nums[fmt.Sprintf("%d;%d", plrs[plr_i].cn, plrs[plr_i].dv)] == 0 {
			dv_nums[fmt.Sprintf("%d;%d", plrs[plr_i].cn, plrs[plr_i].dv)] = v_num
		} else {
			wc_cnt++
			if wc_cnt == wc {
				wc_num = v_num
			} else if wc_num >= 0 {
				wc_diff := wc_num - v_num
				dv_diff := dv_nums[fmt.Sprintf("%d;%d", plrs[plr_i].cn, plrs[plr_i].dv)] - v_num
				if wc_diff > gms_lft && dv_diff > gms_lft {
					elm = true
				}
			}
		}

		if elm {
			elm_str += ";" + plr
			col_indxs[chs_j]++
			continue
		}
		var tempPlyrtmGms model.PlyrtmGMS
		if curplyr_depc[plr] == "" {
			tempPlyrtmGms.Plyrtm = plr
		} else {
			tempPlyrtmGms.Plyrtm = curplyr_depc[plr] + "-" + plr
		}
		tempPlyrtmGms.Wn = v_num

		for tbr_i := 1; tbr_i < len(dv_tbr_arr); tbr_i++ {
			a := dv_tbr_arr[tbr_i]
			if (a == "dv") || (a == "cm") || (a == "cn") || (a == "act") || (a == "actcm") || (a == "actdv") {
				v_num := float64(nme_attr_int[plr+";w;"+a])
				if nme_attr_int[plr+";t;"+a] == 1 {
					v_num += 0.5
				}
				tempPlyrtmGms.Wns = append(tempPlyrtmGms.Wns, v_num)
			}
		}

		tfw := nme_attr_int[plr+";w;o"]
		if nme_attr_int[plr+";t;o"] == 1 {
			tfw++
		}
		tfw += nme_attr_int[plr+";l;o"]

		for wki := wk; wki <= len(plr_extd_arr); wki++ {
			if tfw < wki {
				tempPlyrtmGms.Remain = append(tempPlyrtmGms.Remain, plr_extd_arr[wki-1])
			}
		}
		col_indxs[chs_j]++

		for wlt_i := 0; wlt_i < 3; wlt_i++ {
			var plr_gms_wlt []string
			switch wlt_i {
			case 0:
				plr_gms_wlt = strings.Split(plr_plr_gms_w[plr], ";")
			case 1:
				plr_gms_wlt = strings.Split(plr_plr_gms_l[plr], ";")
			default:
				plr_gms_wlt = strings.Split(plr_plr_gms_t[plr], ";")
			}

			for wki := 1; wki < len(plr_gms_wlt); wki++ {
				plr2_arr := strings.SplitN(plr_gms_wlt[wki], ",", 2)
				gms_lft2 := float64(len(plr_extd_arr) - wk + 1)

				v_num2 := float64(nme_attr_int[plr2_arr[0]+";w;o"])
				if nme_attr_int[plr2_arr[0]+";t;o"] == 1 {
					v_num2 += .5
				}
				if (v_num > v_num2) && ((v_num - v_num2) > gms_lft2) {
					continue
				}
				if (v_num < v_num2) && ((v_num2 - v_num) > gms_lft) {
					continue
				}

				switch wlt_i {
				case 0:
					tempPlyrtmGms.W = append(tempPlyrtmGms.W, plr2_arr[0])
				case 1:
					tempPlyrtmGms.L = append(tempPlyrtmGms.L, plr2_arr[0])
				default:
					tempPlyrtmGms.T = append(tempPlyrtmGms.T, plr2_arr[0])
				}
			}
		}
		tempConfGms.PlyrtmGms = append(tempConfGms.PlyrtmGms, tempPlyrtmGms)
	}
	stnds.ConfGms = append(stnds.ConfGms, tempConfGms)
	return nil
}

func stnd(wk int, cndv_wc string, sorted_cndv []string, cn_tbr_arr []string) {
	var cndvwc []string
	if strings.HasPrefix(cndv_wc, "X") {
		cndvwc = append(cndvwc, "cn", "-42")
	} else {
		cndvwc = strings.SplitN(cndv_wc, " ", 2)
		if len(cndvwc) < 2 {
			cndvwc = append(cndvwc, "2")
		}
	}

	HasT := false
	if strings.Contains(cndvwc[1], "t") {
		HasT = true
		cndvwc[1] = strings.ReplaceAll(cndvwc[1], "t", "")
	}
	wc := StrToInt(cndvwc[1])
	var srt_indxs []int
	prev := ""
	curr := ""

	for cndv_i := 0; cndv_i < len(sorted_cndv); cndv_i++ {
		switch cndvwc[0] {
		case "cn":
			cnanddv := strings.Split(sorted_cndv[cndv_i], ";")
			curr = cnanddv[0]
		case "dv":
			curr = sorted_cndv[cndv_i]
		}
		if curr != prev {
			if prev != "" {
				stnd_cn(wk, wc, HasT, srt_indxs, sorted_cndv, cn_tbr_arr, cndvwc[0], prev)
			}
			prev = curr
			srt_indxs = nil
		}
		srt_indxs = append(srt_indxs, cndv_i)
	}
	if strings.HasPrefix(cndvwc[0], "A") || (cndvwc[0] == "CD") {
		prev = "A"
	}
	if prev != "" {
		stnd_cn(wk, wc, HasT, srt_indxs, sorted_cndv, cn_tbr_arr, cndvwc[0], prev)
		/*var sprsd []string
		sprsd = append(sprsd, running_sht)
		os.MkdirAll(dir+"stdn", os.ModePerm)
		if strings.Contains(dir, "/") {
			Sprsht_out(dir+"stdn/stdn"+strconv.Itoa(wk)+".xlsx", sprsd)
		} else {
			Sprsht_out(dir+"stdn\\stdn"+strconv.Itoa(wk)+".xlsx", sprsd)
		}*/
	}
}

func stnd_cn(wk int, wc int, HasT bool, srt_indxs []int, sorted_cndv []string, cn_tbr_arr []string, dvwc string, conf string) {
	var arr_dv []string
	var arr_wc []string

	for srt_i := 0; srt_i < len(srt_indxs); srt_i++ {
		arr_dv = append(arr_dv, cndv_to_plrs[sorted_cndv[srt_indxs[srt_i]]][0])
		/*if glb_elm == "" {
			l := len(cndv_to_plrs[sorted_cndv[srt_indxs[srt_i]]])
			temp = strings.Join(cndv_to_plrs[sorted_cndv[srt_indxs[srt_i]]][1:l-elm_nm], ";")
		} else {*/
		temp := strings.Join(cndv_to_plrs[sorted_cndv[srt_indxs[srt_i]]][1:], ";")
		//}
		arr_wc = append(arr_wc, temp)
	}

	var tmp_dvwc []string
	if dvwc != "CD" && !strings.HasPrefix(dvwc, "ACN") {
		tmp_dvwc = sort_plrs(arr_dv, cn_tbr_arr, 0, false, true, false, false)
	} else {
		tmp_dvwc = CN(arr_dv, cn_tbr_arr, 0, false, true, false, false)
	}

	var tempConfStnd model.ConfStnd
	tempConfStnd.Conf = conf
	tempType := ""
	if (dvwc != "CD") && !strings.HasPrefix(dvwc, "ACN") {
		tempType = "DV"
	} else {
		tempType = "CN"
	}
	dvw := false

	for x := 0; x < len(tmp_dvwc); x++ {
		if ((dvwc == "CD") || strings.HasPrefix(dvwc, "ACN")) && (x == (len(tmp_dvwc) / 2)) {
			tempType = "DV"
			if HasT && strings.HasSuffix(dvwc, "DVW") {
				dvw = true
			}
		}
		if dvw && (nme_attr_int[tmp_dvwc[x]+";w;o"] <= nme_attr_int[tmp_dvwc[x]+";l;o"]) {
			continue
		}

		var tempRnkStnd model.RnkStnd
		tempRnkStnd.Rnk = x + 1
		tempRnkStnd.Name = tmp_dvwc[x]
		tmp_w := float64(nme_attr_int[tmp_dvwc[x]+";w;o"]) + float64(nme_attr_int[tmp_dvwc[x]+";t;o"])/2
		tmp_l := float64(nme_attr_int[tmp_dvwc[x]+";l;o"]) + float64(nme_attr_int[tmp_dvwc[x]+";t;o"])/2
		tempRnkStnd.Rec = PrettyFormat(tmp_w) + "--" + PrettyFormat(tmp_l)
		tempRnkStnd.GFields = g_string_eval(tmp_dvwc[x])
		tempRnkStnd.DvRsn = strings.Split(plr_rsn_dv[tmp_dvwc[x]][len(plr_rsn_dv[tmp_dvwc[x]])-1], "\n")
		/*if len(plr_rsn_cn[tmp_dvwc[x]]) > 1 {
			tempRnkStnd.CnRsn = strings.Split(plr_rsn_cn[tmp_dvwc[x]][len(plr_rsn_cn[tmp_dvwc[x]])-2], "\n")
		}*/
		if len(plr_rsn_cn[tmp_dvwc[x]]) > 0 {
			tempRnkStnd.CnRsn = strings.Split(plr_rsn_cn[tmp_dvwc[x]][len(plr_rsn_cn[tmp_dvwc[x]])-1], "\n")
		}

		if tempType == "CN" {
			tempConfStnd.CNRnkStnd = append(tempConfStnd.CNRnkStnd, tempRnkStnd)
		} else {
			tempConfStnd.DVRnkStnd = append(tempConfStnd.DVRnkStnd, tempRnkStnd)
		}
	}

	tmp_dvwc = sort_plrs(arr_wc, cn_tbr_arr, wc, false, false, false, false)

	for x := 0; x < len(tmp_dvwc) && ((x < wc) || (wc == -42)); x++ {
		var tempRnkStnd model.RnkStnd
		tempRnkStnd.Rnk = x + 1
		tempRnkStnd.Name = tmp_dvwc[x]
		tmp_w := float64(nme_attr_int[tmp_dvwc[x]+";w;o"]) + float64(nme_attr_int[tmp_dvwc[x]+";t;o"])/2
		tmp_l := float64(nme_attr_int[tmp_dvwc[x]+";l;o"]) + float64(nme_attr_int[tmp_dvwc[x]+";t;o"])/2
		tempRnkStnd.Rec = PrettyFormat(tmp_w) + "--" + PrettyFormat(tmp_l)
		tempRnkStnd.GFields = g_string_eval(tmp_dvwc[x])
		tempRnkStnd.DvRsn = strings.Split(plr_rsn_dv[tmp_dvwc[x]][len(plr_rsn_dv[tmp_dvwc[x]])-1], "\n")
		tempRnkStnd.CnRsn = strings.Split(plr_rsn_cn[tmp_dvwc[x]][len(plr_rsn_cn[tmp_dvwc[x]])-1], "\n")
		tempConfStnd.WCRnkStnd = append(tempConfStnd.WCRnkStnd, tempRnkStnd)
	}
	stnds.ConfStnd = append(stnds.ConfStnd, tempConfStnd)
}

func PrettyFormat(f float64) string {
	s := strconv.FormatFloat(f, 'f', 10, 64)
	if strings.Contains(s, ".") {
		for strings.HasSuffix(s, "0") {
			s = strings.TrimSuffix(s, "0")
		}
		s = strings.TrimSuffix(s, ".")
	}
	return s
}

func g_string_eval(p string) []string {
	rslt := ""
	for g_i := 0; g_i < len(g_string_ord); g_i++ {
		if strings.Contains(g_string[p], ";"+g_string_ord[g_i]+";") && !strings.Contains(rslt, ";"+gme_mp_prnt(g_string_ord[g_i])+"@") {
			if strings.Contains(rslt, ";"+gme_mp_prnt(strings.ReplaceAll(g_string_ord[g_i], "-o", ""))+"@") {
				continue
			}
			rslt += ";" + gme_mp_prnt(strings.ReplaceAll(g_string_ord[g_i], "-o", "")) + "@"
			a := g_string_ord[g_i]
			if (a == "o") || (a == "dv") || (a == "cm") || (a == "cn") {
				t := ""
				if nme_attr_int[p+";t;"+a] == 1 {
					t = ".5"
				}
				v := strconv.Itoa(nme_attr_int[p+";w;"+a]) + t + "--"
				v += strconv.Itoa(nme_attr_int[p+";l;"+a]) + t
				rslt += v
				//} else if (a == "act") || (a == "actcn") || (a == "actdv") {
				//	v := strconv.Itoa(len(act_str[p+";w;"+a])) + "--" + strconv.Itoa(len(act_str[p+";l;"+a]))
				//	rslt += v
			} else if (a == "swp") || (a == "flg") || strings.HasPrefix(a, "hh") {
				continue
			} else if a == "tme" {
				v := nme_attr_tme[p+";"+a+";"]
				v = "-" + v + "-"
				rslt += v
			} else if strings.HasPrefix(a, "tme-") {
				vs := strings.SplitN(a, "-", 2)
				v := nme_attr_tme[p+";"+vs[0]+";"+vs[1]]
				v = "-" + v + "-"
				rslt += v
			} else if strings.HasPrefix(a, "rnk") {
				vs := strings.Split(a, "-")
				v := ""
				if len(vs) > 1 {
					v = plyr_rnk_dsc[p+";"+vs[1]]
				} else {
					v = plyr_rnk_dsc[p+";"]
				}
				pos := strings.Index(v, "#")
				rslt += v[:pos]
			} else if strings.HasPrefix(a, "strk") {
				rslt += strk_hsh[p]
			} else if strings.HasPrefix(a, "clsgm") {
				rslt += clsgm_hsh[p]
			} else if strings.Contains(a, "(") {
				tmp_arr := strings.SplitN(a, "(", 2)
				attr := tmp_arr[0]
				tmp_arr[1] = strings.Trim(tmp_arr[1], ")")
				tmp_arr = strings.Split(tmp_arr[1], "-")
				v := 0
				for tmp_i := 0; tmp_i < len(tmp_arr); tmp_i++ {
					if tmp_arr[tmp_i] == "o" {
						v = nme_attr_int[p+";"+attr+";"]
					} else {
						v = nme_attr_int[p+";"+attr+";"+tmp_arr[tmp_i]]
					}
					rslt += strconv.Itoa(v)
				}
			} else {
				vs := strings.SplitN(a, "-", 2)
				v := ""
				if (len(vs) == 1) || (vs[1] == "o") {
					v = strconv.Itoa(nme_attr_int[p+";"+vs[0]+";"])
				} else {
					v = strconv.Itoa(nme_attr_int[p+";"+vs[0]+";"+vs[1]])
				}
				if strings.HasPrefix(v, "-") {
					v = "(" + v[1:] + ")"
				} else if nme_attr_int[p+";"+a+"-t;"] == 1 {
					v += ".5"
				}
				/*if (len(vs)>1) && (vs[1] == "pt") {
				  for bsc_i := 0; bsc_i<len(dv_tbr_bsc); bsc_i++ {
				  v2 := strconv.Itoa(nme_attr_int[p+";"+vs[0]+";"+dv_tbr_bsc[bsc_i]])
				  if strings.HasPrefix(v2, "-") {
				  v2 = "("+v2[1:]+")"
				  }
				  v += "--"+v2
				  }
				  }*/
				rslt += v
			}
		}
	}
	rslt = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(rslt, "@", " "), ";", "\r\n"))
	return strings.Split(rslt, "\n")
}

/*
	func plr_rsn_eval(arr []string) []string {
		var result []string
		for _, val := range arr {
			if val == "" {
				continue
			}
			result = append(result, val)
		}

		return result
	}

	func gme_mp_strt() {
		dbl_gme_mp("r", "Week ")
		gme_mp["afceast"] = "AFC East"
		gme_mp["afcnorth"] = "AFC North"
		gme_mp["afcsouth"] = "AFC South"
		gme_mp["afcwest"] = "AFC West"
		gme_mp["nfceast"] = "NFC East"
		gme_mp["nfcnorth"] = "NFC North"
		gme_mp["nfcsouth"] = "NFC South"
		gme_mp["nfcwest"] = "NFC West"
		gme_mp["pl"] = "place"
		gme_mp["t1"] = "team"
		gme_mp["o"] = "overall"
		gme_mp["dv"] = "division"
		gme_mp["cm"] = "common (division)"
		gme_mp["CMM"] = "common (conference)"
		gme_mp["cn"] = "conference"
		gme_mp["sov"] = "strength of victory"
		gme_mp["sos"] = "strength of schedule"
		gme_mp["rnk-cn"] = "conference rank points for and against"
		gme_mp["rnk"] = "overall rank points for and against"
		gme_mp["pts-cn"] = "net points conference"
		gme_mp["pts"] = "net points overall"
		gme_mp["nt"] = "net touchdowns"
		gme_mp["sds"] = "last year place"

		dbl_gme_mp("T1", "Away Team")
		dbl_gme_mp("T2", "Home Team")
		dbl_gme_mp("W", "Win")
		dbl_gme_mp("Pt", "Net Points")
		dbl_gme_mp("Pf", "Points For")
		dbl_gme_mp("Pa", "Points Against")
		dbl_gme_mp("Nt", "Net Touchdowns")
		dbl_gme_mp("Rec", "Record")
	}

	func dbl_gme_mp(one string, two string) {
		gme_mp[one] = two
		gme_mp_rev[two] = one
	}
*/
func gme_mp_prnt(key string) string {
	if is_gme_mp && (gme_mp[key] != "") {
		return gme_mp[key]
	} else {
		return key
	}
}

/*
func gme_mp_prnt_rev(key string) string {
	if is_gme_mp && (gme_mp_rev[key] != "") {
		return gme_mp_rev[key]
	} else {
		return key
	}
}

func get_sch_extd(fname string, wh_wk int) {
	var spr_in_file *xlsx.File
	var sheet *xlsx.Sheet
	var cell *xlsx.Cell
	var err error
	spr_in_file, err = xlsx.OpenFile(fname)
	if err != nil {
		fmt.Println("Can not open input spr " + fname)
	}

	//var sprsd [] string
	for sht_i := range spr_in_file.Sheets {
		n := spr_in_file.Sheets[sht_i].Name

		if (n == "msc") || strings.HasPrefix(n, "Sheet") {
			continue
		}
		n = strings.ReplaceAll(n, "r", "")
		var n_i int
		n_i, err := strconv.Atoi(n)
		if (err != nil) || (n_i <= wh_wk) {
			continue
		}

		sheet = spr_in_file.Sheets[sht_i]
	L:
		for row := 1; ; row++ {
			for col := 0; col < 2; col++ {
				cell, _ = sheet.Cell(row, col)
				strings.TrimSpace(cell.Value)
				if strings.TrimSpace(cell.Value) == "" {
					if col == 0 {
						break L
					} else {
						break
					}
				}

				if col == 1 {
					tmp := cell.Value
					cell, _ = sheet.Cell(row, 0)
					other := cell.Value
					plyr_extd_sch[tmp] = other + ";" + plyr_extd_sch[tmp]
					plyr_extd_sch[other] = tmp + ";" + plyr_extd_sch[other]
				}
			}
		}
	}
}

func Re_sch_xls(fname string, wh_wk int) string {
	var spr_in_file *xlsx.File
	var sheet *xlsx.Sheet
	var cell *xlsx.Cell
	var err error
	spr_in_file, err = xlsx.OpenFile(fname)
	if err != nil {
		fmt.Println("Can not open input spr " + fname)
	}

	var sprsd []string
	result := ""
	get_return := false
	for sht_i := range spr_in_file.Sheets {
		n := spr_in_file.Sheets[sht_i].Name
		if (n == "msc") || strings.HasPrefix(n, "Sheet") {
			continue
		}
		if (wh_wk > 9) && strings.HasSuffix(n, strconv.Itoa(wh_wk)) {
			get_return = true
		} else if (wh_wk < 10) && strings.HasSuffix(n, strconv.Itoa(wh_wk)) && !strings.HasSuffix(n, strconv.Itoa(10+wh_wk)) {
			get_return = true
		} else {
			get_return = false
		}

		running_sht := n
		sheet = spr_in_file.Sheets[sht_i]
		tme_col := -1
	L:
		for row := 0; ; row++ {
			for col := 0; ; col++ {
				cell, _ = sheet.Cell(row, col)
				strings.TrimSpace(cell.Value)
				if strings.TrimSpace(cell.Value) == "" {
					if col == 0 {
						break L
					} else {
						break
					}
				}
				tmp := cell.Value
				if (row == 0) && (tmp == "GN") {
					tme_col = col
				} else if (row > 0) && (col == tme_col) {
					f, err2 := strconv.ParseFloat(tmp, 64)
					if err2 != nil {
						return ""
					}
					tmp = FloatToTme(f)
				}
				if _, hasKey := plyr_stnd_rplc[tmp]; hasKey {
					if get_return && (row > 0) && (col < 2) {
						result += ";" + plyr_stnd_rplc[tmp]
					}
					if col == 1 {
						cell, _ = sheet.Cell(row, 0)
						other := plyr_stnd_rplc[cell.Value]
						plyr_extd_sch[plyr_stnd_rplc[tmp]] = other + ";" + plyr_extd_sch[plyr_stnd_rplc[tmp]]
						plyr_extd_sch[other] = plyr_stnd_rplc[tmp] + ";" + plyr_extd_sch[other]
					}
					running_sht += AddRunningSht(row, col, plyr_stnd_rplc[tmp])
				} else {
					running_sht += AddRunningSht(row, col, tmp)
				}
			}
		}
		sprsd = append(sprsd, running_sht)
	}
	Sprsht_out2(strings.ReplaceAll(fname, ".xls", "-1.xls"), sprsd)
	return result
}*/
