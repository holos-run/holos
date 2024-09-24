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
        Focus on building platform features, not maintaining scripts and
        templates. Add type checking to Helm charts and manifests. Define golden
        paths for teams to launch projects independently. Automatically manage
        Namespaces, Certificates, Secrets, and Roles with CUE. Use upstream
        charts and mix in resources that make your platform unique.
        <br />
        <a href="/docs/">Learn More</a>
      </>
    ),
  },
  {
    title: 'For Software Developers',
    Svg: require('@site/static/img/base00/undraw_through_the_park_lxnl.svg').default,
    description: (
      <>
        Move faster with paved paths from your platform and security teams. Spin
        up new projects and environments without tickets. Develop locally or in
        the cloud. Deploy Helm charts, Kubernetes manifests, or containers
        confidently with type checking. Reduce friction when integrating
        services into your organization's platform.
        <br />
        <a href="/docs/">Learn More</a>
      </>
    ),
  },
  {
    title: 'For Security Teams',
    Svg: require('@site/static/img/base00/undraw_security_on_re_e491.svg').default,
    description: (
      <>
        Define security policies as reusable, typed configurations.
        Automatically apply the policy to new projects. Build guardrails for
        secure service development. Get clear visibility into the platform's
        configuration to quickly identify risks and audit security posture.
        Integrate your preferred security tools seamlessly with the platform.
        <br />
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
