// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// HACK: This is file is identical to go.talks/present's playground.js except
// the lines containing "HACK". Also, the window.playground function is removed.

(function() {

  function lineHighlight(error) {
    // HACK: hook back into tour.js.
    if (window.highlightErrors) window.highlightErrors(error); // HACK
  }

  function connectPlayground() {
    var playbackTimeout;

    function playback(pre, events) {
      if (!pre.data("cleared")) pre.empty().data("cleared", true); // HACK
      function show(msg) {
        // ^L clears the screen.
        var msgs = msg.split("\x0c");
        if (msgs.length == 1) {
          pre.text(pre.text() + msg);
          return;
        }
        pre.text(msgs.pop());
      }
      function next() {
        if (events.length === 0) {
          var exit = $('<span class="exit"/>');
          exit.text("\nProgram exited.");
          exit.appendTo(pre);
          return;
        }
        var e = events.shift();
        if (e.Delay === 0) {
          show(e.Message);
          next();
        } else {
          playbackTimeout = setTimeout(function() {
            show(e.Message);
            next();
          }, e.Delay / 1000000);
        }
      }
      next();
    }

    function stopPlayback() {
      clearTimeout(playbackTimeout);
    }

    function setOutput(output, events, error) {
      stopPlayback();
      output.empty();
  
      // Display errors.
      if (error) {
        lineHighlight(error);
        output.addClass("error").text(error);
        return;
      }
  
      // Display image output.
      if (events.length > 0 && events[0].Message.indexOf("IMAGE:") === 0) {
        var out = "";
        for (var i = 0; i < events.length; i++) {
          out += events[i].Message;
        }
        var url = "data:image/png;base64," + out.substr(6);
        $("<img/>").attr("src", url).appendTo(output);
        return;
      }
  
      // Play back events.
      if (events !== null) {
        playback(output, events);
      }
    }

    var seq = 0;
    function runFunc(body, output) {
      output = $(output);
      seq++;
      var cur = seq;
      var data = {
        "version": 2,
        "body": body
      };
      $.ajax("/compile", {
        data: data,
        type: "POST",
        dataType: "json",
        success: function(data) {
          if (seq != cur) {
            return;
          }
          if (!data) {
            return;
          }
          if (data.Errors) {
            setOutput(output, null, data.Errors);
            return;
          }
          setOutput(output, data.Events, false);
        },
        error: function() {
          output.addClass("error").text(
            "Error communicating with remote server."
          );
        }
      });
      return stopPlayback;
    }

    return runFunc;
  }

  window.connectPlayground = connectPlayground;

})();
