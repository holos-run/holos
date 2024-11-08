import { themes as prismThemes } from 'prism-react-renderer';
import type { Config } from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Holos',
  tagline: 'An easier way for platform teams to integrate software into their platform',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://holos.run',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',
  // trailing slash is necessary for Cloudflare pages.
  trailingSlash: true,

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

  // TODO: These redirects don't seem to be working, at least with the `npm run
  // start` dev server.
  plugins: [
    [
      '@docusaurus/plugin-client-redirects',
      {
        redirects: [],
      },
    ],
  ],

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
          // https://docusaurus.io/docs/versioning#configuring-versioning-behavior
          // lastVersion: 'current',
          // versions: {
          //   current: {
          //     label: 'v1alpha6',
          //     path: 'v1alpha6',
          //   }
          // }
        },
        blog: {
          path: "blog",
          blogSidebarCount: "ALL",
          blogSidebarTitle: "All posts",
          feedOptions: {
            type: 'all',
            copyright: `Copyright © ${new Date().getFullYear()}, The Holos Authors`,
          },
          showReadingTime: false,
        },
        theme: {
          customCss: './src/css/custom.css',
        },
        gtag: {
          trackingID: 'G-M00QMB1N05',
          anonymizeIP: true,
        }
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
        { to: '/docs', label: 'Docs', position: 'left' },
        { to: '/blog', label: 'Blog', position: 'left' },
        {
          href: 'https://github.com/holos-run',
          label: 'GitHub',
          position: 'right',
        },
        {
          href: 'https://discord.gg/JgDVbNpye7',
          label: 'Discord',
          position: 'right',
        },
        {
          type: 'docsVersionDropdown',
          position: 'right',
          // dropdownItemsAfter: [{ to: '/versions', label: 'All versions' }],
          dropdownActiveClassDisabled: true,
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
              to: '/docs/tutorial/overview',
            },
            {
              label: 'Topics',
              to: '/docs/topics',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Support',
              href: '/docs/support',
            },
            {
              label: 'Discord',
              href: 'https://discord.gg/JgDVbNpye7',
            },
            {
              label: 'Discussion List',
              href: 'https://groups.google.com/g/holos-discuss',
            },
            {
              label: 'Discussion Forum',
              href: 'https://github.com/holos-run/holos/discussions',
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
            {
              label: 'GoDoc',
              href: 'https://pkg.go.dev/github.com/holos-run/holos?tab=doc',
            }
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
