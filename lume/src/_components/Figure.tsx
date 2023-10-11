export interface FigureProps {
  className?: string;
  path: string;
  desc?: string;
  alt?: string;
}

export default function Figure({ className, path, alt, desc = alt}: FigureProps) {
  return (
    <figure className={`max-w-3xl mx-auto ${className}`}>
      <img src={`https://cdn.xeiaso.net/file/christine-static/${path}`} alt={desc} />
      {desc && <figcaption>{desc}</figcaption>}
    </figure>
  );
}
