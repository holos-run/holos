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
      label: 'Tutorial',
      collapsed: false,
      items: [
        'tutorial/start',
        'tutorial/register',
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
  ],
  api: [
    'cli',
    'api/core/v1alpha2'
  ],
};

export default sidebars;
