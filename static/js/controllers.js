/* Copyright 2012 The Go Authors.   All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
'use strict';

/* Controllers */


angular.module('tour.controllers', []).

// Navigation controller
controller('EditorCtrl', ['$scope', '$routeParams', '$location', 'toc', 'i18n', 'run', 'fmt', 'editor', 'analytics', 'storage',
    function($scope, $routeParams, $location, toc, i18n, run, fmt, editor, analytics, storage) {
        var lessons = [];
        toc.lessons.then(function(v) {
            lessons = v;
            $scope.gotoPage($scope.curPage);

            // Store changes on the current file to local storage.
            $scope.$watch(function() {
                var f = file();
                return f && f.Content;
            }, function(val) {
                var key = $scope.lessonId + '.' + ($scope.curPage - 1) + '.' + $scope.curFile;
                storage.set(key, val);
            });
        });

        $scope.toc = toc;
        $scope.lessonId = $routeParams.lessonId;
        $scope.curPage = parseInt($routeParams.pageNumber);
        $scope.curFile = 0;

        $scope.nextPage = function() {
            $scope.gotoPage($scope.curPage + 1);
        };
        $scope.prevPage = function() {
            $scope.gotoPage($scope.curPage - 1);
        };
        $scope.gotoPage = function(page) {
            var l = $routeParams.lessonId;
            if (page >= 1 && page <= lessons[$scope.lessonId].Pages.length) {
                $scope.curPage = page;
            } else {
                l = (page < 1) ? toc.prevLesson(l) : toc.nextLesson(l);
                if (l === '') { // If there's not previous or next
                    $location.path('/list');
                    return;
                }
                page = (page < 1) ? lessons[l].Pages.length : 1;
            }
            $location.path('/' + l + '/' + page);
            $scope.openFile($scope.curFile);
            analytics.trackView();
        };
        $scope.openFile = function(file) {
            $scope.curFile = file;
            editor.paint();
        };

        function log(mode, text) {
            $('.output.active').html('<pre class="' + mode + '">' + text + '</pre>');
        }

        function clearOutput() {
            $('.output.active').html('');
        }

        function file() {
            return lessons[$scope.lessonId].Pages[$scope.curPage - 1].Files[$scope.curFile];
        }

        $scope.run = function() {
            log('info', i18n.l('waiting'));
            var f = file();
            run(f.Content, $('.output.active > pre')[0], {
                path: f.Name
            });
        };

        $scope.format = function() {
            log('info', i18n.l('waiting'));
            fmt(file().Content).then(
                function(data) {
                    if (data.data.Error !== '') {
                        log('stderr', data.data.Error);
                        return;
                    }
                    clearOutput();
                    file().Content = data.data.Body;
                },
                function(error) {
                    log('stderr', error);
                });
        };

        $scope.reset = function() {
            file().Content = file().OrigContent;
        };
    }
]);
