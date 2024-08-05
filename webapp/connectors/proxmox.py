from app import app, Log
from proxmoxer import ProxmoxAPI
import time
import urllib3
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

connectorSettings   = app.config['config']['connector']
machines            = app.config['config']['machines']

def proxmox_connect():
    return ProxmoxAPI(
        connectorSettings['url'],
        user=connectorSettings['user'],
        password=connectorSettings['password'],
        verify_ssl=False
    )

def vm_action(action):
    try:
        proxmox = proxmox_connect()
        
        for machine in machines:
            vm_name = machine['name']
            Log.section(f"{action.capitalize()}ing VM: {vm_name}")
            
            for node in proxmox.nodes.get():
                for vm in proxmox.nodes(node['node']).qemu.get():
                    if vm['name'] == vm_name:
                        vm_id = vm['vmid']
                        node_name = node['node']
                        
                        if action == 'start':
                            proxmox.nodes(node_name).qemu(vm_id).status.start.post()
                        elif action == 'stop':
                            proxmox.nodes(node_name).qemu(vm_id).status.stop.post()
                        
                        Log.section(f"VM {vm_name} (ID: {vm_id}) {action} command sent successfully")
                        
                        while True:
                            status = proxmox.nodes(node_name).qemu(vm_id).status.current.get()
                            expected_status = 'running' if action == 'start' else 'stopped'
                            if status['status'] == expected_status:
                                Log.info(f"VM {vm_name} is now {expected_status}")
                                break
                            time.sleep(5)  # Wait for 5 seconds before checking again
                        
                        break  # Exit the inner loop once the VM is found and actioned
                else:
                    continue
                break  # Exit the outer loop once the VM is found and actioned
            else:
                Log.error(f"VM {vm_name} not found")

        return True
    except Exception as e:
        Log.error(f"Error {action}ing VMs: {str(e)}")
        return False

def turnOn():
    return vm_action('start')

def turnOff():
    return vm_action('stop')