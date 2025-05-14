from fastapi import APIRouter, HTTPException

from fastapi_demo.models.schemas import UserCreate, User

router = APIRouter(
    prefix="/users",
    tags=["users"],
    responses={404: {"description": "Not found [users]"}},
)

fake_users_db = {}


@router.post("/create_user", response_model=User, status_code=201)
async def create_user(user: UserCreate):
    """创建新用户"""
    user_id = str(len(fake_users_db) + 1)
    new_user = User(id=user_id, **user.dict())
    fake_users_db[user_id] = new_user
    return new_user


@router.get("/uid/{user_id}", response_model=User)
async def read_user(user_id: str):
    """获取用户信息"""
    if user_id not in fake_users_db:
        raise HTTPException(status_code=404, detail="User not found")
    return fake_users_db[user_id]
