from flask import render_template, request, redirect, jsonify, send_file
from app import app
import time
from utils.utils import getFileInfo, getConfig, getCleanResults, requestStatus

# Gets configs and inits semaphore
with app.app_context():
    app.config['uploading'] = False
    app.config['scan_config'] = getConfig()
    app.config['scan_results'] = getCleanResults(app.config['scan_config'])

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
    
    # If scan status is "scanning"
    if app.config['scan_results']['status'] == "scanning":
        # Request the scan status
        app.config['scan_results'] = requestStatus(app.config['scan_results'])

    # If scan status is "done"
    else:
        # Clear the scan results
        app.config['scan_results'] = getCleanResults(app.config['scan_config'])

    # Return the scan status
    return jsonify(app.config['scan_results'])

# Route for errors
@app.errorhandler(404)
def page_not_found(e):
    return render_template('404.html'), 404