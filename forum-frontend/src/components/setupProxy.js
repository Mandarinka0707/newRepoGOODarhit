// // In setupProxy.js
// const { createProxyMiddleware } = require('http-proxy-middleware');

// module.exports = function(app) {
//   app.use(
//     '/api/auth',
//     createProxyMiddleware({
//       target: 'http://localhost:8080',
//       changeOrigin: true,
//       pathRewrite: {
//         '^/api/auth': '',
//       },
//     })
//   );

//   app.use(
//     '/api/forum',
//     createProxyMiddleware({
//       target: 'http://localhost:3000',
//       changeOrigin: true,
//       pathRewrite: {
//         '^/api/forum': '',
//       },
//     })
//   );
// };
