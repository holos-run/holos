import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

// TODO: Consider focusing on the three pillars of Safe, Easy, Consistent.
const FeatureList: FeatureItem[] = [
  {
    title: 'Kustomize Helm',
    Svg: require('@site/static/img/base00/undraw_together_re_a8x4.svg').default,
    description: (
      <>
        Super charge your existing Helm charts by providing well defined,
        validated input values, post-processing the output with Kustomize,
        and mixing in your own custom resources.  All without modifying upstream
        charts to alleviate the pain of upgrades.
      </>
    ),
  },
  {
    title: 'Unified Data Model',
    Svg: require('@site/static/img/base00/undraw_fitting_pieces_re_nss7.svg').default,
    description: (
      <>
        Unify all of your platform components into one well-defined, strongly
        typed data model with CUE.  Holos makes it easier and safer to integrate
        seamlessly with software distributed with current and future tools that
        produce Kubernetes resource manifests.
      </>
    ),
  },
  {
    title: 'Deep Insights',
    Svg: require('@site/static/img/base00/undraw_code_review_re_woeb.svg').default,
    description: (
      <>
        Reduce risk and increase confidence in your configuration changes.
        Holos offers clear visibility into platform wide changes before the
        change is made.
      </>
    ),
  },
  {
    title: 'Interoperable',
    Svg: require('@site/static/img/base00/undraw_version_control_re_mg66.svg').default,
    description: (
      <>
        Holos is designed for compatibility with your preferred tools and
        processes, for example git diff and code reviews.
      </>
    ),
  },
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
