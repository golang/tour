function playground(codeEl, outputEl, runEl, shareEl, shareURLEl) {
	var code = $(codeEl);
	var editor = CodeMirror.fromTextArea(
		code[0],
		{
			lineNumbers: true,
			indentUnit: 8,
			indentWithTabs: true,
			onKeyEvent: function(editor, e) {
				if (e.keyCode == 13 && e.shiftKey) {
					if (e.type == "keydown") {
						run();
					}
					e.stop();
					return true;
				}
			}
		}
	);
	var output = $(outputEl);

	function clearErrors() {
		var lines = editor.lineCount();
		for (var i = 0; i < lines; i++) {
			editor.setLineClass(i, null);
		}
	}
	function highlightErrors(text) {
		var errorRe = /[a-z]+\.go:([0-9]+): /g;
		var result;
		while ((result = errorRe.exec(text)) != null) {
			var line = result[1]*1-1;
			editor.setLineClass(line, "errLine")
		}
	}

	var seq = 0;
	function run() {
		clearErrors();
		output.removeClass("error").html(
			'<div class="loading">Waiting for remote server...</div>'
		);
		seq++;
		var cur = seq;
		$.ajax("/compile", {
			processData: false,
			data: editor.getValue(),
			type: "POST",
			dataType: "json",
			success: function(data) {
				if (seq != cur) {
					return;
				}
				pre = $("<pre/>");
				output.empty().append(pre);
				if (data.compile_errors != "") {
					pre.text(data.compile_errors);
					output.addClass("error");
					highlightErrors(data.compile_errors);
				} else {
					var out = ""+data.output;
					if (out.indexOf("IMAGE:") == 0) {
						var img = $("<img/>");
						var url = "data:image/png;base64,";
						url += out.substr(6)
						img.attr("src", url);
						output.empty().append(img);
					} else {
						pre.text(data.output);
					}
				}
			},
			error: function() {
				output.addClass("error").text(
					"Error communicating with remote server."
				);
			}
		});
	}
	$(runEl).click(run);

	if (shareEl == null || shareURLEl == null) {
		return editor;
	}

	function origin(href) {
		return (""+href).split("/").slice(0, 3).join("/");
	}

	var shareURL = $(shareURLEl).hide();
	var sharing = false;
	$(shareEl).click(function() {
		if (sharing) return;
		sharing = true;
		$.ajax("/share", {
			processData: false,
			data: editor.getValue(),
			type: "POST",
			complete: function(xhr) {
				sharing = false;
				if (xhr.status != 200) {
					alert("Server error; try again.");
					return
				}
				var url = origin(window.location) + "/p/" +
					xhr.responseText;
				shareURL.show().val(url).focus().select();
			}
		});
	});

	return editor;
}
