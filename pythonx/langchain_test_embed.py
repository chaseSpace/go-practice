from langchain_community.document_loaders import TextLoader
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_chroma import Chroma
from langchain_community.embeddings import DashScopeEmbeddings  # tongyi
from zzz.gpt import TONGYI_APIKEY
import os
os.environ["DASHSCOPE_API_KEY"] = TONGYI_APIKEY

loader = TextLoader("./introduction.txt")
docs = loader.load()

text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=200)
splits = text_splitter.split_documents(docs)
vectorstore = Chroma(
    collection_name="ai_learning",
    embedding_function=DashScopeEmbeddings(model="text-embedding-v2"),
    persist_directory="vectordb"
)
vectorstore.add_documents(splits)
