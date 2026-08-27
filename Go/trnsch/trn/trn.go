package trn

import (
	"fmt"
	"strings"
	"trnsch/db"
)

type dv_strct struct {
	pid []int
	sds []int
}

type cn_strct struct {
	dv_map map[int]dv_strct
}

type SchStruct struct {
	cn_map map[int]cn_strct
}

type ExtraStruct struct {
	attr_lines string
}

func Start(trnId int, build int) error {
	rows, err := db.DB.Query("select count(*) cnt from gms where tourn_id = ? and rnd >= ?", trnId, build)
	if err != nil {
		return err
	}
	defer rows.Close()

	cnt := 0
	row := db.DB.QueryRow("SELECT DATABASE()")
	var dbName string
	_ = row.Scan(&dbName)
	fmt.Println("Connected to DB:", dbName)
	for rows.Next() {
		if err = rows.Scan(&cnt); err != nil {
			return err
		}
		if cnt > 0 {
			fmt.Println("Already Existing ", cnt)
			return nil
		}
	}

	err = BuildSch(trnId, build)
	return err
}

func BuildSch(trnId int, build int) error {

	schStruct := SchStruct{
		cn_map: make(map[int]cn_strct),
	}

	extraStruct := ExtraStruct{}
	var odiv map[int][]int
	var err error
	if err = ConfDivi(trnId, &schStruct); err != nil {
		return err
	}
	if odiv, err = Grammar(trnId, &schStruct); err != nil {
		return err
	}

	if err = SchFunc(trnId, &schStruct, odiv, build); err != nil {
		fmt.Println("err", err.Error())
		return err
	}

	if err = Extras(trnId, &extraStruct); err != nil {
		return err
	}

	return nil
}

func ConfDivi(trnId int, sch *SchStruct) error {
	sql := "select c.id cid, d.id did, p.id pid, p.sds sds " +
		"from conf c inner join divi d on d.conf_id = c.id " +
		"inner join plyrtm p on p.divi_id = d.id " +
		"where tourn_id = ? " +
		"order by sds"
	rows, err := db.DB.Query(sql, trnId)
	if err != nil {
		return err
	}
	defer rows.Close()

	var cid, did, pid, sds int
	for rows.Next() {
		if err = rows.Scan(&cid, &did, &pid, &sds); err != nil {
			return err
		}
		if _, haskey := sch.cn_map[cid]; !haskey {
			var cn_tmp cn_strct
			cn_tmp.dv_map = make(map[int]dv_strct)
			sch.cn_map[cid] = cn_tmp
		}
		if _, haskey := sch.cn_map[cid].dv_map[did]; !haskey {
			var dv_tmp dv_strct
			sch.cn_map[cid].dv_map[did] = dv_tmp
		}

		dv := sch.cn_map[cid].dv_map[did]
		dv.pid = append(dv.pid, pid)
		dv.sds = append(dv.sds, sds)
		sch.cn_map[cid].dv_map[did] = dv
	}
	return nil
}

func Grammar(trnId int, sch *SchStruct) (map[int][]int, error) {
	rows, err := db.DB.Query("select id, dv1_id, dv2_id from grammar where tourn_id = ?", trnId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	odiv := make(map[int][]int)
	var gid, divi1, divi2 int
	for rows.Next() {
		if err = rows.Scan(&gid, &divi1, &divi2); err != nil {
			return nil, err
		}

		odiv[gid] = append(odiv[gid], divi1, divi2)
	}
	return odiv, nil
}

func SchFunc(trnId int, sch *SchStruct, odiv map[int][]int, build int) error {
	sql := "select ws.rnd, g.id, ws.rpt, ws.id wsid, " +
		"d1.conf_id cid1, d2.conf_id cid2, g.dv1_id dv1, g.dv2_id dv2, wss.sd1, wss.sd2, " +
		"coalesce(ws.has_stnd,'Y') has_stnd " +
		"from wk_sch ws inner join grammar g on ws.grammar_id = g.id " +
		"left join wk_sch_sds wss on wss.wk_sch_id = ws.id " +
		"left join divi d1 on g.dv1_id = d1.id " +
		"left join divi d2 on g.dv2_id = d2.id " +
		"where g.tourn_id = ? and ws.rnd >= ? " +
		"order by ws.rnd, g.id"
	rows, err := db.DB.Query(sql, trnId, build)
	if err != nil {
		return err
	}
	defer rows.Close()

	var rnd, gid, cn1, cn2, dv1, dv2, wsid, sd1, sd2 int
	var rpt, HasStnd string
	gmsId := 0

	for rows.Next() {
		if err = rows.Scan(&rnd, &gid, &rpt, &wsid, &cn1, &cn2, &dv1, &dv2, &sd1, &sd2, &HasStnd); err != nil {
			return err
		}
		if (rpt != "rvs") && (rpt != "elm") && (rnd > 0) {
			if _, haskey := odiv[gid]; !haskey {
				fmt.Errorf("Undefined line %s %s\n", gid, odiv[gid])
				return err
			}

			var sdArr1, sdArr2 []int
			if rpt == "rr" {
				sdsCnt := len(sch.cn_map[cn1].dv_map[dv1].sds)
				if len(sch.cn_map[cn2].dv_map[dv2].sds) < sdsCnt {
					sdsCnt = len(sch.cn_map[cn2].dv_map[dv2].sds)
				}
				for indx2 := 0; indx2 < sdsCnt; indx2++ {
					sdArr1 = append(sdArr1, sch.cn_map[cn1].dv_map[dv1].sds[indx2])
					sdArr2 = append(sdArr2, sch.cn_map[cn2].dv_map[dv2].sds[indx2])
				}
			} else {
				sdArr1 = append(sdArr1, sd1)
				sdArr2 = append(sdArr2, sd2)
				if rpt == "vv" {
					sdArr1 = append(sdArr1, sd2)
					sdArr2 = append(sdArr2, sd1)
				}
			}

			for indx2 := 0; indx2 < len(sdArr1); indx2++ {
				var plyrgm1, plyrgm2 int
				for indx3 := 0; indx3 < len(sch.cn_map[cn1].dv_map[dv1].sds); indx3++ {
					if sch.cn_map[cn1].dv_map[dv1].sds[indx3] == sdArr1[indx2] {
						plyrgm1 = sch.cn_map[cn1].dv_map[dv1].pid[indx3]
						break
					}
				}
				if plyrgm1 == 0 {
					return fmt.Errorf("No sds found rnd %d, gid %d, wsid %s, sds %d", rnd, gid, wsid, sd1)
				}
				for indx3 := 0; indx3 < len(sch.cn_map[cn2].dv_map[dv2].sds); indx3++ {
					if sch.cn_map[cn2].dv_map[dv2].sds[indx3] == sdArr2[indx2] {
						plyrgm2 = sch.cn_map[cn2].dv_map[dv2].pid[indx3]
						break
					}
				}
				if plyrgm2 == 0 {
					return fmt.Errorf("No sds found rnd %d, gid %d, wsid %s, sds %d", rnd, gid, wsid, sd2)
				}

				//insert gms record, and plyrgm1 and plyrgm2
				sql = "insert into gms (tourn_id, rnd, has_stnd) values (?,?,?)"
				params := []interface{}{trnId, rnd, HasStnd}
				gmsId, err = db.ChangDB(sql, params)
				if err != nil {
					return err
				}
				plyrgm := plyrgm1
				for indxx := 0; indxx < 2; indxx++ {
					if rnd == 1 {
						sql = "insert into gms_plyrtm (gms_id, plyrtm_id, rec) " +
							"values (?,?,?)"
						params := []interface{}{gmsId, plyrgm, 0}
						_, err = db.ChangDB(sql, params)
					} else {
						sql = "insert into gms_plyrtm (gms_id, plyrtm_id) " +
							"values (?,?)"
						params := []interface{}{gmsId, plyrgm}
						_, err = db.ChangDB(sql, params)
					}
					if err != nil {
						return err
					}
					plyrgm = plyrgm2
				}
			}
		}
		/*} else if strings.HasPrefix(l, "rvs;") {
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
		}*/
	}
	return nil
}

func Extras(trnId int, ex *ExtraStruct) error {
	rows, err := db.DB.Query("select description from lll_line where tourn_id = ? order by ordr", trnId)
	if err != nil {
		return err
	}
	defer rows.Close()
	description := ""
	for rows.Next() {
		if err = rows.Scan(&description); err != nil {
			return err
		}
		l := strings.TrimSpace(description)

		if strings.HasPrefix(l, "attr-w;") || strings.HasPrefix(l, "attr-blw;") || strings.HasPrefix(l, "attr_orig+") {
			ex.attr_lines += l + "\r\n"
			/*} else if strings.HasPrefix(l, "ftb::") {
				ex.is_ftb = true
				ex.is_gme_mp = true
				ex.gme_mp = make(map[string]string)
				ex.gme_mp_rev = make(map[string]string)
				gme_mp_strt(ex)
			} else if strings.HasPrefix(l, "tme;") {
				tmp := strings.Split(strings.TrimSpace(l), ";")
				if len(tmp) < 2 {
					return fmt.Errorf("bad err line fmt: %s", l)
				}
				tmp = strings.Split(tmp[1], "-")
				if len(tmp) < 2 {
					return fmt.Errorf("bad err line fmt: %s", l)
				}
				ex.tm_plyr[tmp[1]] = tmp[0]
			} else if strings.HasPrefix(l, "fle;") {
				tmp := strings.Split(strings.TrimSpace(l), ";")
				if len(tmp) < 2 {
					return fmt.Errorf("bad err line fmt: %s", l)
				}
				ex.add_fle_nme = tmp[1]*/
		}
	}

	return nil
}
