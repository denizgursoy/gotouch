import { defineConfig } from 'vitepress'

export default defineConfig({
  base: '/gotouch/',
  title: 'Gotouch',
  description: 'Customizable Project Creator',

  head: [
    ['meta', { name: 'theme-color', content: '#3eaf7c' }],
    ['link', { rel: 'icon', href: '/gotouch/icon.png' }],
  ],

  themeConfig: {
    nav: [
      {
        text: 'GitHub',
        link: 'https://github.com/denizgursoy/gotouch'
      }
    ],

    sidebar: [
      {
        text: 'Get Started',
        items: [
          { text: 'What is Gotouch?', link: '/get-started/' },
          { text: 'Installation', link: '/get-started/installation' },
          { text: 'Execute Gotouch', link: '/get-started/execute' },
          { text: 'How It Works', link: '/get-started/how-it-works' },
          { text: 'Change Log', link: '/get-started/change-log' },
        ],
      },
      {
        text: 'Customize',
        items: [
          { text: 'Start Customizing', link: '/customize/' },
          { text: 'Template Project', link: '/customize/template-project' },
          { text: 'Learn Go Template', link: '/customize/learn-go-template' },
          { text: 'Properties YAML', link: '/customize/properties-yaml' },
          { text: 'Project Structure', link: '/customize/project-structure' },
          { text: 'Question', link: '/customize/question' },
          { text: 'Choice', link: '/customize/choice' },
          { text: 'Dependencies', link: '/customize/dependencies' },
          { text: 'Files', link: '/customize/files' },
          { text: 'Values', link: '/customize/values' },
          { text: 'Init Files', link: '/customize/init' },
          { text: 'Examples', link: '/customize/examples' },
          { text: 'Distributing Templates with Docker', link: '/customize/local-path-docker-example' },
        ],
      },
      {
        text: 'Features',
        items: [
          { text: 'Overview', link: '/features/' },
          { text: 'Authentication', link: '/features/authentication' },
        ],
      },
      { text: 'Share', link: '/share' },
      { text: 'Commands', link: '/commands' },
    ],
  },
})
