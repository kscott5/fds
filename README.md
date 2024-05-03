# AWS Serverless Pattern: Module [Food Delivery Server (fds)](https://catalog.workshops.aws/serverless-patterns/en-US)

## [Project Road Map (Pitch)](https://catalog.workshops.aws/serverless-patterns/en-US/business-scenario)
As more people work remotely, restaurants have experienced huge growth for carry-out orders. Customers want dishes from their favorite restaurants, but they are also more aware and concerned about the environmental impact of delivery services.

A startup company received funding to connect restaurants and nearby customers. The new service will exclusively use electric bike riders to pick up and deliver food to hyperlocal customers. The goal is to create a fast, economical, and environmentally friendly delivery service.

> [!NOTE]
> This is the [self-service](https://catalog.workshops.aws/serverless-patterns/en-US/logistics-self-service) project option. My local laptop has the development tools used with these modules. Don't forget to visit [aws cli install](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)

```shell 
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
rm ./aws -rf
```

## Local development
```shell
docker pull amazon/aws-lambda-go
docker pull amazon/aws-lamdba-python
docker pull amazon/dynamodb-local
```
## Local or remote staging of test data
export AWS_ENPOINT_URL_DYNAMODB={localhost of docker dynamodb-local}
```shell
aws dynamodb batch-write-item --endpoint_url $AWS_ENDPOINT_URL_DYNAMODB  --request file://./data/data.json
```
