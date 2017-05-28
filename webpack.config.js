"use strict";

let path = require('path');
let webpack = require('webpack');
let precss = require('precss');
let autoprefixer = require('autoprefixer');
let ExtractTextPlugin = require("extract-text-webpack-plugin");

let postCssLoader = [
  'css-loader?module',
  '&localIdentName=[name]__[local]___[hash:base64:5]',
  '&disableStructuralMinification',
  '!postcss-loader'
];

let plugins = [
    new webpack.NoErrorsPlugin(),
    new webpack.optimize.DedupePlugin(),
    new ExtractTextPlugin('bundle.css')
];

module.exports = {
    entry:  {
        bundle: path.join(__dirname, './client/routes.js')
    },
    output: {
        path:     './server/static/build',
        publicPath: "/static/build/",
        filename: 'bundle.js',
    },
    plugins: plugins,
    module: {
        loaders: [
            {
                test: /\.jsx?$/,
                include: path.join(__dirname, './client'),
                exclude: /(node_modules|bower_components)/,
                loaders: ['babel']
            },
            {
                test: /\.css/,
                loader: ExtractTextPlugin.extract('style-loader', postCssLoader.join(''))
            },
            {
                test: /\.(png|gif|jpg)$/,
                loader: 'url-loader?name=[name]@[hash].[ext]&limit=5000'
            }
        ]
    },
    resolve: {
        extensions: ['', '.js', '.jsx', '.css'],
        alias: {
            '#components': path.join(__dirname, './client/components'),
            '#containers': path.join(__dirname, './client/containers')
        }
    },
    watch: true
};
