const path = require("path");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");

module.exports = {
  target: "web",
  entry: './src/index.ts',
  // mode: "development",
  // devtool: "source-map", // source map for develop
  output: {
    path: path.resolve(__dirname, "./dist"),
    filename: "lib/identifo-js.js",
    library: "IdentifoJs",
    libraryTarget: "umd",
    globalObject: "this",
    umdNamedDefine: true,
  },
  watchOptions: {
    aggregateTimeout: 600,
    ignored: /node_modules/,
  },
  plugins: [
    new CleanWebpackPlugin({
      cleanStaleWebpackAssets: false,
      cleanOnceBeforeBuildPatterns: [path.resolve(__dirname, "./dist")],
    }),
  ],
  module: {
    rules: [
      {
        test: /\.ts(x?)$/,
        exclude: /node_modules/,
        use: [
          {
            loader: "ts-loader",
            options: {
              configFile: "tsconfig.lib.json"
            }
          },
        ],
      },
    ],
  },
  resolve: {
    extensions: [".tsx", ".ts", ".js"],
  },
};
