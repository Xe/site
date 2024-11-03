export default function PullQuote({ children }) {
  return (
    <blockquote className="border-l-4 border-blue-500 dark:border-blue-300 p-4 my-6 italic bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-gray-100 text-lg">
      <p className="m-0">{children}</p>
    </blockquote>
  );
}
