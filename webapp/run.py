from app import app
from utils.retrieve_config import getConfig
import json

if __name__ == "__main__":
    app.config['config'] = getConfig()
    app.run(host='0.0.0.0', port=8000)