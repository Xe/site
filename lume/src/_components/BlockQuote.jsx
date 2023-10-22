export default function BlockQuote({ children }) {
  return (
    <div className="mx-auto mt-4 mb-2 rounded-lg bg-bg-2 p-4 dark:bg-bgDark-2 md:max-w-lg xe-dont-newline">
      &gt; {children}
    </div>
  );
};
