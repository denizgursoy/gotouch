const {description, name} = require('../../package.json')

module.exports = {
    base: '/gotouch/',
    /**
     * Ref：https://v1.vuepress.vuejs.org/config/#title
     */
    title: name,
    /**
     * Ref：https://v1.vuepress.vuejs.org/config/#description
     */
    description: description,

    /**
     * Extra tags to be injected to the page HTML `<head>`
     *
     * ref：https://v1.vuepress.vuejs.org/config/#head
     */
    head: [
        ['meta', {name: 'theme-color', content: '#3eaf7c'}],
        ['meta', {name: 'apple-mobile-web-app-capable', content: 'yes'}],
        ['meta', {name: 'apple-mobile-web-app-status-bar-style', content: 'black'}]
    ],

    /**
     * Theme configuration, here is the default theme configuration for VuePress.
     *
     * ref：https://v1.vuepress.vuejs.org/theme/default-theme-config.html
     */
    themeConfig: {
        repo: '',
        editLinks: false,
        docsDir: '',
        editLinkText: '',
        lastUpdated: false,
        nav: [
            {
                text: 'GitHub',
                link: 'https://github.com/denizgursoy/gotouch'
            }
        ],
        sidebar: [
            {
                title: 'Get Started',
                path: '/get-started',
                children: [
                    '/',
                    '/installation.md',
                    '/execute.md',
                    '/how-it-works.md',
                    '/change-log.md',
                ],
            },
            {
                title: 'Customize',
                path: '/customize',
                children: [
                    '/template-project.md',
                    '/template-examples.md',
                    '/properties-yaml.md',
                    '/project-structure.md',
                    '/question.md',
                    '/choice.md',
                    '/dependency.md',
                    '/file.md',
                    '/value.md',
                    '/init.md',
                    '/examples.md',
                ],
            },
            "share.md",
            "commands.md",
        ]
    },

    /**
     * Apply plugins，ref：https://v1.vuepress.vuejs.org/zh/plugin/
     */
    plugins: [
        '@vuepress/plugin-back-to-top',
        '@vuepress/plugin-medium-zoom',
    ],
    configureWebpack: {
        resolve: {
            alias: {
                '@images': '../../images'
            }
        }
    }
}
