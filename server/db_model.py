from sqlalchemy import Column, Integer, String, Boolean
from database import Base


class Client(Base):
    __tablename__ = 'clients'
    id = Column(Integer, primary_key=True)
    name = Column(String)
    status = Column(Boolean, default=True)
