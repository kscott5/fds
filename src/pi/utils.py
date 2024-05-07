import boto3
import uuid

def load_test_data(event, context):
    table_name = 'fds.apps.users'
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table(table_name)

    result = None
    people = [
            { 'userid' : 'marivera', 'name' : 'Martha Rivera'},
            { 'userid' : 'nikkwolf', 'name' : 'Nikki Wolf'},
            { 'userid' : 'pasantos', 'name' : 'Paulo Santos'},
        ]

    with table.batch_writer() as batch_writer:
        for person in people:
            item = {
                '_id'     : uuid.uuid4().hex,
                'Userid'  : person['userid'],
                'FullName': person['name']
            }
            print("> batch writing: {}".format(person['userid']) )
            batch_writer.put_item(Item=item)
            
        result = f"Success. Added {len(people)} people to {table_name}."

    return {'message': result}

if __name__ == "__main__" :
    load_test_data({}, {})