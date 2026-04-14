from flask import Flask, request, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app)

state_iter = 0
game_states = {
    0: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
    1: "rnbqkbnr/pppppppp/8/8/8/P7/1PPPPPPP/RNBQKBNR b KQkq - 0 1"
}
registrations_iter = 0
registrations = {}

@app.route('/api/mock-board')
def mock_board():
    return jsonify({"fenString": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"})

@app.route('/api/update-board')
def update_board():
    return jsonify({"arbitraryValue": True})

@app.route('/api/register', methods=['POST'])
def register():
  global registrations
  global registrations_iter

  req = request.get_json()
  req_user = req["username"]

  registrations[req_user] = {
      "username": req_user,
      "email": req["email"],
      "password": req["password"],
      "user_id": registrations_iter
      }

  registrations_iter += 1

  return jsonify(registrations[req_user]), 201


@app.route('/api/login', methods=['POST'])
def login():
  global registrations
  req = request.get_json()
  req_user = req["username"]

  if req_user in registrations.keys() and req["password"] == registrations[req_user]["password"]:
    return jsonify({
      "mesage": "login successful",
      "user": registrations[req_user]})
  else:
    return jsonify({"error": "invalid"})

@app.route('/api/game', methods=['POST'])
def game():
  global state_iter

  req = request.get_json()
  print(req)
  action = req['action']
  if action == 'state':
    return jsonify({
      "user": {
        "state": game_states[min(1, state_iter)],
        "next_moves": ["a2a3"],
        "side": "w",
        "status": "InProgress",
        "white_player_id": 0,
        "black_player_id": 1
      }})
  elif action == 'move':
    state_iter += 1
    return jsonify({
      "user": {
        "state": game_states[min(1, state_iter)],
        "next_moves": ["a3a4"],
        "status": "InProgress",
        "side": "b",
        "white_player_id": 0,
        "black_player_id": 1
      }})

if __name__ == '__main__':
    app.run(debug=True, port=8080)

