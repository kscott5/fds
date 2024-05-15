const process = require('node:process');
console.clear();

console.log("Environment Varaible:");
console.log("************************************************************************\n");
console.log(`export/set AWS_DEFAULT_REGION=${process.env["AWS_DEFAULT_REGION"]}\n\n`);
console.log(`export/set FDS_APPS_USERS_TABLE=${process.env["FDS_APPS_USERS_TABLE"]}`);
console.log(`\n`)

console.log("AWS CLI (access key and secret from IAM center\n")
console.log("https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html")
console.log("**************************************************************************\n");
console.log("aws configure")
