from fastapi import APIRouter
from pydantic import BaseModel
from db import get_db_connection  # Import the database function
from typing import List
import logging

# Create a router
tbrkRtr = APIRouter()

# Define a Pydantic model for request payload
class TbrkRequest(BaseModel):
    tourn_id: int
    league_id: int
    dvCnOth: str
    id: int
    descriptions: List[str]  
    rnk: int
    up: bool

# Define a Pydantic model for request payload
class LTRequest(BaseModel):
    tourn_id: int
    league_id: int
    
@tbrkRtr.post("/tbrk/get")
def get_tbrk(request: TbrkRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    if request.league_id > 0:
        sql = "select * from tbrk where league_id = %s order by dv_cn_other, rnk"
        cursor.execute(sql, (request.league_id,))
    else:
        sql = "select * from tbrk where tourn_id = %s order by dv_cn_other, rnk"
        cursor.execute(sql, (request.tourn_id,))
    records=cursor.fetchall()
    conn.close()
    groups = {}
    for record in records:
        dv_cn_other = record["dv_cn_other"]
        if dv_cn_other not in groups:
            groups[dv_cn_other] = {
                "dv_cn_other": dv_cn_other,
                "descriptions": []
            }
        groups[dv_cn_other]["descriptions"].append({
            "description": record["description"],
            "id": record["id"],
            "rnk": record["rnk"]
        })             
    conn.close()

    result = []
    for group in groups.values():
        result.append(group)

    return result

@tbrkRtr.post("/tbrk/add")
def add_tbrk(request: TbrkRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    added = []
    existing = []
    
    if request.league_id > 0:
        sql = "select COALESCE(max(rnk), 0)+1 rnk from tbrk " \
            "where league_id = %s and dv_cn_other = %s"
        cursor.execute(sql, (request.league_id, request.dvCnOth))
    else:
        sql = "select COALESCE(max(rnk), 0)+1 rnk from tbrk " \
            "where tourn_id = %s and dv_cn_other = %s"
        cursor.execute(sql, (request.tourn_id, request.dvCnOth))
    result = cursor.fetchone() 
    rnk = 0
    if result is None or 'rnk' not in result:
        rnk = 1
    else:
        rnk = result['rnk']

    for indx in range(len(request.descriptions)):
        if request.league_id > 0:
            sql = "select * from tbrk " \
                "where league_id = %s and dv_cn_other = %s and description = %s"
            cursor.execute(sql, (request.league_id,request.dvCnOth,request.descriptions[indx],))
        else:
            sql = "select * from tbrk " \
                "where tourn_id = %s and dv_cn_other = %s and description = %s"
            cursor.execute(sql, (request.tourn_id,request.dvCnOth,request.descriptions[indx],))
        records=cursor.fetchall()
        if len(records) > 0:
            existing.append(request.descriptions[indx])
        else:
            sql = "insert into tbrk (tourn_id, league_id, dv_cn_other, description, rnk) "
            if request.league_id > 0:
                cursor.execute(sql + "values (null, %s, %s, %s, %s)",(request.league_id,request.dvCnOth,request.descriptions[indx],rnk,))
            else:
                cursor.execute(sql + "values (%s, null, %s, %s, %s)",(request.tourn_id,request.dvCnOth,request.descriptions[indx],rnk,))
            conn.commit()
            added.append(request.descriptions[indx])
        rnk = rnk + 1
    conn.close()
    return {"Added": added, "Already existing": existing}

@tbrkRtr.post("/tbrk/change")
def change_tbrk(request: TbrkRequest):
    if len(request.descriptions)<1:
        return "Nothing in descriptions"
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    if request.league_id > 0:
        sql = "select * from tbrk " \
            "where league_id = %s and dv_cn_other = %s and description = %s"
        cursor.execute(sql, (request.league_id,request.dvCnOth,request.descriptions[0],))
    else:
        sql = "select * from tbrk " \
            "where tourn_id = %s and dv_cn_other = %s and description = %s"
        cursor.execute(sql, (request.tourn_id,request.dvCnOth,request.descriptions[0],))
    records=cursor.fetchall()
    if len(records) > 0:
        conn.close()
        return "Already exists: "+request.descriptions[0]
    
    cursor.execute("update tbrk set description = %s where id = %s", (request.descriptions[0],request.id,))
    conn.commit()
    conn.close()
    return "Changed id"+ str(request.id) + ": "+ request.descriptions[0]

@tbrkRtr.post("/tbrk/delete")
def delete_tbrk(request: TbrkRequest):
    if len(request.descriptions)<1:
        return "Nothing in descriptions"
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    
    #Keep commented out until plyr_id has other tables 
    #records={}
    #sql = "select 1 from plyrtm where divi_id = %s"
    #cursor.execute(sql, (request.id,))
    #records=cursor.fetchall()  
    
    #if len(records) > 0:
    #    return "Can not cascade delete "+request.table+": "+request.descriptions[0]
    
    if request.league_id > 0:
        sql = "select league_id parent_id, rnk from tbrk where id = %s"
    else:
        sql = "select tourn_id parent_id, rnk from tbrk where id = %s"
    cursor.execute(sql, (request.id,))
    record = cursor.fetchone() 
    leagOrTourId = 0
    rnk = 0
    if record:
        leagOrTourId = record['parent_id']
        rnk = record['rnk']
    cursor.execute("delete from tbrk where id = %s", (request.id,))
    conn.commit()
    if request.league_id > 0:
        cursor.execute("update tbrk set rnk = rnk- 1 where league_id = %s and dv_cn_other = %s and rnk > %s and id > 0", (leagOrTourId,request.dvCnOth,rnk,))
    else:
        cursor.execute("update tbrk set rnk = rnk- 1 where tourn_id = %s and dv_cn_other = %s and rnk > %s and id > 0", (leagOrTourId,request.dvCnOth,rnk,))
    conn.commit()
    conn.close()
    return "Deleted:"+request.descriptions[0]+" -- "+str(leagOrTourId)

@tbrkRtr.post("/tbrk/rnk")
def rnk_tbrk(request: TbrkRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    
    if request.league_id > 0:
        sql = "select league_id parent_id, rnk from tbrk where id = %s"
    else:
        sql = "select tourn_id parent_id, rnk from tbrk where id = %s"
    cursor.execute(sql, (request.id,))
    record=cursor.fetchone() 
    leagOrTourId = 0
    rnk = 0 
    if record:
        leagOrTourId = record['parent_id']
        rnk = record["rnk"]
    else:
        return print("No record found for ID:", request.id)

    if ((rnk==1) and request.up):
        conn.close()
        return print("Can not move rnk up:", request.id) 
    
    if request.league_id > 0:
        sql = "select max(rnk) rnk from tbrk where league_id = %s and dv_cn_other = %s "
        cursor.execute(sql, (request.league_id,request.dvCnOth))
    else:
        sql = "select max(rnk) rnk from tbrk where tourn_id = %s and dv_cn_other = %s "
        cursor.execute(sql, (request.tourn_id,request.dvCnOth))
    record=cursor.fetchone()
    maxRnk=0
    if record:
        maxRnk = record["rnk"]
    elif request.league_id > 0:
        conn.close()
        return print("No record found for league_id and dv_cn_oth:", request.league_id, request.dvCnOth)
    else: 
        conn.close()
        return print("No record found for tourn_id and dv_cn_oth:", request.tourn_id, request.dvCnOth)
    if ((rnk>=maxRnk) and not request.up):
        conn.close()
        return print("Can not move rnk down:", request.id, rnk) 
   
    if request.up:
        if request.league_id > 0:
            sql = "update tbrk set rnk = %s where league_id = %s and rnk = %s and dv_cn_other = %s"
        else:
            sql = "update tbrk set rnk = %s where tourn_id = %s and rnk = %s and dv_cn_other = %s"
        cursor.execute(sql, (rnk, leagOrTourId, rnk-1,request.dvCnOth,))
        sql = "update tbrk set rnk = %s where id = %s"
        cursor.execute(sql, (rnk-1, request.id))
    else:
        if request.league_id > 0:
            sql = "update tbrk set rnk = %s where league_id = %s and rnk = %s and dv_cn_other = %s"
        else:
            sql = "update tbrk set rnk = %s where tourn_id = %s and rnk = %s and dv_cn_other = %s"
        cursor.execute(sql, (rnk, leagOrTourId, rnk+1,request.dvCnOth,))
        sql = "update tbrk set rnk = %s where id = %s"
        cursor.execute(sql, (rnk+1, request.id,))
    conn.commit()
    conn.close()
    return "Changed plyrs rnk"

@tbrkRtr.post("/tbrk/change-dco")
def change_dco_tbrk(request: TbrkRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    if request.league_id > 0:
        sql = "select * from tbrk " \
            "where league_id = %s and dv_cn_other = %s"
        cursor.execute(sql, (request.league_id,request.descriptions[0],))
    else:
        sql = "select * from tbrk " \
            "where tourn_id = %s and dv_cn_other = %s"
        cursor.execute(sql, (request.tourn_id,request.descriptions[0],))
    records=cursor.fetchall()
    if len(records) > 0:
        return "Already exists: "+request.descriptions[0]
    
    cursor.execute("update tbrk set dv_cn_other = %s where dv_cn_other = %s", (request.descriptions[0],request.dvCnOth,))
    conn.commit()
    conn.close()
    return "Changed id"+ str(request.id) + ": "+ request.descriptions[0]

@tbrkRtr.post("/tbrk/league_to_trn")
def league_to_trn(request: LTRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    sql = "delete from tbrk where tourn_id = %s"
    cursor.execute(sql, (request.tourn_id,))
    sql = "update tbrk set tourn_id = %s, league_id = null where league_id = %s"
    cursor.execute(sql, (request.tourn_id,request.league_id,))
    conn.commit()
    conn.close()
    return "Copied trn to league "+ str(request.tourn_id) + ": "+ str(request.league_id)

@tbrkRtr.post("/tbrk/trn_to_league")
def trn_to_league(request: LTRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    sql = "delete from tbrk where league_id = %s"
    cursor.execute(sql, (request.league_id,))
    sql = "update tbrk set league_id = %s, tourn_id = null where tourn_id = %s"
    cursor.execute(sql, (request.league_id,request.tourn_id,))
    conn.commit()
    conn.close()
    return "Copied trn to league "+ str(request.league_id) + ": " + str(request.tourn_id)
