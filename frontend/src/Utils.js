// Module App.BlogEntry

showdown = require("showdown");

showdown.extension('blog', function() {
  return [{
    type: 'output',
    regex: /<ul>/g,
    replace: '<ul class="browser-default">'
  }];
});

exports.mdify = function(corpus) {
  var converter = new showdown.Converter({ extensions: ['blog'] });
  return converter.makeHtml(corpus);
};
