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

def put(file)
  copy "./#{file}", "/site/#{file}"
end

files = [
  "backend",
  "blog",
  "frontend/package.json",
  "frontend/bower.json",
  "frontend/webpack.production.config.js",
  "frontend/src",
  "static",
  "build.sh",
  "run.sh",
]

files.each { |x| put(x) }

### Build
run %q[ cd /site && bash ./build.sh ]

### Cleanup
run %q[ rm -rf /usr/local/go /usr/local/node /site/frontend/node_modules /site/frontend/bower_components /go /site/backend /tmp/phantomjs /site/frontend /site/static ]
run %q[ apt-get remove -y xz-utils bzip2 git-core && apt-get -y autoremove && apt-get clean ]

### Runtime
entrypoint "/sbin/my_init"
cmd "/site/run.sh"

env "USE_ASAR" => "yes"

flatten
tag "xena/christine.website"
