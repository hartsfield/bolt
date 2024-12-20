package main

import (
	"log"
	"os"
)

func autosplash(params []string) {
	name := params[0]
	wd := "internal/components/"
	createComponent([]string{"autosplash"})
	insertcomponent([]string{"autosplash"})
	tmpl_, err := os.OpenFile(wd+"autosplash/autosplash.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Println(err)
	}
	tmpl_.WriteString(`<div class="section-outer autosplash-outer">
  <div class="splash-inner">
    <div class="splash-logo">{{.AppName}}</div>
  </div>
<style>{{ template "autosplash.css" }}</style>
<script>{{ template "autosplash.js"}}</script>
</div>
`)

	css_, err := os.OpenFile(wd+"autosplash/autosplash.css", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Println(err)
	}
	css_.WriteString(`.splash-inner {
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
}`)
}
