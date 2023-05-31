from flask import Flask, request, Response
import json
import func

app = Flask(__name__)

@app.route('/', methods=['POST'])
def handleTrigger():
    body = request.get_data()
    try:
        params = json.loads(body)
    except json.decoder.JSONDecodeError:
        params = ""
    result = func.main(params)
    return Response(json.dumps(result),  mimetype='application/json')

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8080, debug=True)