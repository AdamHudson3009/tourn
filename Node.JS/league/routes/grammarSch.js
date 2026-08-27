const express = require('express');
const router = express.Router();
const db = require('../db');


// Get a list of tournaments with grammar defined
router.post('/GetCopyFrom', async (req, res) => {
  const { id } = req.body; //id must be the key in json
  try {
    const sql = `SELECT distinct tr.id tournId, tr.description tourn 
                from grammar g inner join tourn tr on tr.id = g.tourn_id
                where tr.league_id in (select t1.league_id from tourn t1 where t1.id = ?)
                order by tr.description `;
    const [rows] = await db.query(sql, [id]);
    res.status(200).json(rows);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

//Copy Grammar and schedule from tournament to tournament
//but only if grammar is not defined for the To tournmanent
//only copy if conf and divi names are the same.
router.post('/CopyFrom', async (req, res) => {
  const { TournIdFrom, TournIdTo } = req.body; //id must be the key in json
  try {
    let sql = "SELECT 1 from grammar where tourn_id = ?";
    let [rows] = await db.query(sql, [TournIdTo]);            
    let found = false;
    for (const row of rows) {
        found = true;
    }
    if (found) {
        return res.status(200).send("Already Exists");
    }

    sql = `SELECT g.id gid, letters, dv1_id d1, dv2_id d2, dnew1.id dnew1, dnew2.id dnew2 
        from grammar g 
        inner join divi d1 on d1.id = g.dv1_id
        inner join conf c1 on c1.id = d1.conf_id
        inner join divi d2 on d2.id = g.dv2_id
        inner join conf c2 on c2.id = d2.conf_id
        inner join conf cnew1 on cnew1.description = c1.description and cnew1.tourn_id = ?
        inner join divi dnew1 on dnew1.description = d1.description and dnew1.conf_id = cnew1.id
        inner join conf cnew2 on cnew2.description = c2.description and cnew2.tourn_id = cnew1.tourn_id
        inner join divi dnew2 on dnew2.description = d2.description and dnew2.conf_id = cnew2.id
        where g.tourn_id = ?`;

    [rows] = await db.query(sql, [TournIdTo, TournIdFrom]);
    let GrammarIdMap = {};  
    let GrammarIdIn = [];

    for (const row of rows) {
      sql = `insert into grammar (tourn_id, letters, dv1_id, dv2_id) 
          values (?, ?, ?, ?)`
      let [result] = await db.query(sql, [TournIdTo, row.letters, row.dnew1, row.dnew2]);  
      GrammarIdMap[row.gid] = result.insertId;
      GrammarIdIn.push(row.gid);
    }
    console.log(GrammarIdMap);

    let placeholders = GrammarIdIn.map(() => '?').join(', ');
    sql = `SELECT id, rnd, grammar_id, rpt  
        from wk_sch  
        where grammar_id in (${placeholders})`;
    
    [rows] = await db.query(sql, GrammarIdIn);
    let WkSchIdMap = {};  
    let WkSchIdIn = [];

    for (const row of rows) {
      sql = `insert into wk_sch (rnd, grammar_id, rpt) 
          values (?, ?, ?)`;
      let [result] = await db.query(sql, [row.rnd, GrammarIdMap[row.grammar_id], row.rpt]);
      WkSchIdMap[row.id] = result.insertId;
      WkSchIdIn.push(row.id);
    }

    placeholders = WkSchIdIn.map(() => '?').join(', ');
    sql = `SELECT wk_sch_id, sd1, sd2
        from wk_sch_sds  
        where wk_sch_id in (${placeholders})`;
    
    [rows] = await db.query(sql, WkSchIdIn);
    
    for (const row of rows) {
      sql = `insert into wk_sch_sds (wk_sch_id, sd1, sd2) 
          values (?, ?, ?)`;
      let [result] = await db.query(sql, [WkSchIdMap[row.wk_sch_id], row.sd1, row.sd2]);
    }

    return res.status(200).send("OK");
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});


// THIS IS REQUIRED
module.exports = router;