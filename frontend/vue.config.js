const StylelintWebpackPlugin = require('stylelint-webpack-plugin');

module.exports = {
  css: {
    sourceMap: true,
  },
  configureWebpack: {
    plugins: [
      new StylelintWebpackPlugin({
        files: '**/*.{vue,scss}',
        // failOnError: true,
      }),
    ],
  },
};
