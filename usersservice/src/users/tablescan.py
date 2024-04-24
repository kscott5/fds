import uuid

import boto3
from boto3.dynamodb.table import TableResource
from boto3.dynamodb.conditions import Attr

def lambda_handler(event, context):
    table_name = 'fds.apps.users'
    db = boto3.resource('dynamodb')

    table = db.Table(table_name)
    results = table.scan(FilterExpression=Attr("userid").exists())
    
    return results


