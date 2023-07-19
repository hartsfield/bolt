package main

import (
	"log"
	"os"
)

func autosplash(name string) {
	wd := "internal/components/"
	createComponent("autosplash")
	insertcomponent("autosplash")
	tmpl_, err := os.OpenFile(wd+"autosplash/autosplash.tmpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Println(err)
	}
	tmpl_.WriteString(`{{ define "autosplash.tmpl" }}
<div class="section-outer autosplash-outer">
  <div class="splash-inner">
    <div class="splash-logo">{{.CompanyName}}</div>
  </div>
</div>
<style>{{ template "autosplash.css" }}</style>
<script>{{ template "autosplash.js"}}</script>
{{end}}`)

	css_, err := os.OpenFile(wd+"autosplash/autosplash.css", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Println(err)
	}
	css_.WriteString(`{{ define "autosplash.css"}}
.splash-inner {
    background-image: url(public/media/` + name + `);
    background-size: cover;
    height: 100vh;
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    align-items: center;
    justify-content: center;
}
.splash-logo {
    font-size: 3em;
    font-weight: bold;
    color: seashell;
    animation: 600ms linear scalein;
}
@keyframes scalein {
  0% {transform: scale(5) rotate(720deg);}
  // 100% {transform: scale(1);}
}
{{end}}`)
}
