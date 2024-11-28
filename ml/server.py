from flask import Flask, request, jsonify
from rag_model import RAGModel

app = Flask(__name__)

@app.route('/train', methods=['POST'])
def train_model():
    data = request.get_json()
    model_name = data['model_name']
    print(model_name)
    txt_files_directory = data['txt_files_directory']

    try:
        # Создаем и обучаем модель
        print(txt_files_directory)
        model = RAGModel(model_name=model_name, database_path=txt_files_directory)
        return jsonify({"status": "OK"}), 200
    except Exception as e:
        print(e)
        return jsonify({"error": str(e)}), 500

@app.route('/ask', methods=['POST'])
def ask_question():
    data = request.get_json()
    model_name = "assistant_" + data['assistant_id']
    question = data['message']

    try:
        # Загружаем модель и получаем ответ
        model = RAGModel(model_name=model_name)
        answer = model.ask_question(question)
        
        return jsonify({"answer": answer}), 200
    except Exception as e:
        print(e)
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(debug=True)
