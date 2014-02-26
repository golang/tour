/* Copyright 2012 The Go Authors.   All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
'use strict';

/* Controllers */


angular.module('tour.controllers', []).

// Navigation controller
controller('EditorCtrl', ['$scope', '$routeParams', '$location', 'TOC', 'I18n', 'Run', 'Fmt',
    function($scope, $routeParams, $location, TOC, I18n, Run, Fmt) {
        var lessons = [];
        TOC.lessons.then(function(v) {
            lessons = v;
        })

        $scope.TOC = TOC;
        $scope.lessonId = $routeParams.lessonId;
        $scope.curPage = parseInt($routeParams.pageNumber);
        $scope.curFile = 0;

        $scope.nextPage = function() {
            $scope.gotoPage($scope.curPage + 1);
        }
        $scope.prevPage = function() {
            $scope.gotoPage($scope.curPage - 1);
        }
        $scope.gotoPage = function(page) {
            var l = $routeParams.lessonId;
            if (page >= 1 && page <= lessons[$scope.lessonId].Pages.length) {
                $scope.curPage = page;
            } else {
                l = (page < 1) ? TOC.prevLesson(l) : TOC.nextLesson(l);
                if (l == '') { // If there's not previous or next
                    $location.path('/list');
                    return;
                }
                page = (page < 1) ? lessons[l].Pages.length : 1;
            }
            $location.path('/' + l + '/' + page);
        }

        function log(mode, text) {
            $(".output.active").html('<pre class="' + mode + '">' + text + '</pre>');
        }

        function clearOutput() {
            $(".output.active").html('');
        }

        function file() {
            return lessons[$scope.lessonId].Pages[$scope.curPage - 1].Files[$scope.curFile];
        }

        $scope.run = function() {
            log('info', I18n.L('waiting'));
            Run(file().Content, $(".output.active > pre")[0]);
        };

        $scope.format = function() {
            log('info', I18n.L('waiting'));
            Fmt(file().Content).then(
                function(data) {
                    if (data.data.Error != '') {
                        log('stderr', data.data.Error);
                        return
                    }
                    clearOutput();
                    file().Content = data.data.Body;
                },
                function(error) {
                    log('stderr', error);
                });
        };
    }
]);
