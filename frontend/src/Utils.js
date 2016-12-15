// Module App.BlogEntry

exports.mdify = function(id) {
  var converter = new showdown.Converter()
  elem = document.getElementById(id);
  md = elem.innerHTML;
  elem.innerHTML = converter.makeHtml(md);
  return "done :)";
}
