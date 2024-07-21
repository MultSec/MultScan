from app import app, Log
from utils.retrieve_config import getConfig
import json
from utils.utils import turnOnMachines, turnOffMachines, cleanUploads

if __name__ == "__main__":
    app.config['config'] = getConfig()

    try:
        turnOnMachines()
        app.run(host='0.0.0.0', port=8000)
    except Exception as e:
        Log.error(f"An error occurred: {str(e)}")
    finally:
        turnOffMachines()
        cleanUploads()
        Log.info("Shutting down server")

