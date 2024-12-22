let np = document.getElementById("nav-portrait");
np.style.position = "absolute";
np.style.right = "-" + np.offsetWidth + "px";
function showNavPortrait() {
    np.style.right = 0;
    setTimeout(function () {
        document.addEventListener('click', tf, false);
    }, 50);
}

function tf() {
    np.style.right = "-" + np.offsetWidth + "px";
    document.removeEventListener('click', tf);
}