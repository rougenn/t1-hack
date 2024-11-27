import os
import json
from loguru import logger
from langchain_community.vectorstores import FAISS
from langchain.docstore.document import Document
from langchain_huggingface import HuggingFaceEmbeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_ollama import ChatOllama


class RAGModel:
    def __init__(self, model_name, database_path=None, base_dir="models", llama_version="3.2-3B"):
        """
        Инициализация RAG модели с выбором версии LLaMA.
        :param model_name: Имя модели.
        :param database_path: Путь к директории с текстовыми файлами.
        :param base_dir: Базовая директория для хранения моделей.
        :param llama_version: Версия LLaMA ('3.2-3B' или '3.2-1B').
        """
        self.model_name = model_name
        self.database_path = database_path
        self.model_dir = os.path.join(base_dir, model_name)
        self.db_dir = os.path.join(self.model_dir, "db")
        self.log_dir = os.path.join(self.model_dir, "log")
        self.config_file = os.path.join(self.model_dir, "config.json")
        self.chunk_size = 512
        self.chunk_overlap = 50
        self.embeddings_model_id = "intfloat/multilingual-e5-large"
        self.llama_version = llama_version

        # Настройка ID модели в зависимости от версии LLaMA
        if self.llama_version == "3.2-3B":
            self.llm_model_id = "llama3.2:3b-instruct-fp16"
        elif self.llama_version == "3.2-1B":
            self.llm_model_id = "llama3.2:1b-instruct-fp16"
        else:
            raise ValueError(f"Неверная версия LLaMA: {llama_version}")

        # Проверяем, существует ли модель
        if not os.path.exists(self.model_dir):
            if database_path is None:
                raise ValueError(f"Модель {model_name} не существует. Укажите путь к базе данных для её создания.")
            os.makedirs(self.db_dir, exist_ok=True)
            os.makedirs(self.log_dir, exist_ok=True)
            self._save_config()
        else:
            self._load_config()

        logger.add(
            os.path.join(self.log_dir, f"{model_name}.log"),
            format="{time} {level} {message}",
            level="DEBUG",
            rotation="100 KB",
            compression="zip"
        )
        logger.info(f"Инициализирована модель {model_name} с LLaMA {self.llama_version} в {self.model_dir}")

        self.embeddings = HuggingFaceEmbeddings(
            model_name=self.embeddings_model_id,
            model_kwargs={"device": "cpu"}
        )

        # Проверяем или создаем базу знаний
        self.db = self._get_index_db()

    def _save_config(self):
        config = {
            "model_name": self.model_name,
            "database_path": self.database_path,
            "chunk_size": self.chunk_size,
            "chunk_overlap": self.chunk_overlap,
            "embeddings_model_id": self.embeddings_model_id,
            "llama_version": self.llama_version,
            "llm_model_id": self.llm_model_id
        }
        with open(self.config_file, "w", encoding="utf-8") as f:
            json.dump(config, f, indent=4)

    def _load_config(self):
        if not os.path.exists(self.config_file):
            raise FileNotFoundError(f"Конфигурация для модели {self.model_name} не найдена.")
        with open(self.config_file, "r", encoding="utf-8") as f:
            config = json.load(f)
        self.database_path = config.get("database_path")
        self.chunk_size = config.get("chunk_size", 512)
        self.chunk_overlap = config.get("chunk_overlap", 50)
        self.embeddings_model_id = config.get("embeddings_model_id", "intfloat/multilingual-e5-large")
        self.llama_version = config.get("llama_version", "3.2-3B")
        self.llm_model_id = config.get("llm_model_id", "llama3.2:3b-instruct-fp16")

    def _load_txt_documents(self):
        documents = []
        for root, dirs, files in os.walk(self.database_path):
            for file in files:
                if file.endswith(".txt"):
                    file_path = os.path.join(root, file)
                    with open(file_path, 'r', encoding='utf-8') as f:
                        content = f.read().strip()
                        if content:
                            documents.append(Document(page_content=content, metadata={"source": file}))
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
        return "\n".join([doc.page_content for doc in docs])

    def ask_question(self, topic):
        message_content = self.get_relevant_chunks(topic, num_chunks=3)
        llm = ChatOllama(model=self.llm_model_id, temperature=0)

        prompt = f"""
        Контекст:
        {message_content}

        Вопрос: {topic}
        Ответ:
        """
        return llm.invoke([{"role": "user", "content": prompt}]).content
