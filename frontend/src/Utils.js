// Module App.BlogEntry

function htmlDecode(value){
  return $('<div/>').html(value).text();
}

exports.mdify = function(id) {
  var converter = new showdown.Converter()
  elem = document.getElementById(id);
  md = elem.innerHTML;
  elem.innerHTML = htmlDecode(converter.makeHtml(md));
  return "done :)";
}
