from fastapi import APIRouter
from pydantic import BaseModel
from db import get_db_connection  # Import the database function
from typing import List
import logging

# Create a router
cmmRtr = APIRouter()

# Define a Pydantic model for request payload
class CmmRequest(BaseModel):
    tourn_id: int
    league_id: int
    divi_id: List[int]
    num1: List[int]
    num2: List[int]  

# Define a Pydantic model for request payload
class LTRequest(BaseModel):
    tourn_id: int
    league_id: int

@cmmRtr.post("/cmm/get")
def cmm_get(request: LTRequest):
    conn = get_db_connection()
    sql1 = "select c.description conf, d.id did, d.description divi, " \
            "COALESCE(conf_num,0) conf_num, COALESCE(i_num,0) i_num " \
            "from divi d inner join conf c on c.id = d.conf_id "
    sql2 = "left join cmm on d.id = cmm.divi_id "
    cursor=conn.cursor(dictionary=True) #json format
    if request.league_id > 0:
        sql = sql1 + "inner join tourn t on c.tourn_id = t.id " + sql2 + "and cmm.league_id = t.league_id " \
            "where t.league_id = %s order by c.description, d.description"
        cursor.execute(sql, (request.league_id,))
    else:
        sql = sql1 + sql2 + "and cmm.tourn_id = c.tourn_id " \
            "where c.tourn_id = %s order by c.description, d.description"
        cursor.execute(sql, (request.tourn_id,))
    records=cursor.fetchall()  
    confs = {}
    for record in records:
        conf = record["conf"]
        if conf not in confs:
            confs[conf] = {
                "conf": conf,
                "descriptions": []
            }
        confs[conf]["descriptions"].append({
            "divi": record["divi"],
            "did": record["did"],
            "conf_num": record["conf_num"],
            "i_num": record["i_num"]
        })       
    conn.close()

    result = []
    for conf in confs.values():
        result.append(conf)

    return result    

@cmmRtr.post("/cmm/change")
def cmm_get(request: CmmRequest):
    returnMsg = ""
    try:
        conn = get_db_connection()
        cursor=conn.cursor(dictionary=True) #json format
 
        for indx in range(len(request.divi_id)):
            sql = "select 1 from cmm where divi_id = %s"
            cursor.execute(sql, (request.divi_id[indx],))
            result = cursor.fetchone() 
            if result is None:
                if request.league_id > 0:
                    sql = "insert into cmm (league_id, divi_id, conf_num, i_num) " \
                        "values (%s, %s, %s, %s)"
                    cursor.execute(sql, (request.league_id,request.divi_id[indx], \
                                         request.num1[indx],request.num2[indx],))
                else:
                    sql = "insert into cmm (tourn_id, divi_id, conf_num, i_num) " \
                        "values (%s, %s, %s, %s)"
                    cursor.execute(sql, (request.tourn_id,request.divi_id[indx], \
                                         request.num1[indx],request.num2[indx],))
            else:
                sql = "update cmm set conf_num = %s, i_num = %s where divi_id = %s"
                cursor.execute(sql, (request.num1[indx],request.num2[indx],request.divi_id[indx],))
        conn.commit()
        returnMsg = "Changes saved"
    except Exception as e:
        returnMsg = f"Insert error: {e}"
    finally:
        conn.close()
    return returnMsg

@cmmRtr.post("/cmm/delete")
def cmm_get(request: LTRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    if request.league_id > 0:
        sql = "delete from cmm where league_id = %s"
        cursor.execute(sql, (request.league_id,))
    else:
        sql = "delete from cmm where tourn_id = %s"
        cursor.execute(sql, (request.tourn_id,))
    records=cursor.fetchall()           
    conn.close()
    return records

@cmmRtr.post("/cmm/league_to_trn")
def league_to_trn(request: LTRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    sql = "delete from cmm where tourn_id = %s"
    cursor.execute(sql, (request.tourn_id,))
    sql = "update cmm set tourn_id = %s, league_id = null where league_id = %s"
    cursor.execute(sql, (request.tourn_id,request.league_id,))
    conn.commit()
    conn.close()
    return "Copied trn to league "+ str(request.tourn_id) + ": "+ str(request.league_id)

@cmmRtr.post("/cmm/trn_to_league")
def trn_to_league(request: LTRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    sql = "delete from cmm where league_id = %s"
    cursor.execute(sql, (request.league_id,))
    sql = "update cmm set league_id = %s, tourn_id = null where tourn_id = %s"
    cursor.execute(sql, (request.league_id,request.tourn_id,))
    conn.commit()
    conn.close()
    return "Copied trn to league "+ str(request.league_id) + ": " + str(request.tourn_id)