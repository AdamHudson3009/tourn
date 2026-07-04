from fastapi import APIRouter
from pydantic import BaseModel
from db import get_db_connection  # Import the database function

# Create a router
leagueRtr = APIRouter()

# Define a Pydantic model for request payload
class LeagueRequest(BaseModel):
    id: int
    league: str

@leagueRtr.get("/")
def read_root():
    return {"message": "Hello, Leagues!"}

@leagueRtr.get("/league/get")
def get_leagues():
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    sql = "select * from league order by start_time desc"
    cursor.execute(sql)
    records=cursor.fetchall()
    conn.close()
    return records

@leagueRtr.post("/league/add")
def add_league(request: LeagueRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    sql = "select * from league where description = %s"
    cursor.execute(sql, (request.league,))
    records=cursor.fetchall()
    if len(records) > 0:
        return "Already exists: "+request.league       
    
    cursor.execute("insert into league (description, start_time, last_time) values (%s, Now(), Now())",(request.league,))
    conn.commit()
    conn.close()
    return "Added:"+request.league

@leagueRtr.post("/league/change")
def change_league(request: LeagueRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    sql = "select * from league where description = %s"
    cursor.execute(sql, (request.league,))
    records=cursor.fetchall()
    if len(records) > 0:
        return "Already exists: "+request.league
    
    cursor.execute("update league set description = %s, last_time = Now() where id = %s", (request.league, request.id,))
    conn.commit()
    conn.close()
    return "Changed id"+ str(request.id) + ": "+ request.league

@leagueRtr.post("/league/delete")
def delete_league(request: LeagueRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    sql = "select * from tourn where league_id = %s"
    cursor.execute(sql, (request.id,))
    records=cursor.fetchall()
    if len(records) > 0:
        return "Can not cascade delete: "+request.league
    
    cursor.execute("delete from league where id = %s", (request.id,))
    conn.commit()
    conn.close()
    return "Deleted:"+request.league