from flask import render_template, request, redirect, jsonify, send_file
from app import app, Log
import os
import shutil
from utils.utils import getFileInfo

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
def plugins():
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
    return jsonify({"message": "File uploaded successfully"})

# Retrieve fileinfo for a given id sample
@app.route('/api/v1/sample/fileinfo/<id>', methods=['GET'])
def getResult(id):

    # Retrieve fileinfo
    fileInfo = getFileInfo(f'./uploads/{id}/sample')

    # Remove directory
    Log.info(f"[\033[34m{id}\033[0m] Removing temp dir")
    shutil.rmtree(f'./uploads/{id}')

    return fileInfo

# Route for errors
@app.errorhandler(404)
def page_not_found(e):
    return render_template('404.html'), 404