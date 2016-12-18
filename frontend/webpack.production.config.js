var path = require('path');
var webpack = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var webpackUglifyJsPlugin = require('webpack-uglify-js-plugin');
var FaviconsWebpackPlugin = require('favicons-webpack-plugin');

module.exports = {
  entry: [ path.join(__dirname, 'support/index.js') ],
  output: {
    path: path.resolve('./static/dist'),
    filename: '[name]-[hash].min.js',
    publicPath: '/dist/'
  },
  module: {
    loaders: [
      {
        test: /\.purs$/,
        loader: 'purs-loader',
        exclude: /node_modules/,
        query: {
          psc: 'psa',
          bundle: true,
          warnings: false
        }
      }
    ],
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify('production')
    }),
    new webpack.optimize.OccurrenceOrderPlugin(true),
    new webpack.LoaderOptionsPlugin({
      minimize: true,
      debug: false
    }),
    new HtmlWebpackPlugin({
      template: 'support/index.html',
      inject: 'body',
      filename: 'index.html'
    }),
    new FaviconsWebpackPlugin('../static/img/avatar.png'),
    new webpack.optimize.DedupePlugin(),
    new webpackUglifyJsPlugin({
      cacheFolder: path.resolve(__dirname, 'output/ug'),
      debug: false,
      minimize: true,
      sourceMap: false,
      output: {
        comments: false
    },
      compressor: {
        warnings: false
      }
    }),
  ],
  resolveLoader: {
    modules: [
      path.join(__dirname, 'node_modules')
    ]
  },
  resolve: {
    modules: [
      'node_modules',
      'bower_components'
    ],
    extensions: ['.js', '.purs']
  }
};
