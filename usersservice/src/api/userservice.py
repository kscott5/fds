import json
import os
import uuid

import boto3
from boto3.dynamodb.table import TableResource
from boto3.dynamodb.conditions import Attr # allows use of FilterExpression=Attr('UserId').Exists() or others

# Prepare DynamoDB client
USERS_TABLE = os.getenv('USERS_TABLE', 'fds.apps.users')
dynamodb = boto3.resource('dynamodb')
ddbTable = dynamodb.Table(USERS_TABLE)

def lambda_handler(event, context):
    route_key = f"{event['httpMethod']} {event['resource']}"

    headers = {
         'Content-Type': 'application/json',
         'Access-Control-Allow-Orgins': '*',
    }
    status_code = 400
    response_body = {}

    try:    
        match route_key:
            case 'GET /users':
                ddbResults = ddbTable.scan(Select='ALL_ATTRIBUTES')
                response_body = ddbResults['Items']
                status_code = 200

            case 'GET /users/{userid}':
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


