import json
import magic
import hashlib
import time
import os

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
    fileInfo['info']['digests'].append("MD5:" + hashlib.md5(open(filename, 'rb').read()).hexdigest())
    fileInfo['info']['digests'].append("SHA-1:" + hashlib.sha1(open(filename, 'rb').read()).hexdigest())
    fileInfo['info']['digests'].append("SHA-256:" + hashlib.sha256(open(filename, 'rb').read()).hexdigest())

    # Get public presence
    fileInfo['info']['public_presence'] = {}
    fileInfo['info']['public_presence']['Virustotal'] = True
    fileInfo['info']['public_presence']['IBM X-Force'] = False

    # Return fileInfo
    return fileInfo

def scan(payload):
    # Print payload
    print(payload)