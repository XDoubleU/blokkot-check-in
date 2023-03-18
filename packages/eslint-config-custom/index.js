const { commonNamingConvention } = require("./shared");

module.exports = {
  root: true,
  parser: "@typescript-eslint/parser",
  plugins: [
    "@typescript-eslint",
    "import",
    "sonarjs",
    "redundant-undefined"
  ],
  extends: [
    "plugin:prettier/recommended",
    "eslint:recommended",
    "turbo",
    "plugin:@typescript-eslint/recommended",
    "plugin:@typescript-eslint/recommended-requiring-type-checking",
    "plugin:@typescript-eslint/strict",
    "plugin:import/typescript",
    "plugin:sonarjs/recommended"
  ],
  ignorePatterns: [
    ".eslintrc.js",
    "**/dist/**",
    "*.config.*"
  ],
  rules: {
    "@typescript-eslint/explicit-function-return-type": "error",
    "@typescript-eslint/explicit-member-accessibility": "error",
    "@typescript-eslint/member-ordering": "error",
    "@typescript-eslint/no-require-imports": "error",
    "@typescript-eslint/parameter-properties": "error",
    "@typescript-eslint/prefer-readonly": "error",
    "max-lines-per-function": "error",
    "no-duplicate-imports": "error",
    "no-warning-comments": "error",
    "redundant-undefined/redundant-undefined": "error",
    "complexity": "error",
    "@typescript-eslint/no-unused-vars": [
      "error",
      {
        "ignoreRestSiblings": true
      }
    ],
    "@typescript-eslint/naming-convention": [
      "error",
      ...commonNamingConvention
    ],
    "@typescript-eslint/consistent-type-imports": [
      "error",
      {
        "fixStyle": "inline-type-imports"
      }
    ],
    "import/consistent-type-specifier-style": [
      "error", 
      "prefer-inline"
    ],
    "prettier/prettier": [
      "error",
      {
        "semi": false,
        "trailingComma": "none",
        "parser": "typescript",
        
      }
    ]
  },
  settings: {
    "import/resolver": {
      "typescript": {
        "alwaysTryTypes": true,
        "project": [
          "apps/*/tsconfig.json",
          "packages/*/tsconfig.json"
        ]
      }
    },
  }
}