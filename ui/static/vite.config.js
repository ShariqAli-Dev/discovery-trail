import { defineConfig } from "vite";
import fs from "fs";
import path from "path";

function getRouteEntries() {
  const routesDirectory = path.resolve(__dirname, "src/routes");
  const entryPoints = {};
  const routeFiles = fs.readdirSync(routesDirectory);

  routeFiles.forEach(
    (
      /**  @type {string} */
      fileName
    ) => {
      const fileAttributes = path.parse(fileName);
      entryPoints[fileAttributes.name] = `src/routes/${fileAttributes.base}`;
    }
  );

  return entryPoints;
}

export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        index: "src/index.js",
        ...getRouteEntries(),
      },
    },
  },
});
