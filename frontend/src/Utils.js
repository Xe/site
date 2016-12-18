// Module App.BlogEntry

exports.mdify = function(corpus) {
  var converter = new showdown.Converter()
  return converter.makeHtml(corpus);
}
