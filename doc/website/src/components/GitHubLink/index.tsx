export default function GitHubLink({ repo, tree, path, children }: { repo: string, commit: string, path: string, children: React.ReactNode }) {
  const href = `https://github.com/${repo}/tree/${tree}/${path}`
  return (
    <a href={href} target="_blank" rel="noopener noreferrer">{children}</a>
  );
};
