from flask import render_template, request, redirect
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

# Define the route that will generate the loader
@app.route('/', methods=['POST'])
def scanFile():
    # Check if the payload file was uploaded
    if 'payload' not in request.files:
        return redirect(request.url)
    payload = request.files['payload']

    # Scan the payload against AVs
    scan(payload)

    # redirect to index
    return redirect(request.url)

# Route for errors
@app.errorhandler(404)
def page_not_found(e):
    return render_template('404.html'), 404