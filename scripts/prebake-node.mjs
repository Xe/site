import { execaCommand } from "execa";

if (process.argv.length === 2) {
  console.error(`usage: node prebake-node.js <path> <video segment count>`);
  process.exit(1);
}

const [_node, _script, basePath, countStr] = process.argv;

const instances = await (async () => {
  try {
    const { stdout } = await execaCommand(
      "flyctl machines list -a xedn --json --state started | jq [.[].ID]",
      {
        shell: true,
      },
    );
    const instances = JSON.parse(stdout);
    return instances;
  } catch (error) {
    console.error(error);
  }
})();

for (const i of Array(parseInt(countStr) + 1).keys()) {
  try {
    const reqs = instances.map((x) =>
      fetch(
        `https://cdn.xeiaso.net/file/christine-static/${basePath}/index${i}.ts`,
        {
          headers: {
            "fly-force-instance": x,
          },
        },
      ).then((resp) => {
        console.log(`${x}: response get: ${resp.status}`);
        return resp;
      })
    );

    const resps = await Promise.all(reqs);
    resps.forEach(async (resp) => {
      if (resp.status !== 200) {
        console.error(
          `failure: ${resp.url}: request ID: ${
            resp.headers["fly-request-id"]
          }, body: ${await resp.text()}`,
        );
      }

      console.log(`ok: ${resp.headers["fly-request-id"]}`);
    });
  } catch (e) {
    console.error(e);
  }
}
