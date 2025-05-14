from typing import Optional

from pydantic import BaseModel


class UserBase(BaseModel):
    username: str
    email: Optional[str] = None


class UserCreate(UserBase):
    password: str


class User(UserBase):
    id: str

    class Config:
        from_attributes = True


class ItemBase(BaseModel):
    name: str
    description: Optional[str] = None
    price: float


class ItemCreate(ItemBase):
    pass


class Item(ItemBase):
    id: str

    class Config:
        from_attributes = True
