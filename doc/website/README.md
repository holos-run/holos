# Website

This website is built using [Docusaurus](https://docusaurus.io/), a modern
static website generator.

## Installation

```shell
npm install
```

### Local Development

```shell
npm run start
```

This command starts a local development server and opens up a browser window.
Most changes are reflected live without having to restart the server.

### Build

```shell
npm run build
```

This command generates static content into the `build` directory and can be
served using any static contents hosting service.

### Deployment

Deployments are made with Cloudflare Pages. Cloudflare deploys on changes to
the main branch, and Pull Requests get comments with links to preview
environments.