const dotenv = require("dotenv-webpack");
const path = require("path");
const webpack = require("webpack");

module.exports = {
  mode: "development",
  entry: ["./src/js/index.js", "./src/js/attention.js"],
  output: {
    path: path.join(__dirname, "static"),
    filename: "bundle.js",
  },
  plugins: [new dotenv()],
};
