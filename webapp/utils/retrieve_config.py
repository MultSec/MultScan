import yaml
from app import Log

def getConfig():
    with open('config.yml', 'r') as ymlfile:
        cfg = yaml.safe_load(ymlfile)

    Log.info(f"Using connector: {cfg['config']['connector']['connector_type']}")
    Log.info(f"Loaded {str(len(cfg['config']['machines']))} machines:")
    for machine in cfg['config']['machines']:
        Log.section(f'{machine['machine_name']} ({machine['machine_ip']})')
    
    return cfg['config']