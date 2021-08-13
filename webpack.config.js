"use strict";

const ExtractTextPlugin = require("extract-text-webpack-plugin");

const sass = new ExtractTextPlugin("style.css");

module.exports = {
  entry: ["./client/js/app.js"],
  output: {
    path: __dirname + "/client/build",
    filename: "app.js"
  },
  optimization: {
    usedExports: true
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
                  useBuiltIns: "entry"
                }
              ]
            ]
          }
        }
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
                minimize: true
              }
            },
            "sass-loader"
          ]
        })
      }
    ]
  },
  plugins: [sass],
  resolve: {
    alias: {
      vue$: "vue/dist/vue.min.js"
    }
  },
  // devtool: "source-map",
  watch: false
};
