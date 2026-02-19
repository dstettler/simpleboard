from flask import Flask, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app)

@app.route('/api/mock-board')
def mock_board():
    return jsonify({"fenString": "mocked"})

if __name__ == '__main__':
    app.run(debug=True)

