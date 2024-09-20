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
    title: 'Faster',
    Svg: require('@site/static/img/base00/undraw_together_re_a8x4.svg').default,
    description: (
      <>
        Cut your time to market.  Something about developer productivity.
      </>
    ),
  },
  {
    title: 'Safer',
    Svg: require('@site/static/img/base00/undraw_fitting_pieces_re_nss7.svg').default,
    description: (
      <>
        Replace manual tasks with workflows that are well structured and strongly typed.
      </>
    ),
  },
  {
    title: 'Secure',
    Svg: require('@site/static/img/base00/undraw_code_review_re_woeb.svg').default,
    description: (
      <>
        Empower your security team to pave smooth roads for dev teams to deploy
        services securely.
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
