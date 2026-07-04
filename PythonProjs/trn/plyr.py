from fastapi import APIRouter
from pydantic import BaseModel
from db import get_db_connection  # Import the database function
from typing import List
import logging

# Create a router
plyrRtr = APIRouter()

# Define a Pydantic model for request payload
class PlyrRequest(BaseModel):
    parent_id: int
    id: int
    descriptions: List[str]  
    descriptions2: List[str]
    sds: int
    up: bool

@plyrRtr.post("/plyr/get")
def get_tcdps(request: PlyrRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    sql = "select * from plyrtm where divi_id = %s order by sds"
    cursor.execute(sql, (request.parent_id,))
    records=cursor.fetchall()
    conn.close()
    return records

@plyrRtr.post("/plyr/add")
def add_plyr(request: PlyrRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    added = []
    existing = []
    
    sql = "select COALESCE(max(sds), 0)+1 sds from plyrtm where divi_id = %s"
    cursor.execute(sql, (request.parent_id,))
    result = cursor.fetchone() 
    if result is None or 'sds' not in result:
        sds = 1
    else:
        sds = result['sds']

    for indx in range(len(request.descriptions)):
        if (request.descriptions2[indx] == ""):
            sql = "select * from plyrtm where divi_id = %s and description = %s"
            cursor.execute(sql, (request.parent_id,request.descriptions[indx],))
        else:
            sql = "select * from plyrtm where divi_id = %s and description = %s and description2 = %s"
            cursor.execute(sql, (request.parent_id,request.descriptions[indx],request.descriptions2[indx],))  
        records=cursor.fetchall()
        if len(records) > 0:
            if (request.descriptions2[indx] == ""):
                existing.append(request.descriptions[indx])
            else:
                existing.append(request.descriptions[indx]+" "+request.descriptions2[indx])
        else:
            if (request.descriptions2[indx] == ""):
                cursor.execute("insert into plyrtm (divi_id, description, sds) values (%s, %s, %s)",(request.parent_id,request.descriptions[indx],sds,))
            else:
                cursor.execute("insert into plyrtm (divi_id, description, description2, sds) values (%s, %s, %s, %s)",(request.parent_id,request.descriptions[indx],request.descriptions2[indx],sds,))
            conn.commit()
            if (request.descriptions2[indx] == ""):
                added.append(request.descriptions[indx])
            else:
                added.append(request.descriptions[indx]+" "+request.descriptions2[indx])
        sds = sds + 1
    conn.close()
    return {"Added": added, "Already existing": existing}

@plyrRtr.post("/plyr/change")
def change_plyr(request: PlyrRequest):
    if len(request.descriptions)<1:
        return "Nothing in descriptions"
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    if request.descriptions2[0] == "":
        sql = "select * from plyrtm where divi_id = %s and description = %s and description2 is null"
        cursor.execute(sql, (request.parent_id,request.descriptions[0],))
    else:
        sql = "select * from plyrtm where divi_id = %s and description = %s and description2 = %s"
        cursor.execute(sql, (request.parent_id,request.descriptions[0],request.descriptions2[0],))
    records=cursor.fetchall()
    if len(records) > 0:
        return "Already exists: "+request.descriptions[0]
    
    if request.descriptions2[0] == "":
        cursor.execute("update plyrtm set description = %s, description2 = null where id = %s", (request.descriptions[0],request.id,))
    else:
        cursor.execute("update plyrtm set description = %s, description2 = %s where id = %s", (request.descriptions[0],request.descriptions2[0],request.id,))
    conn.commit()
    conn.close()
    return "Changed id"+ str(request.id) + ": "+ request.descriptions[0]+ " " +request.descriptions2[0]

@plyrRtr.post("/plyr/delete")
def delete_plyr(request: PlyrRequest):
    if len(request.descriptions)<1:
        return "Nothing in descriptions"
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    
    #Keep commented out until plyr_id has other tables 
    #records={}
    sql = "select 1 from gms g " \
        "inner join gms_plyrtm g2 on g2.gms_id = g.id where g2.plyrtm_id = %s"
    cursor.execute(sql, (request.id,))
    record = cursor.fetchone()   
    if record:
        conn.close()
        return "Tourn started, so cannot delete "+request.descriptions[0]
    
    sql = "select divi_id, sds from plyrtm where id = %s"
    cursor.execute(sql, (request.id,))
    record = cursor.fetchone() 
    diviId = 0
    sds = 0
    if record:
        diviId = record['divi_id']
        sds = record['sds']
    cursor.execute("delete from plyrtm where id = %s", (request.id,))
    conn.commit()
    cursor.execute("update plyrtm set sds = sds - 1 where divi_id = %s and sds > %s and id > 0", (diviId,sds,))
    conn.commit()
    conn.close()
    return "Deleted:"+request.descriptions[0]+" -- "+request.descriptions2[0]+" -- "+str(diviId)

@plyrRtr.post("/plyr/sds")
def sds_plyr(request: PlyrRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    
    sql = "select divi_id, sds from plyrtm where id = %s"
    cursor.execute(sql, (request.id,))
    record=cursor.fetchone() 
    diviId = 0
    sds = 0 
    if record:
        diviId = record["divi_id"]
        sds = record["sds"]
    else:
        return print("No record found for ID:", request.id)

    if ((sds==1) and request.up):
        return print("Can not move sds up:", request.id) 
    
    sql = "select max(sds) sds from plyrtm where divi_id = %s"
    cursor.execute(sql, (diviId,))
    record=cursor.fetchone()
    maxSds=0
    if record:
        maxSds = record["sds"]
    else:
        return print("No record found for id and divi:", request.id, diviId)
    
    if ((sds>=maxSds) and not request.up):
        return print("Can not move sds down:", request.id, sds) 
    
    if request.up:
        sql = "update plyrtm set sds = %s where divi_id = %s and sds = %s"
        cursor.execute(sql, (sds, diviId, sds-1,))
        sql = "update plyrtm set sds = %s where id = %s"
        cursor.execute(sql, (sds-1, request.id,))
    else:
        sql = "update plyrtm set sds = %s where divi_id = %s and sds = %s"
        cursor.execute(sql, (sds, diviId, sds+1,))
        sql = "update plyrtm set sds = %s where id = %s"
        cursor.execute(sql, (sds+1, request.id,))
    conn.commit()
    conn.close()
    return "Changed plyrs sds"

