import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Zero Trust Authorization',
    Svg: require('@site/static/img/base00/undraw_security_on_re_e491.svg').default,
    description: (
      <>
        Spend more time on your business features and less time rebuilding
        authentication and authorization.  Holos provides zero trust security
        with no code necessary.
      </>
    ),
  },
  {
    title: 'Multi-Cloud',
    Svg: require('@site/static/img/base00/undraw_cloud_hosting_7xb1.svg').default,
    description: (
      <>
        Avoid vendor lock in, downtime, and price hikes.  Holos is designed to
        easily deploy workloads into multiple clouds and multiple regions.
      </>
    ),
  },
  {
    title: 'Observability',
    Svg: require('@site/static/img/base00/undraw_data_trends_re_2cdy.svg').default,
    description: (
      <>
        Receive alerts before issues impact your customers.  Holos provides an
        integrated observability system complete with dashboards and
        notification rules.
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
