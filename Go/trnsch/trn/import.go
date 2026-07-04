package trn

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"trnsch/db"
	"trnsch/model"
)

func Import(trnId int, filename string) ([]byte, error) {
	DivHsh := make(map[string]int)
	AbbrHsh := make(map[string]string)
	SdsHsh := make(map[string]int)

	sql := `select c.description conf, d.description divi, d.id did,
			p.description tm, p.sds
			from conf c inner join divi d on d.conf_id = c.id 
			inner join plyrtm p on p.divi_id = d.id 
			where c.tourn_id = ?`
	rows, err := db.DB.Query(sql, trnId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var did, sds int
	var conf, divi, tm string
	for rows.Next() {
		if err = rows.Scan(&conf, &divi, &did, &tm, &sds); err != nil {
			return nil, err
		}
		tmArr := strings.Split(tm, " ")
		tm = strings.ToLower(tmArr[len(tmArr)-1])
		conf = conf[0:1]
		divi = divi[0:1]

		DivHsh[tm] = did
		AbbrHsh[tm] = conf + divi
		SdsHsh[tm] = sds
	}

	prevGid := 0
	sql = `select letters, dv1_id, dv2_id, g.id gid,
		coalesce(wk.rnd,0) rnd, coalesce(wk.id,0) wkid, coalesce(ws.sd1,0), coalesce(ws.sd2,0) 
		from grammar g 
		left join wk_sch wk on wk.grammar_id = g.id
		left join wk_sch_sds ws on ws.wk_sch_id = wk.id
		where g.tourn_id = ?`
	rows, err = db.DB.Query(sql, trnId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	GrmDv1Hsh := make(map[string]int)
	GrmDv2Hsh := make(map[string]int)
	GrmGidHsh := make(map[string]int)
	GidWk := make(map[string]int)
	WkIdSd1 := make(map[int]int)
	WkIdSd2 := make(map[int]int)

	var dv1, dv2, gid, rnd, wkid, sd1, sd2 int
	var letters string
	for rows.Next() {
		if err = rows.Scan(&letters, &dv1, &dv2, &gid, &rnd, &wkid, &sd1, &sd2); err != nil {
			return nil, err
		}
		if prevGid != gid {
			GrmDv1Hsh[letters] = dv1
			GrmDv2Hsh[letters] = dv2
			GrmGidHsh[letters] = gid
			prevGid = gid
		}
		if wkid > 0 {
			GidWk[fmt.Sprintf("%d,%d", gid, rnd)] = wkid
			if (sd1 > 0) && (sd2 > 0) {
				WkIdSd1[wkid] = sd1
				WkIdSd2[wkid] = sd2
			}
		}

	}
	/*fmt.Println(DivHsh)
	fmt.Println()
	fmt.Println(AbbrHsh)
	fmt.Println()
	fmt.Println(SdsHsh)
	*/
	wk := 0
	tm1 := ""
	tm2 := ""
	lines := LinesInFile(filename)
	for _, line := range lines {
		line = strings.ToLower(strings.TrimSpace(line))
		if (strings.ToLower(line) == "thursday night football") && (wk < 17) {
			wk++
			tm1 = ""
			tm2 = ""
		} else if strings.HasPrefix(line, "week") {
			arr := strings.Split(line, " ")
			wk, _ = strconv.Atoi(arr[1])
			tm1 = ""
			tm2 = ""
		} else if DivHsh[line] > 0 {
			if tm1 == "" {
				tm1 = line
			} else {
				tm2 = line

				letters = AbbrHsh[tm1] + "-" + AbbrHsh[tm2]
				if GrmGidHsh[letters] == 0 {
					gid := 0
					sql = `insert into grammar (tourn_id,letters,dv1_id,dv2_id)
						values (?,?,?,?)`
					params := []interface{}{trnId, letters, DivHsh[tm1], DivHsh[tm2]}
					gid, err = db.ChangDB(sql, params)
					if err != nil {
						return nil, err
					}
					GrmDv1Hsh[letters] = DivHsh[tm1]
					GrmDv2Hsh[letters] = DivHsh[tm2]
					GrmGidHsh[letters] = gid
				}

				gidwk := fmt.Sprintf("%d,%d", GrmGidHsh[letters], wk)
				if GidWk[gidwk] == 0 {
					sql = `insert into wk_sch (rnd,grammar_id,rpt,has_stnd)
						values (?,?,'','Y')`
					params := []interface{}{wk, GrmGidHsh[letters]}
					wkid, err = db.ChangDB(sql, params)
					if err != nil {
						return nil, err
					}
					GidWk[gidwk] = wkid
				}

				if (WkIdSd1[GidWk[gidwk]] == 0) || (WkIdSd2[GidWk[gidwk]] == 0) {
					//Get sd1 and sd2
					sql = `insert into wk_sch_sds (wk_sch_id,sd1,sd2)
						values (?,?,?)`
					params := []interface{}{GidWk[gidwk], SdsHsh[tm1], SdsHsh[tm2]}
					_, err = db.ChangDB(sql, params)
					if err != nil {
						return nil, err
					}
				}

				fmt.Println(wk, tm1, tm2, gid, gidwk)
				tm1 = ""
				tm2 = ""
			}
		}
	}

	//Find grammar and create if not existing
	//Find wk_sch rnd and grammar and create if not existing
	//Create seeds if wk_sch_sds
	resp := model.Response{Msg: "ok"}
	jsonOut, _ := json.MarshalIndent(resp, "", "  ")
	return jsonOut, nil
}

func LinesInFile(fname string) []string {
	gen_bfile, err := os.ReadFile(fname)
	gen_is := err == nil
	if !gen_is {
		fmt.Println("file not found: " + fname)
		return nil
	}
	gen_file := fmt.Sprintf("%s", gen_bfile)
	lines := strings.Split(gen_file, "\r\n")
	return lines
}
