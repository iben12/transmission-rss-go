"use strict";

const path = require("path")
const ExtractTextPlugin = require("extract-text-webpack-plugin");
const html = new ExtractTextPlugin("dummy.html");
const sass = new ExtractTextPlugin("style.css");

module.exports = {
  entry: ["./client/js/app.js"],
  output: {
    path: __dirname + "/client/build",
    filename: "app.js",
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: "babel-loader",
          options: {
            presets: [
              [
                "@babel/preset-env",
                {
                  corejs: "2",
                  useBuiltIns: "entry",
                },
              ],
            ],
          },
        },
      },
      {
        test: /\.scss$/,
        use: sass.extract({
          fallback: "style-loader",
          use: [
            {
              loader: "css-loader",
              options: {
                sourceMap: true,
                minimize: true,
              },
            },
            "sass-loader",
          ],
        }),
      },
      {
        test: /\.html$/,
        use: html.extract({ use: "raw-loader" }),
      },
    ],
  },
  plugins: [sass],
  resolve: {
    alias: {
      vue$: "vue/dist/vue.common.js",
    },
  },
  devtool: "source-map",
  devServer: {
    contentBase: path.join(__dirname, 'client/public'),
  }
};
