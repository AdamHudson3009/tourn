package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

const (
	lll_fname  = "lll.txt"
	sch_fsname = "sch"
	attr_fname = "attr"
)

var spr_mp_arr []string

func main2() {
	del := "/"
	notdel := "\\"
	dir, err := os.Getwd()
	trn := "1"
	if err != nil {
		fmt.Println(err)
		return
	}
	if strings.Contains(dir, "\\") {
		del = "\\"
		notdel = "/"
	}
	dir += del
	o_dir := dir

	scanned := ""
	fmt.Scanln(&scanned)
	trnArr := strings.SplitN(scanned, "-", 2)

	lll_bfile, err := os.ReadFile(dir + strings.ReplaceAll(lll_fname, ".txt", "-"+trnArr[0]+".txt"))
	lll_is := err == nil
	nfl_trn := false
	if !lll_is {
		if scanned == "nfl-trn" {
			fmt.Scanln(&scanned)
			lll_bfile, err = os.ReadFile(dir + scanned)
			lll_is = (err == nil)
			nfl_trn = true
		} else {
			lll_bfile, err = os.ReadFile(dir + lll_fname)
			lll_is = (err == nil) && (scanned != "AS")
		}
	}
	lll_file := ""

	if lll_is {
		lll_file = fmt.Sprintf("%s", lll_bfile)
		if nfl_trn {
			schtxt_file, err := os.ReadFile(dir + "\\TMaster\\schM.txt")
			schtxt_is := err == nil
			if !schtxt_is {
				nfltrn(dir+"\\TMaster", lll_file)
				build_attr(dir+attr_fname+".xlsx", "Master", dir, false)
			} else {
				build_master_schres(schtxt_file, dir+"TMaster\\"+sch_fsname, dir)
			}
			return
		} else {
			lll_temp := strings.SplitN(lll_file, "\n", 2)
			if len(lll_temp) < 2 {
				fmt.Println("lll.txt must have at least 2 lines")
				return
			}
			trn = strings.ReplaceAll(lll_temp[0], notdel, del)
			trn = strings.TrimSpace(trn)
			//fmt.Printf(trn)
			lll_file = lll_temp[1]
		}
	} else {
		dir += "T" + scanned + del
		fmt.Print("Rnd:")
		fmt.Scanln(&scanned)
		num1 := StrToInt(scanned)
		if num1 < 2 {
			num1 = 1
		}
		attr_wk_assn2(dir, num1, dir+"as1.xlsx")
		os.Exit(0)
	}
	if strings.HasPrefix(strings.ToLower(trn), "mult"+del) {
		tmpArr := strings.Split(trn, del)
		stmp := tmpArr[1]
		dir += stmp + del
		boolSet := false
		areFiles := false
		for i := 2; i < len(tmpArr); i++ {
			_, err := os.Stat(fmt.Sprintf("%s%s\\%s%s.txt", dir, tmpArr[i], sch_fsname, tmpArr[i]))
			if !boolSet {
				boolSet = true
				areFiles = err == nil
			} else if areFiles != (err == nil) {
				fmt.Println("All files must exist or not exist", i, 1)
				return
			}
			_, err = os.Stat(fmt.Sprintf("%s%s\\%s%s.xlsx", dir, tmpArr[i], sch_fsname, tmpArr[i]))
			if areFiles != (err == nil) {
				fmt.Println("All files must exist or not exist", i, 2)
				return
			}
		}
		if areFiles {
			for i := 2; i < len(tmpArr); i++ {
				schtxt_file, _ := os.ReadFile(dir + tmpArr[i] + "\\" + sch_fsname + tmpArr[i] + ".txt")
				build_schres(schtxt_file, dir+tmpArr[i]+"\\"+sch_fsname+tmpArr[i], dir+tmpArr[i]+"\\")
			}
		} else {
			for i := 2; i < len(tmpArr); i++ {
				build_sch(dir, tmpArr[i], lll_file, true)
				build_attr(o_dir+attr_fname+".xlsx", tmpArr[i], dir, true)
			}
			return
		}
	} else if strings.Contains(trn, del) {
		tmparr := strings.Split(trn, del)
		trn = tmparr[len(tmparr)-1]
		tmparr = tmparr[:1]
		stmp := strings.Join(tmparr, del)
		dir += stmp + del
		os.MkdirAll(dir, os.ModePerm)
	}

	schtxt_file, err := os.ReadFile(dir + sch_fsname + trn + ".txt")
	schtxt_is := err == nil

	_, err = os.ReadFile(dir + sch_fsname + trn + ".xlsx")
	schxls_is := err == nil

	if schtxt_is != schxls_is {
		fmt.Println(sch_fsname + trn + ".txt and " + sch_fsname + trn + ".xls should exist")
	} else if schtxt_is {
		if len(trnArr) > 1 {
			RunGME(dir, trnArr[0], trnArr[1])
		} else {
			build_schres(schtxt_file, dir+sch_fsname+trn, dir)
		}
	} else if !lll_is {
		fmt.Println(lll_fname + " or " + sch_fsname + trn + ".txt or " + sch_fsname + trn + ".xls should exist")
	} else {
		build_sch(dir, trn, lll_file, false)
		build_attr(o_dir+attr_fname+".xlsx", trn, dir, false)
	}
}

func build_sch(dir, trn string, lll_file string, mult bool) {
	//$hsh{cn_str}{dv_str}{sds} = nme
	type dv_strct struct {
		sds   []string
		debug string
	}
	type cn_strct struct {
		dv_map map[string]dv_strct
		debug  string
	}
	cn_map := make(map[string]cn_strct)
	odiv := make(map[string][]string)
	lines := strings.Split(lll_file, "\r\n")
	if len(lines) < 4 {
		fmt.Println(lll_fname + " has less than 4 lines.")
		return
	}

	lv_wc := strings.TrimSpace(lines[0])
	dv_tbrk := strings.TrimSpace(lines[1])
	cn_tbrk := strings.TrimSpace(lines[2])
	hdr_spr := strings.TrimSpace(lines[3])
	fmt.Println("dv_tbrk", dv_tbrk)
	fmt.Println("cn_tbrk", cn_tbrk)
	fmt.Println("hdr_spr", hdr_spr)
	is_ftb := false
	is_cnf_cm_cm := false
	elm_lines := ""
	elm_num := 0
	elm_stnd := 0
	plcD_lines := ""
	plcD_num := 0
	plcD_stnd := 0
	add_fle_nme := ""
	tm_plyr := make(map[string]string)
	cnf_cm_cm := make(map[string]string)
	cndv_plyr := make(map[string]string)
	all_plyrs := "|"
	all_plyrs_cpy := ""

	var wks []string
	var fields []string
	var l string
	plyrs_sch := make(map[string][]string)

	maxwk := 0
	maxsds := 0
	maxnme := 0
	maxcndv := 0
	attr_lines := ""
	arch_wk1 := 0
	arch_wk2 := 0
	rvs_cnf := ""
	expr := "^([a-z]|[0-9])+-([a-z]|[0-9])+;"
	cnfld := 0
	dvfld := 1
	if mult {
		expr = "^([A-Z]|[a-z]|[0-9])+-([a-z]|[0-9])+-([a-z]|[0-9])+;"
		cnfld = 1
		dvfld = 2
	} else {
		wks = append(wks, trn)
	}

	for line_i := 4; line_i < len(lines); line_i++ {
		l = strings.TrimSpace(lines[line_i])
		if l == "" {
			continue
		}
		//Extra
		if strings.HasPrefix(l, "attr-w;") || strings.HasPrefix(l, "attr-blw;") || strings.HasPrefix(l, "attr_orig+") {
			attr_lines += l + "\r\n"
			//Extra
		} else if strings.HasPrefix(l, "ftb::") {
			is_ftb = true
			is_gme_mp = true
			gme_mp = make(map[string]string)
			gme_mp_rev = make(map[string]string)
			gme_mp_strt()
			//Conf-Divi
		} else if regexp.MustCompile(expr).MatchString(l) {
			fields = strings.Split(l, ";")
			trncndv := strings.Split(fields[0], "-")
			if mult && (strings.ToLower(trncndv[0]) != strings.ToLower(trn)) {
				continue
			}
			if len(fields[0]) > maxcndv {
				maxcndv = len(fields[0])
			}
			if _, haskey := cn_map[trncndv[cnfld]]; !haskey {
				var cn_tmp cn_strct
				cn_tmp.debug = "Here is cn"
				cn_tmp.dv_map = make(map[string]dv_strct)
				cn_map[trncndv[cnfld]] = cn_tmp
			}
			if _, haskey := cn_map[trncndv[cnfld]].dv_map[trncndv[dvfld]]; !haskey {
				var dv_tmp dv_strct
				dv_tmp.debug = "Here is dv"
				dv_tmp.sds = make([]string, len(fields))
				cn_map[trncndv[cnfld]].dv_map[trncndv[dvfld]] = dv_tmp
			}

			if (len(fields) - 1) > maxsds {
				maxsds = len(fields) - 1
			}
			for field_i := 1; field_i < len(fields); field_i++ {
				f := strings.TrimSpace(fields[field_i])
				if len(f) > maxnme {
					maxnme = len(f)
				}
				cn_map[trncndv[cnfld]].dv_map[trncndv[dvfld]].sds[field_i-1] = f
				if mult {
					tmp := strings.Split(fields[0], "-")
					cndv_plyr[tmp[1]+"-"+tmp[2]] += "|" + f
				} else {
					cndv_plyr[fields[0]] += "|" + f
				}
				all_plyrs += f + "|"
			}
			//Sch
		} else if strings.HasPrefix(l, "rvs;") {
			rvs_cnf = ";" + l[4:] + ";"
			//Sch
		} else if strings.HasPrefix(l, "elm;") {
			tmp := strings.Split(strings.TrimSpace(l), ";")
			if len(tmp) < 3 {
				fmt.Println("Bad err line fmt: " + l)
				return
			}
			elm_lines += strconv.Itoa(maxwk) + ";" + l + "\r\n"
			itmp := 1
			var err error
			if itmp, err = strconv.Atoi(tmp[1]); err != nil {
				fmt.Println("Bad err line fmt: "+l, itmp)
				return
			}
			maxsds -= itmp
			if elm_num == 0 {
				elm_num = maxwk
			}
			if (elm_stnd == 0) && (tmp[2] == "stnd") {
				elm_stnd = maxwk
			}
			/*if (len(tmp)>3) && (tmp[3] == "c->d") {
			    elm_cd = maxwk
			}*/
			//Later
		} else if strings.HasPrefix(l, "plc-d;") {
			tmp := strings.Split(strings.TrimSpace(l), ";")
			plcD_lines += strconv.Itoa(maxwk) + ";" + l + "\r\n"
			if plcD_num == 0 {
				plcD_num = maxwk
			}
			if (plcD_stnd == 0) && (tmp[1] == "stnd") {
				plcD_stnd = maxwk
			}
			//Later
		} else if strings.HasPrefix(l, "cmm;") {
			is_cnf_cm_cm = true
			tmp := strings.SplitN(strings.TrimSpace(l), ";", 3)
			if len(tmp) < 3 {
				fmt.Println("Bad err line fmt: " + l)
				return
			}
			f2 := strings.Split(cndv_plyr[tmp[1]], "|")
			for f2_i := 1; f2_i < len(f2); f2_i++ {
				cnf_cm_cm[f2[f2_i]] = strings.ReplaceAll(tmp[2], ";", "|")
			}
			//Extra
		} else if strings.HasPrefix(l, "tme;") {
			tmp := strings.Split(strings.TrimSpace(l), ";")
			if len(tmp) < 2 {
				fmt.Println("Bad err line fmt: " + l)
				return
			}
			tmp = strings.Split(tmp[1], "-")
			if len(tmp) < 2 {
				fmt.Println("Bad err line fmt: " + l)
				return
			}
			tm_plyr[tmp[1]] = tmp[0]
			//Extra
		} else if strings.HasPrefix(l, "fle;") {
			tmp := strings.Split(strings.TrimSpace(l), ";")
			if len(tmp) < 2 {
				fmt.Println("Bad err line fmt: " + l)
				return
			}
			add_fle_nme = tmp[1]
			//Grammar
		} else if regexp.MustCompile("^[a-z]+[0-9]*;").MatchString(l) {
			fields = strings.Split(l, ";")
			odiv[fields[0]] = fields[1:]
			//Sch
		} else if regexp.MustCompile(`^\d+\w+\d*;`).MatchString(l) {
			l2 := strings.Split(l, "|")
			wk := 0
			od := ""
			for l2_cnt := 0; l2_cnt < len(l2); l2_cnt++ {
				l = l2[l2_cnt]
				fields = strings.Split(l, ";")
				if l2_cnt == 0 {
					wk, od = intfirst(fields[0])
				} else {
					od = fields[0]
				}
				if wk > maxwk {
					maxwk = wk
					for len(wks) < (maxwk + 1) {
						wks = append(wks, "")
					}
				}
				if wks[wk] == "" {
					wks[wk] = "[" + strconv.Itoa(wk) + "]\r\n"
					if od == "a" {
						if arch_wk1 == 0 {
							arch_wk1 = wk
						} else {
							arch_wk2 = wk
						}
						od = "d"
					}
				}
				odnotd := true
				tmp_fields := ""
				if od == "d" {
					odnotd = false
					for cni := range cn_map {
						for dvi := range cn_map[cni].dv_map {
							for fieldi := 1; fieldi < len(fields); fieldi++ {
								pls := strings.Split(fields[fieldi], "-")
								old_plsi := "0"
								for plsi := 0; plsi < len(pls); plsi++ {
									pls[plsi] = strings.TrimSpace(pls[plsi])
									if len(odiv[pls[plsi]]) > 0 {
										tmp_fields = strconv.Itoa(wk) + pls[plsi] + ";" + old_plsi + "-" + old_plsi
										od = pls[plsi]
										odnotd = true
										break
									}
									itmp := StrToInt(pls[plsi])
									if strings.Contains(rvs_cnf, ";"+cni+";") {
										itmp = maxsds - itmp + 1
									}
									itmp--
									old_plsi = pls[plsi]
									if (elm_num == 0) && (plcD_num == 0) {
										pls[plsi] = cn_map[cni].dv_map[dvi].sds[itmp]
									} else if elm_num > 0 {
										if wk <= elm_num {
											pls[plsi] = cn_map[cni].dv_map[dvi].sds[itmp]
										} else if (elm_stnd == 0) || (wk <= elm_stnd) {
											pls[plsi] = cni + ":" + dvi + ":sds:" + strconv.Itoa(itmp+1)
										} else {
											pls[plsi] = cni + ":" + dvi + ":stnd:" + strconv.Itoa(itmp+1)
										}
									} else { //plcD_num >= 0)
										if wk <= plcD_num {
											pls[plsi] = cn_map[cni].dv_map[dvi].sds[itmp]
										} else if (plcD_stnd == 0) || (wk <= plcD_stnd) {
											pls[plsi] = cni + ":" + dvi + ":sds:" + strconv.Itoa(itmp+1)
										} else {
											pls[plsi] = cni + ":" + dvi + ":stnd:" + strconv.Itoa(itmp+1)
										}
									}
								}
								if (len(odiv[pls[0]]) > 0) || (len(odiv[pls[1]]) > 0) {
									continue
								}
								wks[wk] = wks[wk] + pls[0] + "-" + pls[1] + "\r\n"
								for len(plyrs_sch[pls[0]]) < (maxwk + 1) {
									plyrs_sch[pls[0]] = append(plyrs_sch[pls[0]], "")
								}
								plyrs_sch[pls[0]][wk] = pls[1]
								for len(plyrs_sch[pls[1]]) < (maxwk + 1) {
									plyrs_sch[pls[1]] = append(plyrs_sch[pls[1]], "")
								}
								plyrs_sch[pls[1]][wk] = pls[0]
							}
						}
					}
				}
				if odnotd {
					if _, haskey := odiv[od]; !haskey {
						fmt.Printf("Undefined line %s %s\n", od, l)
						return
					}

					if tmp_fields != "" {
						fields = strings.Split(tmp_fields, ";")
					}
					isvv := fields[1] == "vv"
					isrr := fields[1] == "rr"
					var starti int
					if isvv || isrr {
						starti = 2
					} else {
						starti = 1
					}

					if isrr {
						for len(fields) < (maxsds + 2) {
							fields = append(fields, "")
						}
						for field_i := 1; field_i <= maxsds; field_i++ {
							fields[field_i+1] = fmt.Sprintf("%d-%d", field_i, field_i)
						}
					}
					for odiv_i := 0; odiv_i < len(odiv[od]); odiv_i++ {
						cnanddv := strings.Split(odiv[od][odiv_i], ":")
						cn12 := make([]string, 2)
						dv12 := make([]string, 2)
						cnordv := strings.Split(cnanddv[0], "-")
						a1 := len(cnordv) - 2
						a2 := a1 + 1
						cn12[0] = cnordv[a1]
						dv12[0] = cnordv[a2]
						if len(cnanddv) > 1 {
							cnordv = strings.Split(cnanddv[1], "-")
							a1 = len(cnordv) - 2
							a2 = a1 + 1
							cn12[1] = cnordv[a1]
							dv12[1] = cnordv[a2]
						} else {
							cn12[1] = cnordv[0]
							dv12[1] = cnordv[1]
						}
						for field_i := starti; field_i < len(fields); field_i++ {
							pls := strings.Split(fields[field_i], "-")
							pls2 := make([]string, len(pls))
							for pls_i := 0; pls_i < len(pls); pls_i++ {
								pls[pls_i] = strings.TrimSpace(pls[pls_i])
								itmp := StrToInt(pls[pls_i])
								if strings.Contains(rvs_cnf, ";"+cn12[pls_i]+";") {
									itmp = maxsds - itmp + 1
								}
								itmp--
								pls_eq := (pls[0] == pls[1])
								if (elm_num == 0) || (wk <= elm_num) {
									pls[pls_i] = cn_map[cn12[pls_i]].dv_map[dv12[pls_i]].sds[itmp]
								} else if (elm_stnd == 0) || (wk <= elm_stnd) {
									pls[pls_i] = cn12[pls_i] + ":" + dv12[pls_i] + ":sds:" + strconv.Itoa(itmp+1)
								} else {
									pls[pls_i] = cn12[pls_i] + ":" + dv12[pls_i] + ":stnd:" + strconv.Itoa(itmp+1)
								}
								if isvv && !pls_eq {
									pls_i2 := (pls_i + 1) % 2
									if (elm_num == 0) || (wk <= elm_num) {
										pls2[pls_i] = cn_map[cn12[pls_i2]].dv_map[dv12[pls_i2]].sds[itmp]
									} else if (elm_stnd == 0) || (wk <= elm_stnd) {
										pls2[pls_i] = cn12[pls_i2] + ":" + dv12[pls_i2] + ":sds:" + strconv.Itoa(itmp+1)
									} else {
										pls2[pls_i] = cn12[pls_i2] + ":" + dv12[pls_i2] + ":stnd:" + strconv.Itoa(itmp+1)
									}
								}
							}
							wks[wk] = wks[wk] + pls[0] + "-" + pls[1] + "\r\n"
							for len(plyrs_sch[pls[0]]) < (maxwk + 1) {
								plyrs_sch[pls[0]] = append(plyrs_sch[pls[0]], "")
							}
							plyrs_sch[pls[0]][wk] = pls[1]
							for len(plyrs_sch[pls[1]]) < (maxwk + 1) {
								plyrs_sch[pls[1]] = append(plyrs_sch[pls[1]], "")
							}
							plyrs_sch[pls[1]][wk] = pls[0]
							if isvv && (pls2[0] != "") && (pls2[1] != "") {
								wks[wk] = wks[wk] + pls2[0] + "-" + pls2[1] + "\r\n"
								for len(plyrs_sch[pls2[0]]) < (maxwk + 1) {
									plyrs_sch[pls2[0]] = append(plyrs_sch[pls2[0]], "")
								}
								plyrs_sch[pls2[0]][wk] = pls2[1]
								for len(plyrs_sch[pls2[1]]) < (maxwk + 1) {
									plyrs_sch[pls2[1]] = append(plyrs_sch[pls2[1]], "")
								}
								plyrs_sch[pls2[1]][wk] = pls2[0]
							}
						}
					}
				}
			}
		}
	}

	if add_fle_nme != "" {
		dir_arr := strings.Split(strings.Trim(dir, "\\"), "\\")
		dir_arr[len(dir_arr)-1] = ""
		dir2 := strings.Join(dir_arr, "\\")

		lll_bfile, err := os.ReadFile(dir2 + add_fle_nme)
		lll_is := err == nil
		lll_file := ""
		lll_temp := strings.Split(",,", ",")

		if lll_is {
			lll_file = fmt.Sprintf("%s", lll_bfile)
			lll_temp = strings.Split(lll_file, "\n")
		}
		wk := 0
		cnt1 := 0
		pls := strings.Split(",", ",")
		for ti := 0; ti < len(lll_temp); ti++ {
			l := strings.TrimSpace(lll_temp[ti])
			if strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
				pos := strings.Index(l, " ")
				tmp := strings.Trim(l[pos+1:], "]")
				if wk > 0 {
					all_plrs_arr := strings.Split(all_plyrs_cpy, "|")
					for all_plyrs_i := 1; all_plyrs_i < (len(all_plrs_arr) - 1); all_plyrs_i++ {
						plyrs_sch[all_plrs_arr[all_plyrs_i]] = append(plyrs_sch[all_plrs_arr[all_plyrs_i]], "")
						plyrs_sch[all_plrs_arr[all_plyrs_i]][wk] = "open"
					}
				}
				wk = StrToInt(tmp)
				if wk > maxwk {
					maxwk = wk
					for len(wks) < (maxwk + 1) {
						wks = append(wks, "")
					}
				}
				wks[wk] = "[" + strconv.Itoa(wk) + "]\r\n"
				all_plyrs_cpy = all_plyrs
			} else if strings.HasPrefix(l, "FINAL") || strings.HasSuffix(l, "M CST") || strings.HasSuffix(l, "M MST") || (l == "TBD") || strings.HasSuffix(l, "CDT") || strings.HasSuffix(l, "CET") {
				cnt1 = 0
			} else if strings.Contains(l, " logo ") {
				tmp_arr := strings.Split(l, " logo ")
				if tmp_arr[0] == tmp_arr[1] {
					cnt1++
					if cnt1 == 1 {
						pls[0] = tm_plyr[tmp_arr[0]]
						if len(pls[0]) > maxnme {
							maxnme = len(pls[0])
						}
						all_plyrs_cpy = strings.ReplaceAll(all_plyrs_cpy, "|"+pls[0]+"|", "|")
					} else if cnt1 == 2 {
						pls[1] = tm_plyr[tmp_arr[0]]
						if len(pls[1]) > maxnme {
							maxnme = len(pls[1])
						}
						all_plyrs_cpy = strings.ReplaceAll(all_plyrs_cpy, "|"+pls[1]+"|", "|")
						wks[wk] += pls[0] + "-" + pls[1] + "\r\n"
						for len(plyrs_sch[pls[0]]) < (maxwk + 1) {
							plyrs_sch[pls[0]] = append(plyrs_sch[pls[0]], "")
						}
						plyrs_sch[pls[0]][wk] = pls[1]
						for len(plyrs_sch[pls[1]]) < (maxwk + 1) {
							plyrs_sch[pls[1]] = append(plyrs_sch[pls[1]], "")
						}
						plyrs_sch[pls[1]][wk] = pls[0]
					}
				}
			}
		}
		all_plrs_arr := strings.Split(all_plyrs_cpy, "|")
		for all_plyrs_i := 1; all_plyrs_i < (len(all_plrs_arr) - 1); all_plyrs_i++ {
			plyrs_sch[all_plrs_arr[all_plyrs_i]] = append(plyrs_sch[all_plrs_arr[all_plyrs_i]], "")
			plyrs_sch[all_plrs_arr[all_plyrs_i]][wk] = "open"
		}
	}

	//out array
	var sht_mp_arr []string
	running_sht := ""
	sht_mp_arr = append(sht_mp_arr, "msc")

	if elm_lines != "" {
		elm_lines += "\r\n"
	}
	buffer := lv_wc + "\r\n" + dv_tbrk + "\r\n" + cn_tbrk + "\r\n" + hdr_spr + "\r\n" + attr_lines + elm_lines
	if is_ftb {
		buffer += "ftb::\r\n"
	}
	if arch_wk1 > 0 {
		buffer += "arch;" + strconv.Itoa(arch_wk1) + ";" + strconv.Itoa(arch_wk2) + ";" + "\r\n"
	}
	buffer += "\r\n"
	cp1 := ""

	for wk := 1; wk < len(wks); wk++ {
		buffer += wks[wk] + "\r\n"
		running_sht = gme_mp_prnt("r") + strconv.Itoa(wk)
		lines := strings.Split(hdr_spr, ";")
		rec_c := -1
		for linei := 1; linei < len(lines); linei++ {
			lines[linei] = strings.TrimSpace(lines[linei])
			if (lines[linei] == "Rec") && (wk == 1) {
				rec_c = linei
			}
			running_sht += AddRunningSht(0, linei-1, gme_mp_prnt(lines[linei]))
		}
		if arch_wk2 == wk {
			cur_col := len(lines) - 1
			running_sht += AddRunningSht(0, cur_col, "ARW")
			cur_col++
			running_sht += AddRunningSht(0, cur_col, "PW")
			cur_col++

			dv_tbrk_arr := strings.Split(dv_tbrk, ";")
			for i := 1; i < len(dv_tbrk_arr); i++ {
				if strings.HasPrefix(dv_tbrk_arr[i], "gn") || strings.HasPrefix(dv_tbrk_arr[i], "ot") || strings.HasPrefix(dv_tbrk_arr[i], "pt") {
					running_sht += AddRunningSht(0, cur_col, "P"+dv_tbrk_arr[i][:2])
					cur_col++
				} else if strings.HasPrefix(dv_tbrk_arr[i], "sht") {
					running_sht += AddRunningSht(0, cur_col, "P"+dv_tbrk_arr[i][:3])
					cur_col++
				}
			}
		}
		lines = strings.Split(wks[wk], "\n")
		lines[len(lines)-1] = strings.TrimSpace(lines[len(lines)-1])
		row := 1
		for linei := 0; linei < len(lines); linei++ {
			lines[linei] = strings.TrimSpace(lines[linei])
			if strings.Contains(lines[linei], "-") {
				onetwo := strings.Split(lines[linei], "-")
				running_sht += AddRunningSht(row, 0, onetwo[0])
				running_sht += AddRunningSht(row, 1, onetwo[1])
				if wk == 1 {
					cp1 += onetwo[0] + "\t" + onetwo[1] + "\r\n"
				}
				if rec_c > 1 {
					running_sht += AddRunningSht(row, rec_c-1, "0--0")
				}
				row++
			}
		}
		sht_mp_arr = append(sht_mp_arr, running_sht)
	}
	os.MkdirAll(dir+trn+"\\"+sch_fsname, 544)
	Sprsht_out(dir+trn+"\\"+sch_fsname+trn+".xlsx", sht_mp_arr)

	os.MkdirAll(dir+trn+"\\cp", os.ModePerm)
	err := os.WriteFile(dir+trn+"\\cp\\1-cp.txt", []byte(cp1), 0644)
	if err != nil {
		fmt.Println(err)
	}

	temp := ""
	fmt.Println("all", maxnme, maxcndv)
	maxbegin := 0
	if is_cnf_cm_cm {
		maxbegin = maxnme + maxcndv + 7
		temp += rightPad("cn|dv|cm1|cm2|sds|nme", " ", maxbegin)
	} else {
		maxbegin = maxnme + maxcndv + 3
		temp += rightPad("cn|dv|sds|nme", " ", maxbegin)
	}
	for i := 1; (i <= maxwk) && ((elm_num == 0) || (i <= elm_num)); i++ {
		temp += ";" + leftPad(strconv.Itoa(i), " ", maxnme)
	}
	buffer += temp + "\r\n"
	/*addsp := 0
	  if maxnme<6 {
	      addsp = 6-maxnme
	  }*/

	for cni := range cn_map {
		for dvi := range cn_map[cni].dv_map {
			for i := 0; i < len(cn_map[cni].dv_map[dvi].sds); i++ {
				nme := cn_map[cni].dv_map[dvi].sds[i]
				if nme == "" {
					continue
				}
				temp = cni + "|" + dvi + "|"
				if is_cnf_cm_cm {
					temp += cnf_cm_cm[nme] + "|"
				}
				temp += strconv.Itoa(i+1) + "|" + nme
				temp = rightPad(temp, " ", maxbegin)

				for j := 1; j < len(plyrs_sch[nme]); j++ {
					temp += ";" + leftPad(plyrs_sch[nme][j], " ", maxnme)
				}
				buffer += temp + "\r\n"
			}
			buffer += "\r\n"
		}
	}

	err = os.WriteFile(dir+trn+"\\"+sch_fsname+trn+".txt", []byte(buffer), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func build_attr(afilename string, trn string, dir string, mult bool) {
	var spr_in_file *xlsx.File
	var sheet *xlsx.Sheet
	var cell *xlsx.Cell
	var err error

	spr_in_file, err = xlsx.OpenFile(afilename)
	if err != nil {
		fmt.Println("Can not open input spr " + afilename)
	}

	buffer1 := ""
	count1 := 0
	for sht_i := range spr_in_file.Sheets {
		if trn != "Master" {
			trn_1 := "0"
			trn_2 := "0"
			if trn == "A" {
				trn_1 = "1"
				trn_2 = "4"
			} else if trn == "B" {
				trn_1 = "2"
				trn_2 = "3"
			} else if trn == "C" {
				trn_1 = "5"
				trn_2 = "8"
			} else if trn == "D" {
				trn_1 = "6"
				trn_2 = "7"
			} else if strings.HasPrefix(trn, "A") {
				trn_1 = trn[1:]
				trn_2 = "0"
			}
			if trn_1 == "0" {
				if (spr_in_file.Sheets[sht_i].Name != ("t" + trn + "rnk")) && (strings.ToLower(spr_in_file.Sheets[sht_i].Name) != strings.ToLower(trn+"rnk")) {
					continue
				}
			} else {
				if (spr_in_file.Sheets[sht_i].Name != ("t" + trn_1 + "rnk")) && (spr_in_file.Sheets[sht_i].Name != ("t" + trn_2 + "rnk")) {
					continue
				}
			}
		} else {
			if !strings.HasSuffix(spr_in_file.Sheets[sht_i].Name, "rnk") {
				continue
			}
		}
		sheet = spr_in_file.Sheets[sht_i]
		attr := make(map[string]string)
		var c_arr [5]int

		grp := ""
		prev_blank := false
		for row := 1; ; row++ {
			cell, _ = sheet.Cell(row, 0)
			p := cell.Value
			if p == "" {
				if prev_blank {
					break
				} else {
					prev_blank = true
					grp += ";"
				}
			} else {
				prev_blank = false
				grp += "," + p
			}
		}
		for col := 1; ((c_arr[0] == 0) || (c_arr[1] == 0) || (c_arr[2] == 0) || (c_arr[3] == 0) || (c_arr[4] == 0)) && (col < 20); col++ {
			cell, _ = sheet.Cell(0, col)
			p := cell.Value
			if p == "o" {
				c_arr[0] = col
			} else if p == "st" {
				c_arr[1] = col
			} else if p == "sp" {
				c_arr[2] = col
			} else if p == "d" {
				c_arr[3] = col
			} else if p == "spcl" {
				c_arr[4] = col
				buffer1 += "has_spcl\r\n"
			}
		}
		for c_i := 0; c_i < len(c_arr); c_i++ {
			c_arr_p1 := c_arr[c_i] + 1
			a_cnt := 0
			b_cnt := 0
			c_cnt := 0
			row_end := 0
			for row := 1; ; row++ {
				cell, _ = sheet.Cell(row, c_arr_p1)
				p := cell.Value
				if p == "" {
					row_end = row
					break
				}
				if p == "a" {
					a_cnt++
				} else if p == "b" {
					b_cnt++
				} else if p == "c" {
					c_cnt++
				}
			}
			if a_cnt == 0 {
				break
			}
			a_ind := 1
			b_ind := 1
			c_ind := 1
			base := 10
			for row := 1; row < row_end; row++ {
				cell, _ = sheet.Cell(row, c_arr[c_i])
				nme := cell.Value
				cell, _ = sheet.Cell(row, c_arr_p1)
				p := cell.Value
				if p == "d" {
					attr[nme] += "," + strconv.Itoa(base)
				} else if p == "s" {
					attr[nme] += "," + strconv.Itoa(base+7)
				} else if p == "c" {
					if (c_ind * 2) <= c_cnt {
						attr[nme] += "," + strconv.Itoa(base+2)
					} else {
						attr[nme] += "," + strconv.Itoa(base+1)
					}
					c_ind++
				} else if p == "b" {
					if (b_ind * 2) <= b_cnt {
						attr[nme] += "," + strconv.Itoa(base+4)
					} else {
						attr[nme] += "," + strconv.Itoa(base+3)
					}
					b_ind++
				} else if p == "a" {
					if (a_ind * 2) <= a_cnt {
						attr[nme] += "," + strconv.Itoa(base+6)
					} else {
						attr[nme] += "," + strconv.Itoa(base+5)
					}
					a_ind++
				}
			}
		}

		arr1 := strings.Split(grp, ";")
		for x := 0; x < len(arr1)-1; x++ {
			buffer1 += "[D" + strconv.Itoa(count1+1) + "]\r\n"
			count1++
			buffer1 += "1=" + arr1[x][1:] + "\r\n"
			arr2 := strings.Split(arr1[x], ",")
			for y := 1; y < len(arr2); y++ {
				buffer1 += arr2[y] + "=" + attr[arr2[y]][1:] + "\r\n"
			}
		}
	}
	if buffer1 != "" {
		if mult { //if multiple files
			os.MkdirAll(dir+"\\"+trn+"\\attr", os.ModePerm)
			err = os.WriteFile(dir+"\\"+trn+"\\attr\\attr1.txt", []byte(buffer1), 0644)
		} else if trn != "Master" {
			os.MkdirAll(dir+"\\attr", os.ModePerm)
			err = os.WriteFile(dir+"\\attr\\attr1.txt", []byte(buffer1), 0644)
		} else {
			os.MkdirAll(dir+"\\TMaster\\attr", os.ModePerm)
			err = os.WriteFile(dir+"\\TMaster\\attr\\attr1.txt", []byte(buffer1), 0644)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func intfirst(value string) (int, string) {
	for i := len(value) - 1; i >= 0; i-- {
		first := value[:i]
		second := value[i:]
		if ifirst, err := strconv.Atoi(first); err == nil {
			return ifirst, second
		}
	}
	return 0, value
}

func rightPad(str, pad string, lngth int) string {
	for len(str) < lngth {
		str += pad
	}
	return str
}

func leftPad(str, pad string, lngth int) string {
	for len(str) < lngth {
		str = pad + str
	}
	return str
}

func AddRunningSht(rw, cl int, str string) string {
	return ";" + strconv.Itoa(rw) + ";" + strconv.Itoa(cl) + ";" + str
}

func Sprsht_out(fname string, shts []string) {
	var oFile *xlsx.File
	var oSheet *xlsx.Sheet
	var oCell *xlsx.Cell
	var err error
	var rw int
	var cl int

	oFile = xlsx.NewFile()
	for sht_i := len(shts) - 1; sht_i >= 0; sht_i-- {
		sht_arr := strings.Split(shts[sht_i], ";")
		oSheet, err = oFile.AddSheet(sht_arr[0])
		if err != nil {
			fmt.Printf(err.Error())
		}
		for indx := 1; (indx + 2) < len(sht_arr); indx += 3 {
			rw, err = strconv.Atoi(sht_arr[indx])
			cl, err = strconv.Atoi(sht_arr[indx+1])
			oCell, _ = oSheet.Cell(rw, cl)
			oCell.Value = sht_arr[indx+2]
		}
	}
	err = oFile.Save(fname)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func Sprsht_out2(fname string, shts []string) {
	var oFile *xlsx.File
	var oSheet *xlsx.Sheet
	var oCell *xlsx.Cell
	var err error
	var rw int
	var cl int

	oFile = xlsx.NewFile()
	for sht_i := 0; sht_i < len(shts); sht_i++ {
		sht_arr := strings.Split(shts[sht_i], ";")
		oSheet, err = oFile.AddSheet(sht_arr[0])
		if err != nil {
			fmt.Printf(err.Error())
		}
		for indx := 1; (indx + 2) < len(sht_arr); indx += 3 {
			rw, err = strconv.Atoi(sht_arr[indx])
			cl, err = strconv.Atoi(sht_arr[indx+1])
			oCell, _ = oSheet.Cell(rw, cl)
			oCell.Value = sht_arr[indx+2]
		}
	}

	err = oFile.Save(fname)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func nfltrn(mastertrn, file string) {
	larr := strings.Split(file, "\n")
	trncondiv := make(map[string]string)
	tmesch := make(map[string]string)
	wksch := make(map[int]string)
	maxlntme := 0
	maxlntcd := 0
	mxwk := 0
	nmetrn := make(map[string]string)
	trnStr := ""

	for a := 0; a < len(larr); a++ {
		larr[a] = strings.TrimSpace(larr[a])
		if larr[a] == "" {
			continue
		}
		line := strings.Split(larr[a], "\t")
		if len(line) == 1 {
			line = strings.SplitN(line[0], "|", 2)
			trn := strings.SplitN(line[0], ";", 2)
			if len(line) > 1 {
				trncondiv[line[0]] = line[1]
				if len(line[0]) > maxlntcd {
					maxlntcd = len(line[0])
				}
				tms := strings.Split(line[1], "|")
				for b := 0; b < len(tms); b++ {
					nmetrn[tms[b]] = trn[0]
					if !strings.Contains(trnStr+"|", "|"+trn[0]+"|") {
						trnStr += "|" + trn[0]
					}
				}
			}
		} else {
			if len(line[0]) > maxlntme {
				maxlntme = len(line[0])
			}
			tmparr := strings.Split(line[1], "|")
			tmparr2 := strings.Split("", ";")
			for b := 1; b < len(tmparr); b++ {
				tmparr2 = append(tmparr2, "")
			}
			for b := 0; b < len(tmparr); b++ {
				if tmparr[b] == "" {
					continue
				}
				wknum := 0
				tme := ""
				wknum, tme = intfirst(tmparr[b])
				if mxwk < wknum {
					mxwk = wknum
				}
				if tme != "open" {
					if !(strings.Contains(wksch[wknum]+"|", "|"+tme+"-"+line[0]+"|") || strings.Contains(wksch[wknum]+"|", "|"+line[0]+"-"+tme+"|")) {
						if strings.Contains(wksch[wknum]+"|", "-"+line[0]+"|") || strings.Contains(wksch[wknum], "|"+line[0]+"-") {
							fmt.Println("Already A sch", wknum, line[0], line[0]+"-"+tme)
							return
						}
						if strings.Contains(wksch[wknum]+"|", "-"+tme+"|") || strings.Contains(wksch[wknum], "|"+tme+"-") {
							fmt.Println("Already B sch", wknum, tme, line[0]+"-"+tme)
							return
						}
						wksch[wknum] += "|" + line[0] + "-" + tme
					}
				}
				tmparr2[wknum] = tme
			}
			tmesch[line[0]] = strings.Join(tmparr2, "|")
		}
	}
	header := "hdr-spr-ln;T1;T2;W;Pt;Pf;Pa;Rec"
	buffer := "dv-tbr;o;arch;hhx;dv;sht-pt;pts-pt;arch-d;clsgm\r\n"
	buffer += "cn-tbr;o;hha;sht;pts;clsgm\r\n"
	buffer += header + "\r\n"
	buffer += "attr-w;3;6;9;10;12;15;arch\r\n"
	buffer += "attr-blw;;8;3;10\r\n"
	buffer += "split\r\n"

	//out array
	sht_mp_arr := make(map[string][]string)
	running_sht := make(map[string]string)
	trns := strings.Split(trnStr, "|")
	row := make(map[string]int)

	for b := 1; b < len(trns); b++ {
		sht_mp_arr[trns[b]] = append(sht_mp_arr[trns[b]], "msc")
	}
	sht_mp_arr["inter"] = append(sht_mp_arr["inter"], "msc")
	buffer1 := make(map[string]string)

	for a := 1; a <= len(wksch); a++ {
		if wksch[a] == "" {
			continue
		}
		hasinter := false
		buffer += "\r\n[" + strconv.Itoa(a) + "]\r\n"
		headerArr := strings.Split(header, ";")
		for b := 1; b < len(trns); b++ {
			running_sht[trns[b]] = gme_mp_prnt("r") + strconv.Itoa(a)
			for h_a := 1; h_a < len(headerArr); h_a++ {
				running_sht[trns[b]] += AddRunningSht(0, h_a-1, headerArr[h_a])
			}
		}
		running_sht["inter"] = gme_mp_prnt("r") + strconv.Itoa(a)
		for h_a := 1; h_a < len(headerArr); h_a++ {
			running_sht["inter"] += AddRunningSht(0, h_a-1, headerArr[h_a])
		}
		tmparr := strings.Split(wksch[a], "|")
		for b := 1; b < len(tmparr); b++ {
			onetwo := strings.Split(tmparr[b], "-")
			temptrn := "inter"
			if nmetrn[onetwo[0]] == nmetrn[onetwo[1]] {
				temptrn = nmetrn[onetwo[0]]
			} else {
				hasinter = true
			}
			row[temptrn]++
			running_sht[temptrn] += AddRunningSht(row[temptrn], 0, onetwo[0])
			running_sht[temptrn] += AddRunningSht(row[temptrn], 1, onetwo[1])
			if a == 1 {
				buffer1[temptrn] += onetwo[0] + "\t" + onetwo[1] + "\r\n"
				running_sht[temptrn] += AddRunningSht(row[temptrn], len(headerArr)-2, "0--0")
			}
			buffer += tmparr[b] + "\r\n"
		}
		for b := 1; b < len(trns); b++ {
			sht_mp_arr[trns[b]] = append(sht_mp_arr[trns[b]], running_sht[trns[b]])
			row[trns[b]] = 0
		}
		if hasinter {
			sht_mp_arr["inter"] = append(sht_mp_arr["inter"], running_sht["inter"])
			row["inter"] = 0
		}
	}
	maxlntcd += 3 + maxlntme

	tmpbuffer := rightPad("trn|cn|dv|sds|nme", " ", maxlntcd) + ";"
	for a := 1; a <= mxwk; a++ {
		tmpbuffer += leftPad(strconv.Itoa(a), " ", maxlntme) + ";"
	}

	buffer += "\r\n" + tmpbuffer
	tmpbuffer = ""
	sorted_tcd := SortNamesInMap2(trncondiv)
	for a := 0; a < len(sorted_tcd); a++ {
		tempArr := strings.Split(trncondiv[sorted_tcd[a]], "|")
		for b := 0; b < len(tempArr); b++ {
			tmpbuffer += rightPad(strings.ReplaceAll(sorted_tcd[a], ";", "|")+"|"+strconv.Itoa(b+1)+"|"+tempArr[b], " ", maxlntcd) + ";"
			tempArr2 := strings.Split(tmesch[tempArr[b]], "|")
			for c := 1; c < len(tempArr2); c++ {
				tmpbuffer += leftPad(tempArr2[c], " ", maxlntme) + ";"
			}
			tmpbuffer += "\r\n"
		}
	}
	buffer += "\r\n" + tmpbuffer
	os.MkdirAll(mastertrn, os.ModePerm)
	for a := range sht_mp_arr {
		Sprsht_out(mastertrn+"\\sch"+a+".xlsx", sht_mp_arr[a])
	}
	err := os.WriteFile(mastertrn+"\\schM.txt", []byte(buffer), 0644)
	if err != nil {
		fmt.Println(err)
	}
	for a := range sht_mp_arr {
		a2 := a
		if a == "inter" {
			a2 = "trnI"
		}
		os.MkdirAll(mastertrn+"\\"+a2+"\\cp", os.ModePerm)
		err = os.WriteFile(mastertrn+"\\"+a2+"\\cp\\1-cp.txt", []byte(buffer1[a]), 0644)
	}
}

func SortNamesInMap2(mymap map[string]string) []string {
	var arr_s []string
	for one_s := range mymap {
		arr_s = append(arr_s, one_s)
	}
	sort.Strings(arr_s)
	return arr_s
}

func RunGME(fbase string, trn, wk string) {
	gmePath := strings.ReplaceAll(strings.ToLower(fbase), "trn sch\\t"+strings.ToLower(trn), "gms1")
	os.MkdirAll(fbase+"gme", os.ModePerm)
	FileCopy(gmePath+"trn.exe", fbase+"gme\\trn.exe")
	FileCopy(gmePath+"ht-mss-rlls.txt", fbase+"gme\\ht-mss-rlls.txt")
	FileCopy(gmePath+"sb.txt", fbase+"gme\\sb.txt")
	FileCopy(fbase+attr_fname+"\\"+attr_fname+wk+".txt", fbase+"gme\\"+attr_fname+wk+".txt")
	FileCopy(fbase+"cp\\"+wk+"-cp.txt", fbase+"gme\\sch"+wk+".txt")

	a := exec.Command(fbase+"gme\\trn.exe", wk+"t", fbase+"gme")
	if err := a.Run(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("fb", gmePath, fbase, trn, wk)
	//FileDelete(fbase + "gme\\" + attr_fname + wk + ".txt")
	//FileDelete(fbase + "gme\\" + wk + "-cp.txt")
}

func FileCopy(f1 string, f2 string) {
	lines := LinesInFile(f1)
	b := strings.Join(lines, "\r\n")
	err := os.WriteFile(f2, []byte(b), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func FileDelete(fn string) {
	err := os.Remove(fn)
	if err != nil {
		fmt.Println(err)
	}
}
