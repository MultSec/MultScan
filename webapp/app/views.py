from flask import render_template, request, redirect, jsonify, send_file
from app import app
import time
from utils.utils import getFileInfo, getConfig

# Gets configs and inits semaphore
with app.app_context():
    app.config['uploading'] = False
    app.config['scan_config'] = getConfig()

# Route for favicon
@app.route('/favicon.ico')
def favicon():
    # Return the favicon
    return send_file('./static/MultScan.ico', mimetype='image/vnd.microsoft.icon')

# Route for the homepage
@app.route('/', methods=['GET'])
def index():
    return render_template('index.html')

# Route for file upload
@app.route('/api/v1/payload/upload', methods=['POST'])
def upload():
    # Get the file from the request
    file = request.files['payload']

    # Write to semaphore config that the file is being uploaded
    app.config['uploading'] = True

    # Save the file to the uploads folder as payload
    file.save('./uploads/payload')

    # Write to semaphore config that the file has been uploaded
    app.config['uploading'] = False

    # Return the success message
    return jsonify({"message": "File uploaded successfully"})

# Route for file information
@app.route('/api/v1/payload/info', methods=['GET'])
def fileInfo():
    # Wait for the file to be uploaded
    while app.config['uploading']:
        pass
        
    return jsonify(getFileInfo())

# Route for payload download
@app.route('/api/v1/payload/download', methods=['GET'])
def download():
    return send_file('../uploads/payload')

# Route for scan status
@app.route('/api/v1/payload/scan', methods=['GET'])
def scan():
    # Wait for the file to be uploaded
    while app.config['uploading']:
        pass
    

    result = {
        "status": "done",
        "results": {
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
            }
        }
    }
    
    return jsonify(result)

# Route for errors
@app.errorhandler(404)
def page_not_found(e):
    return render_template('404.html'), 404