import json
import os
import uuid

import boto3
from boto3.dynamodb.table import TableResource
from boto3.dynamodb.conditions import Attr # allows use of FilterExpression=Attr('UserId').Exists() or others

# Prepare DynamoDB client
USERS_TABLE = os.getenv('FDS_APPS_USERS_TABLE', 'fds.apps.users')
dynamodb = boto3.resource('dynamodb')
ddbTable = dynamodb.Table(USERS_TABLE)

def lambda_handler(event, context):
    route_key = f"{event['httpMethod']} {event['resource']}"

    status_code = 400
    response_body = {}
    headers = {
         'content-type': 'application/json'
    }

    try:    
        if route_key == 'GET /users':
            ddbResults = ddbTable.scan(Select='ALL_ATTRIBUTES')
            response_body = ddbResults['Items']
            status_code = 200

    except Exception as err:
            response_body = f"Error: {str(err)}"

    return {
         'statusCode': status_code,
         'headers': headers,
         'body': json.dumps(response_body)
    }


