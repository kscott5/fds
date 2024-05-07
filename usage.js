const process = require('node:process');
console.clear();

console.log("Environment Varaible:");
console.log("**********************************************\n");
console.log(`export/set AWS_ACCESS_KEY_ID=${process.env["AWS_ACCESS_KEY_ID"]}`);
console.log(`export/set AWS_SECRET_ACCESS_KEY=${process.env["AWS_SECRET_ACCESS_KEY"]}`); 
console.log(`export/set AWS_REGION=${process.env["AWS_REGION"]}\n\n`);

