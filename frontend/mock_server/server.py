from flask import Flask, request, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app)

@app.route('/api/mock-board')
def mock_board():
    return jsonify({"fenString": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"})

@app.route('/api/update-board')
def update_board():
    return jsonify({"arbitraryValue": True})

@app.route('/api/game', methods=['POST'])
def game():
  req = request.get_json()
  action = req['action']
  if action == 'state':
    return jsonify({"state": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"})
  else:
    return jsonify({"state": "rnbqkbnr/1ppppppp/p7/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"})

if __name__ == '__main__':
    app.run(debug=True)

