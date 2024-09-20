import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <p className="projectDesc">
          Building and operating an internal development platform is a
          challenge.  Engineering teams glue their tools to the rest.  This glue
          is sticky and hard to work with, frustrating engineering teams.
        </p>
        <p className="projectDesc">
          Holos is a universal adapter that replaces the glue in your platform
          without replacing your tools.  Holos leverages the safe, consistent,
          and well-defined structure of CUE to solve the problem of integrating
          bespoke tools and processes into a centrally managed platform.
        </p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="docs/quickstart">
            Get Started
          </Link>
          <span className={styles.divider}></span>
          <Link
            className="button button--primary button--lg"
            to="docs/">
            Learn More
          </Link>
          <span className={styles.divider}></span>
        </div>
      </div >
    </header >
  );
}

export default function Home(): JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} Platform Manager`}
      description="Holos adds CUE's type safety, unified structure, and strong validation features to your Kubernetes configuration manifests, including Helm and Kustomize.">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
