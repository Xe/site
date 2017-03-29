local sh = require "sh"
local fs = require "fs"

sh { abort = true }

local cd = function(path)
  local ok, err = fs.chdir(path)
  if err ~= nil then
    error(err)
  end
end

cd "frontend"
sh.rm("-rf", "node_modules", "bower_components"):ok()
print "running npm install..."
sh.npm("install"):print()
print "running npm run build..."
sh.npm("run", "build"):print()
print "packing frontend..."
sh.asar("pack", "static", "../frontend.asar"):print()
cd ".."

sh.box("box.rb"):print()
