const express = require('express');
const routerPlf = express.Router();
const db = require('../db');

routerPlf.post('/ListPlfRnd', async (req, res) => {
    const { TournId } = req.body;
    try {
        const sql = `select rnd, rnd_nme from plf_nme where tourn_id = ?`
        [rows] = await db.query(sql, [TournId, rnd]);
        let rnds = [];
        for (const rnd of rnds) {
            const rnd = {
                rnd: rnd.rnd, 
                rnd_nme: rnd_nme
            };
            rnds.push(rnd);
        }
        return res.status(200).send("OK");
    } catch (err) {
        res.status(500).json({ error: err.message });
    } 
});

routerPlf.post('/ListPlfRow', async (req, res) => {
    try {
        const { TournId, rnd, rndNme } = req.body;
        const sql = `select a.*, p1.description plyr1nm, p2.description plyr2nm, 
		    coalesce(gp1.rec,0) rec1, coalesce(gp1.tie,0) tie1, 
		    coalesce(gp2.rec,0) rec2, coalesce(gp2.tie,0) tie2, 
		    coalesce(pw.description, '') plyrwnnm, 
		    row_number() over(order by coalesce(ordr,999999)) real_ordr
		    from (select g.id, min(gp1.plyrtm_id) plyr1, max(gp2.plyrtm_id) plyr2, 
			    coalesce(g.ordr,999999) ordr, coalesce(g.plyrtm_wn,-1) plyrwn, 
			    coalesce(g.pf,0) pf, coalesce(g.pa,0) pa, coalesce(g.notes,'') notes 
			    from gms g inner join gms_plyrtm gp1 on gp1.gms_id = g.id 
			    inner join gms_plyrtm gp2 on gp2.gms_id = g.id and gp1.plyrtm_id <> gp2.plyrtm_id 
			    where tourn_id = ? and rnd = ? 
			    group by g.id, coalesce(g.ordr,0), g.plyrtm_wn, coalesce(g.pf,0), coalesce(g.pa,0), 
			    coalesce(g.notes,0) 
		    ) a inner join plyrtm p1 on plyr1 = p1.id 
		    inner join gms_plyrtm gp1 on gp1.gms_id = a.id and gp1.plyrtm_id = p1.id 
		    inner join plyrtm p2 on plyr2 = p2.id 
		    inner join gms_plyrtm gp2 on gp2.gms_id = a.id and gp2.plyrtm_id = p2.id 
		    left join plyrtm pw on plyrwn = pw.id 
		    order by a.ordr`
        [rows] = await db.query(sql, [TournId, rnd]);
        let mtchs = [];
        for (const row of rows) {
            const mtch = {
                gid: row.gid,
                plyr1: row.plyr1,
                plyr2: row.plyr2,
                ordr: row.ordr,
                plyrwn: row.plyrwn,
                pf: row.pf,
                pa: row.pa,
                notes: row.notes,
                plyr1nm: row.plyr1nm, 
                plyr2nm: row.plyr2nm,
                rec1: row.rec1,
                tie1: row.tie1,
                rec2: row.rec2,
                tie2: row.tie2,
                plyrwnnm: row.plyrwnnm,
                real_ordr: row.real_ordr
            };
            mtchs.push(mtch);
        }
        res.status(201).json(mtchs);
    } catch (err) {
        res.status(500).json({ error: err.message });
    }  
});

routerPlf.post('/AddPlfRnd', async (req, res) => {
    const { TournId, rnd, rndNme } = req.body;
    try {
        const sql = `insert into plf_nme (tourn_id, rnd, rnd_nme) 
                    values (?, ?, ?)`
        let [result] = await db.query(sql, [TournId, rnd, rndNme]); 
        return res.status(200).send("OK");
    } catch (err) {
        res.status(500).json({ error: err.message });
    } 
});

routerPlf.post('/EditPlfRnd', async (req, res) => {
    const { TournId, rnd, rndNme } = req.body;
    try {
        const sql = `update plf_nme set rnd_nme = ?  
                    where tourn_id = ? and rnd = ?`
        let [result] = await db.query(sql, [rndNme, TournId, rnd]);   
        return res.status(200).send("OK");
    } catch (err) {
        res.status(500).json({ error: err.message });
    }          
});

routerPlf.post('DeletePlfRnd', async (req, res) => {
    const { TournId, rnd } = req.body;
    try {
        let sql = `select 1 from gms where turn_id = ? and rnd = ?`
        [rows] = await db.query(sql, [TournId, rnd]);
        found = false
        for (const row of rows) {
            found = true
            break
        }
        if (found) {
            return res.status(200).send("Matches exist, can not delete")
        }
        sql = `delete from plf_nme where tourd_id = ? and rnd = ?`  
        let [result] = await db.query(sql, [TournId, rnd]);   
        return res.status(200).send("OK");
    } catch (err) {
        res.status(500).json({ error: err.message });
    }          
});

routerPlf.post('/AddPlfRow', async (req, res) => {
    const { TournId, rnd, plyrtm1, plyrtm2 } = req.body;
    try {
        sql = `select max(ordr)+1 ordr from gms
            where tourn_id = ? and rnd = ?`
        let ordr = 1
        [rows] = await db.query(sql, [TournId, rnd]);
        for (const row of rows) {
            ordr = row.ordr
        }

        sql = `insert into gms (tourn_id, rnd, ordr) 
            values (?, ?, ?)`
        let [result] = await db.query(sql, [TournId, rnd, ordr]); 
        
        let plyrtm = plyrtm1
		for (let indxx = 0; indxx < 2; indxx++) {
			sql = `insert into gms_plyrtm (gms_id, plyrtm_id) 
				values (?,?)`
            let [result] = await db.query(sql, [gmsId, plyrtm]);   

			plyrtm = plyrtm2
		}

        return res.status(200).send("OK");
    } catch (err) {
        res.status(500).json({ error: err.message });
    } 
});

routerPlf.post('/EditPlfRow', async (req, res) => {
    const { plyrtm_wn, pf, pa, notes, id } = req.body;
    try {
        const sql = `update gms set plyrtm_wn = ?, pf = ?, pa = ?, notes = ? 
                    where id = ?`
        let [result] = await db.query(sql, [plyrtm_wn, pf, pa, notes, id]); 
        return res.status(200).send("OK");
    } catch (err) {
        res.status(500).json({ error: err.message });
    } 
});

routerPlf.post('/DeletePlfRow', async (req, res) => {
    const { id } = req.body;
    try {
        let sql = `delete from gms_plyrtm where gms_id = ?`
        let [result] = await db.query(sql, [id]); 
        sql = `delete from gms where id = ?`
        [result] = await db.query(sql, [id]);
        return res.status(200).send("OK");
    } catch (err) {
        res.status(500).json({ error: err.message });
    } 
});

routerPlf.post('/EditPlfRow', async (req, res) => {
    const { plyrtm_wn, pf, pa, notes, id } = req.body;
    try {
        const sql = `update gms set plyrtm_wn = ?, pf = ?, pa = ?, notes = ? 
                    where id = ?`
        let [result] = await db.query(sql, [plyrtm_wn, pf, pa, notes, id]); 
        return res.status(200).send("OK");
    } catch (err) {
        res.status(500).json({ error: err.message });
    } 
});

routerPlf.post('/PlfSwp', async (req, res) => {
    const { trnId, rnd, gid, up } = req.body;
    try {
        const sql = "select ordr from gms where id = ?"
        [rows] = await db.query(sql, gid);
        let ordr = 0
        for (const row of rows) {
            ordr = row.ordr
        }
        if (ordr > 0) {
            sql = "update gms set ordr = ? where tourn_id = ? and rnd = ? and ordr = ?"
		    if (up) {
                let [result] = await db.query(sql, [ordr, trnId, rnd, ordr - 1]); 
		    } else {
			    let [result] = await db.query(sql, [ordr, trnId, rnd, ordr + 1]);
		    }

		    sql = "update gms set ordr = ? where id = ?"
		    if (up) {
			    let [result] = await db.query(sql, [ordr - 1, gid]);
		    } else {
			    let [result] = await db.query(sql, [ordr + 1, gid]);
		    }
        }
    } catch (err) {
        res.status(500).json({ error: err.message });
    } 
});

// THIS IS REQUIRED
module.exports = routerPlf;