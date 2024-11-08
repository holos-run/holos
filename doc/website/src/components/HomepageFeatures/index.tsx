import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

// We don't focus on features, but rather problems and solutions.
const FeatureList: FeatureItem[] = [
  {
    title: 'For Platform Engineers',
    Svg: require('@site/static/img/base00/undraw_software_engineer_re_tnjc.svg').default,
    description: (
      <>
        <p align="left">
          <ul>
            <li>Provide simple definitions for other teams to use as golden paths.</li>
            <li>Define integrations in <a href="https://cuelang.org/">CUE</a> with strong type checking. No more text templates or bash scripts.</li>
            <li>Reuse your existing Helm charts and Kustomize bases.</li>
          </ul>
        </p>
        <a href="/docs/">Learn More</a>
      </>
    ),
  },
  {
    title: 'For Software Developers',
    Svg: require('@site/static/img/base00/undraw_through_the_park_lxnl.svg').default,
    description: (
      <>
        <p align="left">
          <ul>
            <li>Move faster using paved paths from your platform and security teams.</li>
            <li>Develop locally or in the cloud.</li>
            <li>Spend more time developing software and fewer cycles fighting infrastructure challenges.</li>
          </ul>
        </p>
        <a href="/docs/">Learn More</a>
      </>
    ),
  },
  {
    title: 'For Security Teams',
    Svg: require('@site/static/img/base00/undraw_security_on_re_e491.svg').default,
    description: (
      <>
        <p align="left">
          <ul>
            <li>Define security policy as reusable, typed configurations.</li>
            <li>Automatically enforce security policy on new projects.</li>
            <li>Ensure a consistent security posture cross-platform with fewer code changes.</li>
          </ul>
        </p>
        <a href="/docs/">Learn More</a>
      </>
    ),
  }
];

function Feature({ title, Svg, description }: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
