import DefaultTheme from 'vitepress/theme'
import './custom.css'
import Tabs from './Tabs.vue'

export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    app.component('Tabs', Tabs)
  }
}
