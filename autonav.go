package main

import (
	"fmt"
	"log"
	"os"
)

func autonav(name string) {
	wd := "internal/components/"
	ex, err := exists(wd + name)
	if err != nil || ex {
		fmt.Println("Component already exists", err)
		os.Exit(0)
	}
	os.MkdirAll(wd+name, 0755)
	tmpl_, err := os.Create(wd + name + "/" + name + ".tmpl")
	if err != nil {
		log.Println(err)
	}

	autonav_create(name, tmpl_, []string{"test1"})
}

func autonav_create(name string, tmpl_ *os.File, sections []string) {
	wd := "internal/components/"
	var navListHTML string
	for _, section := range sections {
		navListHTML = navListHTML + `<li onclick="jumpTo('section-` +
			section + `')">` + section + `</li>` + "\n"
	}

	tmpl_.WriteString(`{{ define "` + name + `.tmpl" }}
<div class="navbar-outer">
  <div class="logo-nav" onclick="window.location = window.location.origin">{{ .CompanyName }}</div>

  <div class="nav-landscape">
    <ul>
    ` + "\n      " + navListHTML + `
    </ul>
  </div>

  <div class="nav-portrait" id="nav-portrait">
    <div class="nav-portrait-logo">{{ .CompanyName }}</div>
    <ul>
    ` + "\n" + navListHTML + `
    </ul>
  </div>

  <div class="ham-outer" onclick="showNavPortrait()">
    <div class="hamburger"></div>
    <div class="hamburger"></div>
    <div class="hamburger"></div>
  </div>

</div>

<style>{{ template "` + name + `.css" }}</style>
<script>{{ template "` + name + `.js"}}</script>
{{end}}`)
	js_, err := os.Create(wd + name + "/" + name + ".js")
	if err != nil {
		log.Println(err)
	}
	js_.WriteString(
		`{{ define "` + name + `.js"}}
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
{{end}}`)

	css_, err := os.Create(wd + name + "/" + name + ".css")
	if err != nil {
		log.Println(err)
	}
	css_.WriteString(
		`{{ define "` + name + `.css"}}
.ham-outer {
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    align-content: space-around;
    justify-content: space-evenly;
    align-items: stretch;
    margin-right: 2em;
    padding: 0.5em;
}
.hamburger {
    background-color: seashell;
    height: 0.3em;
    width: 2.5em;
}


.navbar-outer {
    display: inline-flex;
    font-weight: bold;
    font-size: 1.1em;
    /* padding-bottom: 1em; */
    position: fixed;
    background: chocolate;
    top: 0;
    left: 0;
    width: 100vw;
    flex-direction: row;
    flex-wrap: nowrap;
    /* align-content: space-between; */
    justify-content: space-between;
    animation: 0.9s linear navbar-load;;
    z-index: 101010101010;
}

.logo-nav {
    padding: 1em;
    cursor: pointer;
}
@media screen and (orientation:landscape) {
    .ham-outer {
        display: none;
    }
    .nav-landscape > ul {
        display: inline-flex;
        list-style-type: none;
        margin-right: 2em;
    }
    .nav-landscape > ul > li {
        margin: 0 1.5em 0 0;
        padding: 1em;
        cursor: pointer;
    }
    .nav-landscape > ul > li:hover {
        background-color: #919b0b;
    }
    .nav-portrait {
        display: none;
    }
    .navbar-outer {
        display: inline-flex;
        flex-wrap: nowrap;
        justify-content: space-between;
        /* align-items: stretch; */
        /* flex-direction: row; */
    }
}
@media screen and (orientation:portrait) {
    .navbar-outer {
        padding-left: 1em;
    }
    .logo-nav {
        max-width: 75%;
        line-height: 1.5em;
    }
    .nav-landscape {
        display: none;
    }
    .nav-portrait {
        position: absolute;
        background-color: #622b6bc7;
        color: seashell;
        height: 100vh;
        padding: 1em 2em 1em 1em;
        font-size: 1.5em;
    }
    .nav-portrait > ul {
        list-style-type: none;
    }
    .nav-portrait > ul > li {
        margin-top: 1em;
        text-align: right;
        cursor: pointer;
        padding: 0.2em;
        color: #ffd0ec;
    }
    .nav-portrait > ul > li:hover {
        background: #c96100;
    }
    .nav-portrait-logo {
        text-align: center;
        width: 100%;
        margin-bottom: 2em;
        background: #c96100;
    }
}
@keyframes navbar-load {
    0% {transform: translateY(-3em);}
    90% {transform: translateY(-3em);}
    100% {transform: translateY(0);}
}
{{end}}`)
}
