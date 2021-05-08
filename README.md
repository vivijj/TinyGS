## TinyGS
### Tiny gin server -- Simplest restfulAPI webserver 

use Gin as the web framework and Zap for fast logging


### How to run it ?
``` go run server.go ```

### bench
use postman to test it & longest response time is about 3s

### Picmaker
make picture with the message the user set


- Method:
 POST
- URL:
 /picmaker/${useid}
- DATA:
 message
