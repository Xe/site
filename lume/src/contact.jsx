export const title = "Contact";
export const layout = "base.njk";

export default ({ contactLinks }) => (
  <>
    <h1 className="text-3xl mb-4">Contact</h1>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <h2 className="text-2xl mb-2">Email</h2>
        <p className="mb-4">
          <a href="mailto:me@xeiaso.net">me@xeiaso.net</a>
        </p>
        <p className="mb-4">Other useful links:</p>
        <ul className="list-disc list-inside">
          {contactLinks.map(({ title, url }) => (
            <li key={url}>
              <a href={url}>{title}</a>
            </li>
          ))}
        </ul>
      </div>
      <div>
        <h2 className="text-2xl mb-2">Discord</h2>
        <p className="mb-4">
          <code className="p-1 bg-bg-soft dark:bg-bgDark-soft m-1">xeiaso</code> or{" "}
          <code className="p-1 bg-bg-soft dark:bg-bgDark-soft m-1">Cadey~#1337</code>
        </p>
        <p className="mb-4">
          Please note that Discord may reject friend requests if you aren't in a
          mutual server with me. I don't have control over this behavior.
        </p>
      </div>
      <div>
        <h2 className="text-2xl mb-2">Cryptocurrency</h2>
        <p className="mb-4">
          I accept cryptocurrency donations at the following addresses:
        </p>
        <ul className="list-disc list-inside">
          <li>
            Ethereum:{" "}
            <code className="bg-bg-soft dark:bg-bgDark-soft m-1">xeiaso.eth</code> or{" "}
            <code className="bg-bg-soft dark:bg-bgDark-soft m-1">
              0xeA223Ca8968Ca59e0Bc79Ba331c2F6f636A3fB82
            </code>
          </li>
          <li>
            Bitcoin:{" "}
            <code className="bg-bg-soft dark:bg-bgDark-soft m-1">
              bc1qw0pa3zdus94nyehmys6g8td2xfaqtl9pmuv564
            </code>
          </li>
        </ul>
      </div>
    </div>
  </>
);
