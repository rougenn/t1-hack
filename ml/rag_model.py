import os
import json
from loguru import logger
from langchain_community.vectorstores import FAISS
from langchain.docstore.document import Document
from langchain_huggingface import HuggingFaceEmbeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_ollama import ChatOllama

class RAGModel:
    def __init__(self, model_name, database_path=None, base_dir="models", llama_version="3.2-3B", chunk_size=512, chunk_overlap=50, embeddings_model_id="sentence-transformers/all-MiniLM-L6-v2"):
        self.model_name = model_name
        self.database_path = database_path
        self.model_dir = os.path.join(base_dir, model_name)
        self.db_dir = os.path.join(self.model_dir, "db")
        self.log_dir = os.path.join(self.model_dir, "log")
        self.config_file = os.path.join(self.model_dir, "config.json")
        self.chunk_size = chunk_size
        self.chunk_overlap = chunk_overlap
        self.embeddings_model_id = embeddings_model_id
        self.llama_version = llama_version

        if self.llama_version == "3.2-3B":
            self.llm_model_id = "llama3.2:3b-instruct-fp16"
        elif self.llama_version == "3.2-1B":
            self.llm_model_id = "llama3.2:1b-instruct-fp16"
        else:
            raise ValueError(f"Неверная версия LLaMA: {llama_version}")

        # Настройка логирования и эмбеддингов
        logger.add(
            os.path.join(self.log_dir, f"{model_name}.log"),
            format="{time} {level} {message}",
            level="DEBUG",
            rotation="100 KB",
            compression="zip"
        )

        self.embeddings = HuggingFaceEmbeddings(
            model_name=self.embeddings_model_id,
            model_kwargs={"device": "cpu"}
        )

        self.db = self._get_index_db()

    def _load_txt_documents(self):
        documents = []
        for root, dirs, files in os.walk(self.database_path):
            for file in files:
                if file.endswith(".txt"):
                    file_path = os.path.join(root, file)
                    with open(file_path, 'r', encoding='utf-8') as f:
                        content = f.read().strip()
                        if content:
                            # Добавляем документ с метаданными (имя файла без расширения)
                            documents.append(Document(page_content=content, metadata={"source": os.path.splitext(file)[0]}))
        if not documents:
            raise ValueError(f"Папка '{self.database_path}' пуста или файлы не содержат текста.")
        return documents

    def _split_documents_into_chunks(self, documents):
        text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=self.chunk_size,
            chunk_overlap=self.chunk_overlap
        )
        return text_splitter.split_documents(documents)

    def _get_index_db(self):
        db_file = os.path.join(self.db_dir, "index.faiss")
        if os.path.exists(db_file):
            return FAISS.load_local(self.db_dir, self.embeddings, allow_dangerous_deserialization=True)
        else:
            documents = self._load_txt_documents()
            chunks = self._split_documents_into_chunks(documents)
            db = FAISS.from_documents(chunks, self.embeddings)
            db.save_local(self.db_dir)
            return db

    def get_relevant_chunks(self, topic, num_chunks=3):
        docs = self.db.similarity_search(topic, k=num_chunks)
        return "\n".join([doc.page_content for doc in docs]), [doc.metadata["source"] for doc in docs]

    def ask_question(self, topic):
        message_content, sources = self.get_relevant_chunks(topic, num_chunks=3)
        llm = ChatOllama(model=self.llm_model_id, temperature=0)

        prompt = f"""
        Контекст:
        {message_content}

        Вопрос: {topic}
        Ответ:
        """
        answer = llm.invoke([{"role": "user", "content": prompt}]).content
        # Возвращаем только ответ и первый источник
        return answer, sources[0]  # Возвращаем только первый файл, так как источник один

