from "phusion/baseimage:0.9.19"

copy "./runit/", "/etc/system/"

run %q[ curl -o backplane-stable-linux-amd64.tgz https://bin.equinox.io/c/jWahGASjoRq/backplane-stable-linux-amd64.tgz \
     && tar xf backplane-stable-linux-amd64.tgz \
     && mv backplane /usr/bin/backplane \
     && rm backplane-stable-linux-amd64.tgz ]

run %q[ cd /usr/local && curl -o go.tar.gz https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz \
     && tar xf go.tar.gz && rm go.tar.gz && mkdir -p /go/src ]
env "GOPATH" => "/go"
run %q[ apt-get update && apt-get install -y git-core \
     && /usr/local/go/bin/go get github.com/gernest/front && /usr/local/go/bin/go get gopkg.in/yaml.v2 ]

run %q[ apt-get update && apt-get install xz-utils bzip2 ]

run %q[ cd /usr/local && curl -o node.tar.xz https://nodejs.org/dist/v6.9.2/node-v6.9.2-linux-x64.tar.xz \
     && tar xf node.tar.xz && mv node-v6.9.2-linux-x64 node && rm node.tar.xz ]

### Copy files
run "mkdir -p /site"
copy "./backend", "/site/backend"
copy "./blog", "/site/blog"
copy "./frontend/package.json", "/site/frontend/package.json"
copy "./frontend/bower.json", "/site/frontend/bower.json"
copy "./frontend/webpack.production.config.js", "/site/frontend/webpack.production.config.js"
copy "./frontend/src", "/site/frontend/src"
copy "./frontend/support", "/site/frontend/support"
copy "./static", "/site/static"
copy "./build.sh", "/site/build.sh"
copy "./run.sh", "/site/run.sh"

### Build
run %q[ cd /site && bash ./build.sh ]

### Cleanup
run %q[ rm -rf /usr/local/go /usr/local/node /site/frontend/node_modules /site/frontend/bower_components /go /site/backend ]
run %q[ apt-get remove -y xz-utils bzip2 git-core ]

### Runtime
entrypoint "/sbin/my_init"
cmd "/site/run.sh"

flatten
tag "xena/christine.website"
