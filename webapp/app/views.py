from flask import render_template, request, redirect, jsonify, send_file
from app import app

# favicon.ico route
@app.route('/favicon.ico')
def favicon():
    # Return the favicon
    return send_file('./static/MultScan.ico', mimetype='image/vnd.microsoft.icon')

# Define the route that will serve the loader configuration
@app.route('/', methods=['GET'])
def index():
    return render_template('index.html')

# Define the route that scan the files
@app.route('/', methods=['POST'])
def scan_file():
    # Check if the payload file was uploaded
    if 'payload' not in request.files:
        return redirect(request.url)
    payload = request.files['payload']

    # Scan the payload against AVs
    scan(payload)

    # redirect to index
    return redirect(request.url)

# Route for scan status
@app.route('/api/v1.0/scan/status', methods=['GET'])
def status():
    result = {
        "status": "done",
        "data": {
            "file_name": "test.exe",
            "file_size": "14.5MB",
            "digests": [
                "MD5:6a46ba7a9cd4016294e6a713193c2642",
                "SHA-1:12676b985e9d3a422252364195576f5f97b17cc2",
                "SHA-256:78d348f7cefda75dd582a0412b408be8cedf200670e92de89ec442a93d0a1c46"
            ],
            "scan_results": {
                "Avast": {
                    "badBytes": "",
                    "result": "Undetected"
                },
                "Mcafee": {
                    "badBytes": "TWFsaWNpb3VzIGNvbnRlbnQgZm91bmQgYXQgb2Zmc2V0OiAwMDA0OGUzZAowMDAwMDAwMCAgNjUgNzQgNWYgNjEgNjQgNjQgNjkgNzQgIDY5IDZmIDZlIDYxIDZjIDVmIDc0IDY5ICB8ZXRfYWRkaXRpb25hbF90aXwKMDAwMDAwMTAgIDYzIDZiIDY1IDc0IDczIDAwIDY3IDY1ICA3NCA1ZiA3NCA2OSA2MyA2YiA2NSA3NCAgfGNrZXRzLmdldF90aWNrZXR8CjAwMDAwMDIwICA3MyAwMCA3MyA2NSA3NCA1ZiA3NCA2OSAgNjMgNmIgNjUgNzQgNzMgMDAgNTMgNzkgIHxzLnNldF90aWNrZXRzLlN5fAowMDAwMDAzMCAgNzMgNzQgNjUgNmQgMmUgNGUgNjUgNzQgIDJlIDUzIDZmIDYzIDZiIDY1IDc0IDczICB8c3RlbS5OZXQuU29ja2V0c3w=",
                    "result": "Detected"
                },
                "Dev01": {
                    "badBytes": "",
                    "result": "Detected"
                },
            }
        }
    }

    return jsonify(result)

# Route for errors
@app.errorhandler(404)
def page_not_found(e):
    return render_template('404.html'), 404