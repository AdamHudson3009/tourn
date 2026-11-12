from fastapi import APIRouter
from pydantic import BaseModel
from db import get_db_connection

# Define a Pydantic model for request payload
class BookPageRequest(BaseModel):
    book_id: int
    book: str
    page_id: int
    page: str

@leagueRtr.get("/trading/bookpage/get")
def get_book_pages():
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True) #json format
    sql = "select b.id, b.name, p.id, p.name
        "from book b inner join page p on p.book_id = b.id"
        "order by book, page"
    cursor.execute(sql)
    records=cursor.fetchall()
    conn.close()
    return records

@leagueRtr.post("/trading/bookpage/add")
def add_book_pages(request: LeagueRequest):
    conn = get_db_connection()
    cursor=conn.cursor(dictionary=True)
    sql = "select * from book where name = %s"
    cursor.execute(sql, (request.book_name,))
    records=cursor.fetchall()
    if len(records) > 0:
        return "Book Name Already Exists: "+request.book_name       
    book_id = cursor.lastrowid
    cursor.execute("insert into book (name) values (%s)",(request.book_name,))
    conn.commit()
    
    #Now entry the page name for that book
    #Since the book is new, so will be the page
    #the page name does not have to be unique

    cursor.execute("insert into page (name, book_id) values (%s, %s)",(request.page_name, request.book_id,))
    

    conn.close()
    return "Added, book:"+request.book_name+", page:"+request.page_name    

