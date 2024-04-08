module.exports = {
  root: true,
  env: { browser: true, es2020: true },
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:react-hooks/recommended',
    'prettier',
  ],
  overrides: [
    {
      extends: ['plugin:@typescript-eslint/strict-type-checked'],
      files: ['./**/*.{ts,tsx}'],
    },
    {
      extends: ['plugin:@typescript-eslint/disable-type-checked'],
      files: ['./vite.config.ts'],
    }
  ],
  ignorePatterns: ['dist', 'gen', '.eslintrc.cjs'],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
    project: ['./tsconfig.json', './tsconfig.node.json'],
    tsconfigRootDir: __dirname,
  },
  plugins: ['react-refresh'],
  rules: {
    'react-refresh/only-export-components': [
      'warn',
      { allowConstantExport: true },
    ],
  },
}
