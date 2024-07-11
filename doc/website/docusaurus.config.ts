import { themes as prismThemes } from 'prism-react-renderer';
import type { Config } from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Holos',
  tagline: 'The Platform Operating System',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://holos.run',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'holos-run', // Usually your GitHub org/user name.
  projectName: 'holos', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'throw',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  // https://docusaurus.io/docs/markdown-features/diagrams
  markdown: {
    mermaid: true
  },
  themes: ['@docusaurus/theme-mermaid'],

  presets: [
    [
      'classic',
      {
        docs: {
          path: "../md",
          // Remove this to remove the "edit this page" links.
          editUrl: 'https://github.com/holos-run/holos/edit/main/doc/md/',
          showLastUpdateAuthor: true,
          showLastUpdateTime: true,
          sidebarPath: './sidebars.ts',
        },
        blog: {
          path: "blog",
          blogSidebarCount: "ALL",
          blogSidebarTitle: "All posts",
          feedOptions: {
            type: 'all',
            copyright: `Copyright © ${new Date().getFullYear()}, The Holos Authors.`,
          },
          showReadingTime: false,
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Replace with your project's social card
    image: 'img/holos-social-card.png',
    docs: {
      sidebar: {
        autoCollapseCategories: false,
      }
    },
    navbar: {
      title: '',
      logo: {
        src: 'img/logo.svg',
        srcDark: 'img/logo-dark.svg',
      },
      items: [
        {
          type: 'doc',
          docId: 'intro',
          position: 'left',
          label: 'Docs',
        },
        {
          type: 'docSidebar',
          sidebarId: 'api',
          position: 'left',
          label: 'API',
        },
        { to: '/blog', label: 'Blog', position: 'left' },
        {
          "href": "https://pkg.go.dev/github.com/holos-run/holos?tab=doc",
          "label": "GoDoc",
          "position": "left",
          "className": "header-godoc-link",
        },
        {
          href: 'https://github.com/holos-run/holos',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Tutorial',
              to: '/docs/intro',
            },
            {
              label: 'API Reference',
              to: '/docs/api/core/v1alpha2',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Discuss',
              href: 'https://github.com/orgs/holos-run/discussions',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'Blog',
              to: '/blog',
            },
            {
              label: 'GitHub',
              href: 'https://github.com/holos-run/holos',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} The Holos Authors.`,
    },
    prism: {
      // Refer to https://docusaurus.io/docs/api/themes/configuration#theme
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      // Refer to https://docusaurus.io/docs/next/markdown-features/code-blocks#supported-languages
      additionalLanguages: ['protobuf', 'cue', 'bash', 'diff', 'json'],
      magicComments: [
        {
          className: 'theme-code-block-highlighted-line',
          line: 'highlight-next-line',
          block: { start: 'highlight-start', end: 'highlight-end' },
        },
        {
          className: 'code-block-error-message',
          line: 'highlight-next-line-error-message',
        },
        {
          className: 'code-block-info-line',
          line: 'highlight-next-line-info',
          block: { start: 'highlight-info-start', end: 'highlight-info-end' },
        },
      ],
    },
    mermaid: {
      // Refer to https://mermaid.js.org/config/theming.html
      theme: { light: 'neutral', dark: 'dark' },
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
