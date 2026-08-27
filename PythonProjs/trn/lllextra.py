from fastapi import APIRouter
from pydantic import BaseModel
from db import get_db_connection  # Import the database function
from typing import List
import logging

# Create a router
lllExtraRtr = APIRouter()

class lllExtraRequest(BaseModel):
    id: int
    tourn_id: int
    ordr: int
    description: str
    up: bool

@lllExtraRtr.post("/lllextra/get")
def get_lllExtra(request: lllExtraRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    sql = "select id, ordr, description from lll_line " \
        "where tourn_id = %s order by ordr"
    cursor.execute(sql, (request.tourn_id,))
    records=cursor.fetchall()
    
    return records

@lllExtraRtr.post("/lllextra/add")
def add_lllExtra(request: lllExtraRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    
    sql = "select COALESCE(max(ordr), 0)+1 ordr from lll_line where tourn_id = %s"
    cursor.execute(sql, (request.tourn_id,))
    result = cursor.fetchone() 
    ordr = 0
    if result is None or 'ordr' not in result:
        ordr = 1
    else:
        ordr = result['ordr']

    sql = "insert into lll_line (tourn_id,description,ordr) values (%s,%s,%s)"
    cursor.execute(sql, (request.tourn_id,request.description,ordr,))
    
    conn.commit()
    return "record added: "+request.description

@lllExtraRtr.post("/lllextra/change")
def change_lllExtra(request: lllExtraRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    
    sql = "update lll_line set description = %s where id = %s"
    cursor.execute(sql, (request.description,request.id,))
    conn.commit()
    conn.close()
    
    return "Changed: "+request.description

@lllExtraRtr.post("/lllextra/delete")
def delete_lllExtra(request: lllExtraRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format

    sql = "select tourn_id, ordr from lll_line where id = %s"
    cursor.execute(sql, (request.id,))
    record = cursor.fetchone() 
    ordr = 0
    if record:
        ordr = record['ordr']
    
    sql = "delete from lll_line where id = %s"
    cursor.execute(sql, (request.id,))
    conn.commit()
    cursor.execute("update lll_line set ordr = ordr - 1 where tourn_id = %s and ordr > %s", (request.tourn_id,ordr,))
    conn.commit()
    
    return "deleted: "+request.description

@lllExtraRtr.post("/lllextra/rnk")
def rnk_xtra(request: lllExtraRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    
    sql = "select tourn_id, ordr from lll_line where id = %s"
    cursor.execute(sql, (request.id,))
    record=cursor.fetchone() 
    tournId = 0
    ordr = 0 
    if record:
        tournId = record['tourn_id']
        ordr = record["ordr"]
    else:
        return print("No record found for ID:", request.id)

    if ((ordr==1) and request.up):
        conn.close()
        return print("Can not move ordr up:", request.id) 
    
    sql = "select max(ordr) ordr from lll_line where tourn_id = %s"
    cursor.execute(sql, (request.tourn_id,))
    record=cursor.fetchone()
    maxOrdr=0
    if record:
        maxOrdr = record["ordr"]
    else: 
        conn.close()
        return print("No record found for tourn_id:", request.tourn_id,)
    if ((ordr>=maxOrdr) and not request.up):
        conn.close()
        return print("Can not move rnk down:", request.id, rnk) 
   
    if request.up:
        sql = "update lll_line set ordr = %s where tourn_id = %s and ordr = %s"
        cursor.execute(sql, (ordr, tournId, ordr-1,))
        sql = "update lll_line set ordr = %s where id = %s"
        cursor.execute(sql, (ordr-1, request.id,))
    else:
        sql = "update lll_line set ordr = %s where tourn_id = %s and ordr = %s"
        cursor.execute(sql, (ordr, tournId, ordr+1,))
        sql = "update lll_line set ordr = %s where id = %s"
        cursor.execute(sql, (ordr+1, request.id,))
    conn.commit()
    conn.close()
    return "Changed rnk"
