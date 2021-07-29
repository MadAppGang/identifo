const webpack = require('webpack');
const path = require('path');
const HtmlPlugin = require('html-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');

const dotenv = require('dotenv')
  .config({ path: path.resolve(__dirname, '../../.env') });

const env = {
  ...dotenv.parsed,
  MOCK_API: process.env.MOCK_API,
  BASE_URL: process.env.BASE_URL,
  ASSETS_PATH: process.env.ASSETS_PATH,
};

module.exports = {
  mode: 'development',

  entry: {
    polyfill: '@babel/polyfill',
    bundle: './src/index.js',
  },

  output: {
    filename: '[name].[hash].bundle.js',
    publicPath: env.ASSETS_PATH || '/',
  },

  devtool: 'inline-source-map',

  resolve: {
    alias: {
      '~': path.resolve(__dirname, '../../src'),
    },
  },

  module: {
    rules: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        use: 'babel-loader',
      },
      {
        test: /\.css/,
        use: ['style-loader', 'css-loader'],
      },
      {
        test: /\.(eot|woff|woff2|ttf|png|jpg|gif)$/,
        use: 'url-loader?limit=30000&name=[path][name].[ext]',
      },
      {
        test: /\.svg$/,
        use: ['babel-loader', 'react-svg-loader'],
      },
    ],
  },

  plugins: [
    new HtmlPlugin({
      title: 'Identifo Admin',
      template: path.resolve(__dirname, '../..', 'index.template.html'),
    }),
    new CopyPlugin([{
      from: 'src/config.json',
    }]),
    new webpack.DefinePlugin({
      'process.env': JSON.stringify(env),
    }),
  ],

  devServer: {
    historyApiFallback: true,
    contentBase: path.join(__dirname, '../../build'),
    port: 3000,
    hot: true,
  },
};
