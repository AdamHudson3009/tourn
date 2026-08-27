from fastapi import APIRouter
from pydantic import BaseModel
from db import get_db_connection  # Import the database function
from typing import List
import logging

# Create a router
grammarRtr = APIRouter()

# Define a Pydantic model for request payload
class GrammarRequest(BaseModel):
    tourn_id: int
    letters: str
    divi1_id: int
    divi2_id: int
    id: int
    
@grammarRtr.post("/grammar/get")
def get_grammar(request: GrammarRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    sql = "select g.id, letters, " \
        "g.dv1_id, c1.description conf1, d1.description divi1, " \
        "g.dv2_id, c2.description conf2, d2.description divi2 " \
        "from grammar g inner join divi d1 on g.dv1_id = d1.id " \
        "inner join conf c1 on d1.conf_id = c1.id " \
        "inner join divi d2 on g.dv2_id = d2.id " \
        "inner join conf c2 on d2.conf_id = c2.id " \
        "where g.tourn_id = %s order by letters"
    cursor.execute(sql, (request.tourn_id,))
    records=cursor.fetchall()
    groups = {"grammar": []}
    for record in records:
        groups["grammar"].append({
            "id": record["id"],
            "letters": record["letters"],
            "dv1_id": record["dv1_id"],
            "divi1": record["divi1"],
            "conf1": record["conf1"],
            "dv2_id": record["dv2_id"],
            "divi2": record["divi2"],
            "conf2": record["conf2"]
        })
    sql = "select d.id did, d.description divi, c.description conf " \
        "from divi d inner join conf c on d.conf_id = c.id " \
        "where tourn_id = %s order by c.description, d.description"
    cursor.execute(sql, (request.tourn_id,))
    records=cursor.fetchall()
    groups["confdivi"] = []
    for record in records:
        groups["confdivi"].append({
            "did": record["did"],
            "divi": record["divi"],
            "conf": record["conf"]
        })

    conn.commit()
    conn.close()
    return groups

@grammarRtr.post("/grammar/add")
def add_grammar(request: GrammarRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    sql = "select 1 from grammar where tourn_id = %s and letters = %s"
    cursor.execute(sql, (request.tourn_id,request.letters,))
    result = cursor.fetchone() 
    if result is not None:
        conn.close()
        return "id exists: "+request.letters
    sql = "insert into grammar (tourn_id, letters, dv1_id, dv2_id) " \
        "values (%s, %s, %s, %s)"
    cursor.execute(sql, (request.tourn_id,request.letters,request.divi1_id,request.divi2_id,))
    conn.commit()
    conn.close()
    
    return "Added record"

@grammarRtr.post("/grammar/change")
def change_grammar(request: GrammarRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    sql = "select 1 from grammar where tourn_id = %s and letters = %s and id <> %s"
    cursor.execute(sql, (request.tourn_id,request.letters,request.id,))
    result = cursor.fetchone() 
    if result is not None:
        conn.close()
        return "id exists: "+request.letters
    sql = "update grammar set letters = %s, dv1_id = %s, dv2_id = %s where id = %s"
    cursor.execute(sql, (request.letters,request.divi1_id,request.divi2_id,request.id,))
    conn.commit()
    conn.close()
    
    return "Changed record"

@grammarRtr.post("/grammar/delete")
def delete_grammar(request: GrammarRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    
    sql = "delete from grammar where id = %s"
    cursor.execute(sql, (request.id,))
    conn.commit()
    conn.close()
    
    return "Deleted record"