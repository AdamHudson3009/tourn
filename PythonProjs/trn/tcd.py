from fastapi import APIRouter
from pydantic import BaseModel
from db import get_db_connection  # Import the database function
from typing import List
import logging

# Create a router
tcdRtr = APIRouter()

# Define a Pydantic model for request payload
class TCDPRequest(BaseModel):
    table: str
    column: str
    parent_id: int
    id: int
    descriptions: List[str]  


@tcdRtr.get("/")
def read_root():
    return {"message": "Hello, Tourns!"}

@tcdRtr.post("/tcd/get")
def get_tcdps(request: TCDPRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    sql = "select * from " + request.table + " where "+ request.column + " = %s order by description"
    cursor.execute(sql, (request.parent_id,))
    records=cursor.fetchall()
    conn.close()
    return records

@tcdRtr.post("/tcd/add")
def add_tcdp(request: TCDPRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    added = []
    existing = []
    cnt = 1
    for item in request.descriptions:
        sql = "select * from " + request.table + " where "+ request.column + " = %s and description = %s"
        cursor.execute(sql, (request.parent_id,item))  
        records=cursor.fetchall()
        if len(records) > 0:
             existing.append(item)
        else:
            cursor.execute("insert into "+request.table+" ("+request.column+", description) values (%s, %s)",(request.parent_id,item,))
            conn.commit()
            added.append(item)
        cnt = cnt + 1
    conn.close()
    return {"Added": added, "Already existing": existing}

@tcdRtr.post("/tcd/change")
def change_tcdp(request: TCDPRequest):
    if len(request.descriptions)<0:
        return "Nothing in descriptions"
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    sql = "select * from " + request.table + " where "+ request.column + " = %s and description = %s"
    cursor.execute(sql, (request.parent_id,request.descriptions[0],))
    records=cursor.fetchall()
    if len(records) > 0:
        return "Already exists: "+request.descriptions[0]
    
    cursor.execute("update " + request.table + " set description = %s where id = %s", (request.descriptions[0],request.id,))
    conn.commit()
    conn.close()
    return "Changed id"+ str(request.id) + ": "+ request.descriptions[0]

@tcdRtr.post("/tcd/delete")
def delete_tcdp(request: TCDPRequest):
    if len(request.descriptions)<0:
        return "Nothing in descriptions"
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    records={}
    if request.table=="tourn":
        sql = "select 1 from conf where tourn_id = %s union " \
            "select 1 from grammar where tourn_id = %s union " \
            "select 1 from lll_line where tourn_id = %s union " \
            "select 1 from tbrk where tourn_id = %s"
        cursor.execute(sql, (request.id,request.id,request.id,request.id,))
        records=cursor.fetchall()
    elif request.table=="conf":
        sql = "select 1 from divi where conf_id = %s"
        cursor.execute(sql, (request.id,))
        records=cursor.fetchall()
    elif request.table=="divi":
        sql = "select 1 from plyrtm where divi_id = %s"
        cursor.execute(sql, (request.id,))
        records=cursor.fetchall()  
    
    if len(records) > 0:
        return "Can not cascade delete "+request.table+": "+request.descriptions[0]
    
    cursor.execute("delete from "+request.table+" where id = %s", (request.id,))
    conn.commit()
    conn.close()
    return "Deleted:"+request.descriptions[0]