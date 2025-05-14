from fastapi_demo.models.schemas import Item, ItemCreate
from fastapi import APIRouter, HTTPException

router = APIRouter()

fake_items_db = {}


@router.post("/", response_model=Item, status_code=201)
async def create_item(item: ItemCreate):
    """创建新物品"""
    item_id = str(len(fake_items_db) + 1)
    new_item = Item(id=item_id, **item.dict())
    fake_items_db[item_id] = new_item
    return new_item


@router.get("/{item_id}", response_model=Item)
async def read_item(item_id: str):
    """获取物品信息"""
    if item_id not in fake_items_db:
        raise HTTPException(status_code=404, detail="Item not found")
    return fake_items_db[item_id]
