// Module App.BlogEntry

exports.mdify = function(id) {
  var converter = new showdown.Converter()
  elem = document.getElementById(id);
  md = elem.innerHTML;
  elem.innerHTML = unescape(converter.makeHtml(md));
  return "done :)";
}
