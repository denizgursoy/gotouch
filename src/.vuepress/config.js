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
                children: [
                    '/get-started/',
                    '/get-started/installation.md',
                    '/get-started/execute.md',
                    '/get-started/how-it-works.md',
                    '/get-started/change-log.md',
                ],
            },
            {
                title: 'Customize',
                children: [
                    '/customize',
                    '/customize/template-project.md',
                    '/customize/template-examples.md',
                    '/customize/properties-yaml.md',
                    '/customize/project-structure.md',
                    '/customize/question.md',
                    '/customize/choice.md',
                    '/customize/dependency.md',
                    '/customize/file.md',
                    '/customize/value.md',
                    '/customize/init.md',
                    '/customize/examples.md',
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
