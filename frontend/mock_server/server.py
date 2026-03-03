from flask import Flask, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app)

@app.route('/api/mock-board')
def mock_board():
    return jsonify({"fenString": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"})

@app.route('/api/update-board')
def update_board():
    return jsonify({"arbitraryValue": True})

if __name__ == '__main__':
    app.run(debug=True)

