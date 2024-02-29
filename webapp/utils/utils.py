import json
import magic
import hashlib
import time
import os
import requests

def getFileInfo():
    filename = './uploads/payload'
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
    fileInfo['info']['public_presence']['IBM X-Force'] = checkIIBMXForce(digest)

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
    if response.json().get('error'):
        return '❌'
    else:
        return '<a href="https://www.virustotal.com/gui/search/' + hash + '" target="_blank">✅</a>'

def checkIIBMXForce(hash):
    headers = {
        "Host": "exchange.xforce.ibmcloud.com",
        "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0",
        "Accept": "application/json, text/plain, */*",
        "Accept-Language": "en-US,en;q=0.5",
        "Accept-Encoding": "gzip, deflate",
        "X-Ui": "XFE",
        "Sec-Fetch-Dest": "empty",
        "Sec-Fetch-Mode": "cors",
        "Sec-Fetch-Site": "same-origin",
        "Te": "trailers"
    }

    response = requests.get('https://exchange.xforce.ibmcloud.com/api/malware/' + hash, headers=headers)
    
    # Check for response
    if response.json().get('error'):
        return '❌'
    else:
        return '<a href="https://exchange.xforce.ibmcloud.com/malware/' + hash + '" target="_blank">✅</a>'