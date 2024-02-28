from app import app
from utils.get_conf import get_config
    
if __name__ == "__main__":
    app.config['uploading'] = False
    app.config['scan_config'] = get_config()
    app.run(port=5001)