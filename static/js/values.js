/* Copyright 2012 The Go Authors.   All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
'use strict';

angular.module('tour.values', []).

// List of modules with description and lessons in it.
value('TableOfContents', [{
    'id': 'mechanics',
    'title': 'Using the tour',
    'description': '<p>Welcome to a tour of the <a href="http://golang.org">Go programming language</a>. The tour covers the most important features of the language, mainly:</p>',
    'lessons': ['welcome']
}, {
    'id': 'basics',
    'title': 'Basics',
    'description': '<p>The starting point, learn all the basics of the language.</p><p>Declaring variables, calling functions, and all the things you need to know before moving to the next lessons.</p>',
    'lessons': ['basics', 'flowcontrol', 'moretypes']
}, {
    'id': 'methods',
    'title': 'Methods and interfaces',
    'description': '<p>Learn how to define methods on types, how to declare interfaces, and how to put everything together.</p>',
    'lessons': ['methods']
}, {
    'id': 'concurrency',
    'title': 'Concurrency',
    'description': '<p>Go provides concurrency features as part of the core language.</p><p>This module goes over goroutines and channels, and how they are used to implement different concurrency patterns.</p>',
    'lessons': ['concurrency']
}]).

// Translation
value('Translation', {
    "off": "off",
    "on": "on",
    "syntax": "Syntax-Highlighting",
    "lineno": "Line-Numbers",
    "reset": "Reset Slide",
    "format": "Format Source Code",
    "kill": "Kill Program",
    "run": "Run",
    "compile": "Compile and Run",
    "more": "Options",
    "toc": "Table of Contents",
    "prev": "Previous",
    "next": "Next",
    "waiting": "Waiting for remote server...",
    "errcomm": "Error communicating with remote server.",
}).

// Config for codemirror plugin
value('ui.config', {
    codemirror: {
        mode: 'text/x-go',
        matchBrackets: true,
        lineNumbers: true,
        autofocus: true,
        indentWithTabs: true,
        lineWrapping: true,
        extraKeys: {
            "Shift-Enter": function() {
                $('#run').click();
            },
            "PageDown": function() {
                return false;
            },
            "PageUp": function() {
                return false;
            },
            "Shift-Space": function() {
                $('#format').click();
            },
        }
    }
});
