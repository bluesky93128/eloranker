const StylelintWebpackPlugin = require('stylelint-webpack-plugin');
const FontelloWebpackPlugin = require('fontello-webpack-plugin');

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
      new FontelloWebpackPlugin({
        config: require('./fontello.config.json'),
        fonts: ['woff', 'woff2', 'ttf'],
        output: {
          css: 'css/[name].css',
          font: 'fonts/[name].[ext]',
        },
      }),
    ],
  },
};
