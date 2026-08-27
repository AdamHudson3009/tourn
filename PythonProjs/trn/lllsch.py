from fastapi import APIRouter
from pydantic import BaseModel
from db import get_db_connection  # Import the database function
from typing import List
import logging

# Create a router
lllschRtr = APIRouter()

class lllschRequest(BaseModel):
    tourn_id: int
    gid: int
    rnd: int
    rpt: str
    w1id: int
    w2id: List[int]
    sd1: List[int]
    sd2: List[int]
    HasStnd: bool 

class lllschCopy(BaseModel):
    tourn_id_from: int
    tourn_id_to: int

@lllschRtr.post("/lllsch/get")
def get_lllsch(request: lllschRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json formath
    sql = "select w1.id w1id, w1.rnd, g.id gid, w1.rpt, w2.id w2id, w2.sd1, w2.sd2, " \
        "coalesce(w1.has_stnd,'Y') has_stnd " \
        "from wk_sch w1 inner join grammar g on w1.grammar_id = g.id " \
        "left join wk_sch_sds w2 on w2.wk_sch_id = w1.id " \
        "where g.tourn_id = %s " \
        "order by w1.rnd, g.letters, w2.sd1"
    cursor.execute(sql, (request.tourn_id,))
    records=cursor.fetchall()
    groups = {"rnds": []}
    rnd_map = {}  # To track existing rnd objects

    for record in records:
        rnd = record["rnd"]
        w1id = record["w1id"]
        HasStnd = record["has_stnd"] == "Y"

        if rnd not in rnd_map:
            rnd_obj = {"rnd": rnd, "entries": []}
            rnd_map[rnd] = rnd_obj
            groups["rnds"].append(rnd_obj)

        rnd_info = rnd_map[rnd]
        entry_map = {entry["w1id"]: entry for entry in rnd_info["entries"]}
    
        if w1id not in entry_map:
            entry = {
                "w1id": w1id,
                "HasStnd": HasStnd,
                "rpt": record["rpt"],
                "gid": record["gid"],
                "w2Info": []
            }
            rnd_info["entries"].append(entry)
        else:
            entry = entry_map[w1id]

        w2id = record["w2id"]
        if w2id is not None:
            entry["w2Info"].append({
                "w2id": w2id,
                "sd1": record["sd1"],
                "sd2": record["sd2"]
            })

    sql = "select id, letters from grammar " \
        "where tourn_id = %s order by letters"
    cursor.execute(sql, (request.tourn_id,))
    records=cursor.fetchall()
    groups["grammar"] = []
    for record in records:
        groups["grammar"].append({
            "letters": record["letters"],
            "gid": record["id"]
        })

    conn.commit()
    conn.close()

    return groups

@lllschRtr.post("/lllsch/add")
def add_lllsch(request: lllschRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    sql = "select 1 from wk_sch w1 " \
        "inner join grammar g on w1.grammar_id = g.id " \
        "where tourn_id = %s and w1.rnd = %s and g.id = %s"
    cursor.execute(sql, (request.tourn_id,request.rnd,request.gid,))
    result = cursor.fetchone() 
    if result is not None:
        conn.close()
        return "combo exists: "+request.rnd+", "+request.gid
    
    HasStndY = "N"
    if request.HasStnd:
        HasStndY = "Y"
    sql = "insert into wk_sch (rnd, grammar_id, rpt, has_stnd) " \
        "values (%s,%s,%s,%s)"
    cursor.execute(sql, (request.rnd,request.gid,request.rpt,HasStndY,))
    inserted_id = cursor.lastrowid
    for indx in range(len(request.w2id)):
        sql = "insert into wk_sch_sds (wk_sch_id, sd1, sd2) " \
            "values (%s,%s,%s)"
        cursor.execute(sql, (inserted_id,request.sd1[indx],request.sd2[indx],))
    conn.commit()
    conn.close()

@lllschRtr.post("/lllsch/change")
def change_lllsch(request: lllschRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    print("Here" + str(request.rnd))
    sql = "select 1 from wk_sch w " \
        "inner join grammar g on w.grammar_id = g.id " \
        "where g.tourn_id = %s and g.id = %s and w.id <> %s and rnd = %s"
    cursor.execute(sql, (request.tourn_id,request.gid,request.w1id,request.rnd))
    result = cursor.fetchone()
    HasStndNum = "N"
    if request.HasStnd:
        HasStndNum = "Y"
    if result is not None:
        conn.close()
        return "id exists"
    sql = "update wk_sch set grammar_id = %s, rpt = %s, rnd = %s, has_stnd = %s where id = %s"
    cursor.execute(sql, (request.gid,request.rpt,request.rnd,HasStndNum,request.w1id,))
    
    w2idArr = []
    for indx in range(len(request.w2id)):
        w2idArr.append(request.w2id[indx])
        if request.w2id[indx] == 0:
            sql = "insert into wk_sch_sds (wk_sch_id,sd1,sd2) " \
                "values (%s,%s,%s)"
            cursor.execute(sql, (request.w1id,request.sd1[indx],request.sd2[indx],))
            w2idArr[indx] = cursor.lastrowid
        else:
            sql = "update wk_sch_sds set sd1 = %s, sd2 = %s where id = %s"
            cursor.execute(sql, (request.sd1[indx],request.sd2[indx],w2idArr[indx],))

    if w2idArr:
        placeholders = ','.join(['%s'] * len(w2idArr))
        sql = f"DELETE FROM wk_sch_sds WHERE wk_sch_id = %s AND id NOT IN ({placeholders})"
        params = [request.w1id] + w2idArr  # wk_sch_id first, then IDs to exclude
    else:
        sql = "DELETE FROM wk_sch_sds WHERE wk_sch_id = %s"
        params = [request.w1id]
    cursor.execute(sql, params)
    conn.commit()
    conn.close()
    
    return "Changed record"

@lllschRtr.post("/lllsch/delete")
def delete_lllsch(request: lllschRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    
    sql = "delete from wk_sch_sds where wk_sch_id = %s"
    cursor.execute(sql, (request.w1id,))
    sql = "delete from wk_sch where id = %s"
    cursor.execute(sql, (request.w1id,))
    conn.commit()
    conn.close()

    return "Deleted record"

@lllschRtr.post("/lllsch/CopyFromTourn")
def CopyFromTourn_lllsch(request: lllschCopy):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format

    sql = "select id, letters, dv1_id, dv2_id from grammar where tourn_id = %s"
    cursor.execute(sql, (request.tourn_id_to,))
    records=cursor.fetchall()
    for record in records:
        conn.close()
        return "not empty: "+str(request.tourn_id_to) 
    
    sql  = "select c1.id cid1, coalesce(a.cid,0) cid2, d1.id did1, coalesce(a.did,0) did2 " \
        "from conf c1 inner join divi d1 on d1.conf_id = c1.id " \
        "left join (select c2.id cid, c2.description conf, d2.id did, d2.description divi " \
		"           from conf c2 inner join divi d2 on d2.conf_id = c2.id " \
		"           where c2.tourn_id = %s) a on a.conf = c1.description and a.divi = d1.description " \
        "where c1.tourn_id = %s " \
        "order by c1.id, d1.id "
    cursor.execute(sql, (request.tourn_id_to,request.tourn_id_from,))
    records=cursor.fetchall()
    
    confs = {}
    divis = {}
    for record in records:
        cid1 = record["cid1"]
        if cid1 not in confs:
            confs[cid1] = record["cid2"]
        did1 = record["did1"]
        divis[did1] = record["did2"]

    sql = "select id, letters, dv1_id, dv2_id from grammar where tourn_id = %s "
    cursor.execute(sql, (request.tourn_id_from,))
    records=cursor.fetchall()
    recordSet = []
    for record in records:
        id = record["id"]
        letters = record["letters"]
        dv1_id = record["dv1_id"]
        dv2_id = record["dv2_id"]
        
        recordSet.append({"id": id, "letters": letters, "dv1_id": dv1_id, "dv2_id": dv2_id})

    grammars = {}
    for indx in range(recordSet):
        gid = recordSet[indx]["id"]
        letters = recordSet[indx]["letters"]
        dv1_id = recordSet[indx]["dv1_id"]
        dv2_id = recordSet[indx]["dv2_id"]
        sql = "insert into grammar (letters, dv1_id, dv2_id, tourn_id) " \
            "values (%s,%s,%s,%s)"
        cursor.execute(sql, (letters,divis[dv1_id],divis[dv2_id],request.tourn_id_to,))
        grammars[gid] = cursor.lastrowid

    sql = "select id, rnd, grammar_id, rpt from wk_sch " \
        "where grammar_id in (select id from grammar where tourn_id = %s)"
    cursor.execute(sql, (request.tourn_id_from,))
    records=cursor.fetchall()
    recordSet = []
    for record in records:
        id = record["id"]
        rnd = record["rnd"]
        grammar_id = record["grammar_id"]
        rpt = record["rpt"]
        
        recordSet.append({"id": id, "rnd": rnd, "grammar_id": grammar_id, "rpt": rpt})

    wksch = {}
    for indx in range(recordSet):
        swid = recordSet[indx]["id"] 
        gid = recordSet[indx]["grammar_id"]
        rnd = recordSet[indx]["rnd"]
        rpt = recordSet[indx]["rpt"]
        sql = "insert into wk_sch (rnd, grammar_id, rpt) values (%s,%s,%s)"
        cursor.execute(sql, (rnd,grammars[gid],rpt,))
        wksch[swid] = cursor.lastrowid

    sql = "select wk_sch_id, sd1, sd2 from wk_sch_sds " \
        "where wk_sch_id in (select id from wk_sch where grammar_id in " \
         "                  (select id from grammar where tourn_id = %s))"
    cursor.execute(sql, (request.tourn_id_from,))
    records=cursor.fetchall()
    recordSet = []
    for record in records:
        wk_sch_id = record["wk_sch_id"]
        sd2 = record["sd1"]
        sd1 = record["sd2"]
        
        recordSet.append({"wk_sch_id": wk_sch_id, "sd1": sd1, "sd2": sd2})

    for indx in range(recordSet):       
        swid = recordSet[indx]["wk_sch_id"]
        sd1 = recordSet[indx]["sd1"]
        sd2 = recordSet[indx]["sd2"]
        sql = "insert into wk_sch_sds (wk_sch_sds, sd1, sd2) values (%s,%s,%s)"
        cursor.execute(sql, (wksch[swid], sd1, sd2,))

    conn.close()

    return "success: "+str(request.tourn_id_from) + " -- " + str(request.tourn_id_to) 
