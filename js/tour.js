/* Copyright 2012 The Go Authors.  All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
(function() {
"use strict";

var slides, editor, $editor, $output;
var slide = null;
var slidenum = 0;

// manage translations
function L(k) {
	if (tr[k]) {
		return tr[k];
	} else {
		console.log("translation missing for: "+k);
		return "(no translation for "+k+")";
	}
}

function init() {
	var $tocdiv = $('<div id="toc" />').insertBefore('#slides').hide();
	$tocdiv.append($('<h2>'+L('toc')+'</h2>'));
	var $toc = $('<ol />').appendTo($tocdiv);
	$("#tocbtn").click(toggleToc);

	slides = $("div.slide");
	slides.each(function(i, slide) {
		var $s = $(slide).hide();

		var $h2 = $s.find("h2").first();
		var $nav;
		if ($h2.length > 0) {
			$("<div/>").addClass("clear").insertAfter($h2);
			$nav = $("<div/>").addClass("nav");
			if (i > 0) {
				$nav.append($("<a>◀</a>").click(function() {
					show(i-1);
					return false;
				}).attr("href", "#"+(i)).attr("title", L('prev')));
			} else {
				$nav.append($("<span>◀</span>"));
			}
			if (i+1 < slides.length) {
				$nav.append($("<a>▶</a>").click(function() {
					show(i+1);
					return false;
				}).attr("href", "#"+(i+2)).attr("title", L('next')));
			} else {
				$nav.append($("<span>▶</span>"));
			}
			$nav.insertBefore($h2);

			var thisI = i;
			var $entry = $("<a />").text($h2.text()).click(function() {
				show(thisI);
			}).attr('href', '#'+(i+1));
			$toc.append($entry);
			$entry.wrap('<li />');

		}
	});

	// set up playground editor
	editor = CodeMirror.fromTextArea(document.getElementById('editor'), {
		theme: "default",
		matchBrackets: true,
		indentUnit: 4,
		tabSize: 4,
		indentWithTabs: false,
		mode: "text/x-go",
		lineNumbers: true,
		extraKeys: {
			"Shift-Enter": function() {
				run();
			}
		}
	});
	$editor = $(editor.getWrapperElement()).attr('id', 'code');
	$output = $('#output');

	$('#more').click(function() {
		$('.controls').toggleClass('expanded');
		return false;
	});
	$('html').click(function() {
		$('.controls').removeClass('expanded');
	});

	$('#run').click(function() {
		run();
		$('.controls').removeClass('expanded');
		return false;
	});

	$('#reset').click(function() {
		reset();
		$('.controls').removeClass('expanded');
		return false;
	});

	$('#kill').click(function() {
		kill();
		$('.controls').removeClass('expanded');
		return false;
	});

	$('#format').click(function() {
		format();
		$('.controls').removeClass('expanded');
		return false;
	});

	$('#togglesyntax').click(function() {
		if (editor.getOption('theme') === 'default') {
			editor.setOption('theme', 'plain');
			$('#togglesyntax').text(L('syntax')+': '+L('off'));
		} else {
			editor.setOption('theme', 'default');
			$('#togglesyntax').text(L('syntax')+': '+L('on'));
		}
		setcookie('theme', editor.getOption('theme'), 14);
		$('.controls').removeClass('expanded');
		return false;
	});

	$('#togglelineno').click(function() {
		if (editor.getOption('lineNumbers')) {
			editor.setOption('lineNumbers', false);
			$('#togglelineno').text(L('lineno')+': '+L('off'));
		} else {
			editor.setOption('lineNumbers', true);
			$('#togglelineno').text(L('lineno')+': '+L('on'));
		}
		setcookie('lineno', editor.getOption('lineNumbers'), 14);
		$('.controls').removeClass('expanded');
		return false;
	});

	if (getcookie('lineno') != ""+editor.getOption('lineNumbers')) {
		$('#togglelineno').trigger('click');
	} else {
		$('#togglelineno').text(L('lineno')+': '+L('on'));
	}

	if (getcookie('theme') != ""+editor.getOption('theme')) {
		$('#togglesyntax').trigger('click');
	} else {
		$('#togglesyntax').text(L('syntax')+': '+L('on'));
	}

	// set these according to lang.js
	$('#run').text(L('run'));
	$('#reset').text(L('reset'));
	$('#format').text(L('format'));
	$('#kill').text(L('kill'));
	$('#tocbtn').attr('title', L('toc'));
	$('#run').attr('title', L('compile'));
	$('#more').attr('title', L('more'));
}

function toggleToc() {
	if ($('#toc').is(':visible')) {
		show(slidenum);
	} else {
		$('#slides, #workspace, #slidenum').hide();
		$('#toc').show();
	}
	return false;
}

function show(i) {
	if(i < 0 || i >= slides.length) {
		return;
	}

	// if a slide is already onscreen, hide it and store its code
	if(slide !== null) {
		var $oldSlide = $(slide).hide();
		if (!$oldSlide.hasClass("nocode")) {
			save(slidenum);
		}
	}

	$('#toc').hide();
	$('#slidenum, #slides').show();

	// switch to new slide
	slidenum = i;
	$("#slidenum").text(i+1);
	slide = slides[i];
	var $s = $(slide).show();

	// load stored code, or hide code box
	if ($s.hasClass("nocode")) {
		$('#workspace').hide();
		$('#wrap').addClass('full-width');
	} else {
		$('#wrap').removeClass('full-width');
		$('#workspace').show();
		$output.empty();
		editor.setValue(load(i) || $s.find('div.source').text());
		editor.focus();
	}

	// update url fragment
	var url = location.href;
	var j = url.indexOf("#");
	if(j >= 0) {
		url = url.substr(0, j);
	}
	url += "#" + (slidenum+1).toString();
	location.href = url;
}

function reset() {
	editor.setValue($(slide).find('div.source').text());
	save(slidenum);
}

var pageData = {};

function save(page) {
	pageData[page] = editor.getValue();
	return true;
}

function load(page) {
	var data = pageData[page];
	if (data) {
		return data;
	}
	return false;
}

function urlSlideNumber(url) {
	var i = url.indexOf("#");
	if(i < 0) {
		return 0;
	}
	var frag = decodeURIComponent(url.substr(i+1));
	if(/^\d+$/.test(frag)) {
		i = parseInt(frag, 10);
		if(i-1 < 0 || i-1 >= slides.length) {
			return 0;
		}
		return i-1;
	}
	return 0;
}

function pageUpDown(event) {
	var e = window.event || event;
	if (e.keyCode === 33) { // page up
		e.preventDefault();
		show(slidenum-1);
		return false;
	}
	if (e.keyCode === 34) { // page down
		e.preventDefault();
		show(slidenum+1);
		return false;
	}
	return true;
}

$(document).ready(function() {
	init();
	if (location.href.indexOf('#') < 0) {
		show(0);
	} else {
		show(urlSlideNumber(location.href));
	}
	document.onkeydown = pageUpDown;
});

$(window).unload(function() {
	save(slidenum);
});

var runFunc, stopFunc;

function body() {
	return editor.getValue();
}
function loading() {
	$output.html('<pre><span class="loading">'+L('waiting')+'</span></pre>');
}
function run() {
	loading();
	stopFunc = runFunc(body(), $output.find("pre")[0]);
}

function kill() {
	if (stopFunc) stopFunc();
}

var seq = 0;

function format() {
	seq++;
	var cur = seq;
	loading();
	$.ajax("/fmt", {
		data: {"body": body()},
		type: "POST",
		dataType: "json",
		success: function(data) {
			if (seq !== cur) {
				return;
			}
			$output.empty();
			if (data.Error) {
				$('<pre class="error" />').text(data.Error).appendTo($output);
				highlightErrors(data.Error);
			} else {
				editor.setValue(data.Body);
			}
		},
		error: function() {
			$('<pre class="error" />').text(L('errcomm')).appendTo($output);
		}
	});
}

function highlightErrors(text) {
	if (!editor || !text) {
		return;
	}
	var errorRe = /[a-z0-9]+\.go:([0-9]+):/g;
	var result;
	while ((result = errorRe.exec(text)) !== null) {
		var line = result[1]*1-1;
		editor.setLineClass(line, null, 'errLine');
	}
	editor.setOption('onChange', function() {
		for (var i = 0; i < editor.lineCount(); i++) {
			editor.setLineClass(i, null, null);
		}
		editor.setOption('onChange', null);
	});
}

// Nasty hack to make this function available to playground.js and socket.js.
window.highlightErrors = highlightErrors;

function getcookie(name) {
	if (document.cookie.length > 0) {
		var start = document.cookie.indexOf(name + '=');
		if (start >= 0) {
			start += name.length + 1;
			var end = document.cookie.indexOf(';', start);
			if (end < 0) {
				end = document.cookie.length;
			}
			return decodeURIComponent(document.cookie.substring(start, end));
		}
	}
	return null;
}

function setcookie(name, value, expire) {
	var expdate = new Date();
	expdate.setDate(expdate.getDate() + expire);
	document.cookie = name + '=' + encodeURIComponent(value) +
		((expire === undefined) ? '' : ';expires=' + expdate.toGMTString());
}

if (window.connectPlayground) {
	runFunc = window.connectPlayground(window.socketAddr);
} else {
	// If this message is logged,
	// we have neglected to include socket.js or playground.js.
	console.log("No playground transport available.");
}

}());
