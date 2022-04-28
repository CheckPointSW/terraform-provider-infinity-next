#!/usr/bin/env python3
import sys
import time
import requests
import argparse

parser = argparse.ArgumentParser(description='Infinity Next Enforce Utility')
parser.add_argument(
    '--region', help='The Infinity Next region to enforce policy on: eu or us (default: eu)', default="eu")
parser.add_argument(
    '--client-id', help='The Infinity Next client ID: required', required=True)
parser.add_argument(
    '--access-key', help='The Infinity Next access key: required', required=True)
args = parser.parse_args()

if args.region == "eu":
    inext_url = "https://cloudinfra-gw.portal.checkpoint.com"
elif args.region == "us":
    inext_url = "https://cloudinfra-gw-us.portal.checkpoint.com"
else:
    raise argparse.ArgumentError(argument='region', message='got unknown region ' + args.region + ', expected eu or us')

print('authenticating to ' + inext_url)

auth_reponse = requests.post(inext_url+'/auth/external', headers={'Content-Type': 'application/x-www-form-urlencoded'}, data={
                             'clientId': args.client_id, 'accessKey': args.access_key})
if auth_reponse.status_code != 200:
    raise SystemExit('failed authenticating, got response ' + str(auth_reponse.content) + ' exiting...')
auth = auth_reponse.json()

if 'data' not in auth.keys() or 'token' not in auth['data'].keys():
    raise SystemExit('got malformed authentication response ' + str(auth) + ' exiting...')
token = auth['data']['token']

enforce_reponse = requests.post(inext_url + '/app/i2/graphql/V1', headers={'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json'}, json={
                                "query": "mutation {enforcePolicy {id}}", "variables": {}})
if enforce_reponse.status_code != 200:
    raise SystemExit('failed enforcing policy, got response ' + str(enforce_reponse.content) + ' exiting...')
enforce = enforce_reponse.json()
if 'data' not in enforce.keys() or 'enforcePolicy' not in enforce['data'].keys() or 'id' not in enforce['data']['enforcePolicy'].keys():
    raise SystemExit('got malformed enforce policy response ' + str(enforce) + ' exiting...')
task_id = enforce['data']['enforcePolicy']['id']

task_status = 'InProgress'
print('Waiting for task ' + task_id + ' to complete...')
while task_status == 'InProgress':
    task_response = requests.post(inext_url + '/app/i2/graphql/V1', headers={'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json'}, json={
        "query": "query {getTask(id: \""+task_id+"\") {id\r\nstatus}}", "variables": {}})
    if task_response.status_code != 200:
        print('failed polling task ' + task_id + ' status with response: ' + task_response.content)

    task = task_response.json()
    if 'data' not in task.keys() or 'getTask' not in task['data'].keys() or 'status' not in task['data']['getTask'].keys():
        raise SystemExit('got malformed task response ' + str(task) + ' exiting...')
    task_status = task['data']['getTask']['status']
    time.sleep(0.5)

print('Enforce policy task ' + task_id + ' done with status "' + task_status+'"')
sys.exit(0)
