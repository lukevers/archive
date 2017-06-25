function toggle(on, c) {
    [].slice.call(document.getElementsByClassName(on)).map(function(e) {
        e.classList.toggle(c);
    });
}
