let path = require("path");
let webpack = require("webpack");

module.exports = {
  mode: "development",
  entry: ["./src/js/index"],
  output: {
    path: path.join(__dirname, "static"),
    filename: "bundle.js",
  },
};
