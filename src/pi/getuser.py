import json
import os
import uuid

import boto3
from boto3.dynamodb.table import TableResource
from boto3.dynamodb.conditions import Attr # allows use of FilterExpression=Attr('UserId').Exists() or others

# Prepare DynamoDB client
USERS_TABLE = os.getenv('FDS_APPS_USERS_TABLE', 'FDSAppsUsers')
dynamodb = boto3.resource('dynamodb')
ddbTable = dynamodb.Table(USERS_TABLE)

def lambda_handler(event, context):
    route_key = f"{event['httpMethod']} {event['resource']}"

    headers = {
         'conten-type': 'application/json',
         'access-control-allow-orgins': '*',
    }
    status_code = 400
    response_body = {}

    try:    
        if route_key == 'GET /users{userid}':
            ddbResults = ddbTable.get_item(
                Key={'userid': event['pathParameters']['userid']}
            )

            if 'Item' in ddbResults:
                response_body = ddbResults['Item']
            else:
                response_body = {}
            status_code = 200

    except Exception as err:
        response_body = f"Error: {str(err)}"

    return {
         'headers': headers,
         'statuscode': status_code,
         'body': json.dumps(response_body)
    }
