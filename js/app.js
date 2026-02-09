import Alpine from 'alpinejs'
import ajax from '@imacrayon/alpine-ajax'
window.htmx = require('htmx.org');

window.Alpine = Alpine
Alpine.plugin(ajax)


Alpine.store('app', {
    sidebarOpen: true,
    
})

console.log('App initialized test')

Alpine.start()