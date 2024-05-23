const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require("html-webpack-plugin");
const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin');
//const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const { NODE_ENV, CI, WEBPACK_SERVE } = process.env;

const isDevelopment = WEBPACK_SERVE === 'true' && CI == null;
const isProduction = NODE_ENV === 'production';
const isPuppeteer = NODE_ENV === 'puppeteer';

const useReactRefresh = isDevelopment && !isPuppeteer;

function employCache(loaders) {
  if (isDevelopment && !isPuppeteer) {
    return [
      {
        loader: 'cache-loader',
        options: {
          cacheDirectory: path.join(__dirname, '..', '.cache-loader'),
        },
      },
      ...loaders,
    ];
  }

  return loaders;
}


const fontPublicPath = process.env.NODE_ENV === 'production' ? '/ui/' : '';
module.exports = {
  entry: './src/main.tsx',
  output: {
    path: path.resolve(__dirname, './dist'),
    publicPath: '/',
    filename: 'build.js'
  },
  mode: process.env.NODE_ENV || "development",
  resolve: {
    extensions: [".tsx", ".ts", ".js"],
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: ["babel-loader"],
      },
      {
        test: /\.(ts|tsx)$/,
        exclude: /node-modules/,
        use: ["ts-loader"],
      },
      {
        test: /\.(css)$/i,
        use: employCache(['style-loader', 'css-loader']),
      },
      {
        test: /\.(scss_disabled)$/i,
        use: employCache([
          {
            loader: 'style-loader',
            options: {
              injectType: 'lazySingletonStyleTag',
              insert: 'meta[name="sass-styles-compiled"]',
            },
          },
          'css-loader',
          'postcss-loader',
          'sass-loader',
        ]),
      },
      {
        test: /\.(scss)$/i,
        use: employCache([
          'style-loader',
          'css-loader',
          'postcss-loader',
          'sass-loader',
        ]),
      },
      {
        test: /\.(png|jpg|gif|svg|ttf|woff|woff2|eot)\w*/,
        loader: 'file-loader',
        options: {
          publicPath: fontPublicPath,
          name: '[name].[ext]?[hash]'
        }
      },
    ],
  },
  plugins: [
    new HtmlWebpackPlugin({ template: path.join(__dirname, "src", "index.html") }),
    new MonacoWebpackPlugin(),
  ],
  //devServer: {
    ////contentBase: path.join(__dirname, "src")
    //static: path.join(__dirname, "src")
  //},
  devServer: {
    proxy: {
      '/v1': {
        target: 'http://127.0.0.1:7079',
        secure: false
      }
    },
    historyApiFallback: true,
    //noInfo: true
  },
  performance: {
    hints: false
  },
  //devtool: '#eval-source-map'
  devtool: 'eval-source-map'
};

//if (process.env.NODE_ENV === 'production') {
//  module.exports.devtool = '#source-map'
//  // http://vue-loader.vuejs.org/en/workflow/production.html
//  module.exports.plugins = (module.exports.plugins || []).concat([
//    new webpack.DefinePlugin({
//      'process.env': {
//        NODE_ENV: '"production"'
//      }
//    }),
//    new webpack.optimize.UglifyJsPlugin({
//      sourceMap: true,
//      compress: {
//        warnings: false
//      }
//    }),
//    new webpack.LoaderOptionsPlugin({
//      minimize: true
//    })
//  ])
//}
