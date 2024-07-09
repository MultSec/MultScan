from flask import render_template, request, redirect, jsonify, send_file
from app import app

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

# Route for errors
@app.errorhandler(404)
def page_not_found(e):
    return render_template('404.html'), 404