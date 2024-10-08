import magic
import hashlib
import os
import sys
import json
import requests
import importlib
import shutil
from app import app, Log

def turnOnMachines():
    connectorPath = "connectors." + app.config['config']['connector']['type']

    Log.info("Turning On Machines")

    importlib.import_module(connectorPath).turnOn()

def turnOffMachines():
    # Clean line
    sys.stdout.write('\r\033[K')
    sys.stdout.flush()

    connectorPath = "connectors." + app.config['config']['connector']['type']

    Log.info("Turning Off Machines")

    importlib.import_module(connectorPath).turnOff()

def cleanUploads():
    Log.info("Cleaning Up uploads storage")

    parent_dir = './uploads'

    # Iterate over all items in the parent directory
    for item in os.listdir(parent_dir):
        item_path = os.path.join(parent_dir, item)
        # Check if the item is a directory
        if os.path.isdir(item_path):
            try:
                # Remove the directory and all its contents
                shutil.rmtree(item_path)
            except Exception as e:
                Log.error(f"Error removing {item_path}: {e}")

def fileInfo(filename):
    # Create emptyt dictionary to store file info
    fileInfo = {}
    fileInfo['info'] = {}
            
    # Get file size in bytes and convert to MB
    sizeBytes = os.path.getsize(filename)
    sizeMB = sizeBytes / (1024 * 1024)
    fileInfo['info']['size'] = f'{sizeMB:.2f} MB ({sizeBytes} bytes)'

    # Get file type
    fileInfo['info']['type'] = magic.from_file(filename)

    # Get digests
    fileInfo['info']['digests'] = []
    
    with open(filename, 'rb') as f:
        data = f.read()
        fileInfo['info']['digests'].append("MD5:" + hashlib.md5(data).hexdigest())
        fileInfo['info']['digests'].append("SHA-1:" + hashlib.sha1(data).hexdigest())
        fileInfo['info']['digests'].append("SHA-256:" + hashlib.sha256(data).hexdigest())

    # Public presence
    digest = fileInfo['info']['digests'][2].split(':')[1]
    fileInfo['info']['public_presence'] = {}
    fileInfo['info']['public_presence']['Virustotal'] = checkVirusTotal(digest)

    # Return fileInfo
    return fileInfo

def checkVirusTotal(hash):
    headers = {
        'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0',
        'Accept': 'application/json',
        'Accept-Language': 'en-US,en;q=0.5',
        'Accept-Encoding': 'gzip, deflate',
        'Referer': 'https://www.virustotal.com/',
        'Content-Type': 'application/json',
        'X-Tool': 'vt-ui-main',
        'X-App-Version': 'v1x249x0',
        'Accept-Ianguage': 'en-US,en;q=0.9,es;q=0.8',
        'X-Vt-Anti-Abuse-Header': 'a',
        'Sec-Fetch-Dest': 'empty',
        'Sec-Fetch-Mode': 'cors',
        'Sec-Fetch-Site': 'same-origin',
        'Te': 'trailers',
    }

    response = requests.get('https://www.virustotal.com/ui/files/' + hash, headers=headers)

    # Check for response {"error":{"code":"NotFoundError","message":"Resource not found."}} that indicates file not found in VirusTotal

    return not response.json().get('error')

def getSampleStatus(id):
    statusFilePath = f'./uploads/{id}/status'

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

            status["status"][machine["name"]] = {"badBytes": '', "result": 'Scanning'}

        # Write json to file
        with open(statusFilePath, 'w') as statusFile:
            json.dump(status, statusFile, indent=4)

        # Request scan to agents
        Log.info(f"[\033[34m{id}\033[0m] Requesting sample scan")
        for machine in app.config['config']['machines']:
            Log.subsection(f"[\033[34m{id}\033[0m] Requesting {machine['name']} ({machine['ip']}) for scan")
            
            response = requests.get(f"http://{machine['ip']}:9090/task/{id}")

            # Check if the request was successful
            if response.status_code != 200:
                Log.error(f"Request failed with status code: {response.status_code}")
    else:
        # Check status on agents
        Log.info(f"[\033[34m{id}\033[0m] Checking sample scan status")

        # Load existing status
        with open(statusFilePath, 'r') as statusFile:
            status = json.load(statusFile)
        
        # Update each machine status
        for machine in app.config['config']['machines']:
            Log.subsection(f"[\033[34m{id}\033[0m] Requesting {machine['name']} ({machine['ip']}) for scan")
            
            response = requests.get(f"http://{machine['ip']}:9090/task/{id}")

            # Check if the request was successful
            if response.status_code == 200:
                # Update sample status
                status["status"][machine["name"]] = response.json()
            else:
                Log.error(f"Request failed with status code: {response.status_code}")

        # Update status file
        with open(statusFilePath, 'w') as statusFile:
            json.dump(status, statusFile, indent=4)

    # Load status file
    with open(statusFilePath, 'r') as statusFile:
        status = json.load(statusFile)

    # Return status
    return status