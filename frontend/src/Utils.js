// Module App.BlogEntry

function htmlDecode(input) {
  var doc = new DOMParser().parseFromString(input, "text/html");
  return doc.documentElement.textContent;
}

exports.mdify = function(id) {
  var converter = new showdown.Converter()
  elem = document.getElementById(id);
  md = elem.innerHTML;
  elem.innerHTML = htmlDecode(converter.makeHtml(md));
  return "done :)";
}
