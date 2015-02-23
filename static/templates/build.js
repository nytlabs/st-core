this["JST"] = this["JST"] || {};

this["JST"]["templates/block.html"] = function(obj) {
obj || (obj = {});
var __t, __p = '', __e = _.escape;
with (obj) {
__p += '<script id="block-template" type="text/template">\n    <h4>' +
((__t = ( subhead )) == null ? '' : __t) +
'</h4>\n    <h1><span class="head-text">' +
((__t = ( head )) == null ? '' : __t) +
'</span>\n      <span class="subtext">' +
((__t = ( subtext )) == null ? '' : __t) +
'</span>\n    </h1>\n    <p></p>\n</script>\n';

}
return __p
};