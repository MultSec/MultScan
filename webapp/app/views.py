from flask import render_template, request, redirect, jsonify, send_file
from app import app, Log
import os
import shutil
from utils.utils import fileInfo, getSampleStatus

# Route for favicon
@app.route('/favicon.ico')
def favicon():
    # Return the favicon
    return send_file('./static/MultScan.ico', mimetype='image/vnd.microsoft.icon')

# Route for the homepage
@app.route('/', methods=['GET'])
def index():
    return render_template('index.html')

# Get present machines
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
    return jsonify(getSampleStatus(id))

# Request a sample deletion for a given id
@app.route('/api/v1/sample/delete/<id>', methods=['GET'])
def deleteSample(id):
    # Remove directory
    Log.info(f"[\033[34m{id}\033[0m] Removing temp dir")
    shutil.rmtree(f'./uploads/{id}')

    # Return the success message
    return jsonify({"message": "Sample deleted successfully"})

# Request a sample download for a given id
@app.route('/api/v1/sample/download/<id>', methods=['GET'])
def downloadSample(id):
    file_path = os.path.abspath(f'./uploads/{id}/sample')

    # Check if the file exists
    if os.path.exists(file_path):
        # Send the file to the client
        result = send_file(file_path, as_attachment=True)
    else:
        result = jsonify({"error": "File not found"}), 404

    return result

# Route for errors
@app.errorhandler(404)
def page_not_found(e):
    return render_template('404.html'), 404