### bunny

Main goal:

gather information about running rabbitmq instances and update clustering configuration.
For use with auto scaling and cloudformation.

UserData part of CloudFormation template should look something like:

```javascript
"UserData": { "Fn::Base64": { "Fn::Join": [ "", [
    "[ ... more stuff here ...]\n",
    "docker run -d -v /tmp:/tmp tray/bunny:latest -destination /tmp/rabbit-cluster.conf\n",
    "while [ ! -f /tmp/rabbit-cluster.conf ]; do\n",
    "  sleep 2\n",
    "done\n",
    "docker run -d -v /tmp/rabbit-cluster.conf:/tmp/rabbit-cluster.conf rabbitmq:latest rabbitmq-server -config /tmp/rabbit-cluster.conf\n",
    "[ ... mode stuff here ...]\n"
]]}}
```
