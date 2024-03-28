// @ts-check
// `@type` JSDoc annotations allow editor autocompletion and type checking
// (when paired with `@ts-check`).
// There are various equivalent ways to declare your Docusaurus config.
// See: https://docusaurus.io/docs/api/docusaurus-config

import {themes as prismThemes} from 'prism-react-renderer';

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'SpecPipe',
  tagline: 'Data Pipeline for Spectrum',
  favicon: 'img/favicon.ico',
  url: 'https://ml4wireless.github.io/',
  baseUrl: '/specpipe/',

  // Github pages configuration
  organizationName: 'ml4wireless',
  projectName: 'specpipe',
  deploymentBranch: 'gh-pages',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          routeBasePath: '/',
          sidebarPath: './sidebars.js'
        },
        blog: false
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/docusaurus-social-card.jpg',
      navbar: {
        title: 'SpecPipe',
        logo: {
          alt: 'My Site Logo',
          src: 'img/logo.svg',
        },
        items: [
          {
            href: 'https://github.com/ml4wireless/specpipe',
            label: 'SpecPipe Repo',
            position: 'right',
          },
          {
            href: 'https://github.com/ml4wireless/specpipe-sdk-py',
            label: 'SpecPipe Python SDK Repo',
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
                label: 'Home',
                to: '/'
              }
            ],
          },
          {
            title: 'More',
            items: [
              {
                href: 'https://github.com/ml4wireless/specpipe',
                label: 'SpecPipe Repo'
              },
              {
                href: 'https://github.com/ml4wireless/specpipe-sdk-py',
                label: 'SpecPipe Python SDK Repo'
              }
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} ml4wireles.`,
      },
      prism: {
        theme: prismThemes.github,
        darkTheme: prismThemes.dracula,
      },
    }),
};

export default config;
