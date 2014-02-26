/* Copyright 2012 The Go Authors.   All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
'use strict';

angular.module('tour', ['ui', 'tour.services', 'tour.controllers', 'tour.directives', 'tour.values', 'ng']).

config(['$routeProvider', '$locationProvider',
    function($routeProvider, $locationProvider) {
        $routeProvider.
        when('/', {
            redirectTo: '/welcome/1'
        }).
        when('/list', {
            templateUrl: '/static/partials/list.html',
        }).
        when('/:lessonId/:pageNumber', {
            templateUrl: '/static/partials/editor.html',
            controller: 'EditorCtrl'
        }).
        when('/:lessonId', {
            redirectTo: '/:lessonId/1'
        }).
        otherwise({
            redirectTo: '/'
        });

        $locationProvider.html5Mode(true);
    }
]);
