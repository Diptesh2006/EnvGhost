const http = require('http');
const fs = require('fs');
if (!process.env.REQUIRED_TOKEN) { console.error('missing REQUIRED_TOKEN'); process.exit(1); }
if (!fs.existsSync('/app/config.json')) { console.error('missing /app/config.json'); process.exit(1); }
http.createServer((req,res)=>{res.writeHead(200,{'content-type':'application/json'});res.end(JSON.stringify({ok:true}));}).listen(8080);
