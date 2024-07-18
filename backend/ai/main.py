from fastapi import FastAPI, HTTPException
from databases import Database
import os

app = FastAPI()

DATABASE_URL = os.getenv(
    'DATABASE_URL', 'mysql://root:password@mysql:3306/mydb')

database = Database(DATABASE_URL)


@app.on_event("startup")
async def startup():
    await database.connect()


@app.on_event("shutdown")
async def shutdown():
    await database.disconnect()


@app.get("/health")
async def health_check():
    query = "SELECT 1"
    try:
        result = await database.execute(query=query)
        return {"status": "ok", "result": result}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/")
async def read_root():
    return {"message": "Hello, World!"}
