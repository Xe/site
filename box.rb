from "xena/go-mini:1.8.1"

### setup go
run "go1.8.1 download"

### Copy files
run "mkdir -p /site"

def debug?()
  return getenv("DEBUG") == "yes"
end

def debug!()
  run "apk add --no-cache bash"
  debug
  run "apk del bash"

  puts "hint: flatten this image if deploying."
end

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

copy "vendor/", "/root/go/src/"

### Build
run "apk add --no-cache --virtual site-builddep build-base"
run %q[ cd /site && sh ./build.sh ]
debug! if debug?

### Cleanup
run %q[ rm -rf /root/go /site/backend /root/sdk ]
run %q[ apk del git go1.8 site-builddep ]

### Runtime
cmd "/site/run.sh"

env "USE_ASAR" => "yes"

flatten
tag "xena/christine.website"
