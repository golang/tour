/* Copyright 2012 The Go Authors.   All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
'use strict';

/* Services */

angular.module('tour.services', []).

// Google Analytics
factory('analytics', ['$window',
    function(win) {
        var track = win.trackPageview || (function() {});
        return {
            trackView: track
        };
    }
]).

// Internationalization
factory('i18n', ['translation',
    function(translation) {
        return {
            l: function(key) {
                if (translation[key]) return translation[key];
                return '(no translation for ' + key + ')';
            }
        };
    }
]).

// Running code
factory('run', ['$window', 'editor',
    function(win, editor) {
        var writeInterceptor = function(writer) {
            return function(write) {
                if (write.Kind == 'stderr') {
                    var lines = write.Body.split('\n');
                    for (var i in lines) {
                        var match = lines[i].match(/.*\.go:([0-9]+): ([^\n]*)/);
                        if (match !== null) {
                            editor.highlight(match[1], match[2]);
                        }
                    }
                }
                writer(write);
            };
        };
        return function(code, output, options) {
            // PlaygroundOutput is defined in playground.js which is prepended
            // to the generated script.js in gotour/tour.go.
            // The next line removes the jshint warning.
            // global PlaygroundOutput
            win.transport.Run(code, writeInterceptor(PlaygroundOutput(output)), options);
        };
    }
]).

// Formatting code
factory('fmt', ['$http',
    function($http) {
        return function(body) {
            var params = $.param({
                'body': body
            });
            var headers = {
                'Content-Type': 'application/x-www-form-urlencoded'
            };
            return $http.post('/fmt', params, {
                headers: headers
            });
        };
    }
]).

// Local storage, persistent to page refreshing.
factory('storage', ['$window',
    function(win) {
        if (win.localStorage) {
            return {
                get: function(key) {
                    return win.localStorage.getItem(key);
                },
                set: function(key, val) {
                    win.localStorage.setItem(key, val);
                }
            };
        }
        return {
            get: function() {
                return null;
            },
            set: function() {}
        };
    }
]).

// Editor context service, kept through the whole app.
factory('editor', ['$window', 'storage',
    function(win, storage) {
        var ctx = {
            syntax: storage.get('syntax') === 'true',
            toggleSyntax: function() {
                ctx.syntax = !ctx.syntax;
                storage.set('syntax', ctx.syntax);
                ctx.paint();
            },
            paint: function() {
                var mode = ctx.syntax && 'text/x-go' || 'text/x-go-comment';
                // Wait for codemirror to start.
                var set = function() {
                    if ($('.CodeMirror').length > 0) {
                        var cm = $('.CodeMirror')[0].CodeMirror;
                        if (cm.getOption('mode') == mode) {
                            cm.refresh();
                            return;
                        }
                        cm.setOption('mode', mode);
                    }
                    win.setTimeout(set, 10);
                };
                set();
            },
            highlight: function(line, message) {
                $('.CodeMirror-code > div:nth-child(' + line + ')')
                    .addClass('line-error').attr('title', message);
            },
            onChange: function() {
                $('.line-error').removeClass('line-error').attr('title', null);
            }
        };
        // Set in the window so the onChange function in the codemirror config
        // can call it.
        win.codeChanged = ctx.onChange;
        return ctx;
    }
]).

// Table of contents management and navigation
factory('toc', ['$http', '$q', '$log', 'tableOfContents', 'storage',
    function($http, $q, $log, tableOfContents, storage) {
        var modules = tableOfContents;

        var lessons = {};

        var prevLesson = function(id) {
            var mod = lessons[id].module;
            var idx = mod.lessons.indexOf(id);
            if (idx < 0) return '';
            if (idx > 0) return mod.lessons[idx - 1];

            idx = modules.indexOf(mod);
            if (idx <= 0) return '';
            mod = modules[idx - 1];
            return mod.lessons[mod.lessons.length - 1];
        };

        var nextLesson = function(id) {
            var mod = lessons[id].module;
            var idx = mod.lessons.indexOf(id);
            if (idx < 0) return '';
            if (idx + 1 < mod.lessons.length) return mod.lessons[idx + 1];

            idx = modules.indexOf(mod);
            if (idx < 0 || modules.length <= idx + 1) return '';
            mod = modules[idx + 1];
            return mod.lessons[0];
        };

        $http.get('/lesson/').then(
            function(data) {
                lessons = data.data;
                for (var m = 0; m < modules.length; m++) {
                    var module = modules[m];
                    module.lesson = {};
                    for (var l = 0; l < modules[m].lessons.length; l++) {
                        var lessonName = module.lessons[l];
                        var lesson = lessons[lessonName];
                        lesson.module = module;
                        module.lesson[lessonName] = lesson;

                        // replace file contents with locally stored copies.
                        for (var p = 0; p < lesson.Pages.length; p++) {
                            var page = lesson.Pages[p];
                            for (var f = 0; f < page.Files.length; f++) {
                                page.Files[f].OrigContent = page.Files[f].Content;
                                var val = storage.get(lessonName + '.' + p + '.' + f);
                                if (val !== null) {
                                    page.Files[f].Content = val;
                                }
                            }
                        }
                    }
                }
                moduleQ.resolve(modules);
                lessonQ.resolve(lessons);
            },
            function(error) {
                $log.error('error loading lessons : ', error);
                moduleQ.reject(error);
                lessonQ.reject(error);
            }
        );

        var moduleQ = $q.defer();
        var lessonQ = $q.defer();

        return {
            modules: moduleQ.promise,
            lessons: lessonQ.promise,
            prevLesson: prevLesson,
            nextLesson: nextLesson
        };
    }
]);
