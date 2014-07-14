/* Copyright 2012 The Go Authors.   All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
'use strict';

/* Directives */

angular.module('tour.directives', []).

// onpageup executes the given expression when Page Up is released.
directive('onpageup', function() {
    return function(scope, elm, attrs) {
        elm.attr('tabindex', 0);
        elm.keyup(function(evt) {
            var key = evt.key || evt.keyCode;
            if (key == 33) {
                scope.$apply(attrs.onpageup);
                evt.preventDefault();
            }
        });
    };
}).

// onpagedown executes the given expression when Page Down is released.
directive('onpagedown', function() {
    return function(scope, elm, attrs) {
        elm.attr('tabindex', 0);
        elm.keyup(function(evt) {
            var key = evt.key || evt.keyCode;
            if (key == 34) {
                scope.$apply(attrs.onpagedown);
                evt.preventDefault();
            }
        });
    };
}).

// autofocus sets the focus on the given element when the condition is true.
directive('autofocus', function() {
    return function(scope, elm, attrs) {
        elm.attr('tabindex', 0);
        scope.$watch(function() {
            return scope.$eval(attrs.autofocus);
        }, function(val) {
            if (val === true) $(elm).focus();
        });
    };
}).

// syntax-checkbox activates and deactivates
directive('syntaxCheckbox', ['editor',
    function(editor) {
        return function(scope, elm) {
            elm.click(function() {
                editor.toggleSyntax();
                scope.$digest();
            });
            scope.editor = editor;
        };
    }
]).

// verticalSlide creates a sliding separator between the left and right elements.
// e.g.:
// <div id="header">Some content</div>
// <div vertical-slide top="#header" bottom="#footer"></div>
// <div id="footer">Some footer</div>
directive('verticalSlide', ['editor',
    function(editor) {
        return function(scope, elm, attrs) {
            var moveTo = function(x) {
                if (x < 0) {
                    x = 0;
                }
                if (x > $(window).width()) {
                    x = $(window).width();
                }
                elm.css('left', x);
                $(attrs.left).width(x);
                $(attrs.right).offset({
                    left: x
                });
                editor.x = x;
            };

            elm.draggable({
                axis: 'x',
                drag: function(event) {
                    moveTo(event.clientX);
                    return true;
                },
                containment: 'parent',
            });

            if (editor.x !== undefined) {
                moveTo(editor.x);
            }
        };
    }
]).

// horizontalSlide creates a sliding separator between the top and bottom elements.
// <div id="menu">Some menu</div>
// <div vertical-slide left="#menu" bottom="#content"></div>
// <div id="content">Some content</div>
directive('horizontalSlide', ['editor',
    function(editor) {
        return function(scope, elm, attrs) {
            var moveTo = function(y) {
                var top = $(attrs.top).offset().top;
                if (y < top) {
                    y = top;
                }
                elm.css('top', y - top);
                $(attrs.top).height(y - top);
                $(attrs.bottom).offset({
                    top: y,
                    height: 0
                });
                editor.y = y;
            };
            elm.draggable({
                axis: 'y',
                drag: function(event) {
                    moveTo(event.clientY);
                    return true;
                },
                containment: 'parent',
            });

            if (editor.y !== undefined) {
                moveTo(editor.y);
            }
        };
    }
]).

directive('tableOfContentsButton', function() {
    var speed = 250;
    return {
        restrict: 'A',
        templateUrl: '/static/partials/toc-button.html',
        link: function(scope, elm, attrs) {
            elm.on('click', function() {
                var toc = $(attrs.tableOfContentsButton);
                // hide all non active lessons before displaying the toc.
                var visible = toc.css('display') != 'none';
                if (!visible) {
                    toc.find('.toc-lesson:not(.active) .toc-page').hide();
                    toc.find('.toc-lesson.active .toc-page').show();
                }
                toc.toggle('slide', {
                    direction: 'right'
                }, speed);

                // if fullscreen hide the rest of the content when showing the atoc.
                var fullScreen = toc.width() == $(window).width();
                if (fullScreen) $('#editor-container')[visible ? 'show' : 'hide']();
            });
        }
    };
}).

// side bar with dynamic table of contents
directive('tableOfContents', ['$routeParams', 'toc',
    function($routeParams, toc) {
        var speed = 250;
        return {
            restrict: 'A',
            templateUrl: '/static/partials/toc.html',
            link: function(scope, elm) {
                scope.toc = toc;
                scope.params = $routeParams;

                scope.toggleLesson = function(id) {
                    var l = $('#toc-l-' + id + ' .toc-page');
                    l[l.css('display') == 'none' ? 'slideDown' : 'slideUp']();
                };

                scope.$watch(function() {
                    return scope.params.lessonId + scope.params.lessonId;
                }, function() {
                    $('.toc-lesson:not(#toc-l-' + scope.params.lessonId + ') .toc-page').slideUp(speed);
                });

                scope.hideTOC = function(fullScreenOnly) {
                    var fullScreen = elm.find('.toc').width() == $(window).width();
                    if (fullScreenOnly && !fullScreen) {
                        return;
                    }
                    $('.toc').toggle('slide', {
                        direction: 'right'
                    }, speed);
                };
            }
        };
    }
]);