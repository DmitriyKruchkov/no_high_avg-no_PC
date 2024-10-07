from fastapi import APIRouter, Depends, Request
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from db_model import Client
from core import templates
from utils import get_db

router = APIRouter()


@router.get("/register")
async def register(name: str = "randomPC", db: AsyncSession = Depends(get_db)):
    new_client = Client(status=True, name=name)
    db.add(new_client)
    await db.commit()

    result = await db.execute(select(Client).order_by(Client.id.desc()))
    client = result.scalar()

    return {"your-id": client.id}


@router.get("/status/{client_id}")
async def get_status(client_id: int, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Client).where(Client.id == client_id))
    client = result.scalar_one_or_none()
    if client:
        return {"status": client.status}
    else:
        return {"error": "Client not found"}


@router.get("/")
async def list_clients(request: Request, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Client))
    clients = result.scalars().all()
    return templates.TemplateResponse("index.html", {"request": request, "clients": clients})


@router.get("/stop/{client_id}")
async def stopper(request: Request, client_id: int, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Client).where(Client.id == client_id))
    client = result.scalar_one_or_none()
    client.status = False
    await db.commit()
    return templates.TemplateResponse("end.html", {"request": request})
