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
        Spend more time building platform features and less time maintaining
        scripts and text templates.  Add type checking and validation to your
        existing Helm charts and manifests. Define golden paths for teams to
        safely spin up new projects without interrupting you.  Automatically
        manage Namespaces, Certificates, Secrets, Roles, and labels with
        strongly typed CUE definitions. Take unmodified, upstream software and
        mix-in the resources making your platform unique and valuable.
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
        Move faster along paved roads provided by your platform and security
        teams.  Quickly spin up new projects and environments without filing
        tickets.  Develop where it works best for you, locally or in the cloud.
        Deploy your existing Helm charts, Kubernetes manifests, or containers
        safely and confidently with strong type checking. Reduce the friction of
        integrating your services into your organization's platform.
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
        Express your security policies as reusable, well defined, typed
        configuration.   Build guard rails making it easy to develop services
        securely.  Automatically provision and configure SecretStores,
        ExternalSecrets, AuthorizationPolicies, RoleBindings, etc... for project
        teams.  Gain clear visibility into the complete configuration of the
        platform to quickly identify risk and audit your security posture.
        Integrate your preferred security tools with the platform.
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
