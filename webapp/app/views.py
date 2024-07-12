from flask import render_template, request, redirect, jsonify, send_file
from app import app, Log
import os
import shutil
from utils.utils import fileInfo

# Route for favicon
@app.route('/favicon.ico')
def favicon():
    # Return the favicon
    return send_file('./static/MultScan.ico', mimetype='image/vnd.microsoft.icon')

# Route for the homepage
@app.route('/', methods=['GET'])
def index():
    return render_template('index.html')

# Get enabled plugins
@app.route('/api/v1/machines', methods=['GET'])
def machines():
    return app.config['config']['machines']

# Upload sample for a given id
@app.route('/api/v1/sample/upload/<id>', methods=['POST'])
def upload(id):
    Log.info(f"[\033[34m{id}\033[0m] Saving sample")

    # Check if 'sample' file is in the request
    if 'sample' not in request.files:
        Log.error("No sample file part in the request")
        return jsonify({"error": "No sample file part"}), 400

    # Get the file from the request
    file = request.files['sample']

    # Make directory
    Log.info(f"[\033[34m{id}\033[0m] Generating temp dir")
    os.makedirs(f'./uploads/{id}')
    
    # Save the file to the uploads folder as sample
    file.save(f'./uploads/{id}/sample')

    # Return the success message
    return jsonify({"message": "Sample uploaded successfully"})

# Retrieve fileinfo for a given id sample
@app.route('/api/v1/sample/fileinfo/<id>', methods=['GET'])
def getFileInfo(id):
    return fileInfo(f'./uploads/{id}/sample')

# Request a scan status for a given id sample
@app.route('/api/v1/sample/scan/<id>', methods=['GET'])
def getStatus(id):
    statusFilePath = f'./uploads/{id}/sample/status'

    # Check if status file exists
    if not os.path.exists(statusFilePath):
        # Create status file
        Log.info(f"[\033[34m{id}\033[0m] Creating status file")

        status = {"status": {}}

        for machine in app.config['config']['machines']:
            result = {
                "badBytes": '',
                "result": 'Scanning'
            }

            status["status"][machine["name"]] = result

        # Write json to file
        with open(statusFilePath, 'w') as statusFile:
            json.dump(status, statusFile, indent=4)

        # Request scan to agents
        Log.info(f"[\033[34m{id}\033[0m] Requesting sample scan")
        # TODO

    else:
        # Check status on agents
        Log.info(f"[\033[34m{id}\033[0m] Checking sample scan status")

        # Load existing status
        with open(statusFilePath, 'r') as statusFile:
            status = json.load(statusFile)
        
        # Update each machine status
        for machine in app.config['config']['machines']:
            machine_name = machine["name"]
            
            # TODO: Add code to query machine for current status
            # Example mock update
            if machine_name == "machine1":
                status["status"][machine_name]["badBytes"] = "TWFsaWNpb3VzIGNvbnRlbnQgZm91bmQgYXQgb2Zmc2V0OiAwMDA0OGUzZAowMDAwMDAwMCAgNjUgNzQgNWYgNjEgNjQgNjQgNjkgNzQgIDY5IDZmIDZlIDYxIDZjIDVmIDc0IDY5ICB8ZXRfYWRkaXRpb25hbF90aXwKMDAwMDAwMTAgIDYzIDZiIDY1IDc0IDczIDAwIDY3IDY1ICA3NCA1ZiA3NCA2OSA2MyA2YiA2NSA3NCAgfGNrZXRzLmdldF90aWNrZXR8CjAwMDAwMDIwICA3MyAwMCA3MyA2NSA3NCA1ZiA3NCA2OSAgNjMgNmIgNjUgNzQgNzMgMDAgNTMgNzkgIHxzLnNldF90aWNrZXRzLlN5fAowMDAwMDAzMCAgNzMgNzQgNjUgNmQgMmUgNGUgNjUgNzQgIDJlIDUzIDZmIDYzIDZiIDY1IDc0IDczICB8c3RlbS5OZXQuU29ja2V0c3w="
                status["status"][machine_name]["result"] = "Detected"
            else:
                status["status"][machine_name]["result"] = "Undetected"

        # Update status file
        with open(statusFilePath, 'w') as statusFile:
            json.dump(status, statusFile, indent=4)

    # Load status file
    with open(statusFilePath, 'r') as statusFile:
        status = json.load(statusFile)

    # Return status
    return jsonify(status)

@app.route('/api/v1/sample/delete/<id>', methods=['GET'])
def deleteSample(id):
    # Remove directory
    Log.info(f"[\033[34m{id}\033[0m] Removing temp dir")
    shutil.rmtree(f'./uploads/{id}')

    # Return the success message
    return jsonify({"message": "Sample deleted successfully"})

# Route for errors
@app.errorhandler(404)
def page_not_found(e):
    return render_template('404.html'), 404