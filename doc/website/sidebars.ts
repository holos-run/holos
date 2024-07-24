import type { SidebarsConfig } from '@docusaurus/plugin-content-docs';

/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */
const sidebars: SidebarsConfig = {
  doc: [
    'intro',
    {
      type: 'category',
      label: 'Guides',
      collapsed: false,
      items: [
        'guides/install',
        'guides/try-holos/index',
        'guides/try-holos/platform-manifests',
        'guides/argocd/index',
        'guides/backstage/index',
        'guides/observability/index',
      ],
    },
    {
      type: 'category',
      label: 'Design',
      collapsed: false,
      items: [
        'design/rendering',
      ],
    },
    {
      type: 'category',
      label: 'Reference Platform',
      collapsed: false,
      items: [
        'reference-platform/architecture',
      ],
    },
    'glossary',
  ],
  api: [
    'api/core/v1alpha2',
    'cli',
  ],
};

export default sidebars;
