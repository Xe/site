export const title = "Patrons";
export const layout = "base.njk";

const PatronCard = ({ full_name, image_url }) => (
  <div className="bg-bg-1 dark:bg-bgDark-1 rounded-xl m-2 shadow-md w-xs">
    <div className="items-center text-lg p-2">
      {full_name}
    </div>
    <div className="flex items-center justify-center p-2">
      <img className="rounded-xl w-24 h-24" src={image_url} alt={full_name} />
    </div>
  </div>
);

export default ({ patrons }) => {
  const users = patrons.included.Items
    .filter((item) => item.type === "user")
    .map((item) => item.attributes)
    .filter((item) => item.full_name !== "Xe");
  return (
    <>
      <h1 className="text-3xl mb-4">Patrons</h1>
      <p className="mb-4">
        This page is a list of all of my patrons. Thank you all so much for your
        support! If you want to get on this list, please consider{" "}
        <a href="https://patreon.com/cadey">becoming a patron</a>! This page
        will update whenever something changes on Patreon.
      </p>

      <div className="flex flex-wrap items-start justify-center p-2">
        {users.map((user) => <PatronCard key={user.full_name} {...user} />)}
      </div>

      {users.length === 0 && (
        <>
          <p className="text-center">
            No patrons yet!
          </p>
          <p className="text-center">
            <a href="https://patreon.com/cadey">Become a patron</a>{" "}
            to get on this list!
          </p>
        </>
      )}
    </>
  );
};
