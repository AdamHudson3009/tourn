package trn

import (
	"encoding/json"
	"fmt"
	"trnsch/db"
	"trnsch/model"
)

func Get(trnId int) ([]byte, error) {
	sql := "select p.id pid, p.description plyrtm, p.sds, " +
		"d.id did, d.description divi, c.id cid, c.description conf, " +
		"g.rnd, gp2.plyrtm_id " +
		"from plyrtm p inner join divi d on p.divi_id = d.id " +
		"inner join conf c on d.conf_id = c.id " +
		"inner join gms_plyrtm gp1 on gp1.plyrtm_id = p.id " +
		"inner join gms g on gp1.gms_id = g.id " +
		"inner join gms_plyrtm gp2 on gp2.gms_id = g.id and gp1.plyrtm_id <> gp2.plyrtm_id " +
		"where c.tourn_id = ? " +
		"order by c.description, d.description, p.sds, g.rnd, gp2.plyrtm_id"
	rows, err := db.DB.Query(sql, trnId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trnTree := model.TrnTree{
		Confs:   make(map[int]model.Conf),
		Plyrtms: make(map[int]model.PlyrtmInfo),
	}

	var pid, did, sds, rnd, plyrtm2, cid, rndExpec int
	var plyrtm, divi, conf string
	for rows.Next() {
		if err = rows.Scan(&pid, &plyrtm, &sds, &did, &divi, &cid, &conf, &rnd, &plyrtm2); err != nil {
			return nil, err
		}
		if _, haskey := trnTree.Confs[cid]; !haskey {
			var cn_tmp model.Conf
			cn_tmp.ConfId = cid
			cn_tmp.Conf = conf
			cn_tmp.Divis = make(map[int]model.Divi)
			trnTree.Confs[cid] = cn_tmp
		}
		if _, haskey := trnTree.Confs[cid].Divis[did]; !haskey {
			var dv_tmp model.Divi
			dv_tmp.DiviId = did
			dv_tmp.Divi = divi
			dv_tmp.PlyrtmsSch = make(map[int]model.PlyrtmSch)
			trnTree.Confs[cid].Divis[did] = dv_tmp
		}
		if _, haskey := trnTree.Confs[cid].Divis[did].PlyrtmsSch[pid]; !haskey {
			var plt_tmp model.PlyrtmSch
			plt_tmp.PlyrtmId = pid
			trnTree.Confs[cid].Divis[did].PlyrtmsSch[pid] = plt_tmp
			rndExpec = 1
		}
		if _, haskey := trnTree.Plyrtms[pid]; !haskey {
			var plt_info model.PlyrtmInfo
			plt_info.PlyrtmId = pid
			plt_info.Plyrtm = plyrtm
			plt_info.Cid = cid
			plt_info.Conf = conf
			plt_info.Did = did
			plt_info.Divi = divi
			plt_info.Sds = sds
			trnTree.Plyrtms[pid] = plt_info
		}

		arr := trnTree.Confs[cid].Divis[did].PlyrtmsSch[pid].Plyrtms
		for ; rndExpec < rnd; rndExpec++ {
			arr = append(arr, 0)
		}
		arr = append(arr, plyrtm2)
		rndExpec++
		plyrtmSch := trnTree.Confs[cid].Divis[did].PlyrtmsSch[pid]
		plyrtmSch.Plyrtms = arr
		trnTree.Confs[cid].Divis[did].PlyrtmsSch[pid] = plyrtmSch
	}

	jsonOut, _ := json.MarshalIndent(trnTree, "", "  ")
	return jsonOut, nil
}

func GetRnd(trnId int, rnd int) ([]byte, error) {
	mscArr := make(map[int][]string)
	sql := "select plyrtm_id, letters from msc where tourn_id = ? and rnd = ? " +
		"order by plyrtm_id, letters"
	rows, err := db.DB.Query(sql, trnId, rnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plyr1 int
	var letters string

	for rows.Next() {
		if err = rows.Scan(&plyr1, &letters); err != nil {
			return nil, err
		}
		mscArr[plyr1] = append(mscArr[plyr1], letters)
	}

	archGmeArr := make(map[int]string)
	sql = "select a.pl1, a.pl2, a.plyrtm_wn, a.pf, a.pa " +
		"from (select gA.rnd, gA.plyrtm_wn, gA.pf, gA.pa, max(gp1A.plyrtm_id) pl1, min(gp2A.plyrtm_id) pl2 " +
		"	from gms gA " +
		"	inner join gms_plyrtm gp1A on gp1A.gms_id = gA.id " +
		"	inner join gms_plyrtm gp2A on gp2A.gms_id = gA.id and gp2A.plyrtm_id <> gp1A.plyrtm_id " +
		"	where gA.tourn_id = ? and gA.plyrtm_wn > -1 " +
		"	group by gA.rnd, gA.plyrtm_wn, pf, pa) a " +
		"inner join (select gA.rnd, gA.ordr, gA.id, max(gp1A.plyrtm_id) pl1, min(gp2A.plyrtm_id) pl2 " +
		"	from gms gA " +
		"	inner join gms_plyrtm gp1A on gp1A.gms_id = gA.id " +
		"	inner join gms_plyrtm gp2A on gp2A.gms_id = gA.id and gp2A.plyrtm_id <> gp1A.plyrtm_id " +
		"	where gA.tourn_id = ? " +
		"	group by gA.rnd, gA.ordr, gA.id) b " +
		"on a.rnd < b.rnd and a.pl1 = b.pl1 and a.pl2 = b.pl2 and b.rnd = ?"
	rows, err = db.DB.Query(sql, trnId, trnId, rnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plyr2, plyrwn, pf, pa int

	for rows.Next() {
		if err = rows.Scan(&plyr1, &plyr2, &plyrwn, &pf, &pa); err != nil {
			return nil, err
		}
		temp := 0
		if (pf > 0) && (pa == 0) {
			temp = 1
		}
		archGmeArr[plyrwn] = fmt.Sprintf("%d--%d--%d--%d", temp, pf-pa, pf, pa)
	}

	sql = "select a.*, p1.description plyr1nm , p2.description plyr2nm, " +
		"coalesce(gp1.rec,0) rec1, coalesce(gp1.tie,0) tie1, " +
		"coalesce(gp2.rec,0) rec2, coalesce(gp2.tie,0) tie2, " +
		"coalesce(pw.description, '') plyrwnnm, " +
		"row_number() over(order by coalesce(ordr,999999)) real_ordr, " +
		"coalesce(m1.letters,''), coalesce(m2.letters,'') " +
		"from (select g.id, min(gp1.plyrtm_id) plyr1, max(gp2.plyrtm_id) plyr2, " +
		"	coalesce(g.ordr,999999) ordr, coalesce(g.plyrtm_wn,-1) plyrwn, " +
		"	coalesce(g.pf,0) pf, coalesce(g.pa,0) pa, coalesce(g.notes,'') notes " +
		"	from gms g inner join gms_plyrtm gp1 on gp1.gms_id = g.id " +
		"	inner join gms_plyrtm gp2 on gp2.gms_id = g.id and gp1.plyrtm_id <> gp2.plyrtm_id " +
		"	where tourn_id = ? and rnd = ? and gp1.id < gp2.id " +
		"	group by g.id, coalesce(g.ordr,0), g.plyrtm_wn, coalesce(g.pf,0), coalesce(g.pa,0), " +
		"	coalesce(g.notes,0) " +
		"    ) a inner join plyrtm p1 on plyr1 = p1.id " +
		"inner join gms_plyrtm gp1 on gp1.gms_id = a.id and gp1.plyrtm_id = p1.id " +
		"inner join plyrtm p2 on plyr2 = p2.id " +
		"inner join gms_plyrtm gp2 on gp2.gms_id = a.id and gp2.plyrtm_id = p2.id " +
		"left join plyrtm pw on plyrwn = pw.id " +
		"left join lateral (select m.letters " +
		"	from msc m where m.rnd < ? and m.plyrtm_id = p1.id " +
		"	and m.letters in ('h','b','d','p','e','ee') " +
		"	order by case letters when 'h' then 1 when 'b' then 2 when 'd' then 3 when 'p' then 4 when 'e' then 5 when 'ee' then 6 else 100 end " +
		"	limit 1) m1 on true " +
		"	left join lateral (select m.letters " +
		"	from msc m where m.rnd < ? and m.plyrtm_id = gp2.plyrtm_id " +
		"	and m.letters in ('h','b','d','p','e','ee') " +
		"	order by case letters when 'h' then 1 when 'b' then 2 when 'd' then 3 when 'p' then 4 when 'e' then 5 when 'ee' then 6 else 100 end " +
		"	limit 1) m2 on true " +
		"order by a.ordr"
	rows, err = db.DB.Query(sql, trnId, rnd, rnd, rnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gid, ordr, realOrdr, rec1, tie1, rec2, tie2 int
	var notes, plyr1nm, plyr2nm, plyrwnnm, letter1, letter2 string
	assignOrdr := make(map[int]int)
	var sch []model.Sch

	for rows.Next() {
		if err = rows.Scan(&gid, &plyr1, &plyr2, &ordr, &plyrwn, &pf, &pa, &notes,
			&plyr1nm, &plyr2nm, &rec1, &tie1, &rec2, &tie2, &plyrwnnm, &realOrdr, &letter1, &letter2); err != nil {
			return nil, err
		}
		var schTemp model.Sch
		schTemp.Gid = gid
		schTemp.Plyr1 = plyr1
		schTemp.Plyr2 = plyr2
		schTemp.Ordr = ordr
		schTemp.Plyrwn = plyrwn
		schTemp.Pf = pf
		schTemp.Pa = pa
		schTemp.Plyr1nm = plyr1nm
		schTemp.Plyr2nm = plyr2nm
		schTemp.Plyrwnnm = plyrwnnm
		schTemp.Notes = notes
		schTemp.Rec1 = rec1
		schTemp.Tie1 = tie1
		schTemp.Rec2 = rec2
		schTemp.Tie2 = tie2
		schTemp.Msc1 = mscArr[plyr1]
		schTemp.Msc2 = mscArr[plyr2]
		schTemp.Letter1 = letter1
		schTemp.Letter2 = letter2

		if archGmeArr[plyr1] != "" {
			schTemp.ArchNm = plyr1nm
			schTemp.ArchPts = archGmeArr[plyr1]
		} else if archGmeArr[plyr2] != "" {
			schTemp.ArchNm = plyr2nm
			schTemp.ArchPts = archGmeArr[plyr2]
		}

		sch = append(sch, schTemp)
		if ordr != realOrdr {
			assignOrdr[gid] = realOrdr
		}
	}

	for gid, ordr = range assignOrdr {
		sql = "update gms set ordr = ? where id = ?"
		params := []interface{}{ordr, gid}
		_, err = db.ChangDB(sql, params)
		if err != nil {
			return nil, err
		}
	}

	jsonOut, _ := json.MarshalIndent(sch, "", "  ")
	return jsonOut, nil
}

func SchSwap(trnId int, rnd int, gid int, up bool) ([]byte, error) {
	sql := "select ordr from gms where id = ?"
	rows, err := db.DB.Query(sql, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordr int
	for rows.Next() {
		if err = rows.Scan(&ordr); err != nil {
			return nil, err
		}
	}

	if ordr > 0 {
		sql = "update gms set ordr = ? where tourn_id = ? and rnd = ? and ordr = ?"
		if up {
			params := []interface{}{ordr, trnId, rnd, ordr - 1}
			_, err = db.ChangDB(sql, params)
		} else {
			params := []interface{}{ordr, trnId, rnd, ordr + 1}
			_, err = db.ChangDB(sql, params)
		}
		if err != nil {
			return nil, err
		}
		sql = "update gms set ordr = ? where id = ?"
		if up {
			params := []interface{}{ordr - 1, gid}
			_, err = db.ChangDB(sql, params)
		} else {
			params := []interface{}{ordr + 1, gid}
			_, err = db.ChangDB(sql, params)
		}
		if err != nil {
			return nil, err
		}
	}

	resp := model.Response{Msg: "ok"}
	jsonOut, _ := json.MarshalIndent(resp, "", "  ")
	return jsonOut, nil
}

func UpdateRnd(gid int, plyrwm int, pf int, pa int, notes string) ([]byte, error) {
	sql := "select ordr from gms where id = ?"
	rows, err := db.DB.Query(sql, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordr int
	for rows.Next() {
		if err = rows.Scan(&ordr); err != nil {
			return nil, err
		}
	}
	if ordr > 0 {
		sql := "update gms set plyrtm_wn = ?, pf = ?, pa = ?, notes = ? where id = ?"
		params := []interface{}{plyrwm, pf, pa, notes, gid} // accepts both int and string
		_, err = db.ChangDB(sql, params)
		if err != nil {
			return nil, err
		}
	}
	resp := model.Response{Msg: "ok"}
	jsonOut, _ := json.MarshalIndent(resp, "", "  ")
	return jsonOut, nil
}

func AddDeleteMsc(trnId int, rnd int, plyrtm int, letters string, addDelete int) ([]byte, error) {
	if addDelete == 1 {
		sql := "insert into msc (tourn_id, rnd, letters, plyrtm_id) " +
			"select ?, ?, ?, ? " +
			"where not exists (select 1 from msc m2 " +
			"					where tourn_id = ? and rnd = ? " +
			"					and letters = ? and plyrtm_id = ?)"
		params := []interface{}{trnId, rnd, letters, plyrtm, trnId, rnd, letters, plyrtm} // accepts both int and string
		_, err := db.ChangDB(sql, params)
		if err != nil {
			return nil, err
		}
	} else {
		sql := "delete from msc " +
			"where tourn_id = ? and rnd = ? " +
			"and letters = ? and plyrtm_id = ?"
		params := []interface{}{trnId, rnd, letters, plyrtm} // accepts both int and string
		_, err := db.ChangDB(sql, params)
		if err != nil {
			return nil, err
		}
	}

	resp := model.Response{Msg: "ok"}
	jsonOut, _ := json.MarshalIndent(resp, "", "  ")
	return jsonOut, nil
}


func OrdNum(trnId int, rnd int, oldnum int, newnum int) ([]byte, error) {
	if (newnum > oldnum) {
		sql := "update gms set ordr = 0 where tourn_id = ? and rnd = ? and ordr = ?"
		params := []interface{}{trnId, rnd, oldnum}
		_, err := db.ChangDB(sql, params)
		if err != nil {
			return nil, err
		}
		sql = "update gms set ordr = ordr -1 where tourn_id = ? and rnd = ? and ordr > ? and ordr <= ?"
		params = []interface{}{trnId, rnd, oldnum, newnum}
		_, err = db.ChangDB(sql, params)
		if err != nil {
			return nil, err
		}
		sql = "update gms set ordr = ? where tourn_id = ? and rnd = ? and ordr = 0"
		params = []interface{}{newnum, trnId, rnd}
		_, err = db.ChangDB(sql, params)
		if err != nil {
			return nil, err
		}
	} else if (newnum < oldnum) {
		sql := "update gms set ordr = 0 where tourn_id = ? and rnd = ? and ordr = ?"
		params := []interface{}{trnId, rnd, oldnum}
		_, err := db.ChangDB(sql, params)
		if err != nil {
			return nil, err
		}
		sql = "update gms set ordr = ordr +1 where tourn_id = ? and rnd = ? and ordr < ? and ordr >= ?"
		params = []interface{}{trnId, rnd, oldnum, newnum}
		_, err = db.ChangDB(sql, params)
		if err != nil {
			return nil, err
		}
		sql = "update gms set ordr = ? where tourn_id = ? and rnd = ? and ordr = 0"
		params = []interface{}{newnum, trnId, rnd}
		_, err = db.ChangDB(sql, params)
		if err != nil {
			return nil, err
		}
	}

	resp := model.Response{Msg: "ok"}
	jsonOut, _ := json.MarshalIndent(resp, "", "  ")
	return jsonOut, nil
}

func Potw(trnId int) ([]byte, error) {
	potwList := []model.Potw{}

	sql := `select m.rnd, m.letters, p.description
		from msc m inner join plyrtm p on p.id = m.plyrtm_id 
		where m.letters in ('potw','rnk') and m.tourn_id = ?
		order by m.rnd, m.letters, p.description`
	rows, err := db.DB.Query(sql, trnId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	var rnd int
	var letters string
	var plyrtm string
	for rows.Next() {
		if err = rows.Scan(&rnd, &letters, &plyrtm); err != nil {
			return nil, err
		}

		potwList = append(potwList, model.Potw{Rnd: rnd, Letters: letters, Plyr: plyrtm})
	}

	jsonOut, err := json.MarshalIndent(potwList, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonOut, nil
}