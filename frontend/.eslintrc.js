module.exports = {
  env: {
    browser: true,
    es6: true,
    jest: true,
  },
  extends: [
    'airbnb',
    'plugin:jsx-a11y/recommended',
    'prettier',
    'prettier/react'
  ],
  globals: {
    Atomics: 'readonly',
    SharedArrayBuffer: 'readonly',
  },
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: 2018,
    sourceType: 'module',
  },
  plugins: [
    'react',
    'jsx-a11y',
    'prettier'
  ],
  rules: {
    'semi': 1,
    'react/jsx-filename-extension': [1, { 'extensions': ['.js', '.jsx'] }],
    'prettier/prettier': [
      'error', {
        'semi': true
      }
    ]
  },
};
