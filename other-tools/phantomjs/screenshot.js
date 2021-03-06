var system = require('system');
var fs = require('fs');

var args = system.args;
if (args.length < 2) {
    console.log('args', args);
    console.log('phantomjs <url>')
    phantom.exit();
}

var urlSplit = args[1].split(':');
var port = urlSplit[urlSplit.length-1];
var protocol = urlSplit[0];
if (port == '443' && protocol == 'http') {
  console.log("Protocol mismatch");
  phantom.exit();
}

var outstring = args[1].split("://").join(".");
var webpage = require('webpage').create();
webpage.settings.resourceTimeout = 3000;

webpage.open(args[1], function(res) {
  if (res !== 'success') {
    console.log('FAIL:', outstring);
    fs.write("failed.txt", "\n fail_" + outstring + ".txt", 'a');
  } else {
    console.log('SUCCESS:', outstring);
    webpage.render("success_" + outstring + ".png");
  }
  phantom.exit();
});
