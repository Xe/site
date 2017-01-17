from "xena/go:1.7.4"

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

copy "vendor/", "/go/src/"

### Build
run %q[ cd /site && sh ./build.sh ]
debug! if debug?

### Cleanup
run %q[ rm -rf /usr/local/go /go /site/backend ]
run %q[ apk del go git ]

### Runtime
cmd "/site/run.sh"

env "USE_ASAR" => "yes"

flatten
tag "xena/christine.website"
