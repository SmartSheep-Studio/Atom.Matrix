import "./assets/style.css"

import { createApp } from "vue"
import { createPinia } from "pinia"

import wrapper from "./app.vue"
import router from "./router"
import i18n from "./i18n"

import "vfonts/Lato.css"
import "vfonts/FiraCode.css"

import VueMarkdownEditor from "@kangc/v-md-editor"
import "@kangc/v-md-editor/lib/style/base-editor.css"
import githubTheme from "@kangc/v-md-editor/lib/theme/github.js"
import "@kangc/v-md-editor/lib/theme/style/github.css"

import Prism from "prismjs"

VueMarkdownEditor.use(githubTheme, {
  Prism,
})

const app = createApp(wrapper)

app.use(VueMarkdownEditor)
app.use(createPinia())
app.use(i18n)
app.use(router)

app.mount("#app")
