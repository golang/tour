/* Copyright 2012 The Go Authors.   All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
'use strict';

angular.module('tour.values', []).

// List of modules with description and lessons in it.
value('tableOfContents', [{
    'id': 'mechanics',
    'title': 'Používání průvodce',
    'description': '<p>Vítejte u průvodce <a href="https://golang.org">Programovacím jazykem Go</a>. Průvodce pokrývá nejdůležitější vlastnosti jazyka, hlavně:</p>',
    'lessons': ['welcome']
}, {
    'id': 'basics',
    'title': 'Základy',
    'description': '<p>Startovací bod, nauč se všechny základy jazyka</p><p>Deklarace proměnných, volání funkcí, a všechno ostatní co potřebuješ vědět pro postup k dalším lekcím.</p>',
    'lessons': ['basics', 'flowcontrol', 'moretypes']
}, {
    'id': 'methods',
    'title': 'Metody a interface',
    'description': '<p>Nauč se jak definovat metody na typech, jak deklarovat interface, a jak to všechno spojit.</p>',
    'lessons': ['methods']
}, {
    'id': 'concurrency',
    'title': 'Konkurence',
    'description': '<p>Go poskytuje vlastnosti souběžnosti jako jeden ze základů jazyka.</p><p>Tento modul prochází přes go rutiny a kanály a jejich použití v různých vzorech konkurence.</p>',
    'lessons': ['concurrency']
}]).

// translation
value('translation', {
    'off': 'vyp',
    'on': 'zap',
    'syntax': 'Syntax-Highlighting',
    'lineno': 'Line-Numbers',
    'reset': 'Reset Slide',
    'format': 'Naformátuj zdrojový kód',
    'kill': 'Ukonči Program',
    'run': 'Start',
    'compile': 'Zkompiluj a Start',
    'more': 'Volby',
    'toc': 'Obsah',
    'prev': 'Předchozí',
    'next': 'Další',
    'waiting': 'Čekám na vzdálený server...',
    'errcomm': 'Chyba při komunikaci se vzdáleným serverem.',
    'submit-feedback': 'Pošli zpětnou vazbu o téhle stránce',

    // GitHub issue template: update repo and messaging when translating.
    'github-repo': 'github.com/martinkunc/go-tour-cz',
    'issue-title': 'průvodce: [NAHRAĎTE KRÁTKÝM POPISEM]',
    'issue-message': 'Nahoře uveďte titulek a zde poskytněte svůj popis problému, nebo poskytněte zpětnou vazbu. Včetně kódu, pokud je potřeba',
    'context': 'Kontext',
}).

// Config for codemirror plugin
value('ui.config', {
    codemirror: {
        mode: 'text/x-go',
        matchBrackets: true,
        lineNumbers: true,
        autofocus: true,
        indentWithTabs: true,
        indentUnit: 4,
        tabSize: 4,
        lineWrapping: true,
        extraKeys: {
            'Shift-Enter': function() {
                $('#run').click();
            },
            'Ctrl-Enter': function() {
                $('#format').click();
            },
            'PageDown': function() {
                return false;
            },
            'PageUp': function() {
                return false;
            },
        },
        // TODO: is there a better way to do this?
        // AngularJS values can't depend on factories.
        onChange: function() {
            if (window.codeChanged !== null) window.codeChanged();
        }
    }
}).

// mapping from the old paths (#42) to the new organization.
// The values have been generated with the map.sh script in the tools directory.
value('mapping', {
    '#1': '/welcome/1', // Hello, 世界
    '#2': '/welcome/2', // Go local
    '#3': '/basics/1', // Packages
    '#4': '/basics/2', // Imports
    '#5': '/basics/3', // Exported names
    '#6': '/basics/4', // Functions
    '#7': '/basics/5', // Functions continued
    '#8': '/basics/6', // Multiple results
    '#9': undefined, // Named results
    '#10': '/basics/8', // Variables
    '#11': '/basics/9', // Variables with initializers
    '#12': '/basics/10', // Short variable declarations
    '#13': '/basics/11', // Basic types
    '#14': '/basics/13', // Type conversions
    '#15': '/basics/15', // Constants
    '#16': '/basics/16', // Numeric Constants
    '#17': '/flowcontrol/1', // For
    '#18': '/flowcontrol/2', // For continued
    '#19': '/flowcontrol/3', // For is Go's "while"
    '#20': '/flowcontrol/4', // Forever
    '#21': '/flowcontrol/5', // If
    '#22': '/flowcontrol/6', // If with a short statement
    '#23': '/flowcontrol/7', // If and else
    '#24': '/flowcontrol/8', // Exercise: Loops and Functions
    '#25': '/moretypes/2', // Structs
    '#26': '/moretypes/3', // Struct Fields
    '#27': '/moretypes/1', // Pointers
    '#28': '/moretypes/5', // Struct Literals
    '#29': undefined, // The new function
    '#30': '/moretypes/6', // Arrays
    '#31': '/moretypes/7', // Slices
    '#32': '/moretypes/8', // Slicing slices
    '#33': '/moretypes/9', // Making slices
    '#34': '/moretypes/10', // Nil slices
    '#35': '/moretypes/12', // Range
    '#36': '/moretypes/13', // Range continued
    '#37': '/moretypes/14', // Exercise: Slices
    '#38': '/moretypes/15', // Maps
    '#39': '/moretypes/16', // Map literals
    '#40': '/moretypes/17', // Map literals continued
    '#41': '/moretypes/18', // Mutating Maps
    '#42': '/moretypes/19', // Exercise: Maps
    '#43': '/moretypes/20', // Function values
    '#44': '/moretypes/21', // Function closures
    '#45': '/moretypes/22', // Exercise: Fibonacci closure
    '#46': '/flowcontrol/9', // Switch
    '#47': '/flowcontrol/10', // Switch evaluation order
    '#48': '/flowcontrol/11', // Switch with no condition
    '#49': undefined, // Advanced Exercise: Complex cube roots
    '#50': undefined, // Methods and Interfaces
    '#51': '/methods/1', // Methods
    '#52': '/methods/2', // Methods continued
    '#53': '/methods/3', // Methods with pointer receivers
    '#54': '/methods/4', // Interfaces
    '#55': '/methods/5', // Interfaces are satisfied implicitly
    '#56': '/methods/8', // Errors
    '#57': '/methods/9', // Exercise: Errors
    '#58': '/methods/13', // Web servers
    '#59': '/methods/14', // Exercise: HTTP Handlers
    '#60': '/methods/15', // Images
    '#61': '/methods/16', // Exercise: Images
    '#62': undefined, // Exercise: Rot13 Reader
    '#63': undefined, // Concurrency
    '#64': '/concurrency/1', // Goroutines
    '#65': '/concurrency/2', // Channels
    '#66': '/concurrency/3', // Buffered Channels
    '#67': '/concurrency/4', // Range and Close
    '#68': '/concurrency/5', // Select
    '#69': '/concurrency/6', // Default Selection
    '#70': '/concurrency/7', // Exercise: Equivalent Binary Trees
    '#71': '/concurrency/8', // Exercise: Equivalent Binary Trees
    '#72': '/concurrency/9', // Exercise: Web Crawler
    '#73': '/concurrency/10', // Where to Go from here...
});
