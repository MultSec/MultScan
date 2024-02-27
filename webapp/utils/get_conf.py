import yaml

def get_config():
    with open('config.yml', 'r') as ymlfile:
        cfg = yaml.safe_load(ymlfile)

    print(" * Using connector: " + cfg['config']['connector']['connector_type'])
    print(" * Loaded " + str(len(cfg['config']['machines'])) + " machines")
    for machine in cfg['config']['machines']:
        print(" *  - " + machine['machine_name'] + " (" + machine['machine_ip'] + ")")
    
    return cfg