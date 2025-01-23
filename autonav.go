package main

import (
	"fmt"
	"log"
	"os"
)

func autonav(params []string) {
	wd := "internal/components/"
	ex, err := exists(wd + "autonav")
	if err != nil || ex {
		fmt.Println("Component already exists", err)
		os.Exit(0)
	}
	os.MkdirAll(wd+"autonav", 0755)
	tmpl_, err := os.Create(wd + "autonav/autonav.html")
	if err != nil {
		log.Println(err)
	}

	autonav_create("autonav", tmpl_, params)
	insertcomponent([]string{"autonav"})
}

func autonav_create(name string, tmpl_ *os.File, sections []string) {
	wd := "internal/components/"
	var navListHTML string
	for _, section := range sections {
		createComponent([]string{section})
		insertcomponent([]string{section})
		navListHTML = navListHTML + `<li onclick="jumpTo('` +
			section + `-outer')">` + section + `</li>` + "\n"
	}

	tmpl_.WriteString(`<div class="navbar-outer">
  <div class="logo-nav" onclick="window.location = window.location.origin">{{ .AppName }}</div>

  <div class="nav-landscape">
    <ul>
    ` + "\n      " + navListHTML + `
    </ul>
  </div>

  <div class="nav-portrait" id="nav-portrait">
    <div class="nav-portrait-logo">{{ .AppName }}</div>
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
<script>{{ template "` + name + `.js"}}</script>`)
	js_, err := os.Create(wd + name + "/" + name + ".js")
	if err != nil {
		log.Println(err)
	}
	js_.WriteString(
		`let np = document.getElementById("nav-portrait");
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
}`)

	css_, err := os.Create(wd + name + "/" + name + ".css")
	if err != nil {
		log.Println(err)
	}
	css_.WriteString(
		`.ham-outer {
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    justify-content: space-evenly;
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
    font-size: 1.1em;
    position: fixed;
    background: #88d5d8;
    color: white;
    top: 0;
    left: 0;
    width: 100vw;
    flex-direction: row;
    flex-wrap: nowrap;
    justify-content: space-between;
    animation: 0.2s linear navbar-load;;
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
        background-color: #7fc3e1;
        color: #ffffff;
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
    }
    .nav-portrait > ul > li:hover {
        background: #ffffff;
        color: #e16f6f;
    }
    .nav-portrait-logo {
        text-align: center;
        width: 100%;
        margin-bottom: 2em;
        background: #e16f6f;
        padding: 0.7em;
    }
}
@keyframes navbar-load {
    0% {transform: translateY(-3em);}
    90% {transform: translateY(-3em);}
    100% {transform: translateY(0);}
}`)
}
