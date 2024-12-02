import styles from './styles.module.css';

//Pulled from: https://gaudion.dev/blog/mdx-youtube-embed
//components/mdx/YouTube.tsx
export default function YouTube ({ id } : { id : string }){
    return (
      <div class={styles.videoWrapper}>
        <iframe
          className="aspect-video w-full"
          src={"https://www.youtube.com/embed/" + id}
          title="YouTube Video Player"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
        ></iframe>
      </div>
    );
  };
