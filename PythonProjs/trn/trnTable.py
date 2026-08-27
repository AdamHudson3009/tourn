from fastapi import APIRouter
from pydantic import BaseModel
from db import get_db_connection  # Import the database function
from typing import List
import logging

class TrnTableRequest(BaseModel):
    tourn_id: int

# Create a router
trnTableRtr = APIRouter()

@trnTableRtr.post("/trn-table/get")
def get_trnTable(request: TrnTableRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    sql = "select c.id cid, c.description conf, d.id did, d.description divi, " \
        "p.id pid, p.description plyrtm, coalesce(p.description2,'') plyrtm2, sds " \
        "from conf c left join divi d on d.conf_id = c.id " \
        "left join plyrtm p on p.divi_id = d.id " \
        "where c.tourn_id = %s " \
        "order by c.description, d.description, p.sds"
    cursor.execute(sql, (request.tourn_id,))
    records=cursor.fetchall()
    confs = {}
    for record in records:
        cid = record["cid"]
        if cid not in confs:
            confs[cid] = {
                "cid": cid,
                "conf": record["conf"],
                "divis": {}
            }
        did = record["did"]
        if did is None:
            continue
        if did not in confs[cid]["divis"]:
            confs[cid]["divis"][did] = {
                "did": did,
                "divi": record["divi"],
                "plyrtms": []
            }
        pid = record["pid"]
        if pid is None:
            continue
        confs[cid]["divis"][did]["plyrtms"].append({
            "pid": pid,
            "plyrtm": record["plyrtm"],
            "plyrtm2": record["plyrtm2"],
            "sds": record["sds"]
        })    
    conn.close()

    # Convert nested dicts to lists for JSON compatibility
    result = []
    for conf in confs.values():
        conf["divis"] = list(conf["divis"].values())
        result.append(conf)

    return result