(function() {

var slides, editor, $editor, $output;
var slide = null;
var slidenum = 0;
var codebox = null;

function init() {
	if (tourMode == 'local') {
		$('.appengineMode').remove();
	} else {
		$('.localMode').remove();
	}

	var $toc = $("#toc").hide();
	$("#tocbtn").click(function() {
		if ($("#toc").is(":visible")) {
			hideToc();
		} else {
			showToc();
		}
	});

	slides = $("div.slide");
	slides.each(function(i, slide) {
		var $s = $(slide).hide();

		var $sdiv = $s.find("div");
		if (!$s.hasClass("nocode") && $sdiv.length > 0) {
			var $div = $sdiv.last();
			$div.remove();
			$s.data("code", $div.text().trim());
		}

		var $content = $('<div class="content"/>');
		$content.html($s.html());
		$s.empty().append($content);

		var $h2 = $content.find("h2").first();
		var $nav;
		if ($h2.length > 0) {
			$("<div/>").addClass("clear").insertAfter($h2);
			$nav = $("<div/>").addClass("nav")
			if (i > 0) {
				$nav.append($("<button>").click(function() {
					show(i-1);
				}).text("PREV").addClass("prev"));
			}
			if (i+1 < slides.length) {
				$nav.append($("<button>").click(function() {
					show(i+1);
				}).text("NEXT").addClass("next"));
			}
			$nav.insertBefore($h2);

			var thisI = i;
			var $entry = $("<li/>").text($h2.text()).click(
				function() { hideToc(); show(thisI); });
			$toc.append($entry);

		}
		if ($s.hasClass("nocode")) {
			$h2.addClass("nocode");
		}
	});

	// set up playground editor
	$editor = $('<div id="code"><button id="run">RUN</button><textarea/></div>');
	$editor.insertBefore("#slides");
	$output = $('<div id="output"/>').insertBefore("#slides");
	editor = playground({
		codeEl: "#code textarea",
		outputEl: "#output",
		runEl: "#run"
	});
}

function showToc() {
	$("#toc").show();
	$("#slides, #code, #output").hide();
	$("#tocbtn").text("SLIDES");
}

function hideToc() {
	$("#toc").hide();
	$("#slides, #code, #output").show();
	$("#tocbtn").text("INDEX");
}

function show(i) {
	if(i < 0 || i >= slides.length)
		return;
		
	// if a slide is already onscreen, hide it and store its code
	if(slide != null) {
		var $oldSlide = $(slide).hide();
		if (!$oldSlide.hasClass("nocode")) {
			$oldSlide.data("code", editor.getValue());
		}
	}

	// switch to new slide
	slidenum = i;
	$("#slidenum").text(i+1);
	slide = slides[i];
	var $s = $(slide).show();

	// load stored code, or hide code box
	if ($s.hasClass("nocode")) {
		$editor.hide();
		$output.hide();
	} else {
		$editor.show();
		$output.show().empty();
		editor.setValue($s.data("code"));
		editor.focus();
	}

	// update url fragment
	var url = location.href;
	var j = url.indexOf("#");
	if(j >= 0)
		url = url.substr(0, j);
	url += "#" + (slidenum+1).toString();
	location.href = url;
}

function urlSlideNumber(url) {
	var i = url.indexOf("#");
	if(i < 0)
		return 0;
	var frag = unescape(url.substr(i+1));
	if(/^\d+$/.test(frag)) {
		i = parseInt(frag);
		if(i-1 < 0 || i-1 >= slides.length)
			return 0;
		return i-1;
	}
	return 0;
}

function pageUpDown(event) {
	var e = window.event || event;
	if (e.keyCode == 33) { // page up
		e.preventDefault();
		show(slidenum-1);
		return false;
	}
	if (e.keyCode == 34) { // page down
		e.preventDefault();
		show(slidenum+1);
		return false;
	}
	return true;
}

$(document).ready(function() {
	init();
	$('body').removeClass('loading');
	show(urlSlideNumber(location.href));
	document.onkeydown = pageUpDown;
});

}());
