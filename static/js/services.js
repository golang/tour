/* Copyright 2012 The Go Authors.   All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
'use strict';

/* Services */

angular.module('tour.services', []).

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
factory('run', ['$window',
    function(win) {
        return function(code, output, options) {
            // PlaygroundOutput is defined in playground.js which is prepended
            // to the generated script.js in gotour/tour.go.
            // The next line removes the jshint warning.
            // global PlaygroundOutput
            win.transport.Run(code, PlaygroundOutput(output), options);
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

// Editor context service, kept through the whole app.
factory('editor', ['$window',
    function(win) {
        var ctx = {
            syntax: false,
            toggleSyntax: function() {
                ctx.syntax = !ctx.syntax;
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
        };
        return ctx;
    }
]).

// Table of contents management and navigation
factory('toc', ['$http', '$q', '$log', 'tableOfContents',
    function($http, $q, $log, tableOfContents) {
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
                    }
                }
                moduleQ.resolve(modules);
                lessonQ.resolve(lessons);
            }, function(error) {
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