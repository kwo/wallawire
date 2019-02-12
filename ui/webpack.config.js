const fs = require('fs');
const path = require('path');
const CleanWebpackPlugin = require('clean-webpack-plugin');
const GitRevisionPlugin = require('git-revision-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const InterpolateHtmlPlugin = require('interpolate-html-plugin');
const TSLintPlugin = require('tslint-webpack-plugin');

const getConfig = function (e, env) {

    const paths = {
        SRC: path.resolve(__dirname, 'src'),
        DST: path.resolve(__dirname, 'public'),
        FONTS: 'fonts/',
        CERTS: 'certs/',
    };

    const config = {

        bail: true,
        devtool: 'source-map',

        entry: {
            app: path.join(paths.SRC, 'index.tsx'),
        },
        output: {
            path: paths.DST,
            publicPath: '/',
            filename: '[hash:8].js',
            chunkFilename: '[chunkhash:8].js',
        },

        optimization: {
            splitChunks: {
                cacheGroups: {
                    vendors: {
                        test: /[\\/]node_modules[\\/]/,
                        chunks: 'all',
                        priority: 1
                    },
                },
            },
        },

        resolve: {
            // Add '.ts' and '.tsx' as resolvable extensions.
            extensions: [".ts", ".tsx", ".js", ".json"]
        },

        module: {
            rules: [
                // All files with a '.ts' or '.tsx' extension will be handled by 'awesome-typescript-loader'.
                {
                    test: /\.tsx?$/,
                    loader: "awesome-typescript-loader"
                },

                {
                    test: /\.css$/,
                    use: ['style-loader', 'css-loader']
                },

                {
                    test: /\.(woff(2)?|ttf|eot|svg)(\?v=\d+\.\d+\.\d+)?$/,
                    use: [{
                        loader: 'file-loader',
                        options: {
                            name: '[name].[ext]',
                            outputPath: paths.FONTS
                        }
                    }]
                }
            ]
        },

        plugins: [
            new CleanWebpackPlugin([paths.DST]),
            new TSLintPlugin({
                files: [paths.SRC + '/**/*.ts', paths.SRC + '/**/*.tsx'],

            }),
            new HtmlWebpackPlugin({
                template: path.join(paths.SRC, 'index.html'),
            }),
            new InterpolateHtmlPlugin({
                'VERSION': new GitRevisionPlugin({
                    versionCommand: 'describe --always --tags --dirty="*"'
                }).version(),
            })
        ]
    };

    if (env.mode !== 'production') {
        config.devServer = {
            host: 'localhost',
            port: 4444,
            contentBase: 'public/',
            compress: true,
            historyApiFallback: true,
            stats: 'minimal',
            https: {
                key: fs.readFileSync(`${paths.CERTS}/server.key`),
                cert: fs.readFileSync(`${paths.CERTS}/server.crt`),
                ca: fs.readFileSync(`${paths.CERTS}/ca.crt`),
            },
            before: (app) => {
                const morgan = require('morgan');
                app.use(morgan('dev'));
            },
            proxy: [
                {
                    context: ['/api/**'],
                    target: 'https://localhost:8888',
                    secure: false
                }
            ]
        };
    }

    return config;

};

module.exports = getConfig;
