from flask import render_template, request, redirect
from app import app
import importlib

# favicon.ico route
@app.route('/favicon.ico')
def favicon():
    # Return the favicon
    return send_file('./static/MultScan.ico', mimetype='image/vnd.microsoft.icon')

# Define the route that will serve the loader configuration
@app.route('/', methods=['GET'])
def serve_conf():
    return render_template('index.html')

# Route for errors
@app.errorhandler(404)
def page_not_found(e):
    return render_template('404.html'), 404