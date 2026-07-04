import mysql.connector

# Function to get a database connection
def get_db_connection():
    return mysql.connector.connect(
        host="localhost",
        user="root",
        passwd="MySQL=1",
        database="league"
    )