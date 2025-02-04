export default function Reflection({ timeTaken, children, open = false }) {
  return (
    <details className="border border-gray-700 rounded bg-gray-800 my-2 transition-all duration-200 ease-linear open:bg-gray-700" open={open}>
      <summary className="cursor-pointer flex justify-between items-center p-2">
        Thought for {(timeTaken || 0).toFixed(0)} seconds
        <span className="ml-2 inline-block transform transition-transform duration-200 ease-out group-open:rotate-180">
          ▼
        </span>
      </summary>
      <div className="p-2 border-t border-gray-700 text-gray-100 text-sm leading-6">
        {children}
      </div>
    </details>
  );
}