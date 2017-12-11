from "xena/go-mini:1.9.2"

### setup go
run "go1.9.2 download"

### Copy files
run "mkdir -p /site"

def debug?()
  getenv("DEBUG") == "yes"
end

def debug!()
  run "apk add --no-cache bash"
  debug
end

def put(file)
  copy "./#{file}", "/site/#{file}"
end

files = [
  "blog",
  "templates",
  "gops.go",
  "hash.go",
  "html.go",
  "main.go",
  "rice-box.go",
  "rss.go",
  "run.sh"
]

files.each { |x| put x }

copy "vendor/", "/root/go/src/"

### Build
run "cd /site && go1.9.2 build -v"

### Cleanup
run %q[ rm -rf /root/go /site/backend /root/sdk /site/*.go ]
run %q[ apk del git go1.9.2 ]

cmd "/site/run.sh"

flatten
tag "xena/christine.website"
