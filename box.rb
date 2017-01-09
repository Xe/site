from "xena/go"

### Copy files
run "mkdir -p /site"

def put(file)
  copy "./#{file}", "/site/#{file}"
end

files = [
  "backend",
  "blog",
  "frontend.asar",
  "static",
  "build.sh",
  "run.sh",

  # This file is packaged in the asar file, but the go app relies on being
  # able to read it so it can cache the contents in ram.
  "frontend/static/dist/index.html",
]

files.each { |x| put x }

### Build
run "apk add --no-cache git"
run %q[ cd /site && sh ./build.sh ]

### Cleanup
run %q[ rm -rf /usr/local/go /usr/local/node /site/frontend/node_modules /site/frontend/bower_components /go /site/backend /tmp/phantomjs ]
run %q[ apk del go git ]

### Runtime
cmd "/site/run.sh"

env "USE_ASAR" => "yes"

flatten
tag "xena/christine.website"
