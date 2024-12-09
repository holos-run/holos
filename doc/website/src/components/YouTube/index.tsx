import styles from './styles.module.css';

//Pulled from: https://gaudion.dev/blog/mdx-youtube-embed
//components/mdx/YouTube.tsx
export default function YouTube({ id }: { id: string }) {
  return (
    <div className={styles.videoWrapper}>
      <iframe
        className="aspect-video w-full"
        src={"https://www.youtube.com/embed/" + id + "?rel=0"}
        title="YouTube Video Player"
        allow="picture-in-picture; fullscreen; accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope;"
      ></iframe>
    </div>
  );
};
