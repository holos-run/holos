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
    'quickstart/index',
    'concepts',
    'install',
    'comparison',
  ],
  api: [
    {
      label: 'Core API',
      type: 'category',
      collapsed: true,
      items: [
        'api/core/v1alpha3',
        'api/core/v1alpha2',
      ],
    },
    'cli',
  ],
};

export default sidebars;
